package interfaces

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"github.com/stretchr/testify/assert"
	"git.humbkr.com/jgalletta/alba-player/internal/alba/domain"
	"git.humbkr.com/jgalletta/alba-player/internal/alba/business"
	"github.com/spf13/viper"
	"os"
	"log"
)

type LocalFSRepoTestSuite struct {
	suite.Suite
	LocalFSRepository LocalFilesystemRepository
}

// Go testing framework entry point.
func TestLocalFSRepoTestSuite(t *testing.T) {
	suite.Run(t, new(LocalFSRepoTestSuite))
}

func (suite *LocalFSRepoTestSuite) SetupSuite() {
	coversDir := os.TempDir() + string(os.PathSeparator) + "covers"
	if _, err := os.Stat(coversDir); os.IsNotExist(err) {
		os.Mkdir(coversDir, 0755)
	}
	viper.Set("Covers.Directory", coversDir)

	ds, err := createTestDatasource()
	if err != nil {
		log.Fatal(err)
	}
	appContext := AppContext{DB: ds}
	suite.LocalFSRepository = LocalFilesystemRepository{AppContext: &appContext}
}

func (suite *LocalFSRepoTestSuite) TearDownSuite() {
	coversDir := os.TempDir() + string(os.PathSeparator) + "covers"
	if _, err := os.Stat(coversDir); err == nil {
		os.Remove(coversDir)
	}
}

func (suite *LocalFSRepoTestSuite) SetupTest() {}

/*
Blackbox tests.
 */

func (suite *LocalFSRepoTestSuite) TestScanMediaFiles() {
	processed, added := suite.LocalFSRepository.ScanMediaFiles(TestFSLibDir)
	assert.Equal(suite.T(), 7, processed)
	assert.Equal(suite.T(), 7, added)
}

func (suite *LocalFSRepoTestSuite) TestMediaFileExists() {
	// Test with an existing media file.
	exists := suite.LocalFSRepository.MediaFileExists(TestFSLibDir + "/no artist - no album - no title.mp3")
	assert.Equal(suite.T(), true, exists)

	// Test with an existing media file.
	exists = suite.LocalFSRepository.MediaFileExists(TestFSLibDir + "/whatever.mp3")
	assert.Equal(suite.T(), false, exists)
}

// TODO test LocalFSRepository.WriteCoverFile.
// TODO test LocalFSRepository.RemoveCoverFile.
// TODO test LocalFSRepository.DeleteCovers.

/*
Below are whitebox (internal) tests.
 */

// TODO Cannot test these functions because gorp.Transaction is not abstracted.
func (suite *LocalFSRepoTestSuite) TestProcessArtist() {}
func (suite *LocalFSRepoTestSuite) TestProcessAlbum() {}
func (suite *LocalFSRepoTestSuite) TestProcessTrack() {}
func (suite *LocalFSRepoTestSuite) TestProcessCover() {}

// Test the preferred source selection only.
func (suite *LocalFSRepoTestSuite) TestGetMediaCoverFile() {
	// Test to get the cover from an image file when a cover tag is present.
	track := domain.Track{Path: TestFSLibDir + "/no artist - album 1 - Track 1.mp3"}
	cover, source, err := getMediaCoverFile(&track, business.CoverPreferredSourceImgFile)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), business.CoverPreferredSourceImgFile, source)
	assert.NotEmpty(suite.T(), cover.Hash)
	assert.NotEmpty(suite.T(), cover.Content)
	assert.Equal(suite.T(), ".jpg", cover.Ext)

	// Test to get the cover from metatags when an image file is present.
	track = domain.Track{Path: TestFSLibDir + "/no artist - album 1 - Track 1.mp3"}
	cover, source, err = getMediaCoverFile(&track, business.CoverPreferredSourceMeta)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), business.CoverPreferredSourceMeta, source)
	assert.NotEmpty(suite.T(), cover.Hash)
	assert.NotEmpty(suite.T(), cover.Content)
	assert.Equal(suite.T(), ".jpg", cover.Ext)

	return
}

func (suite *LocalFSRepoTestSuite) TestGetMediaCoverFromFolder() {
	// Test to get a jpg cover.
	track := domain.Track{Path: TestFSLibDir + "/artist 1/artist 1 - album 1/Artist 1 - Album 1 - Track 1.mp3"}
	cover, err := getMediaCoverFromFolder(track.Path)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), cover.Hash)
	assert.NotEmpty(suite.T(), cover.Content)
	assert.Equal(suite.T(), ".jpg", cover.Ext)

	// Test to get a png cover.
	track = domain.Track{Path: TestFSLibDir + "/artist 1/artist 1 - album 2/Artist 1 - Album 2 - Track 1.mp3"}
	cover, err = getMediaCoverFromFolder(track.Path)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), cover.Hash)
	assert.NotEmpty(suite.T(), cover.Content)
	assert.Equal(suite.T(), ".png", cover.Ext)

	// Test with no cover available
	track = domain.Track{Path: TestFSLibDir + "/artist 2/Artist 2 - Album 1 - Track 1.mp3"}
	cover, err = getMediaCoverFromFolder(track.Path)
	assert.NotNil(suite.T(), err)

	return
}

