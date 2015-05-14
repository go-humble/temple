package temple

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

var _templates_generated_go_tmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xb4\x90\x4d\x6a\xc3\x30\x10\x85\xd7\x9a\x53\x4c\xbc\xb2\xa1\x8d\x0f\x50\xb2\x28\x74\x59\x4a\xa1\x39\x40\xc6\xf2\xc4\x16\xd5\x8f\x91\xa5\x96\x60\x7c\xf7\x4a\xfe\x81\x6e\xb2\x6b\x17\xc6\x1a\xcd\x1b\x7d\xef\xcd\x40\xf2\x93\x3a\xc6\x69\xc2\xe3\xfb\x7a\x7e\x23\xc3\x38\xcf\x00\x75\x8d\xe7\x5e\x8d\xb8\x6b\x7a\x1a\xb1\x61\xb6\x48\x31\x38\x43\x41\x49\xd2\xfa\x86\x1d\x5b\xf6\x14\xb8\xc5\x6f\x15\x7a\x0c\x6c\x06\xcd\xc7\x3c\xfd\xe2\xd0\xba\x80\xdc\xaa\x80\x86\x6c\xcc\xf2\x03\x80\x32\x83\xf3\x01\x4b\x10\x45\x97\x26\x62\x73\x94\xce\xd4\x9d\x7b\xec\xa3\x69\x34\xd7\xeb\x0b\xdb\xaf\x80\x0a\xe0\x1a\xad\x44\x65\x55\x28\x2b\x9c\x40\x7c\x91\x47\xf6\xcb\xe7\x3c\x80\x48\xe6\x3d\xd9\xe4\x30\x45\xf0\x41\x91\x1e\xb3\x7f\xa1\xae\x8b\xea\xb4\x5b\x7a\x6e\xdb\xad\x5f\x16\x39\xef\x16\xb4\x78\xc0\x4b\x2e\x3f\xbc\x4c\xd5\xa5\x7a\x5a\xa6\x0e\x27\xb4\x4a\x67\x9a\x18\xc8\x2a\x59\xa6\xcb\x0a\xc4\xbc\xd0\xd8\xb6\xcb\x86\x7e\x91\x5f\xe9\xe6\x62\xb8\x0b\x5e\xdb\xff\xc0\x3d\x67\x44\xda\xfe\x5d\xf2\x2e\xf8\x23\xf6\xfc\x13\x00\x00\xff\xff\x34\xa7\xfa\xfc\x32\x02\x00\x00")

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

	info := bindata_file_info{name: "templates/generated.go.tmpl", size: 562, mode: os.FileMode(420), modTime: time.Unix(1431577181, 0)}
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

