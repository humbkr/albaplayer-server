package main

import (
	"fmt"
	"log"

	"git.humbkr.com/jgalletta/alba-player/business"
	"git.humbkr.com/jgalletta/alba-player/interfaces"
	"github.com/mnmtanish/go-graphiql"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"net/http"
	gqlHandler "github.com/graphql-go/handler"
	//"github.com/mnmtanish/go-graphiql"
	"github.com/rs/cors"
)

func main() {
	// Load app configuration.
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
/*
	libraryInteractor.EraseLibrary()
	t := time.Now()
	fmt.Println(t.Format("15:04:05"))
	libraryInteractor.UpdateLibrary()
	t2 := time.Now()
	fmt.Println(t2.Format("15:04:05"))
*/

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
