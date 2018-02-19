package main

import (
	"fmt"
	"log"

	"git.humbkr.com/jgalletta/alba-player/internal/alba/business"
	"git.humbkr.com/jgalletta/alba-player/internal/alba/interfaces"
	"github.com/mnmtanish/go-graphiql"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"net/http"
	gqlHandler "github.com/graphql-go/handler"
	//"github.com/mnmtanish/go-graphiql"
	"github.com/rs/cors"
)

func main() {
	// Set default configuration.
	// Database.
	viper.SetDefault("DB.Driver", "sqlite3")
	viper.SetDefault("DB.File", "./alba.db")
	// Covers.
	viper.SetDefault("Covers.Directory", "./covers")
	viper.SetDefault("Covers.PreferredSource", "folder")
	// Logging.
	viper.SetDefault("Log.Enabled", true)
	viper.SetDefault("Log.File", "app.log")
	viper.SetDefault("Log.Path", "./")
	// Webserver.
	viper.SetDefault("Server.Port", "8888")
	// Library.
	viper.SetDefault("Library.Path", "")
	// Dev mode.
	viper.SetDefault("DevMode.Enabled", false)


	// Load app configuration from file.
	viper.SetConfigName("alba")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error loading config file: %s \n", err))
	}

	// Initialize logging system.
	log.SetOutput(&lumberjack.Logger{
		Filename:   viper.GetString("Log.Path") + viper.GetString("Log.File"),
		MaxSize:    10, // Megabytes.
		MaxBackups: 3,
		MaxAge:     15, // Days.
	})

	// Create app context.
	var appContext interfaces.AppContext
	datasource, err := interfaces.InitAlbaDatasource()
	if err != nil {
		panic(fmt.Errorf("Error during the application context creation: %s \n", err))
	}
	appContext.DB = datasource

	// Instanciate all we need to work on the media library.
	libraryInteractor := business.LibraryInteractor{}
	libraryInteractor.ArtistRepository = interfaces.ArtistDbRepository{AppContext: &appContext}
	libraryInteractor.AlbumRepository = interfaces.AlbumDbRepository{AppContext: &appContext}
	libraryInteractor.TrackRepository = interfaces.TrackDbRepository{AppContext: &appContext}
	libraryInteractor.CoverRepository = interfaces.CoverDbRepository{AppContext: &appContext}
	libraryInteractor.LibraryRepository = interfaces.LibraryDbRepository{AppContext: &appContext}
	libraryInteractor.MediaFileRepository = interfaces.LocalFilesystemRepository{AppContext: &appContext}

	// STUB: instanciate the database for tests.
	// libraryInteractor.EraseLibrary()
	// libraryInteractor.UpdateLibrary()

	// Initialize GraphQL stuff.
	graphQLInteractor := interfaces.NewGraphQLInteractor(&libraryInteractor)

	// Create a graphl-go HTTP handler with our previously defined schema
	// and set it to return pretty JSON output.
	graphQLHandler := gqlHandler.New(&gqlHandler.Config{
		Schema: &graphQLInteractor.Schema,
		Pretty: true,
	})

	// Serve a GraphQL endpoint at `/graphql`.
	// Make the server handle cross-domain requests.
	http.Handle("/graphql", cors.Default().Handler(graphQLHandler))

	// Serve media files streaming endpoint.
	// Makes the server handle cross-domain requests.
	mediaFilesHandler := interfaces.NewMediaStreamHandler(&libraryInteractor)
	http.Handle("/stream/", http.StripPrefix("/stream/", cors.Default().Handler(mediaFilesHandler)))

	// Serve media files streaming endpoint.
	// Makes the server handle cross-domain requests.
	coverFilesHandler := interfaces.NewCoverStreamHandler(&libraryInteractor)
	http.Handle("/covers/", http.StripPrefix("/covers/", cors.Default().Handler(coverFilesHandler)))

    if viper.GetBool("DevMode.Enabled") {
		// Serve graphiql.
		http.HandleFunc("/graphiql", graphiql.ServeGraphiQL)
	}

	// Launch the server.
	log.Printf("Server is up: http://localhost:%s/graphql\n", viper.GetString("Server.Port"))
	http.ListenAndServe(":" + viper.GetString("Server.Port"), nil)
}
