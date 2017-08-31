package main

import (
	"fmt"
	"log"

	"git.humbkr.com/jgalletta/alba-player/business"
	"git.humbkr.com/jgalletta/alba-player/interfaces"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"net/http"
	gqlHandler "github.com/graphql-go/handler"
)

func main() {
	// Load app configuration.
	viper.SetConfigName("app_config")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error loading config file: %s \n", err))
	}

	// Initialize logging system.
	log.SetOutput(&lumberjack.Logger{
		Filename:   viper.GetString("Log.Path") + viper.GetString("Log.File"),
		MaxSize:    500, // Megabytes.
		MaxBackups: 3,
		MaxAge:     28, // Days.
	})

	// Create app context.
	appContext, err := interfaces.NewAppContext()
	if err != nil {
		panic(fmt.Errorf("Error during the application context creation: %s \n", err))
	}

	// Instanciate all we need to work on the media library.
	collectionInteractor := new(business.CollectionInteractor)
	collectionInteractor.ArtistRepository = interfaces.ArtistRepository{AppContext: appContext}
	collectionInteractor.AlbumRepository = interfaces.AlbumRepository{AppContext: appContext}
	collectionInteractor.TrackRepository = interfaces.TrackRepository{AppContext: appContext}
	collectionInteractor.LibraryRepository = interfaces.LibraryRepository{AppContext: appContext}

	// Instanciate the main Queue.
	nowPlaying := new(business.Queue)
	nowPlaying.Library = collectionInteractor

	// STUB: instanciate the database for tests.
	//collectionInteractor.EraseLibrary()
	//collectionInteractor.UpdateLibrary()

	// Initialize GraphQL stuff.
	graphQLInteractor := interfaces.NewGraphQLInteractor(collectionInteractor)

	// Create a graphl-go HTTP handler with our previously defined schema
	// and we also set it to return pretty JSON output.
	apiHandler := gqlHandler.New(&gqlHandler.Config{
		Schema: &graphQLInteractor.Schema,
		Pretty: true,
	})

	// Serve a GraphQL endpoint at `/graphql`.
	http.Handle("/graphql", apiHandler)

	// Launch the server.
	http.ListenAndServe(":" + viper.GetString("Server.Port"), nil)
}
