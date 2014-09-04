/*
 * filename   : model_test.go
 * created at : 2014-08-06 18:10:58
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package model

import (
    "testing"

	_ "github.com/jianingy/omnistore/utils/log"
)

func TestBuildPropertyMap(t *testing.T) {
    sample := `
hostname:
  type.string:
    - match: ".cn[0-9]"
    - min: 4
    - unique
    - not_null

rackname:
  type.string:
    - min: 4
    - not_null

sitename:
  type.string:
    - min: 4
    - not_null`

    m := &Manifest{"testing", make(map[string]Property), nil}
    if err := m.buildPropertyMap([]byte(sample)); err != nil {
        t.Errorf("cannot build property map")
    }
}

func TestNewModel(t *testing.T) {
    property := `
hostname:
  type.string:
    - match: ".cn[0-9]"
    - min: 4
    - unique
    - not_null

rackname:
  type.string:
    - min: 4
    - not_null

sitename:
  type.string:
    - min: 4
    - not_null`

    model := `
site:
  displayname: Site
  identifier: "{{ sitename }}"
  properties:
    - sitename

rack:
  displayname: rack
  identifier: "{{ rackname }}"
  properties:
    - rackname
  references:
    - site: site

rackserver:
  displayname: Rack Server
  identifier: "{{ hostname }}"
  properties:
    - hostname
  references:
    - rack: rack`

    _, err := NewManifest("testing", []byte(property), []byte(model))
    if err != nil {
        t.Errorf("cannot build property map: %v", err)
    }
}

func TestNewModelWithMissingProperty(t *testing.T) {
    property := `
hostname:
  type.string:
    - match: ".cn[0-9]"
    - min: 4
    - unique
    - not_null

rackname:
  type.string:
    - min: 4
    - not_null

sitename:
  type.string:
    - min: 4
    - not_null`

    model := `
site:
  displayname: Site
  identifier: "{{ sitename }}"
  properties:
    - sitename

rack:
  displayname: rack
  identifier: "{{ rackname }}"
  properties:
    - rackname
  references:
    - site: site

rackserver:
  displayname: Rack Server
  identifier: "{{ hostname }}"
  properties:
    - hostname
    - height
  references:
    - rack: rack`

    _, err := NewManifest("testing", []byte(property), []byte(model))
    if err == nil {
        t.Errorf("missing property not properly detected")
    }
}

func TestNewModelWithMissingReference(t *testing.T) {
    property := `
hostname:
  type.string:
    - match: ".cn[0-9]"
    - min: 4
    - unique
    - not_null

rackname:
  type.string:
    - min: 4
    - not_null

sitename:
  type.string:
    - min: 4
    - not_null`

    model := `
site:
  displayname: Site
  identifier: "{{ sitename }}"
  properties:
    - sitename
  references:
    - missing: missing

rack:
  displayname: rack
  identifier: "{{ rackname }}"
  properties:
    - rackname
  references:
    - site: site

rackserver:
  displayname: Rack Server
  identifier: "{{ hostname }}"
  properties:
    - hostname
  references:
    - rack: rack`

    _, err := NewManifest("testing", []byte(property), []byte(model))
    if err == nil {
        t.Errorf("missing reference not properly detected")
    }
}
