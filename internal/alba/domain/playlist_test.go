package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PlaylistTestSuite struct {
	suite.Suite
}

// Go testing framework entry point.
func TestPlaylistTestSuite(t *testing.T) {
	suite.Run(t, new(PlaylistTestSuite))
}

func (suite *PlaylistTestSuite) TestAll() {
	playlist := Playlist{Title: "My playlist"}

	track := Track{Title: "Track 01"}

	playlist.Add(track)

	fmt.Println(playlist)

	playlist.Remove(0)

	fmt.Println(playlist)

	playlist.Add(track)
	playlist.Add(track)
	playlist.Add(track)

	fmt.Println(playlist)

	playlist.Clear()

	fmt.Println(playlist)

	playlist.Add(track)

	fmt.Println(playlist)
}
