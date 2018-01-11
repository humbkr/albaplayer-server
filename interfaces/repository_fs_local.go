package interfaces

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"git.humbkr.com/jgalletta/alba-player/business"
	"git.humbkr.com/jgalletta/alba-player/domain"
	"github.com/dhowden/tag"
	mp3info "github.com/xhenner/mp3-go"
	"strings"
	"log"
	"github.com/go-gorp/gorp"
	"github.com/spf13/viper"
)

/**
Stores media metadata retrieved from different sources.
*/
type mediaMetadata struct {
	Format  	string
	Title   	string
	Album   	string
	Artist  	string
	AlbumArtist string
	Genre   	string
	Year    	string
	Track   	int
	Disc    	string // Format: <number>/<total>
	Picture 	*tag.Picture

	Duration int

	Path string
}

// Implements business.MediaFileRepository.
type LocalFilesystemRepository struct{
	AppContext *AppContext
}

func (r LocalFilesystemRepository) ScanMediaFiles(path string) (processed int, added int) {
	log.Println("Scanning files in " + path)

	// TODO Find a way to not have to get the datasource implementation.
	gorpDbMap, ok := r.AppContext.DB.(*gorp.DbMap)
	if !ok {
		log.Fatal("Cannot get underlying gorp dbmap")
	}

	dbTransaction, _ := gorpDbMap.Begin()
	scanDirectory(path, dbTransaction)
	dbTransaction.Commit()

	return
}

// Recursively browses a directory and import / update all the audio files in the database.
func scanDirectory(path string, dbTransaction *gorp.Transaction) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	}

	currentDir := filepath.Clean(path) + string(os.PathSeparator)

	// Collection of tracks found in the directory indexed by album.
	mediaFiles := make(map[string][]mediaMetadata)

	// Get all the entries in the current directory.
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	cpt := 0

	for _, file := range files {
		filePath := currentDir + file.Name()

		if file.IsDir() {
			// Recursion.
			scanDirectory(filePath, dbTransaction)
		} else if matched, _ := filepath.Match("*.mp3", file.Name()); matched {
			// Get ID3 metadata and add it to an array.
			metadata, err := getMetadataFromFile(filePath)
			if err == nil {
				// Add metadata info to the list of media files, sorting by albums.
				if len(metadata.Album) > 0 {
					mediaFiles[metadata.Album] = append(mediaFiles[metadata.Album], metadata)
				} else {
					mediaFiles[business.LibraryDefaultArtist] = append(mediaFiles[business.LibraryDefaultArtist], metadata)
				}
			}

			cpt++
		} else if matched, _ := filepath.Match("*.jpg", file.Name()); matched {
			// It's a good candidate for an album cover, so keep it somewhere.
		}
	}

	processMediaFiles(mediaFiles, dbTransaction)
}

func processMediaFiles(mediaFiles map[string][]mediaMetadata, dbTransaction *gorp.Transaction) {
	var err error

	// Process the media files per album.
	for _, album := range mediaFiles {

		// Here we try to figure out if the album is a compilation or not.
		// If at least 2 of the tracks have different artists, this must be a compilation.
		compilation := false

		currentArtist := album[0].Artist
		for index, metadataTrack := range album {
			if index > 0 && metadataTrack.Artist != currentArtist {
				compilation = true
				break
			}
		}

		// Now we process the metadata to populate the library.
		for _, metadataTrack := range album {
			if compilation {
				metadataTrack.AlbumArtist = business.LibraryDefaultCompilationArtist
			}

			var artistId int
			var albumId int

			artistId, err = processArtist(dbTransaction, &metadataTrack)
			if err != nil {
				// TODO devise a decent logging system.
				log.Println(err)
			}

			albumArtistId := artistId
			if compilation {
				// Get the artist id of "Various artists".
				var entities domain.Artists
				_, transErr := dbTransaction.Select(&entities, "SELECT * FROM artists WHERE name = ?", business.LibraryDefaultCompilationArtist)
				if transErr == nil {
					if len(entities) > 0 {
						albumArtistId = entities[0].Id
					}
				}
			}

			albumId, err = processAlbum(dbTransaction, &metadataTrack, albumArtistId)
			if err != nil {
				// TODO devise a decent logging system.
				log.Println(err)
			}

			_, err := processTrack(dbTransaction, &metadataTrack, artistId, albumId)
			if err != nil {
				// TODO devise a decent logging system.
				log.Println(err)
			}
		}
	}
}

// Checks if a media file physically exists.
func (r LocalFilesystemRepository) MediaFileExists(filepath string) bool {
	return fileExists(filepath)
}

// Writes a cover image.
func (r LocalFilesystemRepository) WriteCoverFile(file *domain.Cover, directory string) error {
	return writeCoverFile(file, directory)
}

