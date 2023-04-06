// Code generated for package bindata by go-bindata DO NOT EDIT. (@generated)
// sources:
// docker/Dockerfile
// manifests/knative-app.yaml
package bindata

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _dockerDockerfile = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\x8f\x41\x6b\xf2\x40\x10\x86\xef\xfb\x2b\x86\xf8\x1d\xdd\xe4\x2e\x7c\x87\x54\x53\x91\xaa\x91\x68\x29\xa5\x14\x19\xd7\x21\x59\xba\xd9\x59\x76\x27\x14\x29\xfd\xef\x45\x6d\xed\xa1\xc7\x79\xdf\x79\x86\x67\xee\x9b\x7a\x05\x9e\x8f\x34\x71\x92\x34\xba\x60\x3d\x29\x35\x82\x29\x7b\x4f\x46\x40\x3a\x82\xc3\x60\x9d\x80\x61\x2f\x68\x3d\x45\x10\x06\x84\x48\x81\x93\x15\x8e\x27\x78\xef\xc8\x43\x18\x52\x77\x6e\xce\xc0\xdc\x4a\x37\x1c\x60\x83\xe6\x0d\x5b\x82\x48\xad\x4d\x12\x4f\x63\xc0\x41\xb8\x47\xb1\x06\x9d\x3b\xa9\x11\x6c\x89\xa0\x13\x09\x69\x52\x14\x47\x36\x29\x6f\x2f\x64\x6e\xb8\x2f\xc8\x17\xe1\x7a\x20\x15\x8e\x30\x7a\x7d\x2d\xf5\x2d\x35\x57\x47\xeb\x5b\x8d\xfa\xd7\x47\x0b\x6b\xfc\xd9\x52\x65\x33\x87\x65\x79\x57\x2d\xf7\x8b\x55\x39\xaf\xf6\xdb\xfa\xb1\x99\x56\xea\x12\x01\xc7\x36\xe7\x40\xfe\xf6\x5b\xca\x6d\x8f\x2d\xe5\x89\x87\x68\xe8\xff\xbf\x8f\xbf\xe8\xa7\x52\xe5\x6c\x06\x39\x14\x18\x82\x52\x4f\x75\xf3\x30\x5b\x34\xdf\x53\xb5\xde\x35\xcf\x9b\x7a\xb1\xde\xc1\x0b\x64\x3e\xf4\xd9\x18\xb2\x24\x18\x25\x83\x57\xf5\x15\x00\x00\xff\xff\xa2\xfd\x42\x22\x6d\x01\x00\x00")

func dockerDockerfileBytes() ([]byte, error) {
	return bindataRead(
		_dockerDockerfile,
		"docker/Dockerfile",
	)
}

func dockerDockerfile() (*asset, error) {
	bytes, err := dockerDockerfileBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "docker/Dockerfile", size: 365, mode: os.FileMode(420), modTime: time.Unix(1680689653, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _manifestsKnativeAppYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\xcd\x3f\x8b\x83\x40\x10\x05\xf0\x7e\x3f\xc5\x14\xd7\x9e\x72\xed\x76\x16\x72\x04\x49\xb2\xa8\x04\x52\xc9\xa0\x83\x2c\x66\xff\xb0\x3b\x6c\x13\xfc\xee\x61\x09\x8a\xc9\xab\x86\xf7\xe0\x37\xe8\xf5\x8d\x42\xd4\xce\x4a\x88\x14\x92\xb6\x73\xb1\x58\x64\x9d\xa8\x98\x28\x95\xe9\x4f\x2c\xda\x4e\x12\xba\x3c\x8e\x24\x0c\x31\x4e\xc8\x28\x05\x80\x45\x43\x12\x7e\x9e\x4d\x53\x0d\x95\x52\xc3\xa5\x3a\xd7\xab\x88\x9e\xc6\xbc\x32\x19\xff\x40\xa6\x7c\x03\x6c\x6d\xce\xe8\x2c\xa3\xb6\x14\xe2\xd6\x00\xfc\x82\x36\x38\x7f\x70\x6d\xfd\x7f\xea\xfa\xf6\xbe\x96\xdf\x2f\x60\x8f\x77\x81\x0f\xca\x5b\xda\x7d\xe5\x02\x1f\x45\x75\x6d\xfb\x55\xbc\x02\x00\x00\xff\xff\x77\x33\xff\xda\xf4\x00\x00\x00")

func manifestsKnativeAppYamlBytes() ([]byte, error) {
	return bindataRead(
		_manifestsKnativeAppYaml,
		"manifests/knative-app.yaml",
	)
}

func manifestsKnativeAppYaml() (*asset, error) {
	bytes, err := manifestsKnativeAppYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "manifests/knative-app.yaml", size: 244, mode: os.FileMode(420), modTime: time.Unix(1680689609, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"docker/Dockerfile":          dockerDockerfile,
	"manifests/knative-app.yaml": manifestsKnativeAppYaml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"docker": &bintree{nil, map[string]*bintree{
		"Dockerfile": &bintree{dockerDockerfile, map[string]*bintree{}},
	}},
	"manifests": &bintree{nil, map[string]*bintree{
		"knative-app.yaml": &bintree{manifestsKnativeAppYaml, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
