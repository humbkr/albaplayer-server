package repositories_test

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"git.humbkr.com/jgalletta/alba-player/interfaces"
	"github.com/stretchr/testify/assert"
	"log"
	"git.humbkr.com/jgalletta/alba-player/domain"
)

type TrackRepoTestSuite struct {
	suite.Suite
	TrackRepository interfaces.TrackDbRepository
}

func (suite *TrackRepoTestSuite) SetupSuite() {
	ds, err := CreateTestDatasource()
	if err != nil {
		log.Fatal(err)
	}
	appContext := interfaces.AppContext{DB: ds}
	suite.TrackRepository = interfaces.TrackDbRepository{AppContext: &appContext}
}

func (suite *TrackRepoTestSuite) TearDownSuite() {
	if err := CloseTestDataSource(suite.TrackRepository.AppContext.DB); err != nil {
		log.Fatal(err)
	}
}

func (suite *TrackRepoTestSuite) SetupTest() {
	ResetTestDataSource(suite.TrackRepository.AppContext.DB)
}

func (suite *TrackRepoTestSuite) TestFind() {
	// Test track retrieval.
	track, err := suite.TrackRepository.Find(1)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, track.Id)
	assert.Equal(suite.T(), 1, track.AlbumId)
	assert.Equal(suite.T(), 1, track.ArtistId)
	assert.Equal(suite.T(), "Stinkfist", track.Title)
	// We do not input anything in the db if there's only one disc, even if the tag is "1/1".
	assert.Empty(suite.T(), track.Disc)
	assert.Equal(suite.T(), 1, track.Number)
	assert.Equal(suite.T(), 311, track.Duration)
	assert.Equal(suite.T(), "Progressive Metal", track.Genre)
	assert.Equal(suite.T(), "/home/test/music/tool/aenima/01 - Stkinfist.mp3", track.Path)

	// TODO Test double disc albums.

	// Test to get a non existing track.
	track, err = suite.TrackRepository.Find(99)
	assert.NotNil(suite.T(), err)
}

func (suite *TrackRepoTestSuite) TestFindAll() {
	tracks, err := suite.TrackRepository.FindAll()
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), tracks)
	for _, track := range tracks {
		assert.NotEmpty(suite.T(), track.Id)
		assert.NotEmpty(suite.T(), track.Title)
		assert.NotEmpty(suite.T(), track.Path)
	}
}

func (suite *TrackRepoTestSuite) TestFindByName() {
	// Test track retrieval.
	track, err := suite.TrackRepository.FindByName("Forty Six & 2", 1, 1)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 5, track.Id)
	assert.Equal(suite.T(), 1, track.AlbumId)
	assert.Equal(suite.T(), 1, track.ArtistId)
	assert.Equal(suite.T(), "Forty Six & 2", track.Title)
	// We do not input anything in the db if there's only one disc, even if the tag is "1/1".
	assert.Empty(suite.T(), track.Disc)
	assert.Equal(suite.T(), 5, track.Number)
	assert.Equal(suite.T(), 364, track.Duration)
	assert.Equal(suite.T(), "Progressive Metal", track.Genre)
	assert.Equal(suite.T(), "/home/test/music/tool/aenima/05 - Forty Six & 2.mp3", track.Path)

	// Test to get a track with non existant name.
	_, err = suite.TrackRepository.FindByName("Bogus", 1, 1)
	assert.NotNil(suite.T(), err)

	// Test to get a track with wrong artist id.
	_, err = suite.TrackRepository.FindByName("Forty Six & 2", 2, 1)
	assert.NotNil(suite.T(), err)

	// Test to get a track with wrong album id.
	_, err = suite.TrackRepository.FindByName("Forty Six & 2", 1, 2)
	assert.NotNil(suite.T(), err)
}

func (suite *TrackRepoTestSuite) TestFindTracksForAlbum() {
	tracks, err := suite.TrackRepository.FindTracksForAlbum(1)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), tracks)
	assert.Equal(suite.T(), 15, len(tracks))
	for _, track := range tracks {
		assert.NotEmpty(suite.T(), track.Id)
		assert.Equal(suite.T(), 1, track.AlbumId)
		assert.NotEmpty(suite.T(), track.Title)
		assert.NotEmpty(suite.T(), track.Path)
	}
}

