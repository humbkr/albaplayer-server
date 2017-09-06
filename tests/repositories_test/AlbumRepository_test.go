package repositories_test

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"git.humbkr.com/jgalletta/alba-player/interfaces"
	"github.com/stretchr/testify/assert"
	"log"
	"git.humbkr.com/jgalletta/alba-player/domain"
)

type AlbumRepoTestSuite struct {
	suite.Suite
	AlbumRepository interfaces.AlbumDbRepository
}

func (suite *AlbumRepoTestSuite) SetupSuite() {
	ds, err := CreateTestDatasource()
	if err != nil {
		log.Fatal(err)
	}
	appContext := interfaces.AppContext{DB: ds}
	suite.AlbumRepository = interfaces.AlbumDbRepository{AppContext: &appContext}
}

func (suite *AlbumRepoTestSuite) TearDownSuite() {
	if err := CloseTestDataSource(suite.AlbumRepository.AppContext.DB); err != nil {
		log.Fatal(err)
	}
}

func (suite *AlbumRepoTestSuite) SetupTest() {
	ResetTestDataSource(suite.AlbumRepository.AppContext.DB)
}

func (suite *AlbumRepoTestSuite) TestFind() {
	// Test album retrieval.
	album, err := suite.AlbumRepository.Find(1)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, album.Id)
	assert.Equal(suite.T(), 1, album.ArtistId)
	assert.Equal(suite.T(), "Ænima", album.Title)
	assert.Equal(suite.T(), "1996", album.Year)
	assert.NotEmpty(suite.T(), album.Tracks)
	assert.Len(suite.T(), album.Tracks, 15)

	// Test to get a non existing album.
	album, err = suite.AlbumRepository.Find(99)
	assert.NotNil(suite.T(), err)
}

func (suite *AlbumRepoTestSuite) TestFindAll() {
	// Test to get albums without tracks.
	albums, err := suite.AlbumRepository.FindAll(false)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), albums)
	for _, album := range albums {
		assert.NotEmpty(suite.T(), album.Id)
		assert.NotEmpty(suite.T(), album.Title)
		assert.Empty(suite.T(), album.Tracks)
	}

	// Test to get albums with tracks.
	albums, err = suite.AlbumRepository.FindAll(true)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), albums)
	for _, album := range albums {
		assert.NotEmpty(suite.T(), album.Id)
		assert.NotEmpty(suite.T(), album.Title)
		assert.NotEmpty(suite.T(), album.Tracks)
	}
}

func (suite *AlbumRepoTestSuite) TestFindByName() {
	// Test album retrieval.
	album, err := suite.AlbumRepository.FindByName("Ænima", 1)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, album.Id)
	assert.Equal(suite.T(), 1, album.ArtistId)
	assert.Equal(suite.T(), "Ænima", album.Title)
	assert.Equal(suite.T(), "1996", album.Year)
	assert.Empty(suite.T(), album.Tracks)

	// Test to get an album with non existant name.
	_, err = suite.AlbumRepository.FindByName("Bogus", 1)
	assert.NotNil(suite.T(), err)

	// Test to get an album with wrong artist id.
	_, err = suite.AlbumRepository.FindByName("Ænima", 2)
	assert.NotNil(suite.T(), err)
}

func (suite *AlbumRepoTestSuite) TestFindAlbumsForArtist() {
	// Test to get albums without tracks.
	albums, err := suite.AlbumRepository.FindAlbumsForArtist(1, false)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), albums)
	for _, album := range albums {
		assert.NotEmpty(suite.T(), album.Id)
		assert.Equal(suite.T(), 1, album.ArtistId)
		assert.NotEmpty(suite.T(), album.Title)
		assert.Empty(suite.T(), album.Tracks)
	}

	// Test to get albums with tracks.
	albums, err = suite.AlbumRepository.FindAlbumsForArtist(2, true)
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

	insertedNewAlbum, errInsert := suite.AlbumRepository.Find(newAlbum.Id)
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

	updatedAlbum, errGetMod := suite.AlbumRepository.Find(newAlbum.Id)
	assert.Nil(suite.T(), errGetMod)
	assert.Equal(suite.T(), newAlbum.Id, updatedAlbum.Id)
	assert.Equal(suite.T(), 1, updatedAlbum.ArtistId)
	assert.Equal(suite.T(), "Update album test", updatedAlbum.Title)
	assert.Equal(suite.T(), "1988", updatedAlbum.Year)
	assert.Empty(suite.T(), updatedAlbum.Tracks)

	// Test to insert an album without title.
	newAlbumNoTitle := &domain.Album{
		ArtistId: 2,
		Year: "2017",
	}

	errNoTitle := suite.AlbumRepository.Save(newAlbumNoTitle)
	assert.NotNil(suite.T(), errNoTitle)

	// Test to insert a new album with a prepopulated albumId (= update a non existant album).
	// Note: it seems gorp.Dbmap.Update() fails silently.
	newAlbumWithId := &domain.Album{
		Id: 44,
		Title: "New album bogus id",
		Year: "2017",
	}

	errBogusId := suite.AlbumRepository.Save(newAlbumWithId)
	assert.Nil(suite.T(), errBogusId)

	// Test to update an album with an empty title.
	updatedAlbum.Title = ""
	errUpdateEmptyTitle := suite.AlbumRepository.Save(&updatedAlbum)
	assert.NotNil(suite.T(), errUpdateEmptyTitle)

	// TODO test to insert an album with a non existant artist id.
}

func (suite *AlbumRepoTestSuite) TestDelete() {
	var albumId = 1

	// Get album to delete.
	album, err := suite.AlbumRepository.Find(albumId)
	assert.Nil(suite.T(), err)

	// Delete album.
	err = suite.AlbumRepository.Delete(&album)
	assert.Nil(suite.T(), err)

	// Check album has been removed from the database.
	_, err = suite.AlbumRepository.Find(albumId)
	assert.NotNil(suite.T(), err)

	// Check album tracks have been removed too.
	trackRepo := interfaces.TrackDbRepository{AppContext: suite.AlbumRepository.AppContext}
	tracks, err := trackRepo.FindTracksForAlbum(albumId)
	assert.Nil(suite.T(), err)
	assert.Empty(suite.T(), tracks)
}

/**
Go testing framework entry point.
 */
func TestAlbumRepoTestSuite(t *testing.T) {
	suite.Run(t, new(AlbumRepoTestSuite))
}