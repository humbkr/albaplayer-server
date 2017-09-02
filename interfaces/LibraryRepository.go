package interfaces

import (
	"os"

	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"math"

	"git.humbkr.com/jgalletta/alba-player/domain"
	"github.com/spf13/viper"
	mp3info "github.com/xhenner/mp3-go"
	"github.com/dhowden/tag"
)

type LibraryRepository struct {
	AppContext *AppContext
}

/**
Stores media metadata retrieved from different sources.
 */
type mediaMetadata struct{
	Format string
	Title string
	Album string
	Artist string
	Genre string
	Year string
	Track int
	Disc string // Format: <number>/<total>

	Duration int
}

func (lr LibraryRepository) Update() {
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

TODO: Use checksum from tag.Sum() to search for existing track for an artist + album.
*/
func (lr LibraryRepository) scanFolder(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return
	}

	filepath.Walk(filePath, func(filePath string, fileInfo os.FileInfo, err error) error {
		if matched, _ := filepath.Match("*.mp3", fileInfo.Name()); matched {
			metadata, err := getMetadataFromFile(filePath)

			if err == nil {
				var artistId int
				var albumId int

				// Process artist if any.
				if metadata.Artist != "" {
					// See if the artist exists.
					artistRepo := ArtistRepository{AppContext: lr.AppContext}
					artist, err := artistRepo.FindByName(metadata.Artist)
					if err != nil {
						artist = domain.Artist{}
					}

					artist.Name = metadata.Artist

					errInsert := artistRepo.Save(&artist)
					if errInsert != nil {
						fmt.Println(errInsert)
					}

					artistId = artist.Id
				}

				// Process album if any.
				if metadata.Album != "" {
					// See if the album exists.
					albumRepo := AlbumRepository{AppContext: lr.AppContext}
					album, err := albumRepo.FindByName(metadata.Album, artistId)
					if err != nil {
						album = domain.Album{}
					}

					album.Title = metadata.Album
					album.ArtistId = artistId
					// TODO Track all the years from an album tracks and compute the final value.
					album.Year = metadata.Year

					errInsert := albumRepo.Save(&album)
					if errInsert != nil {
						fmt.Println(errInsert)
					}

					albumId = album.Id
				}

				// Process the track (there is always a track).
				// See if the track exists.
				trackRepo := TrackRepository{AppContext: lr.AppContext}
				track, err := trackRepo.FindByName(metadata.Title, artistId, albumId)

				if err != nil {
					track = domain.Track{}
				}

				track.Title = metadata.Title
				track.ArtistId = artistId
				track.AlbumId = albumId
				track.Number = metadata.Track
				track.Disc = metadata.Disc
				track.Genre = metadata.Genre
				track.Duration = metadata.Duration
				track.Path = filePath

				errInsert := trackRepo.Save(&track)
				if errInsert != nil {
					fmt.Println(errInsert)
				}
			}
		}

		return nil
	})
}

/**
Get media matadata from a file.

Uses multiple libraries to get a maximum of info depending on the format.
 */
func getMetadataFromFile(filePath string) (info mediaMetadata, err error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	defer file.Close()

	if err != nil {
		return
	}
	tags, err := tag.ReadFrom(file)
	if err != nil {
		return
	}

	// Get all we can from the common tags.
	info.Title = tags.Title()
	info.Album = tags.Album()
	info.Artist = tags.Artist()
	info.Genre = tags.Genre()
	info.Track, _ = tags.Track()

	number, total := tags.Disc()
	// Don't store disc info if there's only one disc.
	if total > 1 {
		info.Disc = strconv.Itoa(number) + "/" + strconv.Itoa(total)
	}

	// If the track has no title, fallback to the filename.
	if  info.Title == "" {
		_, f := path.Split(filePath)
		var extension = filepath.Ext(f)
		info.Title = f[0 : len(f)-len(extension)]
	}

	// If the file is an mp3, get some more info.
	if tags.FileType() == tag.MP3 {
		mp3Info, err := mp3info.Examine(filePath, false)
		if err == nil {
			info.Duration = int(math.Floor(mp3Info.Length + .5))
		}
	}

	return
}
