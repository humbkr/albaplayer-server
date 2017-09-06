package repositories_test

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"git.humbkr.com/jgalletta/alba-player/interfaces"
	"github.com/stretchr/testify/assert"
	"log"
	"git.humbkr.com/jgalletta/alba-player/domain"
)

type ArtistRepoTestSuite struct {
	suite.Suite
	ArtistRepository interfaces.ArtistDbRepository
}

func (suite *ArtistRepoTestSuite) SetupSuite() {
	ds, err := CreateTestDatasource()
	if err != nil {
		log.Fatal(err)
	}
	appContext := interfaces.AppContext{DB: ds}
	suite.ArtistRepository = interfaces.ArtistDbRepository{AppContext: &appContext}
}

func (suite *ArtistRepoTestSuite) TearDownSuite() {
	if err := CloseTestDataSource(suite.ArtistRepository.AppContext.DB); err != nil {
		log.Fatal(err)
	}
}

func (suite *ArtistRepoTestSuite) SetupTest() {
	ResetTestDataSource(suite.ArtistRepository.AppContext.DB)
}

func (suite *ArtistRepoTestSuite) TestFind() {
	// Test artist retrieval.
	artist, err := suite.ArtistRepository.Find(1)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, artist.Id)
	assert.Equal(suite.T(), "Tool", artist.Name)
	assert.NotEmpty(suite.T(), artist.Albums)
	assert.Len(suite.T(), artist.Albums, 1)

	// Test to get a non existing artist.
	artist, err = suite.ArtistRepository.Find(99)
	assert.NotNil(suite.T(), err)
}

func (suite *ArtistRepoTestSuite) TestFindAll() {
	// Test to get artist without albums.
	artists, err := suite.ArtistRepository.FindAll(false)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), artists)
	for _, artist := range artists {
		assert.NotEmpty(suite.T(), artist.Id)
		assert.NotEmpty(suite.T(), artist.Name)
		assert.Empty(suite.T(), artist.Albums)
	}

	// Test to get artist with albums.
	artists, err = suite.ArtistRepository.FindAll(true)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), artists)
	for _, artist := range artists {
		assert.NotEmpty(suite.T(), artist.Id)
		assert.NotEmpty(suite.T(), artist.Name)
		assert.NotEmpty(suite.T(), artist.Albums)
	}
}

func (suite *ArtistRepoTestSuite) TestFindByName() {
	// Test artist retrieval.
	artist, err := suite.ArtistRepository.FindByName("Tool")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, artist.Id)
	assert.Equal(suite.T(), "Tool", artist.Name)
	assert.Empty(suite.T(), artist.Albums)

	// Test to get an artist with non existant name.
	_, err = suite.ArtistRepository.FindByName("Bogus")
	assert.NotNil(suite.T(), err)
}

func (suite *ArtistRepoTestSuite) TestSave() {
	// Note: we do not save embedded objects for the time being.
	// Test to save a new artist.
	newArtist := &domain.Artist{
		Name: "Insert new artist test",
	}

	err := suite.ArtistRepository.Save(newArtist)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), newArtist.Id)

	insertedNewArtist, errInsert := suite.ArtistRepository.Find(newArtist.Id)
	assert.Nil(suite.T(), errInsert)
	assert.Equal(suite.T(), newArtist.Id, insertedNewArtist.Id)
	assert.Equal(suite.T(), "Insert new artist test", insertedNewArtist.Name)
	assert.Empty(suite.T(), insertedNewArtist.Albums)

	// Test to update the artist.
	insertedNewArtist.Name = "Update artist test"
	errUpdate := suite.ArtistRepository.Save(&insertedNewArtist)
	assert.Nil(suite.T(), errUpdate)
	assert.NotEmpty(suite.T(), insertedNewArtist.Id)

	updatedArtist, errGetMod := suite.ArtistRepository.Find(newArtist.Id)
	assert.Nil(suite.T(), errGetMod)
	assert.Equal(suite.T(), newArtist.Id, updatedArtist.Id)
	assert.Equal(suite.T(), "Update artist test", updatedArtist.Name)
	assert.Empty(suite.T(), updatedArtist.Albums)

	// Test to insert an artist without name.
	newArtistNoName := &domain.Artist{}

	errNoTitle := suite.ArtistRepository.Save(newArtistNoName)
	assert.NotNil(suite.T(), errNoTitle)

	// Test to insert a new artist with a prepopulated Id (= update a non existant artist).
	// Note: it seems gorp.Dbmap.Update() fails silently.
	newArtistWithId := &domain.Artist{
		Id: 55,
		Name: "New artist bogus id",
	}

	errBogusId := suite.ArtistRepository.Save(newArtistWithId)
	assert.Nil(suite.T(), errBogusId)

	// Test to update an album with an empty title.
	updatedArtist.Name = ""
	errUpdateEmptyTitle := suite.ArtistRepository.Save(&updatedArtist)
	assert.NotNil(suite.T(), errUpdateEmptyTitle)
 }

func (suite *ArtistRepoTestSuite) TestDelete() {
	var artistId = 1

	//Get artist to delete.
	artist, err := suite.ArtistRepository.Find(artistId)
	assert.Nil(suite.T(), err)

	// Delete artist.
	err = suite.ArtistRepository.Delete(&artist)
	assert.Nil(suite.T(), err)

	// Check artist has been removed from the database.
	_, err = suite.ArtistRepository.Find(artistId)
	assert.NotNil(suite.T(), err)

	// Check artist's albums have been removed too.
	albumRepo := interfaces.AlbumDbRepository{AppContext: suite.ArtistRepository.AppContext}
	albums, err := albumRepo.FindAlbumsForArtist(artistId, false)
	assert.Nil(suite.T(), err)
	assert.Empty(suite.T(), albums)

	// Check that tracks have been removed too is already done in albumRepository tests.
}

/**
Go testing framework entry point.
 */
func TestArtistRepoTestSuite(t *testing.T) {
	suite.Run(t, new(AlbumRepoTestSuite))
}