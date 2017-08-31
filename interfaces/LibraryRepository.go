package interfaces

import (
	"os"

	"fmt"
	"path"
	"path/filepath"

	"strings"

	"strconv"

	"math"

	"git.humbkr.com/jgalletta/alba-player/domain"
	id3 "github.com/mikkyang/id3-go"
	"github.com/spf13/viper"
	mp3info "github.com/xhenner/mp3-go"
)

type LibraryRepository struct {
	AppContext *AppContext
}

func (lr LibraryRepository) Update() {
	// TODO make this configurable.
	lr.scanFolder(viper.GetString("Library.Folder"))
}

/*
Erase all collection data.
*/
func (lr LibraryRepository) Erase() {
	// We have to delete the tables content AND reset the sequences for ID columns.
	lr.AppContext.DB.Exec("DELETE FROM tracks")
	lr.AppContext.DB.Exec("DELETE FROM sqlite_sequence WHERE name = 'tracks'")
	lr.AppContext.DB.Exec("DELETE FROM albums")
	lr.AppContext.DB.Exec("DELETE FROM sqlite_sequence WHERE name = 'albums'")
	lr.AppContext.DB.Exec("DELETE FROM artists")
	lr.AppContext.DB.Exec("DELETE FROM sqlite_sequence WHERE name = 'artists'")
}

/**
Recursively browses a directory and import / update all the audio files in the database.

TODO: Manage other filetypes than mp3.
*/
func (lr LibraryRepository) scanFolder(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return
	}

	filepath.Walk(filePath, func(filePath string, fileInfo os.FileInfo, err error) error {
		if matched, _ := filepath.Match("*.mp3", fileInfo.Name()); matched {
			// Read idTag info.
			mp3Tags, err := id3.Open(filePath)
			// TODO: err is not checked.
			mp3Info, err := mp3info.Examine(filePath, false)

			if err == nil {
				var artistId int
				var albumId int

				// Process artist if any.
				if artistName := mp3Tags.Artist(); artistName != "" {
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

				// Process album if any.
				if albumName := mp3Tags.Album(); albumName != "" {
					// See if the album exists.
					albumRepo := AlbumRepository{AppContext: lr.AppContext}
					album, err := albumRepo.FindByName(albumName, artistId)

					if err != nil {
						album = domain.Album{}
					}

					album.Name = albumName
					album.ArtistId = artistId
					// TODO Track all the years from an album tracks and compute the final value.
					album.Year = mp3Tags.Year()
					albumRepo.Save(&album)

					albumId = album.Id
				}

				// There is always a track, if the title metatag is not present, build one with the filename.
				trackName := mp3Tags.Title()
				if trackName == "" {
					_, f := path.Split(filePath)

					var extension = filepath.Ext(f)
					trackName = f[0 : len(f)-len(extension)]
				}

				//fmt.Println("Processing file: " + mp3Tags.Artist() + " - " + trackName)

				// See if the track exists.
				trackRepo := TrackRepository{AppContext: lr.AppContext}
				track, err := trackRepo.FindByName(trackName, artistId, albumId)

				if err != nil {
					track = domain.Track{}
				}

				track.Name = trackName
				track.ArtistId = artistId
				track.AlbumId = albumId
				track.Number = getTrackNumber(mp3Tags)
				track.Disc = getTrackDisc(mp3Tags)
				track.Duration = getDuration(mp3Info.Length)
				track.Path = filePath

				trackRepo.Save(&track)

			}

			mp3Tags.Close()
		}

		return nil
	})

}

/*
Returns the track number if present, else an empty string.
*/
func getTrackNumber(mp3Tags *id3.File) int {
	trackNumber := 0
	var number string

	newFrame := mp3Tags.Frame("TRCK")
	if newFrame != nil {
		number = newFrame.String()
	} else {
		oldFrame := mp3Tags.Frame("TRK")
		if oldFrame != nil {
			number = oldFrame.String()
		}
	}

	if number != "" {
		shards := strings.Split(number, "/")
		trackNumber, _ = strconv.Atoi(shards[0])
	}

	return trackNumber
}

/*
Returns the disc number if present, else an empty string.
*/
func getTrackDisc(mp3Tags *id3.File) string {
	tlenFrame := mp3Tags.Frame("MCDI")
	if tlenFrame != nil {
		shards := strings.Split(tlenFrame.String(), "/")
		return shards[0]
	} else {
		tleFrame := mp3Tags.Frame("MCI")
		if tleFrame != nil {
			shards := strings.Split(tleFrame.String(), "/")
			return shards[0]
		} else {
			return ""
		}
	}
}

// Get duration in seconds.
func getDuration(f float64) int {
	return int(math.Floor(f + .5))
}
