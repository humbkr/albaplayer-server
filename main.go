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
	"github.com/mnmtanish/go-graphiql"
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
	libraryInteractor := new(business.LibraryInteractor)
	libraryInteractor.ArtistRepository = interfaces.ArtistRepository{AppContext: appContext}
	libraryInteractor.AlbumRepository = interfaces.AlbumRepository{AppContext: appContext}
	libraryInteractor.TrackRepository = interfaces.TrackRepository{AppContext: appContext}
	libraryInteractor.LibraryRepository = interfaces.LibraryRepository{AppContext: appContext}

	// Instanciate the main Queue.
	nowPlaying := new(business.Queue)
	nowPlaying.Library = libraryInteractor

	// STUB: instanciate the database for tests.
	libraryInteractor.EraseLibrary()
	libraryInteractor.UpdateLibrary()

	// Initialize GraphQL stuff.
	graphQLInteractor := interfaces.NewGraphQLInteractor(libraryInteractor)

	// Create a graphl-go HTTP handler with our previously defined schema
	// and we also set it to return pretty JSON output.
	apiHandler := gqlHandler.New(&gqlHandler.Config{
		Schema: &graphQLInteractor.Schema,
		Pretty: true,
	})

	// Serve a GraphQL endpoint at `/graphql`.
	http.Handle("/graphql", apiHandler)

	// Serve graphiql.
	http.HandleFunc("/", graphiql.ServeGraphiQL)

	// Launch the server.
	http.ListenAndServe(":" + viper.GetString("Server.Port"), nil)
	fmt.Printf("Server is up: http://localhost:%s/graphql", viper.GetString("Server.Port"))
}
