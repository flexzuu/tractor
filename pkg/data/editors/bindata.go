// Code generated for package editors by go-bindata DO NOT EDIT. (@generated)
// sources:
// studio/editors/object/index.html
// studio/editors/object/main.css
// studio/editors/object/main.mjs
// studio/editors/table/index.html
// studio/editors/table/main.css
// studio/editors/table/main.mjs
package editors

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

var _studioEditorsObjectIndexHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4c\x8f\xb1\xce\x83\x30\x0c\x84\x77\x9e\xc2\xbf\x1f\xe0\x8f\xd8\x93\x0c\xad\xba\x74\xe9\xd2\x17\xa0\x89\x91\x29\x0e\xa0\x38\x1d\x78\xfb\x0a\x42\x25\x26\xfb\xa4\xfb\xee\x74\xf6\x2f\xce\xa1\xac\x0b\x01\x97\x24\xbe\xb1\xf5\x00\x00\x58\xa6\x2e\xd6\x77\x97\x1a\xf2\xb0\x14\xd8\xcc\x0e\xd3\x1c\x3f\x42\x08\x9a\x83\xc3\xd4\x0d\xd3\x7f\x7a\x2b\x7a\x6b\xaa\xeb\x84\xc9\x30\x8d\x90\x49\x1c\x6a\x59\x85\x94\x89\x0a\x02\x67\xea\x0f\x30\xa8\x22\x98\xa3\xd3\xec\xa5\x00\x55\xbd\xe6\xb8\x82\x2e\x24\x12\x98\xc2\xe8\xb0\xef\x44\x09\x4f\xe9\xdc\xfa\xc7\xe5\x7e\xbb\x3e\xad\xe1\xf6\x97\xb1\x61\xbe\xb1\xa6\x6e\xf9\x06\x00\x00\xff\xff\x14\xbb\xe1\xa1\xe3\x00\x00\x00")

func studioEditorsObjectIndexHtmlBytes() ([]byte, error) {
	return bindataRead(
		_studioEditorsObjectIndexHtml,
		"studio/editors/object/index.html",
	)
}

