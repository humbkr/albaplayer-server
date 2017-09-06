package data

import (
	"database/sql"
	"git.humbkr.com/jgalletta/alba-player/interfaces"
	"git.humbkr.com/jgalletta/alba-player/domain"
	"log"
	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
	"encoding/csv"
	"os"
	"io"
)

/*
@file
Common stuff for repositories tests.
 */

const TestDatasourceFile = "../data/test.db"

// Initialise the application test datasource.
func CreateTestDatasource() (ds interfaces.Datasource, err error) {
	connection, err := sql.Open("sqlite3", TestDatasourceFile)
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

func ResetTestDataSource(ds interfaces.Datasource) error {
	if dbmap, ok := ds.(*gorp.DbMap); ok == true {
		dbmap.Exec("DELETE FROM tracks")
		dbmap.Exec("DELETE FROM sqlite_sequence WHERE name = 'tracks'")
		dbmap.Exec("DELETE FROM albums")
		dbmap.Exec("DELETE FROM sqlite_sequence WHERE name = 'albums'")
		dbmap.Exec("DELETE FROM artists")
		dbmap.Exec("DELETE FROM sqlite_sequence WHERE name = 'artists'")

		initTestDataSource(ds)
	}
	return nil
}

func CloseTestDataSource(ds interfaces.Datasource) error {
	if dbmap, ok := ds.(*gorp.DbMap); ok == true {
		dbmap.Db.Close()

		return os.Remove(TestDatasourceFile)
	}
	return nil
}

// Populate the database with test data from csv.
func initTestDataSource(ds interfaces.Datasource) (err error) {
	if dbmap, ok := ds.(*gorp.DbMap); ok == true {
		// Artists.
		file, _ := os.OpenFile("../data/artists.csv", os.O_RDONLY, 0666)
		r := csv.NewReader(file)
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			// Insert the row in database.
			dbmap.Exec("INSERT INTO artists(id, name) VALUES(?, ?)", record[0], record[1])
		}
		file.Close()

		// Albums.
		file, _ = os.OpenFile("../data/albums.csv", os.O_RDONLY, 0666)
		r = csv.NewReader(file)
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			// Insert the row in database.
			dbmap.Exec(
				"INSERT INTO albums(id, artist_id, title, year) VALUES(?, ?, ?, ?)",
				record[0],
				record[1],
				record[2],
				record[3],
			)
		}
		file.Close()

		// Tracks.
		file, _ = os.OpenFile("../data/tracks.csv", os.O_RDONLY, 0666)
		r = csv.NewReader(file)
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			// Insert the row in database.
			dbmap.Exec(
				"INSERT INTO tracks(id, album_id, artist_id, title, disc, number, duration, genre, path) VALUES(?, ?, ?, ?, ? ,? ,?, ? , ?)",
				record[0],
				record[1],
				record[2],
				record[3],
				record[4],
				record[5],
				record[6],
				record[7],
				record[8],
			)
		}
		file.Close()
	}

	return nil
}
