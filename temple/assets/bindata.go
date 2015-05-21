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

var _templates_generated_go_tmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xb4\x52\xcb\x6e\xdb\x30\x10\x3c\x93\x5f\xb1\xf1\x49\x0e\x52\xa9\x2d\xd0\x4b\x8a\x1c\xd2\x34\x28\x02\xa4\x6e\x80\xb8\xa7\xa2\x40\x28\x72\x25\xb1\x95\x48\x81\x0f\x3b\x86\xe1\x7f\xef\x52\x96\x6c\x03\x09\x7a\x6a\x0f\x86\xc5\x9d\xe1\xce\xec\x2c\xb7\xdb\xe2\x9c\xb3\x1b\xdb\x6f\x9c\xae\x9b\x00\xef\xdf\xbe\xfb\x00\xd7\x2d\x3e\xc3\x27\x67\xd7\x06\x73\xce\xae\xdb\x16\x06\xd0\x83\x43\x8f\x6e\x85\x2a\x87\xef\x1e\xc1\x56\x10\x1a\xed\xc1\xdb\xe8\x24\x82\xb4\x0a\x41\x7b\xce\x6a\xbb\x42\x67\x50\x41\xb9\x21\x02\xc2\xd7\xbb\x25\xb4\x5a\xa2\xf1\x78\x01\xeb\x46\xcb\x06\xa4\x30\x50\x22\x54\x36\x1a\xc5\x99\x36\x03\xef\xfe\xee\xe6\x76\xf1\x78\x0b\x95\x6e\x49\x97\xb3\xc5\xb7\xe5\xed\xe5\x5e\x22\x95\xa8\x37\x60\x57\xa2\x52\xd4\x7a\xa5\x05\xd4\xf6\x4d\xa9\x8d\x12\x41\x40\xd6\x84\xd0\xfb\xcb\xa2\xa8\x75\x68\x62\x99\x4b\xdb\x15\xbf\x02\x62\x5c\xa3\x29\x8e\xbc\x39\x67\x77\x15\x6c\x6c\x04\xd9\x08\x53\x53\xcb\x70\x91\x7c\xf8\xe8\x10\x82\x05\x17\x0d\x75\x85\x1a\x0d\x3a\x11\x10\xf2\x22\xcf\xf3\xc3\x1d\x83\x24\x4c\x2c\x6d\x7c\x10\x14\x4a\xf2\x7c\xe2\x01\x9f\x51\xc6\x20\xca\x96\xa6\x3c\x34\x0a\xf0\x77\x47\xfc\xbc\xd8\xed\x78\x2f\xe4\x6f\x41\x76\xb6\x5b\xc8\x1f\xf6\xdf\x0b\xd1\x21\x10\xc4\x8b\x02\x96\x29\x82\x89\xd3\x08\x4f\x96\xd1\x80\x88\xc1\x76\x22\x68\x49\x5e\x36\x07\xcf\x0a\xd6\x24\x08\x01\xbb\x3e\xa5\x48\xb7\x3f\x5b\x30\x36\x00\x2a\x1d\xa0\x13\x26\x26\xfa\x19\xe7\xba\xeb\xad\x0b\x90\x71\x36\x3b\xb1\x48\xce\x9a\xd8\xd1\x0c\xc5\xbe\xc3\xf8\x37\xe3\x73\xce\x57\xc2\x25\xfa\x32\x55\x48\xc9\x53\xb7\xfe\x87\x0f\x4e\x9b\xfa\xe7\xa8\x37\x61\x9c\x3d\x08\x17\xb4\x68\x3d\xbc\xc2\x1a\x31\xce\xee\x05\xe5\x4a\xef\xea\x35\xd2\x1e\x4b\xba\x55\x34\x92\x42\xd7\x21\x9b\xc3\x96\xb3\x64\x03\xdd\xf0\xb3\x8e\x5e\x1b\x5c\x5e\x4d\xe3\x2e\x70\xfd\xc5\xd9\xd8\x67\xb4\x68\xca\xd2\x0d\x3b\xce\x0f\x56\x28\x4e\xa6\xab\xe1\xf2\x15\xd4\xf9\xb5\x52\x23\x94\xcd\x52\xf2\x63\xe4\xb3\x0b\x78\x4a\xc7\x47\x27\xe9\xf4\x34\xff\x38\x5c\x38\xbb\x02\xa3\xdb\xa4\xcf\x7a\x61\xb4\xcc\xa8\x48\x2a\xbb\x41\x08\x8d\x1a\x76\x75\x22\x3a\x8d\xf6\x52\x73\x8f\xfc\x07\xc9\xe3\x62\x5e\x8a\x4e\xd8\x3f\x92\x3d\x79\x04\x49\xe0\x70\x3a\xd9\x7b\xaa\x4f\x87\xe3\xa6\x53\x75\xfc\xe6\xbb\x3f\x01\x00\x00\xff\xff\x11\xf2\x8a\x0e\x7d\x04\x00\x00")

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

	info := bindata_file_info{name: "templates/generated.go.tmpl", size: 1149, mode: os.FileMode(420), modTime: time.Unix(1432099076, 0)}
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

