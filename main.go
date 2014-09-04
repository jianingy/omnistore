/*
 * filename   : main.go
 * created at : 2014-08-04 23:11:30
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */
package main

import (
    "os"

    "github.com/codegangsta/cli"
    "github.com/go-martini/martini"

    "github.com/jianingy/omnistore/db"
    "github.com/jianingy/omnistore/manifest"
    CONF "github.com/jianingy/omnistore/utils/config"
    LOG "github.com/jianingy/omnistore/utils/log"
)

type Server struct {
    Database struct {
        Connection string
    }
}

func RunServer() {
    app := cli.NewApp()
    app.Name = "omnistore-server"
    app.Flags = []cli.Flag {
        cli.StringFlag{
            Name: "config-file",
            Value: "conf/development.toml",
            Usage: "configuration file",
        },
    }
    app.Action = func(c *cli.Context) {
        configFile := c.String("config-file")
        LOG.Info("loading configuration %s ...\n", configFile)
        CONF.MustLoadFile(configFile)
        LOG.Info("starting server\n")
        m := martini.Classic()
        m.Map(m)
        m.Use(db.GetMartiniMiddleware())
        m.Get("/", func() string { return "It works!" })
        m.Run()
        // serverListen := CONF.GetString("server.listen", "127.0.0.1:3000")
        // LOG.Info("listening on %s\n", serverListen)
        // LOG.Fatal(http.ListenAndServe(serverListen, m))
    }
    app.Run(os.Args)
}


func main() {
    err := manifest.LoadManifests("conf/manifests/*")
    if err != nil { panic(err) }
    LOG.Info(manifest.Manifests)
    // RunServer()
}
