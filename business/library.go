package business

import (
	"git.humbkr.com/jgalletta/alba-player/domain"
	"github.com/dhowden/tag"
	"path"
	"path/filepath"
	"math"
	"os"
	"fmt"
	"strconv"
	mp3info "github.com/xhenner/mp3-go"
	"github.com/spf13/viper"
	"errors"
)

type LibraryRepository interface {
	Erase()
}

type LibraryInteractor struct {
	ArtistRepository  domain.ArtistRepository
	AlbumRepository   domain.AlbumRepository
	TrackRepository   domain.TrackRepository
	// TODO Check if the library repo should be an interface here.
	LibraryRepository LibraryRepository
}

// Gets an artist from its id.
//
// If no artist found, returns an error.
func (interactor LibraryInteractor) GetArtist(artistId int) (domain.Artist, error) {
	return interactor.ArtistRepository.Get(artistId)
}

// Gets all artists.
//
// If no artists found, returns an empty collection.
func (interactor LibraryInteractor) GetAllArtists(hydrate bool) (domain.Artists, error) {
	return interactor.ArtistRepository.GetAll(hydrate)
}

// Saves an artist.
//
// Returns an error if the artist's name is empty.
func (interactor LibraryInteractor) SaveArtist(artist *domain.Artist) error {
	if artist.Name == "" {
		return errors.New("cannot save artist: empty name")
	}

	return interactor.ArtistRepository.Save(artist)
}

// Deletes an artist.
//
// Returns an error if no artistId provided.
func (interactor LibraryInteractor) DeleteArtist(artist *domain.Artist) error {
	if artist.Id == 0 {
		return errors.New("cannot delete artist: id not provided")
	}
	return interactor.ArtistRepository.Delete(artist)
}

// Checks if an artist exists or not.
func (interactor LibraryInteractor) ArtistExists(artistId int) bool {
	return interactor.ArtistRepository.Exists(artistId)
}

// Gets an album from its id.
//
// If no album found, returns an error.
func (interactor LibraryInteractor) GetAlbum(albumId int) (domain.Album, error) {
	return interactor.AlbumRepository.Get(albumId)
}

// Gets all albums.
//
// If no albums found, returns an empty collection.
func (interactor LibraryInteractor) GetAllAlbums(hydrate bool) (domain.Albums, error) {
	return interactor.AlbumRepository.GetAll(hydrate)
}

// Get all albums for a given artist.
//
// If hydrate == true, the albums tracks will be populated.
// If the artist doesnt exists, return an error.
func (interactor LibraryInteractor) GetAlbumsForArtist(artistId int, hydrate bool) (domain.Albums, error) {
	if !interactor.ArtistExists(artistId) {
		return domain.Albums{}, errors.New("cannot get albums: invalid artist ID")
	}
	return interactor.AlbumRepository.GetAlbumsForArtist(artistId, hydrate)
}

// Saves an album.
//
// An album cannot be saved without a title or if the related artist, if any, doesn't exists.
func (interactor LibraryInteractor) SaveAlbum(album *domain.Album) error {
	invalid := false
	var message string
	if album.Title == "" {
		invalid = true
		message = "cannot save album: empty title"
	}
	if album.ArtistId != 0 {
		if !interactor.ArtistExists(album.ArtistId) {
			invalid = true
			message = "cannot save album: invalid artist ID"
		}
	}

	if invalid {
		return errors.New(message)
	}

	return interactor.AlbumRepository.Save(album)
}

// Delete an album.
//
// Returns an error if no albumId provided.
func (interactor LibraryInteractor) DeleteAlbum(album *domain.Album) error {
	if album.Id == 0 {
		return errors.New("cannot delete album: id not provided")
	}

	return interactor.AlbumRepository.Delete(album)
}

// Checks if an album exists or not.
func (interactor LibraryInteractor) AlbumExists(albumId int) bool {
	return interactor.AlbumRepository.Exists(albumId)
}

// Gets atrack from its id.
//
// If no track found, returns an error.
func (interactor LibraryInteractor) GetTrack(trackId int) (domain.Track, error) {
	return interactor.TrackRepository.Get(trackId)
}

// Gets all tracks.
//
// If no tracks found, returns an empty collection.
func (interactor LibraryInteractor) GetAllTracks() (domain.Tracks, error) {
	return interactor.TrackRepository.GetAll()
}

