package business_test

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"
	"git.humbkr.com/jgalletta/alba-player/domain"
	"git.humbkr.com/jgalletta/alba-player/business"
)

type ArtistInteractorTestSuite struct {
	suite.Suite
	// LibraryInteractor where is located what to test.
	Library *business.LibraryInteractor
}

/**
Go testing framework entry point.
 */
func TestArtistRepoTestSuite(t *testing.T) {
	suite.Run(t, new(ArtistInteractorTestSuite))
}

func (suite *ArtistInteractorTestSuite) SetupSuite() {
	suite.Library = createMockLibraryInteractor()
}

func (suite *ArtistInteractorTestSuite) TestGetArtist() {
	// Test artist retrieval.
	artist, err := suite.Library.GetArtist(1)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), artist.Id)
	// Make sure artist has been hydrated.
	assert.NotEmpty(suite.T(), artist.Albums)

	// Test to get a non existing artist.
	artist, err = suite.Library.GetArtist(99)
	assert.NotNil(suite.T(), err)
}

func (suite *ArtistInteractorTestSuite) TestGetAllArtists() {
	// Test to get artists excluding albums.
	artists, err := suite.Library.GetAllArtists(false)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), artists)
	for _, artist := range artists {
		assert.NotEmpty(suite.T(), artist.Id)
		assert.NotEmpty(suite.T(), artist.Name)
		assert.Empty(suite.T(), artist.Albums)
	}

	// Test to get artists including albums.
	artists, err = suite.Library.GetAllArtists(true)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), artists)
	for _, artist := range artists {
		assert.NotEmpty(suite.T(), artist.Id)
		assert.NotEmpty(suite.T(), artist.Name)
		assert.NotEmpty(suite.T(), artist.Albums)
	}
}

func (suite *ArtistInteractorTestSuite) TestSaveArtist() {
	// Test to save a new artist.
	newArtist := &domain.Artist{
		Name: "Insert new artist test",
	}

	err := suite.Library.SaveArtist(newArtist)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), newArtist.Id)

	// Test to update the artist.
	newArtistId := newArtist.Id
	newArtist.Name = "Update artist test"
	errUpdate := suite.Library.SaveArtist(newArtist)
	assert.Nil(suite.T(), errUpdate)
	assert.Equal(suite.T(), newArtist.Id, newArtistId)
	assert.Equal(suite.T(), "Update artist test", newArtist.Name)
	assert.Empty(suite.T(), newArtist.Albums)

	// Test to insert an artist without name.
	newArtistNoName := &domain.Artist{}

	errNoTitle := suite.Library.SaveArtist(newArtistNoName)
	assert.NotNil(suite.T(), errNoTitle)

	// Test to update an album with an empty title.
	newArtist.Name = ""
	errUpdateEmptyTitle := suite.Library.SaveArtist(newArtist)
	assert.NotNil(suite.T(), errUpdateEmptyTitle)
}

func (suite *ArtistInteractorTestSuite) TestDeleteArtist() {
	// Delete artist.
	artist := &domain.Artist{Id: 1}
	err := suite.Library.DeleteArtist(artist)
	assert.Nil(suite.T(), err)

	// Delete non existant album.
	artistFake := &domain.Artist{Id: 55}
	errFake := suite.Library.DeleteArtist(artistFake)
	assert.Nil(suite.T(), errFake)

	// Try to Delete an album which id is not provided.
	artistNoId := &domain.Artist{}
	errNoId := suite.Library.DeleteArtist(artistNoId)
	assert.NotNil(suite.T(), errNoId)
}
