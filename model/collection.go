/*
 * filename   : collection.go
 * created at : 2014-08-06 19:10:57
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package model

import (
	LOG "github.com/jianingy/omnistore/utils/log"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ManifestCollection map[string]*Manifest

func NewManifestCollection(root string) (ManifestCollection, error) {
	var err error
	files, err := filepath.Glob(filepath.Join(root, "*"))
	if err != nil {
		return nil, err
	}

	manifests := make(map[string]*Manifest)

    LOG.Debug("collection manifests: %v from %v", files, root)

	for _, file := range files {
		var property, model []byte
		var manifest *Manifest

		propYAML := filepath.Join(file, "properties.yaml")
		modelYAML := filepath.Join(file, "models.yaml")
		name := filepath.Base(file)

		stat, err := os.Stat(file)
		if err != nil || !stat.IsDir() {
			goto ERROR_OUT
		}

		if property, err = ioutil.ReadFile(propYAML); err != nil {
			goto ERROR_OUT
		}

		if model, err = ioutil.ReadFile(modelYAML); err != nil {
			goto ERROR_OUT
		}

		if manifest, err = NewManifest(name, property, model); err != nil {
			goto ERROR_OUT
		}

		manifests[name] = manifest
		continue

	ERROR_OUT:
		LOG.Warn(err)
	}

	return manifests, nil
}

func (mc ManifestCollection) Get(name string) (*Manifest, bool) {
	manifest, found := mc[name]
	return manifest, found
}

func (mc ManifestCollection) Iterate(fn func(string, *Manifest) error) error {
	for name, value := range mc {
		if err := fn(name, value); err != nil {
			return err
		}
	}
    return nil
}
