// Code generated for package bindata by go-bindata DO NOT EDIT. (@generated)
// sources:
// docker/Dockerfile
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

	info := bindataFileInfo{name: "docker/Dockerfile", size: 365, mode: os.FileMode(420), modTime: time.Unix(1679074884, 0)}
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
	"docker/Dockerfile": dockerDockerfile,
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
