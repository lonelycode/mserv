package bundle

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/TykTechnologies/goverify"
	"github.com/TykTechnologies/tyk/apidef"

	config "github.com/TykTechnologies/mserv/conf"
	"github.com/TykTechnologies/mserv/util/logger"
)

var (
	moduleName = "mserv.bundle"
	log        = logger.GetLogger(moduleName)
)

// Bundle is the basic bundle data structure, it holds the bundle name and the data.
type Bundle struct {
	Name     string
	Data     []byte
	Path     string
	Manifest apidef.BundleManifest
}

// Verify performs a signature verification on the bundle file.
func (b *Bundle) Verify() error {
	log.Info("verifying bundle: ", b.Name)

	var useSignature bool
	var bundleVerifier goverify.Verifier

	// Perform signature verification if a public key path is set:
	pKeyPath := config.GetConf().Mserv.PublicKeyPath
	if pKeyPath != "" {
		if b.Manifest.Signature == "" {
			// Error: A public key is set, but the bundle isn't signed.
			return errors.New("bundle isn't signed")
		}

		var err error
		bundleVerifier, err = goverify.LoadPublicKeyFromFile(pKeyPath)
		if err != nil {
			return err
		}

		useSignature = true
	}

	var bundleData bytes.Buffer

	for _, f := range b.Manifest.FileList {
		extractedPath := filepath.Join(b.Path, f)

		f, err := os.Open(extractedPath)
		if err != nil {
			return err
		}
		_, err = io.Copy(&bundleData, f)
		f.Close()
		if err != nil {
			return err
		}
	}

	checksum := fmt.Sprintf("%x", md5.Sum(bundleData.Bytes()))

	if checksum != b.Manifest.Checksum {
		return errors.New("invalid checksum")
	}

	if useSignature {
		signed, err := base64.StdEncoding.DecodeString(b.Manifest.Signature)
		if err != nil {
			return err
		}
		if err := bundleVerifier.Verify(bundleData.Bytes(), signed); err != nil {
			return err
		}

	}

	return nil
}

// Getter is used for downloading bundle data, see HttpBundleGetter for reference.
type Getter interface {
	Get() ([]byte, error)
}

// HTTPBundleGetter is a simple HTTP Getter.
type HTTPBundleGetter struct {
	URL string
}

// Get performs an HTTP GET request.
func (g *HTTPBundleGetter) Get() ([]byte, error) {
	resp, err := http.Get(g.URL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		httpError := fmt.Sprintf("HTTP Error, got status code %d", resp.StatusCode)
		return nil, errors.New(httpError)
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// Saver is an interface used by bundle saver structures.
type Saver interface {
	Save(*Bundle, string) error
}

// ZipBundleSaver is a Saver for ZIP files.
type ZipBundleSaver struct{}

// Save implements the main method of the Saver interface. It makes use of archive/zip.
func (ZipBundleSaver) Save(bundle *Bundle, bundlePath string) error {
	buf := bytes.NewReader(bundle.Data)
	reader, err := zip.NewReader(buf, int64(len(bundle.Data)))
	if err != nil {
		return err
	}

	for _, f := range reader.File {
		destPath := filepath.Join(bundlePath, f.Name)

		if f.FileHeader.Mode().IsDir() {
			if err := os.Mkdir(destPath, 0700); err != nil {
				return err
			}
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		newFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		if _, err = io.Copy(newFile, rc); err != nil {
			return err
		}
		rc.Close()
		if err := newFile.Close(); err != nil {
			return err
		}
	}
	return nil
}

// saveBundle will save a bundle to the disk, see ZipBundleSaver methods for reference.
func SaveBundle(bundle *Bundle, destPath string) error {
	bundleFormat := "zip"

	var bundleSaver Saver

	// TODO: use enums?
	switch bundleFormat {
	case "zip":
		bundleSaver = ZipBundleSaver{}
	}
	bundleSaver.Save(bundle, destPath)

	return nil
}

// loadBundleManifest will parse the manifest file and return the bundle parameters.
func LoadBundleManifest(bundle *Bundle, skipVerification bool) error {
	log.Info("loading bundle: ", bundle.Name)

	manifestPath := filepath.Join(bundle.Path, "manifest.json")
	f, err := os.Open(manifestPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&bundle.Manifest); err != nil {
		log.Info("couldn't unmarshal the manifest file for bundle: ", bundle.Name)
		return err
	}

	if skipVerification {
		return nil
	}

	if err := bundle.Verify(); err != nil {
		log.Info("bundle verification failed: ", bundle.Name)
	}

	return nil
}

func CreateBundlePath(bundleName, apiID string) string {
	tykBundlePath := filepath.Join(config.GetConf().Mserv.MiddlewarePath, "plugins")
	bundleNameHash := md5.New()
	io.WriteString(bundleNameHash, bundleName)
	bundlePath := fmt.Sprintf("%s_%x", apiID, bundleNameHash.Sum(nil))
	return filepath.Join(tykBundlePath, bundlePath)
}

func LoadBundle(apiID, bundleName string) (*Bundle, error) {
	destPath := CreateBundlePath(bundleName, apiID)

	_, err := os.Stat(destPath)
	if err != nil {
		return nil, err
	}

	// The bundle exists, load and return:
	log.Info("Loading existing bundle: ", bundleName)

	bundle := &Bundle{
		Name: bundleName,
		Path: destPath,
	}

	err = LoadBundleManifest(bundle, false)
	if err != nil {
		log.Error("couldn't load bundle: ", bundleName, " ", err)
		return nil, err
	}

	log.Info("using bundle: ", bundleName)

	return bundle, nil
}

func SaveBundleZip(bundle *Bundle, apiID, bundleName string) error {
	destPath := CreateBundlePath(bundleName, apiID)
	if err := os.MkdirAll(destPath, 0700); err != nil {
		return errors.New("couldn't create bundle directory")
	}

	if err := SaveBundle(bundle, destPath); err != nil {
		return err
	}

	// Set the destination path
	bundle.Path = destPath

	if loadErr := LoadBundleManifest(bundle, false); loadErr != nil {
		if err := os.RemoveAll(bundle.Path); err != nil {
			return fmt.Errorf("%s, %s", loadErr.Error(), err.Error())
		}

		return loadErr
	}

	log.Info("bundle is valid")
	return nil
}
