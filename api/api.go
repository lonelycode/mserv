// Package api provides handlers for mserv's various endpoints.
package api

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/graymeta/stow"
	"github.com/graymeta/stow/local"
	"github.com/graymeta/stow/s3"
	uuid "github.com/satori/go.uuid"

	"github.com/TykTechnologies/mserv/bundle"
	config "github.com/TykTechnologies/mserv/conf"
	coprocess "github.com/TykTechnologies/mserv/coprocess/bindings/go"
	"github.com/TykTechnologies/mserv/storage"
	"github.com/TykTechnologies/mserv/util/logger"
	"github.com/TykTechnologies/mserv/util/storage/errors"
)

const (
	// FmtPluginContainer is a format string for the layout of the container names.
	FmtPluginContainer = "mserv-plugin-%s"

	moduleName = "mserv.api"
)

var log = logger.GetLogger(moduleName)

func NewAPI(store storage.MservStore) *API {
	return &API{store: store}
}

type API struct {
	store storage.MservStore
}

func (a *API) HandleUpdateBundle(filePath string, bundleName string) (*storage.MW, error) {
	mw, err := a.store.GetMWByID(bundleName)
	if err != nil {
		return nil, err
	}

	err = a.store.DeleteMW(mw.UID)
	if err != nil {
		return nil, err
	}

	return a.HandleNewBundle(filePath, mw.APIID, bundleName)
}

func (a *API) HandleDeleteBundle(bundleName string) error {
	mw, err := a.store.GetMWByID(bundleName)
	if err != nil {
		return err
	}

	fStore, err := GetFileStore()
	if err != nil {
		log.WithError(err).Error("failed to get file handle")

		return err
	}

	defer func() {
		if errFC := fStore.Close(); errFC != nil {
			log.WithError(errFC).Error("error while closing file store")
		}
	}()

	pluginContainerID := fmt.Sprintf(FmtPluginContainer, bundleName)

	fCont, err := fStore.Container(pluginContainerID)
	if err != nil {
		return fmt.Errorf("could not get container: %w", err)
	}

	if errWalk := stow.Walk(fCont, "", 100, func(i stow.Item, e error) error {
		if e != nil {
			return fmt.Errorf("error getting item while walking container: %w", e)
		}

		return fCont.RemoveItem(i.ID())
	}); errWalk != nil {
		return fmt.Errorf("error while walking container to delete contents: %w", errWalk)
	}

	// HACK: workaround for https://github.com/graymeta/stow/issues/239 - vvv
	//
	// (stow.Location).RemoveContainer doesn't currently take the full path into account for Kind "local".
	// It merely calls "os.RemoveAll" with the _relative_ path, so we need to change to the parent path, and then defer
	// changing back until after the misbehaving RemoveContainer call.
	//
	// Maybe swap out Stow for the Go CDK one day? https://gocloud.dev/howto/blob/

	fsCfg := config.GetConf().Mserv.FileStore

	if fsCfg.Kind == local.Kind {
		prevWD, errWD := os.Getwd()
		if errWD != nil {
			return fmt.Errorf("could not get current working directory: %w", errWD)
		}

		if errCD := os.Chdir(fsCfg.Local.ConfigKeyPath); errCD != nil {
			return fmt.Errorf("could not change current working directory: %w", errCD)
		}

		defer func() {
			if errPD := os.Chdir(prevWD); errPD != nil {
				log.WithError(errPD).WithField("dir", prevWD).Error("could not revert to previous working directory")
			}
		}()
	}

	// HACK: workaround for https://github.com/graymeta/stow/issues/239 - ^^^

	if errRC := fStore.RemoveContainer(pluginContainerID); errRC != nil {
		return fmt.Errorf("could not remove container '%s': %w", pluginContainerID, errRC)
	}

	return a.store.DeleteMW(mw.UID)
}

