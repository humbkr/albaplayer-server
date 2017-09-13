package interfaces

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"
	"log"
	"git.humbkr.com/jgalletta/alba-player/domain"
)

type AlbumRepoTestSuite struct {
	suite.Suite
	AlbumRepository AlbumDbRepository
}

// Go testing framework entry point.
func TestAlbumRepoTestSuite(t *testing.T) {
	suite.Run(t, new(AlbumRepoTestSuite))
}

func (suite *AlbumRepoTestSuite) SetupSuite() {
	ds, err := createTestDatasource()
	if err != nil {
		log.Fatal(err)
	}
	appContext := AppContext{DB: ds}
	suite.AlbumRepository = AlbumDbRepository{AppContext: &appContext}
}

func (suite *AlbumRepoTestSuite) TearDownSuite() {
	if err := closeTestDataSource(suite.AlbumRepository.AppContext.DB); err != nil {
		log.Fatal(err)
	}
}

func (suite *AlbumRepoTestSuite) SetupTest() {
	resetTestDataSource(suite.AlbumRepository.AppContext.DB)
}

func (suite *AlbumRepoTestSuite) TestGet() {
	// Test album retrieval.
	album, err := suite.AlbumRepository.Get(1)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, album.Id)
	assert.Equal(suite.T(), 1, album.ArtistId)
	assert.Equal(suite.T(), "Ænima", album.Title)
	assert.Equal(suite.T(), "1996", album.Year)
	assert.NotEmpty(suite.T(), album.Tracks)
	assert.Len(suite.T(), album.Tracks, 15)

	// Test to get a non existing album.
	album, err = suite.AlbumRepository.Get(99)
	assert.NotNil(suite.T(), err)
}

func (suite *AlbumRepoTestSuite) TestGetAll() {
	// Test to get albums without tracks.
	albums, err := suite.AlbumRepository.GetAll(false)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), albums)
	for _, album := range albums {
		assert.NotEmpty(suite.T(), album.Id)
		assert.NotEmpty(suite.T(), album.Title)
		assert.Empty(suite.T(), album.Tracks)
	}

	// Test to get albums with tracks.
	albums, err = suite.AlbumRepository.GetAll(true)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), albums)
	for _, album := range albums {
		assert.NotEmpty(suite.T(), album.Id)
		assert.NotEmpty(suite.T(), album.Title)
		assert.NotEmpty(suite.T(), album.Tracks)
	}
}

func (suite *AlbumRepoTestSuite) TestGetByName() {
	// Test album retrieval.
	album, err := suite.AlbumRepository.GetByName("Ænima", 1)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, album.Id)
	assert.Equal(suite.T(), 1, album.ArtistId)
	assert.Equal(suite.T(), "Ænima", album.Title)
	assert.Equal(suite.T(), "1996", album.Year)
	assert.Empty(suite.T(), album.Tracks)

	// Test to get an album with non existant name.
	_, err = suite.AlbumRepository.GetByName("Bogus", 1)
	assert.NotNil(suite.T(), err)

	// Test to get an album with wrong artist id.
	_, err = suite.AlbumRepository.GetByName("Ænima", 2)
	assert.NotNil(suite.T(), err)
}

func (suite *AlbumRepoTestSuite) TestGetAlbumsForArtist() {
	// Test to get albums without tracks.
	albums, err := suite.AlbumRepository.GetAlbumsForArtist(1, false)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), albums)
	for _, album := range albums {
		assert.NotEmpty(suite.T(), album.Id)
		assert.Equal(suite.T(), 1, album.ArtistId)
		assert.NotEmpty(suite.T(), album.Title)
		assert.Empty(suite.T(), album.Tracks)
	}

	// Test to get albums with tracks.
	albums, err = suite.AlbumRepository.GetAlbumsForArtist(2, true)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), albums)
	for _, album := range albums {
		assert.NotEmpty(suite.T(), album.Id)
		assert.Equal(suite.T(), 2, album.ArtistId)
		assert.NotEmpty(suite.T(), album.Title)
		assert.NotEmpty(suite.T(), album.Tracks)

		for _, track := range album.Tracks {
			assert.NotEmpty(suite.T(), track.Id)
			assert.Equal(suite.T(), album.Id, track.AlbumId)
			assert.Equal(suite.T(), album.ArtistId, track.ArtistId)
			assert.NotEmpty(suite.T(), track.Title)
		}
	}
}

func (suite *AlbumRepoTestSuite) TestSave() {
	// Note: we do not save embedded objects for the time being.
	// Test to save a new album.
	newAlbum := &domain.Album{
		ArtistId: 2,
		Title: "Insert new album test",
		Year: "2017",
	}

	err := suite.AlbumRepository.Save(newAlbum)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), newAlbum.Id)

	insertedNewAlbum, errInsert := suite.AlbumRepository.Get(newAlbum.Id)
	assert.Nil(suite.T(), errInsert)
	assert.Equal(suite.T(), newAlbum.Id, insertedNewAlbum.Id)
	assert.Equal(suite.T(), 2, insertedNewAlbum.ArtistId)
	assert.Equal(suite.T(), "Insert new album test", insertedNewAlbum.Title)
	assert.Equal(suite.T(), "2017", insertedNewAlbum.Year)
	assert.Empty(suite.T(), insertedNewAlbum.Tracks)

	// Test to update the album.
	insertedNewAlbum.Title = "Update album test"
	insertedNewAlbum.Year = "1988"
	insertedNewAlbum.ArtistId = 1
	errUpdate := suite.AlbumRepository.Save(&insertedNewAlbum)
	assert.Nil(suite.T(), errUpdate)
	assert.NotEmpty(suite.T(), insertedNewAlbum.Id)

	updatedAlbum, errGetMod := suite.AlbumRepository.Get(newAlbum.Id)
	assert.Nil(suite.T(), errGetMod)
	assert.Equal(suite.T(), newAlbum.Id, updatedAlbum.Id)
	assert.Equal(suite.T(), 1, updatedAlbum.ArtistId)
	assert.Equal(suite.T(), "Update album test", updatedAlbum.Title)
	assert.Equal(suite.T(), "1988", updatedAlbum.Year)
	assert.Empty(suite.T(), updatedAlbum.Tracks)

	// Test to insert a new album with a prepopulated albumId (= update a non existant album).
	// Note: it seems gorp.Dbmap.Update() fails silently.
	newAlbumWithId := &domain.Album{
		Id: 44,
		Title: "New album bogus id",
		Year: "2017",
	}

	errBogusId := suite.AlbumRepository.Save(newAlbumWithId)
	assert.Nil(suite.T(), errBogusId)
}

func (suite *AlbumRepoTestSuite) TestDelete() {
	var albumId = 1

	// Get album to delete.
	album, err := suite.AlbumRepository.Get(albumId)
	assert.Nil(suite.T(), err)

	// Delete album.
	err = suite.AlbumRepository.Delete(&album)
	assert.Nil(suite.T(), err)

	// Check album has been removed from the database.
	_, err = suite.AlbumRepository.Get(albumId)
	assert.NotNil(suite.T(), err)

	// Check album tracks have been removed too.
	trackRepo := TrackDbRepository{AppContext: suite.AlbumRepository.AppContext}
	tracks, err := trackRepo.GetTracksForAlbum(albumId)
	assert.Nil(suite.T(), err)
	assert.Empty(suite.T(), tracks)
}
