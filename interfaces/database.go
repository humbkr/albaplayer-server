package interfaces

import (
	"database/sql"

	"log"

	"git.humbkr.com/jgalletta/alba-player/domain"
	"github.com/go-gorp/gorp"
	"github.com/spf13/viper"
)

func InitDb() (dbmap *gorp.DbMap, err error) {
	connection, err := sql.Open(viper.GetString("DB.driver"), viper.GetString("DB.file"))
	if err != nil {
		return
	}

	// Check database is reachable.
	if err = connection.Ping(); err != nil {
		return
	}

	// Construct a gorp DbMap.
	dbmap = &gorp.DbMap{Db: connection, Dialect: gorp.SqliteDialect{}}

	// Bind tables to objects.
	dbmap.AddTableWithName(domain.Artist{}, "artists").SetKeys(true, "Id")
	dbmap.AddTableWithName(domain.Album{}, "albums").SetKeys(true, "Id")
	dbmap.AddTableWithName(domain.Track{}, "tracks").SetKeys(true, "Id")

	// Create the tables.
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		log.Fatalln("Create tables failed", err)
	}

	return
}
