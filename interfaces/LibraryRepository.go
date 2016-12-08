package interfaces

import (
	"os"
	"path/filepath"

	"fmt"

	"git.humbkr.com/jgalletta/alba-player/domain"
	id3 "github.com/mikkyang/id3-go"
)

type LibraryRepository struct {
	AppContext *AppContext
}

func (lr LibraryRepository) Update() {
	// TODO stub
	lr.scanFolder("/home/humbkr/Music")
}

/*
Erase all collection data.
*/
func (lr LibraryRepository) Erase() {
	lr.AppContext.DB.Exec("DELETE FROM tracks")
	lr.AppContext.DB.Exec("DELETE FROM sqlite_sequence WHERE name = 'tracks'")
	lr.AppContext.DB.Exec("DELETE FROM albums")
	lr.AppContext.DB.Exec("DELETE FROM sqlite_sequence WHERE name = 'albums'")
	lr.AppContext.DB.Exec("DELETE FROM artists")
	lr.AppContext.DB.Exec("DELETE FROM sqlite_sequence WHERE name = 'artists'")
}

/**
Recursively browse a directory and import / update all the audio files in the database.
*/
func (lr LibraryRepository) scanFolder(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	}

	filepath.Walk(path, func(path string, fileInfo os.FileInfo, err error) error {
		if matched, _ := filepath.Match("*.mp3", fileInfo.Name()); matched {
			// Read idTag info.
			mp3File, err := id3.Open(path)

			if err == nil {
				var artistId int
				var albumId int

				if artistName := mp3File.Artist(); artistName != "" {
					// See if the artist exists.
					artistRepo := ArtistRepository{AppContext: lr.AppContext}
					artist, err := artistRepo.FindByName(artistName)

					if err != nil {
						artist = domain.Artist{}
					}

					artist.Name = artistName
					errInset := artistRepo.Save(&artist)
					if errInset != nil {
						fmt.Println(errInset)
					}

					artistId = artist.Id
				}

				if albumName := mp3File.Album(); albumName != "" {
					// See if the album exists.
					albumRepo := AlbumRepository{AppContext: lr.AppContext}
					album, err := albumRepo.FindByName(albumName, artistId)

					if err != nil {
						album = domain.Album{}
					}

					album.Name = albumName
					album.ArtistId = artistId
					album.Year = mp3File.Year()
					albumRepo.Save(&album)

					albumId = album.Id
				}

				if trackName := mp3File.Title(); trackName != "" {
					// See if the album exists.
					trackRepo := TrackRepository{AppContext: lr.AppContext}
					track, err := trackRepo.FindByName(trackName, artistId, albumId)

					if err != nil {
						track = domain.Track{}
					}

					track.Name = trackName
					track.ArtistId = artistId
					track.AlbumId = albumId
					fmt.Println(mp3File.Frame("TRCK"))
					fmt.Println(mp3File.Frame("TLEN"))
					fmt.Println(mp3File.Frame("MCDI"))
					trackRepo.Save(&track)
				}
			}

			mp3File.Close()
		}

		return nil
	})

}
