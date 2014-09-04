/*
 * filename   : martini.go
 * created at : 2014-08-04 12:26:23
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */
package db

import (
    _ "github.com/lib/pq"

    "github.com/go-martini/martini"

	CONF "github.com/jianingy/omnistore/utils/config"
	LOG "github.com/jianingy/omnistore/utils/log"
)


func GetMartiniMiddleware() martini.Handler {

    dialect := CONF.GetString("database.dialect", "")
    connection := CONF.GetString("database.connection", "")

    dbapi, err := NewDBAPI(dialect, connection)

    if err != nil {
        panic(err)
    }

    return func(c martini.Context) {
        c.Map(dbapi)
        dbapi.Begin()

        // NOTE(jianingy): 研究下 rollback 处理的时间点
        defer func() {
            if r := recover(); r != nil {
                LOG.Warn("failed to commit transaction")
                err := dbapi.Rollback()
                if err != nil {
                    LOG.Warn("failed to rollback transaction: %s", err.Error)
                }
            }
        }()
        c.Next()
        err := dbapi.Commit()
        if err != nil {
            panic("dbapi.Commit returns false")
        }
    }
}