func (suite *LocalFSRepoTestSuite) TestGetMediaCoverFromMetadata() {
	// Test to get a jpg cover.
	track := domain.Track{Path: TestFSLibDir + "/no artist - album 1 - Track 1.mp3"}
	cover, err := getMediaCoverFromMetadata(track.Path)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), cover.Hash)
	assert.NotEmpty(suite.T(), cover.Content)
	assert.Equal(suite.T(), ".jpg", cover.Ext)

	// Test with no cover available
	track = domain.Track{Path: TestFSLibDir + "/no artist - no album - no title.mp3"}
	cover, err = getMediaCoverFromMetadata(track.Path)
	assert.NotNil(suite.T(), err)

	return
}

func (suite *LocalFSRepoTestSuite) TestGetMetadataFromFile() {
	// Test with almost full metadata.
	track := domain.Track{Path: TestFSLibDir + "/artist 1/artist 1 - album 1/Artist 1 - Album 1 - Track 1.mp3"}
	meta, err := getMetadataFromFile(track.Path)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "MP3", meta.Format)
	assert.Equal(suite.T(), "Artist #1 - Album #1 - Track #1", meta.Title)
	assert.Equal(suite.T(), "Artist #1 - Album #1", meta.Album)
	assert.Equal(suite.T(), "Artist #1", meta.Artist)
	assert.Equal(suite.T(), "Genre #1", meta.Genre)
	assert.Equal(suite.T(), "2017", meta.Year)
	assert.Equal(suite.T(), 1, meta.Track)
	assert.Empty(suite.T(), meta.Disc)
	assert.Empty(suite.T(), meta.Picture)
	// TODO Cannot test duration with the test file.
	assert.Equal(suite.T(), 0, meta.Duration)
	// Path will be different on each platform so we can only test it's not empty.
	assert.NotEmpty(suite.T(), meta.Path)

	// Test without artist nor album.
	track = domain.Track{Path: TestFSLibDir + "/no artist - no album - Track 1.mp3"}
	meta, err = getMetadataFromFile(track.Path)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "MP3", meta.Format)
	assert.Equal(suite.T(), "No artist - no album - Track #1", meta.Title)
	assert.Empty(suite.T(), meta.Album)
	assert.Empty(suite.T(), meta.Artist)
	assert.Equal(suite.T(), "Genre #5", meta.Genre)
	assert.Equal(suite.T(), "2017", meta.Year)
	assert.Equal(suite.T(), 1, meta.Track)
	assert.Empty(suite.T(), meta.Disc)
	assert.Empty(suite.T(), meta.Picture)
	// TODO Cannot test duration with the test file.
	assert.Equal(suite.T(), 0, meta.Duration)
	// Path will be different on each platform so we can only test it's not empty.
	assert.NotEmpty(suite.T(), meta.Path)

	// Test without any tag.
	track = domain.Track{Path: TestFSLibDir + "/no artist - no album - no title.mp3"}
	meta, err = getMetadataFromFile(track.Path)
	assert.Nil(suite.T(), err)
	// Should set the file name as a title.
	assert.Equal(suite.T(), "no artist - no album - no title", meta.Title)
	assert.Empty(suite.T(), meta.Album)
	assert.Empty(suite.T(), meta.Artist)
	assert.Empty(suite.T(), meta.Genre)
	assert.Empty(suite.T(), meta.Year)
	assert.Empty(suite.T(), meta.Track)
	assert.Empty(suite.T(), meta.Disc)
	assert.Empty(suite.T(), meta.Picture)
	// TODO Cannot test duration with the test file.
	assert.Equal(suite.T(), 0, meta.Duration)
	// Path will be different on each platform so we can only test it's not empty.
	assert.NotEmpty(suite.T(), meta.Path)

	// Test with multiple discs.
	track = domain.Track{Path: TestFSLibDir + "/artist 2/Artist 2 - Album 1 - Track 1.mp3"}
	meta, err = getMetadataFromFile(track.Path)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "1/2", meta.Disc)

	// Test with cover image in tags.
	track = domain.Track{Path: TestFSLibDir + "/no artist - album 1 - Track 1.mp3"}
	meta, err = getMetadataFromFile(track.Path)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "jpg", meta.Picture.Ext)
	assert.NotEmpty(suite.T(), meta.Picture.Data)

	// Test with non existant file.
	meta, err = getMetadataFromFile("non/existant/file.mp3")
	assert.NotNil(suite.T(), err)

	return
}

func (suite *LocalFSRepoTestSuite) TestWriteCoverFileInternal() {
	tmpCoversDirectory := viper.GetString("Covers.Directory")
	cover, _ := getMediaCoverFromFolder(TestFSLibDir + "/artist 1/artist 1 - album 1/Artist 1 - Album 1 - Track 1.mp3")

	err := writeCoverFile(&cover, tmpCoversDirectory)
	assert.Nil(suite.T(), err)

	// Try to read the created file on the filesystem.
	exists := fileExists(tmpCoversDirectory + string(os.PathSeparator) + cover.Path)
	assert.True(suite.T(), exists)
}
