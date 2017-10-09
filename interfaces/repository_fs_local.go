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
	"github.com/spf13/viper"
	mp3info "github.com/xhenner/mp3-go"
	"strings"
	"log"
)

/**
Stores media metadata retrieved from different sources.
*/
type mediaMetadata struct {
	Format  string
	Title   string
	Album   string
	Artist  string
	Genre   string
	Year    string
	Track   int
	Disc    string // Format: <number>/<total>
	Picture *tag.Picture

	Duration int

	Path string
}

// Implements business.MediaFileRepository.
type LocalFilesystemRepository struct{}

// Recursively browses a directory and import / update all the audio files in the database.

// Returns the number of items processed and added.
func (r LocalFilesystemRepository) ScanMediaFiles(path string, interactor *business.LibraryInteractor) (processed int, added int) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	}

	filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
		if matched, _ := filepath.Match("*.mp3", fileInfo.Name()); matched {
			metadata, err := getMetadataFromFile(filePath)
			if err == nil {
				var artistId int
				var albumId int

				artistId, err = processArtist(interactor, &metadata)
				if err != nil {
					// TODO devise a decent logging system.
					log.Println(err)
				}

				albumId, err = processAlbum(interactor, &metadata, artistId)
				if err != nil {
					// TODO devise a decent logging system.
					log.Println(err)
				}

				_, err := processTrack(interactor, &metadata, artistId, albumId)
				if err != nil {
					// TODO devise a decent logging system.
					log.Println(err)
				} else {
					added++
				}
			}

			processed++
		}

		return nil
	})

	return
}

// Checks if a media file physically exists.
func (r LocalFilesystemRepository) MediaFileExists(filepath string) bool {
	return fileExists(filepath)
}

// Writes a cover image.
func (r LocalFilesystemRepository) WriteCoverFile(file *domain.Cover, directory string) error {
	destFileName := directory + string(os.PathSeparator) + file.Hash + file.Ext
	return ioutil.WriteFile(destFileName, file.Content, 755)
}

// Deletes a cover image.
func (r LocalFilesystemRepository) RemoveCoverFile(file *domain.Cover, directory string) error {
	srcFileName := directory + string(os.PathSeparator) + file.Hash + file.Ext
	return os.Remove(srcFileName)
}

// Saves an artist info in the database.
//
// Returns a artist id.
func processArtist(interactor *business.LibraryInteractor, metadata *mediaMetadata) (id int, err error) {
	// Process artist if any.
	if metadata.Artist != "" {
		artist := domain.Artist{}
		// See if the artist exists and if so instanciate it with existing data.
		artist, _ = interactor.ArtistRepository.GetByName(metadata.Artist)

		artist.Name = metadata.Artist

		err = interactor.SaveArtist(&artist)
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
func processAlbum(interactor *business.LibraryInteractor, metadata *mediaMetadata, artistId int) (id int, err error) {
	if metadata.Album != "" {
		album := domain.Album{}
		// See if the album exists and if so instanciate it with existing data.
		album, _ = interactor.AlbumRepository.GetByName(metadata.Album, artistId)

		album.Title = metadata.Album
		album.ArtistId = artistId
		// TODO Track all the years from an album tracks and compute the final value.
		album.Year = metadata.Year

		err = interactor.SaveAlbum(&album)
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
func processTrack(interactor *business.LibraryInteractor, metadata *mediaMetadata, artistId int, albumId int) (id int, err error) {
	if metadata.Title == "" {
		return 0, errors.New("no track title provided")
	}

	track := domain.Track{}
	// See if the track exists and if so instanciate it with existing data.
	track, _ = interactor.TrackRepository.GetByName(metadata.Title, artistId, albumId)

	track.Title = metadata.Title
	track.ArtistId = artistId
	track.AlbumId = albumId
	track.Number = metadata.Track
	track.Disc = metadata.Disc
	track.Genre = metadata.Genre
	track.Duration = metadata.Duration
	track.Path = metadata.Path

	// Process cover.
	coverFile, _, err := getMediaCoverFile(track, viper.GetString("Library.CoverPreferredSource"))
	if err == nil {
		cover := coverFile
		cover.Path = coverFile.Hash + coverFile.Ext
		errSave := interactor.SaveCover(&cover)
		if errSave == nil {
			// Link track to cover.
			track.CoverId = cover.Id
		}
	}

	err = interactor.SaveTrack(&track)
	if err == nil {
		id = track.Id
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
		info.Title = tags.Title()
		info.Album = tags.Album()
		info.Artist = tags.Artist()
		info.Genre = tags.Genre()
		info.Year = strconv.Itoa(tags.Year())
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
		var extension = filepath.Ext(f)
		info.Title = f[0 : len(f)-len(extension)]
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
func getMediaCoverFile(track domain.Track, preferredSource string) (cover domain.Cover, source string, err error) {
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

		cover, err = getMediaCoverFromMetadata(track.Path)
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
			cover.Ext = metadata.Picture.Ext
			cover.Hash = hash
			cover.Content = metadata.Picture.Data

			return
		}
	}

	return cover, errors.New("no cover found in metadata")
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
