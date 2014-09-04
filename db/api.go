/*
 * filename   : api.go
 * created at : 2014-08-04 12:31:38
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package db

import (
    "database/sql"
	"fmt"

	"github.com/lib/pq"
	. "github.com/lib/pq/hstore"
	"github.com/jmoiron/sqlx"

	"github.com/jianingy/omnistore/model"
	"github.com/jianingy/omnistore/utils"
	LOG "github.com/jianingy/omnistore/utils/log"

)

type DBAPI struct {
	db        *sqlx.DB
	tx        *sqlx.Tx
	manifests model.ManifestCollection
}

type NullString struct {
    sql.NullString
}

type ModelItem struct {
    ID          int64
    UUID        string
    Identifier  string
    Model       string
    Value       Hstore
}

// TODO: different name maybe
type StringMap map[string]string

/*
 * SQLs
 */

const (
    SQL_CRETAE_EXNTENSION = `CREATE EXTENSION IF NOT EXISTS "%s"`
	SQL_CREATE_TABLE = `CREATE TABLE IF NOT EXISTS {{ .Table }} (
id         SERIAL PRIMARY KEY,
identifier VARCHAR(124) NOT NULL UNIQUE,
model      VARCHAR(124) NOT NULL,
uuid       UUID NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
value      HSTORE
)
`
	SQL_GET_TABLES     = "SELECT table_name FROM information_schema.tables WHERE table_schema='public'"
	SQL_GET_EXTENSIONS = "SELECT extname FROM pg_extension"

    SQL_CREATE_ITEM    = `INSERT INTO {{ .Table }}(identifier, model, value)
VALUES(:identifier, :model, :value)`

    SQL_GET_ITEM_BY_IDENTIFIER = `SELECT * FROM {{ .Table }} WHERE identifier = $1`
    SQL_GET_ITEM_BY_MODEL = `SELECT * FROM {{ .Table }} WHERE model = $1`
    SQL_DELETE_ITEM_BY_IDENTIFIER = `DELETE FROM {{ .Table }} WHERE identifier = $1`
    SQL_UPDATE_ITEM    = `UPDATE {{ .Table }} SET identifier = :identifier, value = :value
WHERE uuid = :uuid`
)

/*
 * constructors
 */

func NewDBAPI(dialect, connection string) (*DBAPI, error) {
	LOG.Info("initializing database [%s] %s\n", dialect, connection)
	if db, err := sqlx.Connect(dialect, connection); err != nil {
		return nil, err
	} else {
		return &DBAPI{db, nil, nil}, nil
	}
}

/*
 * data structures
 */

func (smap *StringMap) Hstore() Hstore {
    hval := make(map[string]sql.NullString)
    for key, value := range *smap {
        hval[key] = sql.NullString{
            String: value,
            Valid: true,
        }
    }
    return Hstore{Map: hval}
}

/*
 * transaction helpers
 */

func (dbapi *DBAPI) Begin() error {
	if dbapi.tx == nil {
		tx, err := dbapi.db.Beginx()
		dbapi.tx = tx
		return err
	}
	return nil
}

func (dbapi *DBAPI) Commit() error {
	err := dbapi.tx.Commit()
	defer func() {
		dbapi.tx = nil
	}()
	return err
}

func (dbapi *DBAPI) Rollback() error {
	err := dbapi.tx.Rollback()
	defer func() {
		dbapi.tx = nil
	}()
	return err

}

/*
 * Helpers
 */

func GetManifestTable(name string) string {
	return fmt.Sprintf("data_%s", name)
}

func RenderSQL(text, model string) string {
    table := GetManifestTable(model)
    sql, err := utils.RenderTemplate(text, struct { Table string }{ table })
    if err != nil {
        panic(err)
    }

    return sql
}

/*
 * DBAPIs
 */

// Create Datatables if nonexists
func (dbapi *DBAPI) createModelTables() error {
	for name, _ := range dbapi.manifests {
		if _, err := dbapi.db.Exec(RenderSQL(SQL_CREATE_TABLE, name)); err != nil {
			return err
		}
		LOG.Info("creating table for model %s ...", name)
	}
	return nil
}

