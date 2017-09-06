package business_test

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"
	"git.humbkr.com/jgalletta/alba-player/domain"
	"git.humbkr.com/jgalletta/alba-player/business"
)

type AlbumInteractorTestSuite struct {
	suite.Suite
	// LibraryInteractor where is located what to test.
	Library *business.LibraryInteractor
}

/**
Go testing framework entry point.
 */
func TestAlbumRepoTestSuite(t *testing.T) {
	suite.Run(t, new(AlbumInteractorTestSuite))
}

func (suite *AlbumInteractorTestSuite) SetupSuite() {
	suite.Library = createMockLibraryInteractor()
}

func (suite *AlbumInteractorTestSuite) TestGetAlbum() {
	// Test album retrieval.
	album, err := suite.Library.GetAlbum(1)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), album.Id)
	assert.NotEmpty(suite.T(), album.Tracks)

	// Test to get a non existing album.
	album, err = suite.Library.GetAlbum(99)
	assert.NotNil(suite.T(), err)
}

func (suite *AlbumInteractorTestSuite) TestGetAllAlbums() {
	// Test to get albums without tracks.
	albums, err := suite.Library.GetAllAlbums(false)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), albums)
	for _, album := range albums {
		assert.NotEmpty(suite.T(), album.Id)
		assert.NotEmpty(suite.T(), album.Title)
		assert.Empty(suite.T(), album.Tracks)
	}

	// Test to get albums with tracks.
	albums, err = suite.Library.GetAllAlbums(true)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), albums)
	for _, album := range albums {
		assert.NotEmpty(suite.T(), album.Id)
		assert.NotEmpty(suite.T(), album.Title)
		assert.NotEmpty(suite.T(), album.Tracks)
	}
}

func (suite *AlbumInteractorTestSuite) TestGetAlbumsForArtist() {
	// Test to get albums without tracks.
	albums, err := suite.Library.GetAlbumsForArtist(1, false)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), albums)
	for _, album := range albums {
		assert.NotEmpty(suite.T(), album.Id)
		assert.Equal(suite.T(), 1, album.ArtistId)
		assert.NotEmpty(suite.T(), album.Title)
		assert.Empty(suite.T(), album.Tracks)
	}

	// Test to get albums with tracks.
	albums, err = suite.Library.GetAlbumsForArtist(1, true)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), albums)
	for _, album := range albums {
		assert.NotEmpty(suite.T(), album.Id)
		assert.Equal(suite.T(), 1, album.ArtistId)
		assert.NotEmpty(suite.T(), album.Title)
		assert.NotEmpty(suite.T(), album.Tracks)

		for _, track := range album.Tracks {
			assert.NotEmpty(suite.T(), track.Id)
			assert.Equal(suite.T(), album.Id, track.AlbumId)
			assert.Equal(suite.T(), album.ArtistId, track.ArtistId)
			assert.NotEmpty(suite.T(), track.Title)
		}
	}

	// Test to get albums for a non existant artist (or an artist without album).
	albums, err = suite.Library.GetAlbumsForArtist(34, false)
	assert.NotNil(suite.T(), err)
}

func (suite *AlbumInteractorTestSuite) TestSaveAlbum() {
	// Test to save a new album.
	newAlbum := &domain.Album{
		ArtistId: 1,
		Title: "Insert new album test",
		Year: "2017",
	}

	err := suite.Library.SaveAlbum(newAlbum)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), newAlbum.Id)

	// Test to update the album.
	newAlbum.Title = "Update album test"
	newAlbum.Year = "1988"
	newAlbum.ArtistId = 1
	errUpdate := suite.Library.SaveAlbum(newAlbum)
	assert.Nil(suite.T(), errUpdate)
	assert.NotEmpty(suite.T(), newAlbum.Id)

	// Test to insert an album without title.
	newAlbumNoTitle := &domain.Album{
		ArtistId: 2,
		Year: "2017",
	}

	errNoTitle := suite.Library.SaveAlbum(newAlbumNoTitle)
	assert.NotNil(suite.T(), errNoTitle)

	// Test to insert an album with a non existant artist id.
	newAlbumFakeArtistId := &domain.Album{
		ArtistId: 77,
		Title: "Test invalid artist id",
		Year: "2017",
	}

	errInvalidArtist := suite.Library.SaveAlbum(newAlbumFakeArtistId)
	assert.NotNil(suite.T(), errInvalidArtist)

	// Test to update an album with an empty title.
	newAlbum.Title = ""
	errUpdateEmptyTitle := suite.Library.SaveAlbum(newAlbum)
	assert.NotNil(suite.T(), errUpdateEmptyTitle)
}

func (suite *AlbumInteractorTestSuite) TestDeleteAlbum() {
	// Delete album.
	album := &domain.Album{Id: 1}
	err := suite.Library.DeleteAlbum(album)
	assert.Nil(suite.T(), err)

	// Delete non existant album.
	albumFake := &domain.Album{Id: 55}
	errFake := suite.Library.DeleteAlbum(albumFake)
	assert.Nil(suite.T(), errFake)

	// Try to Delete an album which id is not provided.
	albumNoId := &domain.Album{}
	errNoId := suite.Library.DeleteAlbum(albumNoId)
	assert.NotNil(suite.T(), errNoId)
}