func studioEditorsObjectIndexHtml() (*asset, error) {
	bytes, err := studioEditorsObjectIndexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "studio/editors/object/index.html", size: 227, mode: os.FileMode(420), modTime: time.Unix(1585434213, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _studioEditorsObjectMainCss = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xca\x28\xc9\xcd\x51\xa8\xe6\x52\x50\x50\x50\x48\xce\xcf\xc9\x2f\xb2\x52\x28\xcf\xc8\x2c\x49\xb5\xe6\xaa\x05\x04\x00\x00\xff\xff\x92\x19\x6b\xc6\x1a\x00\x00\x00")

func studioEditorsObjectMainCssBytes() ([]byte, error) {
	return bindataRead(
		_studioEditorsObjectMainCss,
		"studio/editors/object/main.css",
	)
}

func studioEditorsObjectMainCss() (*asset, error) {
	bytes, err := studioEditorsObjectMainCssBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "studio/editors/object/main.css", size: 26, mode: os.FileMode(420), modTime: time.Unix(1585434223, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _studioEditorsObjectMainMjs = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x8f\x4f\x4b\xf3\x40\x10\x87\xef\xf9\x14\x43\x2f\xdd\x7d\x79\x33\xdb\xc4\x5e\xb4\x28\x88\x17\x0f\x0a\x42\xbd\xcb\x26\xd9\x92\x2d\xdd\x3f\x4e\x76\x1b\x4b\xe9\x77\x97\x34\x29\x16\x0d\xd6\xbd\xed\xcc\x33\xbf\x79\x46\x08\xd0\xc6\x3b\x0a\x30\x11\x5b\x5b\x09\xa3\x43\x4d\x7a\x93\xe6\x38\xc3\x39\x1a\x6d\x71\xdd\x4c\x16\xc9\x77\x6c\xfd\x1e\x15\xed\xd2\x2b\x9c\x63\x76\x91\x8a\x3a\xcd\x30\xcb\xff\x46\x7a\xd7\xe8\xa0\x9d\xbd\xc8\x96\xce\x06\xf5\x11\x8c\xb2\x31\xcd\xf1\x1a\x67\xbf\x4c\x34\x2f\x9b\x68\x8a\x34\xc7\x2c\xc3\xfc\x8b\x4b\x06\xea\x1f\xc8\x06\x6a\x17\x5a\x55\xc0\x8a\x9c\x81\xa9\x78\xeb\xbf\x68\xd6\xcd\x34\x11\xe2\x1c\x94\xde\x9f\xa8\x8d\x2e\x84\xf4\xfe\x48\x2d\x92\x6e\xef\x2a\xda\xb2\xd3\x87\x96\xa4\x67\x65\xc1\x61\xdf\x95\xbb\x47\x2a\x44\xb2\xb0\xdf\x6a\xd5\xde\x00\xe3\x70\x7b\x07\x86\x95\x05\xe3\xfc\x70\x74\x3e\x1c\x13\x06\x59\x24\x25\xab\x1d\x3b\xe5\xb1\xb3\xa0\x41\xad\x95\xa1\xac\x1f\x96\x4b\xc6\x17\x63\x9d\xc7\xd7\xe7\xa7\x9f\x2d\x52\x2b\x52\x4d\xcd\x86\xf5\x48\xaa\x22\xd9\x32\xce\x4f\x9c\x41\xe3\xa2\x0d\xac\x72\x65\x34\xca\x06\x2c\x5c\xb5\xfb\xdf\x9f\xd3\x0f\x75\x07\xdf\x7b\xcf\xfb\xec\x03\x4f\x92\x11\xa1\x11\x93\xcf\x00\x00\x00\xff\xff\x61\x02\x12\xa4\x6b\x02\x00\x00")

func studioEditorsObjectMainMjsBytes() ([]byte, error) {
	return bindataRead(
		_studioEditorsObjectMainMjs,
		"studio/editors/object/main.mjs",
	)
}

func studioEditorsObjectMainMjs() (*asset, error) {
	bytes, err := studioEditorsObjectMainMjsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "studio/editors/object/main.mjs", size: 619, mode: os.FileMode(420), modTime: time.Unix(1585434233, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _studioEditorsTableIndexHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4c\x8f\x51\x0e\x83\x30\x0c\x43\xff\x39\x45\x96\x03\xac\xe2\xbf\xed\x7e\x77\x0d\xd6\x06\x85\x91\x52\xd4\x74\x9a\xb8\xfd\x04\x65\x12\x5f\x89\x25\x3f\x5b\xb6\xb7\x98\x43\xdd\x56\x02\xae\x49\x7c\x67\xdb\x01\x00\xb0\x4c\x43\x6c\xef\x21\x35\x94\x69\xad\xb0\x9b\x1d\xa6\x1c\x3f\x42\x08\x5a\x82\xc3\x34\x4c\xcb\x3d\xbd\x15\xbd\x35\xcd\x75\xc1\x64\x5a\x66\x28\x24\x0e\xb5\x6e\x42\xca\x44\x15\x81\x0b\x8d\x27\x18\x54\x11\xcc\xd9\x69\x8e\x52\x80\xa6\x5e\x39\x6e\xa0\x2b\x89\x04\xa6\x30\x3b\x1c\x07\x51\xc2\x4b\x3a\xf7\xfe\x49\x22\x19\xbe\xb9\x48\x7c\x58\xc3\xfd\x3f\x69\x87\x7d\x67\x4d\x5b\xf4\x0b\x00\x00\xff\xff\x2b\xac\x9b\x90\xe9\x00\x00\x00")

func studioEditorsTableIndexHtmlBytes() ([]byte, error) {
	return bindataRead(
		_studioEditorsTableIndexHtml,
		"studio/editors/table/index.html",
	)
}

func studioEditorsTableIndexHtml() (*asset, error) {
	bytes, err := studioEditorsTableIndexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "studio/editors/table/index.html", size: 233, mode: os.FileMode(420), modTime: time.Unix(1585868570, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _studioEditorsTableMainCss = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xca\x28\xc9\xcd\x51\xa8\xe6\x52\x50\x50\x50\x48\xce\xcf\xc9\x2f\xb2\x52\x28\xcf\xc8\x2c\x49\xb5\xe6\xaa\x05\x04\x00\x00\xff\xff\x92\x19\x6b\xc6\x1a\x00\x00\x00")

func studioEditorsTableMainCssBytes() ([]byte, error) {
	return bindataRead(
		_studioEditorsTableMainCss,
		"studio/editors/table/main.css",
	)
}

func studioEditorsTableMainCss() (*asset, error) {
	bytes, err := studioEditorsTableMainCssBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "studio/editors/table/main.css", size: 26, mode: os.FileMode(420), modTime: time.Unix(1585433982, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _studioEditorsTableMainMjs = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x8f\x4f\x4b\xf3\x40\x10\x87\xef\xf9\x14\x43\x2f\xdd\x7d\x79\x33\xdb\xc4\x5e\xb4\x28\x88\x17\x0f\x0a\x42\xbd\xcb\x26\xd9\x92\x2d\xdd\x3f\x4e\x76\x1b\x4b\xe9\x77\x97\x34\x29\x16\x0d\xd6\xbd\xed\xcc\x33\xbf\x79\x46\x08\xd0\xc6\x3b\x0a\x30\x11\x5b\x5b\x09\xa3\x43\x4d\x7a\x93\xe6\x38\xc3\x39\x1a\x6d\x71\xdd\x4c\x16\xc9\x77\x6c\xfd\x1e\x15\xed\xd2\x2b\x9c\x63\x76\x91\x8a\x3a\xcd\x30\xcb\xff\x46\x7a\xd7\xe8\xa0\x9d\xbd\xc8\x96\xce\x06\xf5\x11\x8c\xb2\x31\xcd\xf1\x1a\x67\xbf\x4c\x34\x2f\x9b\x68\x8a\x34\xc7\x2c\xc3\xfc\x8b\x4b\x06\xea\x1f\xc8\x06\x6a\x17\x5a\x55\xc0\x8a\x9c\x81\xa9\x78\xeb\xbf\x68\xd6\xcd\x34\x11\xe2\x1c\x94\xde\x9f\xa8\x8d\x2e\x84\xf4\xfe\x48\x2d\x92\x6e\xef\x2a\xda\xb2\xd3\x87\x96\xa4\x67\x65\xc1\x61\xdf\x95\xbb\x47\x2a\x44\xb2\xb0\xdf\x6a\xd5\xde\x00\xe3\x70\x7b\x07\x86\x95\x05\xe3\xfc\x70\x74\x3e\x1c\x13\x06\x59\x24\x25\xab\x1d\x3b\xe5\xb1\xb3\xa0\x41\xad\x95\xa1\xac\x1f\x96\x4b\xc6\x17\x63\x9d\xc7\xd7\xe7\xa7\x9f\x2d\x52\x2b\x52\x4d\xcd\x86\xf5\x48\xaa\x22\xd9\x32\xce\x4f\x9c\x41\xe3\xa2\x0d\xac\x72\x65\x34\xca\x06\x2c\x5c\xb5\xfb\xdf\x9f\xd3\x0f\x75\x07\xdf\x7b\xcf\xfb\xec\x03\x4f\x92\x11\xa1\x11\x93\xcf\x00\x00\x00\xff\xff\x61\x02\x12\xa4\x6b\x02\x00\x00")

func studioEditorsTableMainMjsBytes() ([]byte, error) {
	return bindataRead(
		_studioEditorsTableMainMjs,
		"studio/editors/table/main.mjs",
	)
}

func studioEditorsTableMainMjs() (*asset, error) {
	bytes, err := studioEditorsTableMainMjsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "studio/editors/table/main.mjs", size: 619, mode: os.FileMode(420), modTime: time.Unix(1585434053, 0)}
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
	"studio/editors/object/index.html": studioEditorsObjectIndexHtml,
	"studio/editors/object/main.css":   studioEditorsObjectMainCss,
	"studio/editors/object/main.mjs":   studioEditorsObjectMainMjs,
	"studio/editors/table/index.html":  studioEditorsTableIndexHtml,
	"studio/editors/table/main.css":    studioEditorsTableMainCss,
	"studio/editors/table/main.mjs":    studioEditorsTableMainMjs,
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
	"studio": &bintree{nil, map[string]*bintree{
		"editors": &bintree{nil, map[string]*bintree{
			"object": &bintree{nil, map[string]*bintree{
				"index.html": &bintree{studioEditorsObjectIndexHtml, map[string]*bintree{}},
				"main.css":   &bintree{studioEditorsObjectMainCss, map[string]*bintree{}},
				"main.mjs":   &bintree{studioEditorsObjectMainMjs, map[string]*bintree{}},
			}},
			"table": &bintree{nil, map[string]*bintree{
				"index.html": &bintree{studioEditorsTableIndexHtml, map[string]*bintree{}},
				"main.css":   &bintree{studioEditorsTableMainCss, map[string]*bintree{}},
				"main.mjs":   &bintree{studioEditorsTableMainMjs, map[string]*bintree{}},
			}},
		}},
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
