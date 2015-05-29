package assets

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path"
	"path/filepath"
)

func bindata_read(data []byte, name string) ([]byte, error) {
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

type bindata_file_info struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _templates_generated_go_tmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xb4\x53\x4d\x6f\xdb\x38\x10\x3d\x8b\xbf\x62\xe2\x93\x1c\x78\xa5\xdd\x05\xf6\x92\x45\x0e\x69\x1a\x04\x01\x12\x37\x40\xdc\x7b\x68\x6a\x2c\xb1\x95\x48\x81\x1c\xda\x31\x0c\xff\xf7\x0e\x25\xd9\xb1\xf3\xe1\x53\x7b\xb2\xc9\x79\x6f\xde\xe3\xcc\xd3\x66\x93\x9f\x8b\xe4\xda\xb6\x6b\xa7\xcb\x8a\xe0\xdf\xbf\xff\xf9\x0f\xae\x6a\x7c\x81\x2f\xce\xae\x0c\x66\x22\xb9\xaa\x6b\xe8\x8a\x1e\x1c\x7a\x74\x4b\x2c\x32\xf8\xee\x11\xec\x02\xa8\xd2\x1e\xbc\x0d\x4e\x21\x28\x5b\x20\x68\x2f\x92\xd2\x2e\xd1\x19\x2c\x60\xbe\x66\x00\xc2\xc3\xdd\x0c\x6a\xad\xd0\x78\x9c\xc0\xaa\xd2\xaa\x02\x25\x0d\xcc\x11\x16\x36\x98\x42\x24\xda\x74\xb8\xfb\xbb\xeb\x9b\xe9\xd3\x0d\x2c\x74\xcd\xba\x22\x99\x7e\x9b\xdd\x5c\xf4\x12\xf1\x8a\x7b\x03\x36\x73\x2c\x0a\x6e\xbd\xd4\x12\x4a\xfb\xd7\x5c\x9b\x42\x92\x84\xb4\x22\x6a\xfd\x45\x9e\x97\x9a\xaa\x30\xcf\x94\x6d\xf2\x1f\x84\x18\x56\x68\xf2\x57\xdc\x58\x24\x77\x0b\x58\xdb\x00\xaa\x92\xa6\xe4\x96\x34\x89\x3e\x7c\x70\x08\x64\xc1\x05\xc3\x5d\xa1\x44\x83\x4e\x12\x42\x96\x67\x59\xb6\xe7\x18\x64\x61\x46\x69\xe3\x49\xf2\x50\xa2\xe7\x03\x0f\xf8\x82\x2a\x90\x9c\xd7\xfc\xca\x7d\x23\x82\xd3\x8e\xc4\x79\xbe\xdd\x8a\x56\xaa\x9f\x92\xed\x6c\x36\x90\x3d\xf6\xff\xa7\xb2\x41\xe0\x92\xc8\x73\x98\xc5\x11\xec\x30\x95\xf4\x6c\x19\x0d\xc8\x40\xb6\x91\xa4\x15\x7b\x59\xef\x3d\x17\xb0\x62\x41\x20\x6c\xda\x38\x45\x66\x7f\xb5\x60\x2c\x01\x16\x9a\xa0\x91\x26\x44\xf8\x99\x10\xba\x69\xad\x23\x48\x45\x32\x3a\xb0\xc8\xce\xaa\xd0\xf0\x1b\xf2\xbe\xc3\xf0\x33\x12\x63\x21\x96\xd2\x45\xf8\x2d\xd2\x2c\x5e\xc6\x01\x2d\x82\x51\xa9\x89\x56\x3d\x39\x6d\xca\x31\xa4\xe7\x83\xf4\x0e\x33\x01\x74\xce\xba\x71\x47\x7c\x94\x8e\xb4\xac\x4f\xf1\x06\xc8\x11\xed\x5e\xf2\xfc\xe9\x14\xab\x47\xbc\x92\x1e\x82\xa7\xd3\x46\xdf\xfa\xdc\x73\x3e\xf7\xf8\xc6\xe2\x9e\xf1\xa9\xbd\x63\x77\x71\x86\x11\xc3\x01\xd2\x94\x8e\x61\x23\x92\x38\x52\xb6\xdc\xdb\xe6\x2f\x07\x2e\x2e\x77\xab\x9b\xe2\xea\xd6\xd9\xd0\xa6\xfc\x1a\xce\x85\xeb\xf2\xba\x93\xf6\x31\x1a\x89\x5e\x74\xe4\x4b\x28\xb3\xab\xa2\x18\x4a\xe9\x28\xa6\x68\x88\xcf\x68\x02\xcf\xf1\xf8\xe4\x14\x9f\x9e\xc7\xff\x77\x84\xb3\x4b\x30\xba\x8e\xfa\x49\x2b\x8d\x56\x29\x5f\xb2\xca\xb6\x13\x42\x53\x74\xb9\x3b\x10\xed\xed\x7f\xa4\xd9\x57\xfe\x80\xe4\x6e\x2b\x1f\x89\xee\x6a\xbf\x49\xf6\x28\xd0\x51\xe2\xe0\x7c\x14\xda\xa1\xb6\xdf\xfe\xeb\xe6\x87\xca\xb0\xe6\x77\xe1\x8b\xe5\x37\x77\xef\xd2\x76\x80\xf9\x24\x5e\x07\x88\x41\x68\xfb\x2b\x00\x00\xff\xff\x1e\x4e\x7e\x4b\xbe\x05\x00\x00")

func templates_generated_go_tmpl_bytes() ([]byte, error) {
	return bindata_read(
		_templates_generated_go_tmpl,
		"templates/generated.go.tmpl",
	)
}

func templates_generated_go_tmpl() (*asset, error) {
	bytes, err := templates_generated_go_tmpl_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "templates/generated.go.tmpl", size: 1470, mode: os.FileMode(420), modTime: time.Unix(1432862983, 0)}
	a := &asset{bytes: bytes, info:  info}
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
	if (err != nil) {
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
	"templates/generated.go.tmpl": templates_generated_go_tmpl,
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

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"templates": &_bintree_t{nil, map[string]*_bintree_t{
		"generated.go.tmpl": &_bintree_t{templates_generated_go_tmpl, map[string]*_bintree_t{
		}},
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
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

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