func (dbapi *DBAPI) InitDB(root string) error {
	var err error
	if dbapi.manifests, err = model.NewManifestCollection(root); err != nil {
		return err
	}
    // Create extensions if nonexists
	if _, err = dbapi.db.Exec(fmt.Sprintf(SQL_CRETAE_EXNTENSION, "hstore")); err != nil {
		return err
	}
	if _, err = dbapi.db.Exec(fmt.Sprintf(SQL_CRETAE_EXNTENSION, "uuid-ossp")); err != nil {
		return err
	}
	if err = dbapi.createModelTables(); err != nil {
		return err
	}
	return nil
}

func (dbapi *DBAPI) CreateItem(manifest, model string, values StringMap) error {
    manifestd, found := dbapi.manifests[manifest]
    if !found {
        return ManifestNotFound{manifest}
    }
    modeld, found := manifestd.Models[model]
    if !found {
        return ModelNotFound{model}
    }

    /* filter valid items */
    filtered := make(StringMap)
    for _, name := range modeld.Properties {
        if value, found := values[name]; found {
            filtered[name] = value
        }
    }

    /* generate identifier */
    identifier, err := utils.RenderTemplate(modeld.Identifier, filtered)
    if err != nil {
        panic(err)
    }
    item := ModelItem {
        Identifier: identifier,
        Model: model,
        Value: filtered.Hstore(),
    }

    if _, err := dbapi.db.NamedExec(RenderSQL(SQL_CREATE_ITEM, manifest), item); err != nil {
        if err.(*pq.Error).Code.Name() == "unique_violation" {
            return DuplicatedIdentifier{}
        }
        return err
    }

    return nil
}

func (dbapi *DBAPI) DeleteItemByIdentifier(manifest, identifier string) error {
    // XXX: check if manifest exists
    q := RenderSQL(SQL_DELETE_ITEM_BY_IDENTIFIER, manifest)
    if _, err := dbapi.db.Exec(q, identifier); err != nil {
        return err
    } else {
        return nil
    }
}

func (dbapi *DBAPI) GetItemsByModel(manifest, model string) ([]ModelItem, error) {
    var items []ModelItem
    // XXX: check if manifest exists
    q := RenderSQL(SQL_GET_ITEM_BY_MODEL, manifest)
    if err := dbapi.db.Select(&items, q, model); err != nil {
        return nil, err
    } else {
        return items, nil
    }
}

func (dbapi *DBAPI) GetItemByIdentifier(manifest, identifier string) (*ModelItem, error) {
    var item ModelItem
    // XXX: check if manifest exists
    q := RenderSQL(SQL_GET_ITEM_BY_IDENTIFIER, manifest)
    if err := dbapi.db.Get(&item, q, identifier); err != nil {
        return nil, err
    } else {
        return &item, nil
    }
}

func (dbapi *DBAPI) UpdateItem(manifest, identifier string, values StringMap) error {
    /* model cannot be changed during update */

    manifestd, found := dbapi.manifests[manifest]
    if !found {
        return ManifestNotFound{manifest}
    }

    var item ModelItem
    q := RenderSQL(SQL_GET_ITEM_BY_IDENTIFIER, manifest)
    if err := dbapi.db.Get(&item, q, identifier); err != nil {
        return err
    }

    modeld, found := manifestd.Models[item.Model]
    if !found {
        return ModelNotFound{item.Model}
    }

    /* filter valid items */
    filtered := make(StringMap)
    for _, name := range modeld.Properties {
        if value, found := values[name]; found {
            filtered[name] = value
        }
    }

    /* update */
    for name, _ := range item.Value.Map {
        item.Value.Map[name] = sql.NullString{String: filtered[name], Valid: true}
    }

    /* generate identifier */
    identifier, err := utils.RenderTemplate(modeld.Identifier, filtered)
    if err != nil {
        panic(err)
    }
    item.Identifier = identifier

    if _, err := dbapi.db.NamedExec(RenderSQL(SQL_UPDATE_ITEM, manifest), item); err != nil {
        if err.(*pq.Error).Code.Name() == "unique_violation" {
            return DuplicatedIdentifier{}
        }
        return err
    }

    return nil
}
