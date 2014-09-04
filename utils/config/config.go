/*
 * filename   : config.go
 * created at : 2014-07-20 17:17:57
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package config

import (
    "github.com/pelletier/go-toml"
)

var cfg *toml.TomlTree

func MustLoadFile(filename string) {
    var err error
    cfg, err = toml.LoadFile(filename)
    if err != nil {
        panic(err)
    }
}

func Get(key string, defvalue interface{}) interface{} {
    if cfg.Has(key) {
        return cfg.Get(key)
    } else {
        return defvalue
    }
}

func GetString(key string, defvalue string) string {
    return Get(key, defvalue).(string)
}

func GetInt(key string, defvalue int) int {
    return Get(key, defvalue).(int)
}

func GetBool(key string, defvalue bool) bool {
    return Get(key, defvalue).(bool)
}
