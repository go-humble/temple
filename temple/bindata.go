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

var _templates_generated_go_tmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xb4\x52\xc1\x8a\xdb\x30\x10\x3d\x5b\x5f\x31\x9b\x93\xb3\xa4\xd6\x3d\x25\x87\x42\x4b\x29\x94\x74\xa1\xb9\x95\xc2\x8e\xed\x89\xad\xd6\x96\x8c\x3c\xda\x74\x31\xf9\xf7\x8e\x1c\x3b\x31\xec\xd2\x53\xf7\x60\xa4\x99\x79\x33\xef\xe9\x8d\x87\x41\xdf\x83\x4a\xf6\xdf\x0e\x9f\xb6\xc0\xb5\xe9\xe1\x68\x1a\x02\x39\xa9\xcd\xa9\x2c\xa9\x84\x27\x83\x50\xb9\x77\xb9\xb1\x25\x32\x42\x5a\x33\x77\xfd\x56\xeb\xca\x70\x1d\xf2\xac\x70\xad\xfe\xc5\x44\xe1\x44\x56\xdf\x70\x6b\x95\x7c\x39\xc2\xb3\x0b\x50\xd4\x68\x2b\x19\xc9\x1b\xc8\x09\xfa\xe0\x09\xd8\x81\x0f\x56\xa6\x42\x45\x96\x3c\x32\x41\xa6\xb3\x2c\xbb\xf6\x58\x12\x62\x41\x19\xdb\x33\x36\x8d\x28\xa3\xa5\x06\xfa\x43\x45\x60\xcc\x1b\xda\xdc\x06\x31\xfc\x5b\x91\xba\xd7\xe7\xb3\xea\xb0\xf8\x8d\x22\x67\x18\x20\x7b\xb8\xdc\xf7\xd8\x12\x48\x49\x69\x0d\x87\x68\xc1\x8c\xa9\xb1\x17\xc9\x64\x01\x03\xbb\x16\xd9\x14\xa2\xe5\xf9\xaa\xb9\x84\x93\x10\x02\x53\xdb\x35\x94\xc5\xee\x8f\x0e\xac\x63\xa0\xd2\x30\xb4\x68\x43\x84\xdf\x29\x65\xda\xce\x79\x86\x54\x25\xab\x85\x44\x51\x56\x87\x56\xde\xa0\x2f\x13\xa6\x63\xa5\xd6\x4a\x3d\xa1\x8f\xf0\x43\xcc\x08\x53\x2f\xd3\xba\x1f\x3d\x7b\x63\xab\x9f\x13\xdf\x5c\x53\xc9\x03\x7a\x36\xd8\xf4\xf0\x0a\x6a\xaa\xa9\xe4\x2b\x8a\xaf\x2c\x98\x57\x40\x97\x5a\xe4\x3d\x06\x5b\x88\xe9\x86\xd3\x35\x0c\x2a\x89\x32\xc8\x8f\x9f\xf3\x2a\xa9\x60\xbb\x9b\x9f\xbb\xa7\xd3\x67\xef\x42\x97\xca\xa2\xc5\x4b\x3f\xee\x38\xbb\x4a\x11\x3b\x13\x73\x1c\x9b\x77\x50\x65\x1f\xca\x72\x2a\xa5\xab\xe8\xfc\x64\xf9\x6a\x03\x8f\x31\xfc\xee\x0b\x89\x1e\xd7\xef\xc7\x86\xbb\x1d\x58\xd3\x44\xfe\xa4\x43\x6b\x8a\x54\x92\xc2\x72\x1e\x89\xc8\x96\xe3\xae\x16\xa4\xf3\xd3\x5e\x72\x5e\x2a\x6f\x40\x79\x5b\xcc\x4b\xd2\xb9\xf6\x9f\x68\x17\x3f\x41\x24\xb8\x46\x8b\xbd\xc7\xfc\x1c\xdc\x36\x1d\xb3\xd3\x5d\x9d\xff\x06\x00\x00\xff\xff\xb6\x40\x6e\xe8\xe8\x03\x00\x00")

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

	info := bindata_file_info{name: "templates/generated.go.tmpl", size: 1000, mode: os.FileMode(420), modTime: time.Unix(1431581189, 0)}
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

