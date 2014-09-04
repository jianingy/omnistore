/*
 * filename   : utils.go
 * created at : 2014-07-20 17:38:04
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package utils

import (
    "bytes"
    "fmt"
    "os"

    "text/template"

)

func GetLocalFilePath(path string) string {
    return fmt.Sprintf("%s/src/github.com/jianingy/omnistore/%s",
        os.Getenv("GOPATH"), path)
}


func RenderTemplate(tmpl string, data interface{}) (string, error) {
    var buf bytes.Buffer
    if tmpl, err := template.New(tmpl).Parse(tmpl); err == nil {
        if err := tmpl.Execute(&buf, data); err != nil {
            return "", err
        }
        return buf.String(), nil
    } else {
        return "", err
    }
}
