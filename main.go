package main

import (
	"fmt"
	"log"

	"git.humbkr.com/jgalletta/alba-player/business"
	"git.humbkr.com/jgalletta/alba-player/interfaces"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
)

func main() {
	// Load app configuration.
	viper.SetConfigName("app_config")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// Initialize log system.
	log.SetOutput(&lumberjack.Logger{
		Filename:   viper.GetString("Log.Path") + viper.GetString("Log.File"),
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	})

	// Create app context.
	appContext, err := interfaces.NewAppContext()
	if err != nil {
		panic(fmt.Errorf("Error during the application context creation: %s \n", err))
	}

	collectionInteractor := new(business.CollectionInteractor)
	collectionInteractor.ArtistRepository = interfaces.ArtistRepository{AppContext: appContext}
	collectionInteractor.AlbumRepository = interfaces.AlbumRepository{AppContext: appContext}
	collectionInteractor.TrackRepository = interfaces.TrackRepository{AppContext: appContext}
	collectionInteractor.LibraryRepository = interfaces.LibraryRepository{AppContext: appContext}

	collectionInteractor.EraseLibrary()
	collectionInteractor.UpdateLibrary()
}
