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
	"github.com/rs/cors"
	"time"
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
	var appContext interfaces.AppContext
	datasource, err := interfaces.InitAlbaDatasource()
	if err != nil {
		panic(fmt.Errorf("Error during the application context creation: %s \n", err))
	}
	appContext.DB = datasource

	// Instanciate all we need to work on the media library.
	libraryInteractor := new(business.LibraryInteractor)
	libraryInteractor.ArtistRepository = interfaces.ArtistDbRepository{AppContext: &appContext}
	libraryInteractor.AlbumRepository = interfaces.AlbumDbRepository{AppContext: &appContext}
	libraryInteractor.TrackRepository = interfaces.TrackDbRepository{AppContext: &appContext}
	libraryInteractor.CoverRepository = interfaces.CoverDbRepository{AppContext: &appContext}
	libraryInteractor.LibraryRepository = interfaces.LibraryDbRepository{AppContext: &appContext}
	libraryInteractor.MediaFileRepository = interfaces.LocalFilesystemRepository{AppContext: &appContext}

	// STUB: instanciate the database for tests.
	//libraryInteractor.EraseLibrary()
	//libraryInteractor.UpdateLibrary()


	// Instanciate the main Queue.
	// TODO warning, only works for one user.
	//queue := business.GetQueueInstance()
	//queue.Library = libraryInteractor


	//queue.AppendAlbum(1)

	// Initialize GraphQL stuff.
	graphQLInteractor := interfaces.NewGraphQLInteractor(libraryInteractor)

	// Create a graphl-go HTTP handler with our previously defined schema
	// and set it to return pretty JSON output.
	graphQLHandler := gqlHandler.New(&gqlHandler.Config{
		Schema: &graphQLInteractor.Schema,
		Pretty: true,
	})

	// Make the server handle cross-domain requests.
	apiHandler := cors.Default().Handler(graphQLHandler)

	// Serve a GraphQL endpoint at `/graphql`.
	http.Handle("/graphql", apiHandler)

	// Serve graphiql.
	http.HandleFunc("/", graphiql.ServeGraphiQL)

	// Launch the server.
	fmt.Printf("Server is up: http://localhost:%s/graphql", viper.GetString("Server.Port"))
	http.ListenAndServe(":" + viper.GetString("Server.Port"), nil)

}
