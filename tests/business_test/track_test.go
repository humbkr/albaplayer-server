package business_test

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"
	"git.humbkr.com/jgalletta/alba-player/domain"
	"git.humbkr.com/jgalletta/alba-player/business"
)

type TrackInteractorTestSuite struct {
	suite.Suite
	// LibraryInteractor where is located what to test.
	Library *business.LibraryInteractor
}

/**
Go testing framework entry point.
 */
func TestTrackRepoTestSuite(t *testing.T) {
	suite.Run(t, new(TrackInteractorTestSuite))
}

func (suite *TrackInteractorTestSuite) SetupSuite() {
	suite.Library = createMockLibraryInteractor()
}

func (suite *TrackInteractorTestSuite) TestGetTrack() {
	// Test track retrieval.
	track, err := suite.Library.GetTrack(1)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, track.Id)
	assert.NotEmpty(suite.T(), track.Title)
	assert.NotEmpty(suite.T(), track.Path)

	// Test to get a non existing track.
	track, err = suite.Library.GetTrack(99)
	assert.NotNil(suite.T(), err)
}

func (suite *TrackInteractorTestSuite) TestGetAllTracks() {
	tracks, err := suite.Library.GetAllTracks()
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), tracks)
	for _, track := range tracks {
		assert.NotEmpty(suite.T(), track.Id)
		assert.NotEmpty(suite.T(), track.Title)
		assert.NotEmpty(suite.T(), track.Path)
	}
}

func (suite *TrackInteractorTestSuite) TestGetTracksForAlbum() {
	tracks, err := suite.Library.GetTracksForAlbum(1)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), tracks)
	for _, track := range tracks {
		assert.NotEmpty(suite.T(), track.Id)
		assert.Equal(suite.T(), 1, track.AlbumId)
		assert.NotEmpty(suite.T(), track.Title)
		assert.NotEmpty(suite.T(), track.Path)
	}
}

func (suite *TrackInteractorTestSuite) TestSaveTrack() {
	// Test to save a new track.
	newTrack := &domain.Track{
		Title: "Insert new track test",
		Path: "/home/test/music/artist test/album test/05 - Insert new track test.mp3",
	}

	err := suite.Library.SaveTrack(newTrack)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), newTrack.Id)


	// Test to update the track with valid data.
	newTrack.Title = "Update track test"
	newTrack.Path = "/home/test/music/artist test/album test/05 - Update track test.mp3"
	errUpdate := suite.Library.SaveTrack(newTrack)
	assert.Nil(suite.T(), errUpdate)
	assert.NotEmpty(suite.T(), newTrack.Id)

	// Test to insert a track without title.
	newTrackNoTitle := &domain.Track{
		Path: "/test insert no title.mp3",
	}

	errNoTitle := suite.Library.SaveTrack(newTrackNoTitle)
	assert.NotNil(suite.T(), errNoTitle)

	// Test to insert a track without path.
	newTrackNoPath := &domain.Track{
		Title: "Test insert track no path",
	}

	errNoPath := suite.Library.SaveTrack(newTrackNoPath)
	assert.NotNil(suite.T(), errNoPath)

	// Test to insert a track with a non existant artist id.
	newTrackInvalidArtist := &domain.Track{
		ArtistId:765,
		Title: "Test insert track invalid artist",
		Path: "/home/test/music/artist test/album test/05 - Insert new track invalid artist.mp3",
	}

	errInvalidArtist := suite.Library.SaveTrack(newTrackInvalidArtist)
	assert.NotNil(suite.T(), errInvalidArtist)

	// Test to insert a track with a non existant album id.
	newTrackInvalidAlbum := &domain.Track{
		ArtistId:765,
		Title: "Test insert track invalid artist",
		Path: "/home/test/music/artist test/album test/05 - Insert new track invalid artist.mp3",
	}

	errInvalidAlbum := suite.Library.SaveTrack(newTrackInvalidAlbum)
	assert.NotNil(suite.T(), errInvalidAlbum)

	// Test to update a track with an empty title.
	newTrack.Title = ""
	errUpdateEmptyTitle := suite.Library.SaveTrack(newTrack)
	assert.NotNil(suite.T(), errUpdateEmptyTitle)
}

func (suite *TrackInteractorTestSuite) TestDeleteTrack() {
	// Delete track.
	track := &domain.Track{Id: 1}
	err := suite.Library.DeleteTrack(track)
	assert.Nil(suite.T(), err)

	// Delete non existant album.
	trackFake := &domain.Track{Id: 55}
	errFake := suite.Library.DeleteTrack(trackFake)
	assert.Nil(suite.T(), errFake)

	// Try to Delete an album which id is not provided.
	trackNoId := &domain.Track{}
	errNoId := suite.Library.DeleteTrack(trackNoId)
	assert.NotNil(suite.T(), errNoId)
}
