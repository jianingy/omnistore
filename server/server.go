/*
 * filename   : server.go
 * created at : 2014-08-06 17:32:15
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package server

import (
    "io/ioutil"

    "gopkg.in/yaml.v1"

    LOG "github.com/jianingy/omnistore/utils/log"
)

type ServerSetting struct {
    Database struct {
        Connection string
    }
}

var Setting ServerSetting

func Configure(filename string) error {
    LOG.Info("reading configuration %s\n", filename)
    if raw, err := ioutil.ReadFile(filename); err != nil {
        return err
    } else if err := yaml.Unmarshal(raw, &Setting); err != nil {
        return nil
    }
    return nil
}
