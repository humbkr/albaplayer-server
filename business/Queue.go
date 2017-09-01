package business

import "git.humbkr.com/jgalletta/alba-player/domain"

/**
@file
Contains code relative to the queue system (currently playing list of songs).
*/

type Queue struct {
	Tracklist domain.Tracks
	Library   *LibraryInteractor
}

func (q Queue) AppendTrack(trackId int) {
	track, err := q.Library.TrackRepository.Find(trackId)
	if err != nil {
		q.Tracklist = append(q.Tracklist, track)
	}
}

func (q Queue) AppendAlbum(albumId int) {
	album, err := q.Library.AlbumRepository.Find(albumId)
	if err != nil {
		for _, track := range album.Tracks {
			q.Tracklist = append(q.Tracklist, track)
		}
	}
}

func (q Queue) AppendArtist(artistId int) {
	// TODO code this.
}
