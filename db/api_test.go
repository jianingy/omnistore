/*
 * filename   : api_test.go
 * created at : 2014-08-09 10:07:00
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package db

import (
	. "gopkg.in/check.v1"
	"testing"

	"github.com/jianingy/omnistore/utils"
	CONF "github.com/jianingy/omnistore/utils/config"
	LOG "github.com/jianingy/omnistore/utils/log"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type DBSuite struct {
	dbapi *DBAPI
}

var _ = Suite(&DBSuite{})

func (s *DBSuite) SetUpSuite(c *C) {
	// Use s.dir to prepare some data.
	var err error

	configFile := utils.GetLocalFilePath("conf/testing.toml")
	LOG.Info("loading configuration %s ...\n", configFile)
	CONF.MustLoadFile(configFile)

	dialect := CONF.GetString("database.dialect", "")
	connection := CONF.GetString("database.connection", "")
	// clean testing database
	s.dbapi, err = NewDBAPI(dialect, connection)
	if err != nil {
		c.Error(err)
	}
	s.dbapi.db.MustExec("DROP SCHEMA IF EXISTS public CASCADE")
}

func (s *DBSuite) SetUpTest(c *C) {
	s.dbapi.db.MustExec("CREATE SCHEMA public")
	err := s.dbapi.InitDB(utils.GetLocalFilePath("conf/tests"))
	if err != nil {
		c.Error(err)
	}
}

func (s *DBSuite) TearDownTest(c *C) {
	s.dbapi.db.MustExec("DROP SCHEMA public CASCADE")
}

/*
func (s *DBSuite) TestInitDB(c *C) {
	err := s.dbapi.InitDB(utils.GetLocalFilePath("conf/tests"))
	if err != nil {
		c.Error(err)
	}
	LOG.Debug("manifests: %v", s.dbapi.manifests)
}
*/

func (s *DBSuite) TestCreateItem(c *C) {
	var err error
	values := make(map[string]string)
	values["sitename"] = "cn8"
	err = s.dbapi.CreateItem("colo", "site", values)
	if err != nil {
		c.Error(err)
	}
}

func (s *DBSuite) TestCreateDupItem(c *C) {
	var err error
	values := make(map[string]string)
	values["sitename"] = "cn8"
	err = s.dbapi.CreateItem("colo", "site", values)
	if err != nil {
		c.Error(err)
	}
	err = s.dbapi.CreateItem("colo", "site", values)
	if err != nil {
		switch err.(type) {
		case DuplicatedIdentifier:
			return
		default:
			c.Error("an error but not duplicated identifier raised")
		}
	} else {
		c.Error("duplicated identifier error not raise")
	}
}

func (s *DBSuite) TestGetItemsByModel(c *C) {
	var err error
	values := make(map[string]string)
	values["sitename"] = "cn8"
	err = s.dbapi.CreateItem("colo", "site", values)
	if err != nil {
		c.Error(err)
	}

	if items, err := s.dbapi.GetItemsByModel("colo", "site"); err != nil {
		c.Error(err)
	} else {
		if len(items) < 1 {
			c.Error("cannot get created item")
		}
	}
}

func (s *DBSuite) TestGetItemByIdentifier(c *C) {
	var err error
	values := make(map[string]string)
	values["sitename"] = "cn8"
	err = s.dbapi.CreateItem("colo", "site", values)
	if err != nil {
		c.Error(err)
	}

	if item, err := s.dbapi.GetItemByIdentifier("colo", "cn8"); err != nil {
		c.Error(err)
	} else {
		if item.Identifier != "cn8" {
			c.Error("cannot get correct item")
		}
	}
}

func (s *DBSuite) TestDeleteItemByIdentifier(c *C) {
	var err error
	values := make(map[string]string)
	values["sitename"] = "cn8"
	err = s.dbapi.CreateItem("colo", "site", values)
	if err != nil {
		c.Error(err)
	}

	if err := s.dbapi.DeleteItemByIdentifier("colo", "cn8"); err != nil {
		c.Error(err)
	} else {
		if _, err := s.dbapi.GetItemByIdentifier("colo", "cn8"); err == nil {
			c.Error("item cannot be deleted")
		}
	}
}

func (s *DBSuite) TestUpdateItem(c *C) {
	var err error
	values := make(map[string]string)
	values["sitename"] = "cn1"
	if s.dbapi.CreateItem("colo", "site", values) != nil {
		c.Error(err)
	}
	values["sitename"] = "cn8"
	if s.dbapi.CreateItem("colo", "site", values) != nil {
		c.Error(err)
	}

	values["sitename"] = "cn6"

	if err := s.dbapi.UpdateItem("colo", "cn8", values); err != nil {
		c.Error(err)
	}
	if _, err := s.dbapi.GetItemByIdentifier("colo", "cn6"); err != nil {
		c.Error("item cannot be updated")
	}
	if _, err := s.dbapi.GetItemByIdentifier("colo", "cn8"); err == nil {
		c.Error("item updated, but old item remained")
	}
	if _, err := s.dbapi.GetItemByIdentifier("colo", "cn1"); err != nil {
		c.Error("other items been affected by updating")
    }
}