func (suite *TrackRepoTestSuite) TestSave() {
	// Test to save a new track.
	newTrack := &domain.Track{
		AlbumId: 2,
		ArtistId: 2,
		Title: "Insert new track test",
		Disc: "1/2",
		Number: 5,
		Duration: 321,
		Genre: "Grunge",
		Path: "/home/test/music/artist test/album test/05 - Insert new track test.mp3",
	}

	err := suite.TrackRepository.Save(newTrack)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), newTrack.Id)

	insertednewTrack, errInsert := suite.TrackRepository.Find(newTrack.Id)
	assert.Nil(suite.T(), errInsert)
	assert.Equal(suite.T(), newTrack.Id, insertednewTrack.Id)
	assert.Equal(suite.T(), 2, insertednewTrack.AlbumId)
	assert.Equal(suite.T(), 2, insertednewTrack.ArtistId)
	assert.Equal(suite.T(), "Insert new track test", insertednewTrack.Title)
	assert.Equal(suite.T(), "1/2", insertednewTrack.Disc)
	assert.Equal(suite.T(), 5, insertednewTrack.Number)
	assert.Equal(suite.T(), 321, insertednewTrack.Duration)
	assert.Equal(suite.T(),  "Grunge", insertednewTrack.Genre)
	assert.Equal(suite.T(), "/home/test/music/artist test/album test/05 - Insert new track test.mp3", insertednewTrack.Path)

	// Test to update the track with valid data.
	insertednewTrack.Title = "Update track test"
	insertednewTrack.AlbumId = 1
	insertednewTrack.ArtistId = 1
	insertednewTrack.Disc = "2/2"
	insertednewTrack.Number = 6
	insertednewTrack.Duration = 123
	insertednewTrack.Genre = "Thrash Metal"
	insertednewTrack.Path = "/home/test/music/artist test/album test/05 - Update track test.mp3"
	errUpdate := suite.TrackRepository.Save(&insertednewTrack)
	assert.Nil(suite.T(), errUpdate)
	assert.NotEmpty(suite.T(), insertednewTrack.Id)

	updatedTrack, errGetMod := suite.TrackRepository.Find(newTrack.Id)
	assert.Nil(suite.T(), errGetMod)
	assert.Equal(suite.T(), newTrack.Id, updatedTrack.Id)
	assert.Equal(suite.T(), 1, updatedTrack.AlbumId)
	assert.Equal(suite.T(), 1, updatedTrack.ArtistId)
	assert.Equal(suite.T(), "Update track test", updatedTrack.Title)
	assert.Equal(suite.T(), "2/2", updatedTrack.Disc)
	assert.Equal(suite.T(), 6, updatedTrack.Number)
	assert.Equal(suite.T(), 123, updatedTrack.Duration)
	assert.Equal(suite.T(),  "Thrash Metal", updatedTrack.Genre)
	assert.Equal(suite.T(), "/home/test/music/artist test/album test/05 - Update track test.mp3", updatedTrack.Path)

	// Test to insert a track without title.
	newTrackNoTitle := &domain.Track{
		Path: "/test insert no title.mp3",
	}

	errNoTitle := suite.TrackRepository.Save(newTrackNoTitle)
	assert.NotNil(suite.T(), errNoTitle)

	// Test to insert a track without path.
	newTrackNoPath := &domain.Track{
		Title: "Test insert track no path",
	}

	errNoPath := suite.TrackRepository.Save(newTrackNoPath)
	assert.NotNil(suite.T(), errNoPath)

	// Test to insert a new track with a prepopulated trackId (= update a non existant track).
	// Note: it seems gorp.Dbmap.Update() fails silently.
	newTrackWithId := &domain.Track{
		Id: 88,
		Title: "New track bogus id",
		Path: "/new bogus id.mp3",
	}

	errBogusId := suite.TrackRepository.Save(newTrackWithId)
	assert.Nil(suite.T(), errBogusId)

	// Test to update a track with an empty title.
	updatedTrack.Title = ""
	errUpdateEmptyTitle := suite.TrackRepository.Save(&updatedTrack)
	assert.NotNil(suite.T(), errUpdateEmptyTitle)

	// TODO test to insert a track with a non existant artist id.
	// TODO test to insert a track with a non existant album id.
}

func (suite *TrackRepoTestSuite) TestDelete() {
	var trackId = 1

	// Get track to delete.
	track, err := suite.TrackRepository.Find(trackId)
	assert.Nil(suite.T(), err)

	// Delete track.
	err = suite.TrackRepository.Delete(&track)
	assert.Nil(suite.T(), err)

	// Check track has been removed from the database.
	_, err = suite.TrackRepository.Find(trackId)
	assert.NotNil(suite.T(), err)
}

/**
Go testing framework entry point.
 */
func TestTrackRepoTestSuite(t *testing.T) {
	suite.Run(t, new(TrackRepoTestSuite))
}