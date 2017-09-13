package business

import (
	"git.humbkr.com/jgalletta/alba-player/domain"
	"github.com/spf13/viper"
	"errors"
)

const CoverPreferredSourceImgFile = "file"
const CoverPreferredSourceMeta = "tag"

type LibraryRepository interface {
	Erase()
}

// Interface describing the storage mecanism for media.
type MediaFileRepository interface {
	// TODO Not abstract enough yet, we should not need a path but a reader or something.
	ScanMediaFiles(path string, interactor LibraryInteractor)
	MediaFileExists(filepath string) bool
	WriteCoverFile(file *domain.Cover, directory string) error
	RemoveCoverFile(file *domain.Cover, directory string) error
}

type LibraryInteractor struct {
	ArtistRepository  domain.ArtistRepository
	AlbumRepository domain.AlbumRepository
	TrackRepository domain.TrackRepository
	CoverRepository domain.CoverRepository
	// TODO Check if the library repo should be an interface here.
	LibraryRepository LibraryRepository
	MediaFileRepository MediaFileRepository
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

// Saves a cover.
func (interactor LibraryInteractor) SaveCover(cover *domain.Cover) error {
	coverId := interactor.CoverHashExists(cover.Hash)
	if coverId != 0 {
		// It becomes an update.
		cover.Id = coverId
	}

	// Save cover info to database.
	err := interactor.CoverRepository.Save(cover)
	if err == nil && coverId != 0 {
		// Save image file.
		err = interactor.MediaFileRepository.WriteCoverFile(cover, viper.GetString("Covers.Directory"))
	}

	return err
}

// Deletes a cover.
func (interactor LibraryInteractor) DeleteCover(cover *domain.Cover) error {
	// Save cover info to database.
	err := interactor.CoverRepository.Delete(cover)
	if err == nil {
		// Save image file.
		err = interactor.MediaFileRepository.RemoveCoverFile(cover, viper.GetString("Covers.Directory"))
	}

	return err
}

// Checks if a cover exists or not.
func (interactor LibraryInteractor) CoverExists(coverId int) bool {
	return interactor.CoverRepository.Exists(coverId)
}

// Checks if a cover exists or not by hash.
//
// Returns cover.Id if exists, else 0.
func (interactor LibraryInteractor) CoverHashExists(hash string) int {
	return interactor.CoverRepository.ExistsByHash(hash)
}

// TODO How to unit test this?
func (interactor LibraryInteractor) UpdateLibrary() {
	interactor.MediaFileRepository.ScanMediaFiles(viper.GetString("Library.Folder"), interactor)
}

// TODO How to unit test this?
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
			if !interactor.MediaFileRepository.MediaFileExists(track.Path) {
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
