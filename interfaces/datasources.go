package interfaces

import (
	"database/sql"

	"log"

	"git.humbkr.com/jgalletta/alba-player/domain"
	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

/**
Interface to a datasource to abstract the underlying storage mecanism.
 */
type Datasource interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	SelectOne(holder interface{}, query string, args ...interface{}) error
	Select(i interface{}, query string, args ...interface{}) ([]interface{}, error)
	Get(i interface{}, keys ...interface{}) (interface{}, error)
	Insert(list ...interface{}) error
	Update(list ...interface{}) (int64, error)
	Delete(list ...interface{}) (int64, error)
}

/**
Initialiase the application main datasource.
 */
func InitAlbaDatasource() (ds Datasource, err error) {
	connection, err := sql.Open(viper.GetString("DB.driver"), viper.GetString("DB.file"))
	if err != nil {
		return
	}

	// Check database is reachable.
	if err = connection.Ping(); err != nil {
		return
	}

	// Construct a gorp DbMap.
	dbmap := &gorp.DbMap{Db: connection, Dialect: gorp.SqliteDialect{}}

	// Bind tables to objects.
	dbmap.AddTableWithName(domain.Artist{}, "artists").SetKeys(true, "Id")
	dbmap.AddTableWithName(domain.Album{}, "albums").SetKeys(true, "Id")
	dbmap.AddTableWithName(domain.Track{}, "tracks").SetKeys(true, "Id")

	// Create the tables.
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		log.Fatalln("Create tables failed", err)
	}

	return dbmap, nil
}