// Deletes a cover image.
func (r LocalFilesystemRepository) RemoveCoverFile(file *domain.Cover, directory string) error {
	srcFileName := directory + string(os.PathSeparator) + file.Hash + file.Ext
	return os.Remove(srcFileName)
}

// Deletes all covers
func (r LocalFilesystemRepository) DeleteCovers() error {
	return os.RemoveAll(viper.GetString("Covers.Directory"))
}

// Saves an artist info in the database.
//
// Returns a artist id.
func processArtist(dbTransaction *gorp.Transaction, metadata *mediaMetadata) (id int, err error) {
	// Process artist if any.
	if metadata.Artist != "" {
		artist := domain.Artist{}
		// See if the artist exists and if so instanciate it with existing data.
		var entities domain.Artists
		_, transErr := dbTransaction.Select(&entities, "SELECT * FROM artists WHERE name = ?", metadata.Artist)
		if transErr == nil {
			if len(entities) > 0 {
				artist = entities[0]
			}
		}

		artist.Name = metadata.Artist

		if artist.Id != 0 {
			// Update.
			_, err = dbTransaction.Update(&artist)
		} else {
			// Insert new entity.
			err = dbTransaction.Insert(&artist)

		}
		if err == nil {
			id = artist.Id
		}

		return
	}

	return 0, errors.New("no artist to process")
}

// Saves an album info in the database.
//
// Returns a album id.
func processAlbum(dbTransaction *gorp.Transaction, metadata *mediaMetadata, artistId int) (id int, err error) {
	if metadata.Album != "" {
		album := domain.Album{}
		// See if the album exists and if so instanciate it with existing data.
		var entities domain.Albums
		_, transErr := dbTransaction.Select(&entities, "SELECT * FROM albums WHERE title = ? AND artist_id = ?", metadata.Album, artistId)
		if transErr == nil {
			if len(entities) > 0 {
				album = entities[0]
			}
		}

		album.Title = metadata.Album
		album.ArtistId = artistId
		// TODO Track all the years from an album tracks and compute the final value.
		album.Year = metadata.Year

		if album.Id != 0 {
			// Update.
			_, err = dbTransaction.Update(&album)
		} else {
			// Insert new entity.
			err = dbTransaction.Insert(&album)
		}

		if err == nil {
			id = album.Id
		}

		return
	}

	return 0, errors.New("no album to process")
}

// Saves a track info in the database.
//
// Returns a track id.
//
// TODO: Use checksum from tag.Sum() to search for existing track for an artist + album.
func processTrack(dbTransaction *gorp.Transaction, metadata *mediaMetadata, artistId int, albumId int) (id int, err error) {
	if metadata.Title == "" {
		return 0, errors.New("no track title provided")
	}

	track := domain.Track{}
	// See if the track exists and if so instanciate it with existing data.
	var entities domain.Tracks
	_, transErr := dbTransaction.Select(&entities, "SELECT * FROM tracks WHERE title = ? AND artist_id = ? AND album_id = ?", metadata.Title, artistId, albumId)
	if transErr == nil {
		if len(entities) > 0 {
			track = entities[0]
		}
	}

	track.Title = metadata.Title
	track.ArtistId = artistId
	track.AlbumId = albumId
	track.Number = metadata.Track
	track.Disc = metadata.Disc
	track.Genre = metadata.Genre
	track.Duration = metadata.Duration
	track.Path = metadata.Path

	coverId, err := processCover(dbTransaction, &track)
	if err == nil {
		track.CoverId = coverId
	}

	if track.Id != 0 {
		// Update.
		_, err = dbTransaction.Update(&track)
	} else {
		// Insert new entity.
		err = dbTransaction.Insert(&track)
	}

	if err == nil {
		id = track.Id
	}

	return
}

