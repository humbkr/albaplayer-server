package business

import (
	"git.humbkr.com/jgalletta/alba-player/domain"
	"errors"
	"sync"
)

/*
@file
Contains code relative to the queue system (currently playing list of songs).
*/

const QueueOptionRepeatSingle = 1
const QueueOptionRepeatAll = 2

type Queue struct {
	// Internal ID of the queue.
	Id int
	// List of tracks in the queue. This is not ordered.
	Tracklist map[int]domain.Track
	// Order in which to play the tracks. A same track id can
	// be in here multiple times.
	PlayingOrder []int
	// Flag to see if at least one track has been set as playing.
	startedPlaying bool
	// Index of the currently playing track from PlayinOrder.
	NowPlaying int
	// Queue playing options.
	Options struct{
		Repeat int
		Random bool
	}

	Library   *LibraryInteractor
}

// TODO remove this when coding for multi users.
var instance *Queue
var once sync.Once

func GetQueueInstance() *Queue {
	once.Do(func() {
		instance = &Queue{}
		instance.Tracklist = make(map[int]domain.Track)
	})
	return instance
}
// end todo.


func (q Queue) Current() (domain.Track, error) {
	if len(q.Tracklist) == 0 {
		return domain.Track{}, errors.New("no tracks in the queue")
	}

	return q.Tracklist[q.PlayingOrder[q.NowPlaying]], nil
}

func (q *Queue) Previous() (domain.Track, error) {
	if len(q.Tracklist) == 0 {
		return domain.Track{}, errors.New("no tracks in the queue")
	}

	// TODO not sure if startedPlaying is the best method.
	if !q.startedPlaying {
		q.startedPlaying = true
		return q.Tracklist[q.PlayingOrder[0]], nil
	} else {
		// If repeat is set to one song return the track currently playing.
		if q.Options.Repeat == QueueOptionRepeatSingle {
			return q.Tracklist[q.PlayingOrder[q.NowPlaying]], nil
		}

		// Try to get the previous track.
		if q.NowPlaying - 1 >= 0 {
			// There is a previous track.
			q.NowPlaying--
			return q.Tracklist[q.PlayingOrder[q.NowPlaying]], nil
		} else if q.Options.Repeat == QueueOptionRepeatAll {
			// There is no previous track but queue is set to loop over tracks, go to the last track.
			q.NowPlaying = len(q.PlayingOrder) - 1
			return q.Tracklist[q.PlayingOrder[q.NowPlaying]], nil
		} else {
			// Beginning of the queue.
			return domain.Track{}, errors.New("no previous track in the queue")
		}
	}
}

func (q *Queue) Next() (domain.Track, error) {
	if len(q.Tracklist) == 0 {
		return domain.Track{}, errors.New("no tracks in the queue")
	}

	// TODO not sure if startedPlaying is the best method.
	if !q.startedPlaying {
		q.startedPlaying = true
		return q.Tracklist[q.PlayingOrder[0]], nil
	} else {
		// If repeat is set to one song return the track currently playing.
		if q.Options.Repeat == QueueOptionRepeatSingle {
			return q.Tracklist[q.PlayingOrder[q.NowPlaying]], nil
		}

		// Try to get the next track.
		if q.NowPlaying + 1 < len(q.PlayingOrder) {
			// There is a next track.
			q.NowPlaying++
			return q.Tracklist[q.PlayingOrder[q.NowPlaying]], nil
		} else if q.Options.Repeat == QueueOptionRepeatAll {
			// There is no next track but queue is set to loop over tracks, reinit the playing cursor.
			q.NowPlaying = 0
			return q.Tracklist[q.PlayingOrder[q.NowPlaying]], nil
		} else {
			// End of the queue.
			return domain.Track{}, errors.New("no next track in the queue")
		}
	}
}

func (q *Queue) AppendTrack(trackId int) {
	track, err := q.Library.TrackRepository.Get(trackId)
	if err == nil {
		q.Tracklist[track.Id] = track
		q.PlayingOrder = append(q.PlayingOrder, track.Id)
	}
}

func (q *Queue) PlayTrack(trackId int) {
	q.Clear()
	q.AppendTrack(trackId)
}

func (q *Queue) AppendAlbum(albumId int) {
	album, err := q.Library.AlbumRepository.Get(albumId)
	if err == nil {
		for _, track := range album.Tracks {
			q.Tracklist[track.Id] = track
			q.PlayingOrder = append(q.PlayingOrder, track.Id)
		}
	}
}

func (q *Queue) PlayAlbum(albumId int) {
	q.Clear()
	q.AppendAlbum(albumId)
}

func (q *Queue) AppendArtist(artistId int) {
	albums, err := q.Library.AlbumRepository.GetAlbumsForArtist(artistId, true)
	if err == nil {
		for _, album := range albums {
			for _, track := range album.Tracks {
				q.Tracklist[track.Id] = track
				q.PlayingOrder = append(q.PlayingOrder, track.Id)
			}
		}
	}
}

func (q *Queue) PlayArtist(artistId int) {
	q.Clear()
	q.AppendAlbum(artistId)
}

func (q *Queue) Clear() {
	// Erase tracklist.
	q.Tracklist = make(map[int]domain.Track)
	// Erase playing order.
	q.PlayingOrder = nil
	// Reset Current.
	q.NowPlaying = 0
	// Reset play flag.
	q.startedPlaying = false
}