func (a *API) HandleNewBundle(filePath string, apiID, bundleName string) (*storage.MW, error) {
	// Read the zip file raw data
	bData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	log.WithField("path", filePath).Info("read bundle")

	// Create a bundle object and provide a name
	bdl := &bundle.Bundle{
		Data: bData,
		Name: bundleName,
	}

	// Unzip and verify the bundle
	err = bundle.SaveBundleZip(bdl, apiID, bundleName)
	if err != nil {
		return nil, err
	}

	log.WithField("bundle-path", bdl.Path).Info("saved zip")

	// create DB record of the bundle
	mw := &storage.MW{
		UID:      bdl.Name,
		APIID:    apiID,
		Manifest: &bdl.Manifest,
		Active:   true,
		Added:    time.Now(),
	}

	if len(bdl.Manifest.FileList) != 1 {
		return nil, errors.New("only one plugin file file allowed per bundle")
	}

	log.Info("attempting to get file handle")

	// upload
	fStore, err := GetFileStore()
	if err != nil {
		log.WithError(err).Error("failed to get file handle")
		return nil, err
	}
	defer fStore.Close()

	log.Info("file store handle opened")

	fName := bdl.Manifest.FileList[0]
	pluginPath := path.Join(bdl.Path, fName)

	log.Info("storing bundle in asset repo")

	pluginContainerID := fmt.Sprintf(FmtPluginContainer, bundleName)
	fCont, getErr := fStore.Container(pluginContainerID)
	if getErr != nil {
		log.WithField("container-id", pluginContainerID).Warning("container not found, creating")
		fCont, err = fStore.CreateContainer(pluginContainerID)
		if err != nil {
			return nil, fmt.Errorf("couldn't fetch container: %s, couldn't create container: %s", getErr.Error(), err.Error())
		}
	}

	f, err := os.Open(pluginPath)
	if err != nil {
		return nil, err
	}

	fInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(f)

	data, err := fCont.Put(fInfo.Name(), r, fInfo.Size(), nil)
	if err != nil {
		return nil, err
	}

	// This is an internal URL, must be interpreted by Stow
	ref := data.URL().String()

	// Store the bundle zip file too, because we can use it again
	bF, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	bfInfo, err := bF.Stat()
	if err != nil {
		return nil, err
	}

	bundleData, err := fCont.Put(bfInfo.Name(), bufio.NewReader(bF), bfInfo.Size(), nil)
	if err != nil {
		return nil, err
	}

	// This is an internal URL, must be interpreted by Stow
	mw.BundleRef = bundleData.URL().String()

	log.Info("completed storage")

	for _, f := range bdl.Manifest.CustomMiddleware.Pre {
		p := &storage.Plugin{
			UID:      uuid.NewV4().String(),
			FileName: fName,
			FileRef:  ref,
			Name:     f.Name,
			Type:     coprocess.HookType_Pre,
		}

		mw.Plugins = append(mw.Plugins, p)
	}

	for _, f := range bdl.Manifest.CustomMiddleware.Post {
		p := &storage.Plugin{
			UID:      uuid.NewV4().String(),
			FileName: fName,
			FileRef:  ref,
			Name:     f.Name,
			Type:     coprocess.HookType_Post,
		}

		mw.Plugins = append(mw.Plugins, p)
	}

	for _, f := range bdl.Manifest.CustomMiddleware.PostKeyAuth {
		p := &storage.Plugin{
			UID:      uuid.NewV4().String(),
			FileName: fName,
			FileRef:  ref,
			Name:     f.Name,
			Type:     coprocess.HookType_PostKeyAuth,
		}

		mw.Plugins = append(mw.Plugins, p)
	}

	if bdl.Manifest.CustomMiddleware.AuthCheck.Name != "" {
		p := &storage.Plugin{
			UID:      uuid.NewV4().String(),
			FileName: fName,
			FileRef:  ref,
			Name:     bdl.Manifest.CustomMiddleware.AuthCheck.Name,
			Type:     coprocess.HookType_CustomKeyCheck,
		}

		mw.Plugins = append(mw.Plugins, p)
	}

	log.Warning("not loading into dispatcher")
	// a.LoadMWIntoDispatcher(mw, bdl.Path)

	// store in mongo
	_, err = a.store.CreateMW(mw)
	if err != nil {
		return mw, err
	}

	// clean up
	if err := os.Remove(filePath); err != nil {
		return nil, err
	}

	if !config.GetConf().Mserv.RetainUploads {
		if err := os.RemoveAll(bdl.Path); err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("could not clean up uploaded bundle: %w", err)
		}
	}

	return mw, nil
}

