/*
 * filename   : errors.go
 * created at : 2014-08-11 15:17:11
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package db

import (
    "fmt"
)

type ManifestNotFound struct {
    ManifestName string
}

func (e ManifestNotFound) Error() string {
	return fmt.Sprintf("manifest %s not found", e.ManifestName)
}


type ModelNotFound struct {
    ModelName string
}

func (e ModelNotFound) Error() string {
	return fmt.Sprintf("model %s not found", e.ModelName)
}


type DuplicatedIdentifier struct {}

func (e DuplicatedIdentifier) Error() string {
	return fmt.Sprintf("record with the same identifier already exists")
}