// Saves a cover info in the database and filesystem.
//
// Returns a cover id.
func processCover(dbTransaction *gorp.Transaction, track *domain.Track) (id int, err error) {
	coverFile, _, err := getMediaCoverFile(track, viper.GetString("Library.CoverPreferredSource"))
	if err == nil {
		cover := coverFile
		cover.Path = coverFile.Hash + coverFile.Ext

		var coverFromDb domain.Cover
		coverExistsErr := dbTransaction.SelectOne(&coverFromDb, "SELECT * FROM covers WHERE hash = ?", coverFile.Hash)
		if coverExistsErr == nil && coverFromDb.Id != 0 {
			// Nothing to do about the cover, just return the cover id to be used to link it to the track.
			id = coverFromDb.Id
			return
		}

		// Else we have to save a new cover to database.
		err = dbTransaction.Insert(&cover)
		// And to filesystem.
		if err == nil && cover.Id != 0 {
			// Save image file.
			err = writeCoverFile(&cover, viper.GetString("Covers.Directory"))
		}
	}

	return
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

	tags, errTags := tag.ReadFrom(file)
	if errTags == nil {
		// Get all we can from the common tags.
		info.Format = string(tags.FileType())
		info.Title = sanitizeString(tags.Title())
		info.Album = sanitizeString(tags.Album())
		info.AlbumArtist = sanitizeString(tags.AlbumArtist())
		info.Artist = sanitizeString(tags.Artist())
		info.Genre = sanitizeString(tags.Genre())
		if tags.Year() != 0 {
			info.Year = strconv.Itoa(tags.Year())
		}
		info.Track, _ = tags.Track()
		info.Picture = tags.Picture()

		number, total := tags.Disc()
		// Don't store disc info if there's only one disc.
		if total > 1 {
			info.Disc = strconv.Itoa(number) + "/" + strconv.Itoa(total)
		}
	}

	// If the track has no title, fallback to the filename.
	if info.Title == "" {
		_, f := path.Split(filePath)
		extension := filepath.Ext(f)
		filename := filepath.Base(f)
		info.Title = filename[0 : len(filename) - len(extension)]
	}

	// If the file is an mp3, get some more info.
	if strings.ToLower(filepath.Ext(filePath)) == "mp3" {
		mp3Info, err := mp3info.Examine(filePath, false)
		if err == nil {
			info.Duration = int(math.Floor(mp3Info.Length + .5))
		}
	}

	// Set the filepath.
	info.Path = filePath

	return
}

// Gets media cover image.
//
// Get data from media metadata and / or image file located in the same directory.
func getMediaCoverFile(track *domain.Track, preferredSource string) (cover domain.Cover, source string, err error) {
	if preferredSource == "file" {
		cover, err = getMediaCoverFromFolder(track.Path)
		if err == nil {
			source = preferredSource
			return
		}

		cover, err = getMediaCoverFromMetadata(track.Path)
		source = "tag"
		return
	} else {
		cover, err = getMediaCoverFromMetadata(track.Path)
		if err == nil {
			source = preferredSource
			return
		}

		cover, err = getMediaCoverFromFolder(track.Path)
		source = "file"
		return
	}

	return cover, "", errors.New("no cover image found")
}

// Get media cover from file located in the media file folder.
//
// Returns the info for the first image file that matches.
func getMediaCoverFromFolder(mediaPath string) (cover domain.Cover, err error) {
	var validCoverExtensions = []string{
		".jpg",
		".jpeg",
		".png",
		".gif",
	}

	var validCoverNames = []string{
		"cover",
		"artwork",
		"album",
		"front",
	}

	directory := filepath.Dir(mediaPath)
	for _, name := range validCoverNames {
		for _, ext := range validCoverExtensions {
			fileToLookFor := directory + string(os.PathSeparator) + name + ext
			if fileExists(fileToLookFor) {
				fileContent, errRead := ioutil.ReadFile(fileToLookFor)
				if errRead == nil {
					reader := bytes.NewReader(fileContent)
					hash, errSum := md5Checksum(reader)
					if errSum == nil {
						cover.Ext = ext
						cover.Hash = hash
						cover.Content = fileContent

						return cover, nil
					}
				}
			}
		}
	}

	return cover, errors.New("no cover found in folder")
}

// Get media cover from media file metadata.
//
// Returns cover info if found.
func getMediaCoverFromMetadata(mediaPath string) (cover domain.Cover, err error) {
	metadata, errMeta := getMetadataFromFile(mediaPath)
	if errMeta == nil && metadata.Picture != nil {
		reader := bytes.NewReader(metadata.Picture.Data)
		hash, errSum := md5Checksum(reader)
		if errSum == nil {
			// TODO Change this "hack".
			if metadata.Picture.Ext == "" {
				cover.Ext = ".jpg"
			} else {
				cover.Ext = "." + metadata.Picture.Ext
			}
			cover.Hash = hash
			cover.Content = metadata.Picture.Data

			return
		}
	}

	return cover, errors.New("no cover found in metadata")
}

// Writes a cover image to disk.
func writeCoverFile(file *domain.Cover, directory string) error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.MkdirAll(directory, 0755)
	}

	destFileName := directory + string(os.PathSeparator) + file.Hash + file.Ext
	return ioutil.WriteFile(destFileName, file.Content, 755)
}

// Checks if a file exists on disk.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Computes the checksum of a stream of content.
func md5Checksum(reader io.Reader) (hash string, err error) {
	hasher := md5.New()
	_, err = io.Copy(hasher, reader)
	if err != nil {
		return
	}

	hash = hex.EncodeToString(hasher.Sum(nil))
	return
}

// Removes spaces and nul character from a string.
func sanitizeString(s string) string {
	return strings.Trim(strings.TrimSpace(s), "\x00")
}