// Will only store the bundle file into our store so we can pull it from a gateway if necessary
func (a *API) StoreBundleOnly(filePath string, apiID, bundleName string) (*storage.MW, error) {
	// create DB record of the bundle
	mw := &storage.MW{
		UID:          bundleName,
		APIID:        apiID,
		Active:       true,
		Added:        time.Now(),
		DownloadOnly: true,
	}

	log.Info("attempting to get file handle")

	// upload
	fStore, err := GetFileStore()
	if err != nil {
		log.WithError(err).Error("failed to get file handle")
		return nil, err
	}

	defer fStore.Close()

	log.Info("file store handle opened, storing bundle in asset repo")

	pluginContainerID := fmt.Sprintf(FmtPluginContainer, bundleName)
	fCont, getErr := fStore.Container(pluginContainerID)
	if getErr != nil {
		log.WithField("container-id", pluginContainerID).Warning("container not found, creating")
		fCont, err = fStore.CreateContainer(pluginContainerID)
		if err != nil {
			return nil, fmt.Errorf("couldn't fetch container: %s, couldn't create container: %s", getErr.Error(), err.Error())
		}
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	fInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(f)

	data, err := fCont.Put(fInfo.Name(), r, fInfo.Size(), nil)
	if err != nil {
		return nil, err
	}

	// This is an internal URL, must be interpreted by Stow
	mw.BundleRef = data.URL().String()
	log.Info("completed storage")

	// store in mongo
	_, err = a.store.CreateMW(mw)
	if err != nil {
		return mw, err
	}

	// clean up
	if err := os.Remove(filePath); err != nil {
		return nil, err
	}

	return mw, nil
}

func (a *API) GetMWByID(id string) (*storage.MW, error) {
	return a.store.GetMWByID(id)
}

func (a *API) GetAllActiveMW() ([]*storage.MW, error) {
	return a.store.GetAllActive()
}

func (a *API) LoadMWIntoDispatcher(mw *storage.MW, pluginPath string) (*storage.MW, error) {
	for _, plug := range mw.Plugins {
		// Load the plugin function into memory so we can call it
		hFunc, err := LoadPlugin(plug.Name, pluginPath, plug.FileName)
		if err != nil {
			log.Error("failed to load plugin file: ", plug.FileName)
		}

		// Store a reference
		hookKey := storage.GenerateStoreKey(mw.OrgID, mw.APIID, plug.Type.String(), plug.Name)
		updated, err := storage.GlobalRtStore.UpdateOrStoreHook(hookKey, hFunc)
		if err != nil {
			return nil, err
		}

		msg := "added"
		if updated {
			msg = "updated"
		}

		log.Infof("status: %s plugin %s", msg, hookKey)
	}

	return mw, nil
}

func (a *API) FetchAndServeBundleFile(mw *storage.MW) (string, error) {
	location, err := GetFileStore()
	if err != nil {
		return "", err
	}
	defer location.Close()

	bundleDir := path.Join(config.GetConf().Mserv.PluginDir, mw.UID)
	checkSumDir := path.Join(bundleDir, mw.Manifest.Checksum)
	filePath := path.Join(checkSumDir, "bundle.zip")

	log.Info("fetching bundle from storage backend")

	// if file already exist, then nothing has to be done
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		_, bundleErr := os.Stat(bundleDir)
		if bundleErr == nil {
			errRemove := os.RemoveAll(bundleDir)
			if errRemove != nil {
				log.WithError(errRemove).Error("failed to delete old directory")
			}
		}

		createErr := os.MkdirAll(checkSumDir, os.ModePerm)
		if createErr != nil {
			return "", err
		}

		fUrl, err := url.Parse(mw.BundleRef)
		if err != nil {
			return "", err
		}

		item, err := location.ItemByURL(fUrl)

		f, err := os.Create(filePath)
		if err != nil {
			return "", err
		}

		rc, err := item.Open()
		_, err = io.Copy(f, rc)
		if err != nil {
			return "", err
		}
		rc.Close()
	}

	return filePath, nil
}

func GetFileStore() (stow.Location, error) {
	fsCfg := config.GetConf().Mserv.FileStore

	if fsCfg == nil {
		return nil, ErrNoFSConfig
	}

	switch fsCfg.Kind {
	case local.Kind:
		log.WithField("path", fsCfg.Local.ConfigKeyPath).Info("detected local store")

		// Dialling stow/local will fail if the base directory doesn't already exist
		if err := os.MkdirAll(fsCfg.Local.ConfigKeyPath, 0o750); err != nil && !os.IsExist(err) {
			return nil, fmt.Errorf("%w: %s", ErrCreateLocal, fsCfg.Local.ConfigKeyPath)
		}

		return stow.Dial(local.Kind, stow.ConfigMap{
			local.ConfigKeyPath: fsCfg.Local.ConfigKeyPath,
		})
	case s3.Kind:
		log.Info("detected s3 store")

		return stow.Dial(s3.Kind, stow.ConfigMap{
			s3.ConfigAccessKeyID: fsCfg.S3.ConfigAccessKeyID,
			s3.ConfigRegion:      fsCfg.S3.ConfigRegion,
			s3.ConfigSecretKey:   fsCfg.S3.ConfigSecretKey,
		})
	}

	return nil, fmt.Errorf("%w: %s", ErrFSKind, fsCfg.Kind)
}