// Get all tracks for a given album.
//
// If the album doesn't exists, return an error
func (interactor LibraryInteractor) GetTracksForAlbum(albumId int) (domain.Tracks, error) {
	if !interactor.AlbumExists(albumId) {
		return domain.Tracks{}, errors.New("cannot get tracks: invalid album ID")
	}

	return interactor.TrackRepository.GetTracksForAlbum(albumId)
}

// Saves a track.
//
// A track cannot be saved without a title or if the related artist or album, if any, doesn't exists.
func (interactor LibraryInteractor) SaveTrack(track *domain.Track) error {
	invalid := false
	var message string
	if track.Title == "" {
		invalid = true
		message = "cannot save track: empty title"
	}
	if track.Path == "" {
		invalid = true
		message = "cannot save track: empty path"
	}
	if track.ArtistId != 0 {
		if _, err := interactor.GetArtist(track.ArtistId); err != nil {
			invalid = true
			message = "cannot save track: invalid artist ID"
		}
	}
	if track.AlbumId != 0 {
		if _, err := interactor.GetAlbum(track.AlbumId); err != nil {
			invalid = true
			message = "cannot save track: invalid album ID"
		}
	}

	if invalid {
		return errors.New(message)
	}

	return interactor.TrackRepository.Save(track)
}

// Deletes a track.
//
// Returns an error if no trackId provided.
func (interactor LibraryInteractor) DeleteTrack(track *domain.Track) error {
	if track.Id == 0 {
		return errors.New("cannot delete track: id not provided")
	}

	return interactor.TrackRepository.Delete(track)
}

// Checks if a track exists or not.
func (interactor LibraryInteractor) TrackExists(trackId int) bool {
	return interactor.TrackRepository.Exists(trackId)
}

func (interactor LibraryInteractor) UpdateLibrary() {
	interactor.scanFolder(viper.GetString("Library.Folder"))
}

func (interactor LibraryInteractor) EraseLibrary() {
	interactor.LibraryRepository.Erase()
}

/**
Remove all dead files from library.
 */
func (interactor LibraryInteractor) CleanDeadFiles() {
	// Keep a trace of albums and artists to check after tracks deletion.
	var relatedAlbums map[int]int
	var relatedArtists map[int]int

	tracks, err := interactor.GetAllTracks()
	if err == nil {
		// Delete non existant tracks.
		for _, track := range tracks {
			if _, err := os.Stat(track.Path); os.IsNotExist(err) {
				interactor.DeleteTrack(&track)
				relatedAlbums[track.AlbumId]++
				relatedArtists[track.ArtistId]++
			}
		}

		// Delete albums if no more tracks in it.
		for albumId := range relatedAlbums {
			album, err := interactor.GetAlbum(albumId)
			if err == nil && len(album.Tracks) == 0 {
				interactor.DeleteAlbum(&album)
			}
		}

		// Delete artists if no more albums from them.
		for artistId := range relatedArtists {
			artist, err := interactor.GetArtist(artistId)
			if err == nil && len(artist.Albums) == 0 {
				interactor.DeleteArtist(&artist)
			}
		}
	}
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

/**
Recursively browses a directory and import / update all the audio files in the database.

TODO: Use checksum from tag.Sum() to search for existing track for an artist + album.
*/
func (interactor LibraryInteractor) scanFolder(filePath string) {
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
					artist, err := interactor.ArtistRepository.GetByName(metadata.Artist)
					if err != nil {
						artist = domain.Artist{}
					}

					artist.Name = metadata.Artist

					errInsert := interactor.SaveArtist(&artist)
					if errInsert != nil {
						fmt.Println(errInsert)
					}

					artistId = artist.Id
				}

				// Process album if any.
				if metadata.Album != "" {
					// See if the album exists.
					album, err := interactor.AlbumRepository.GetByName(metadata.Album, artistId)
					if err != nil {
						album = domain.Album{}
					}

					album.Title = metadata.Album
					album.ArtistId = artistId
					// TODO Track all the years from an album tracks and compute the final value.
					album.Year = metadata.Year

					errInsert := interactor.SaveAlbum(&album)
					if errInsert != nil {
						fmt.Println(errInsert)
					}

					albumId = album.Id
				}

				// Process the track (there is always a track).
				// See if the track exists.
				track, err := interactor.TrackRepository.GetByName(metadata.Title, artistId, albumId)

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

				errInsert := interactor.SaveTrack(&track)
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