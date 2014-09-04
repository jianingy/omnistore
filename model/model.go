/*
 * filename   : model.go
 * created at : 2014-08-05 11:32:04
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package model

import (
    "fmt"
	"strings"

	"gopkg.in/yaml.v1"

	LOG "github.com/jianingy/omnistore/utils/log"
)

type Property struct {
	ConstructorName string
	ConstructorArgs map[string]interface{}
}

type Model struct {
    DisplayName  string
    Identifier   string
    Properties   []string
    References   map[string]string
}

type Manifest struct {
	Name          string
    Properties    map[string]Property
    Models        map[string]Model
}

func NewManifest(name string, property, model []byte) (*Manifest, error) {
	manifest := &Manifest{
        name,
        make(map[string]Property),
        make(map[string]Model),
    }
	if err := manifest.makePropertyMap(property); err != nil {
		return nil, err
	}
	if err := manifest.makeModelMap(model); err != nil {
		return nil, err
	}
	return manifest, nil
}

func (mm *Manifest) makePropertyMap(cfg []byte) error {

	m := make(map[interface{}]interface{})
	if err := yaml.Unmarshal(cfg, &m); err != nil {
		return err
	}

	// parse property definitions
	for name, properties := range m {
		var prop Property
		for key, value := range properties.(map[interface{}]interface{}) {
			if name := key.(string); strings.HasPrefix(name, "type.") {
				prop.ConstructorName = strings.TrimPrefix(name, "type.")
				prop.ConstructorArgs = make(map[string]interface{})
				for _, item := range value.([]interface{}) {
					switch item.(type) {
					case string:
						prop.ConstructorArgs[item.(string)] = true
					case map[interface{}]interface{}:
						for k, v := range item.(map[interface{}]interface{}) {
							prop.ConstructorArgs[k.(string)] = v
						}
					default:
						LOG.Warn("unknown item type of property %v", item)
					}
				}
			}
		}
        mm.Properties[name.(string)] = prop
	}

	return nil
}

func (mm *Manifest) makeModelMap(cfg []byte) error {
	m := make(map[interface{}]interface{})
	if err := yaml.Unmarshal(cfg, &m); err != nil {
		return err
	}
	for name, models := range m {
        model := Model{"", "", nil, make(map[string]string)}
        for key, value := range models.(map[interface{}]interface{}) {

            switch key.(string) {
            case "displayname":
                model.DisplayName = value.(string)
            case "identifier":
                model.Identifier = value.(string)
            case "properties":
                for _, prop := range value.([]interface{}) {
                    if _, found := mm.Properties[prop.(string)]; !found {
                        return fmt.Errorf("property %s not defined", prop.(string))
                    }
                    model.Properties = append(model.Properties, prop.(string))
                }
            case "references":
                for _, refer := range value.([]interface{}) {
					switch refer.(type) {
					case string:
						model.References[refer.(string)] = refer.(string)
					case map[interface{}]interface{}:
						for k, v := range refer.(map[interface{}]interface{}) {
							model.References[k.(string)] = v.(string)
						}
					default:
						LOG.Warn("unknown item type of references %v", refer)
					}
                }
            default:
            }
        }
        mm.Models[name.(string)] = model
    }

    // check if references all exists
	for _, model := range mm.Models {
        for _, refer := range model.References {
            if _, found := mm.Models[refer]; !found {
                return fmt.Errorf("model %v referred cannot be found", refer)
            }
        }
    }

    return nil
}
