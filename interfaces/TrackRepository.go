package interfaces

import (
	"errors"

	"git.humbkr.com/jgalletta/alba-player/domain"
)

type TrackRepository struct {
	AppContext *AppContext
}

/*
Fetches a track from the database.
*/
func (tr TrackRepository) Find(id int) (entity domain.Track, err error) {
	err = tr.AppContext.DB.SelectOne(&entity, "SELECT * FROM tracks WHERE id=?", id)

	return
}

/**
Fetches a track from database by name, artist id, and album id.

If several tracks are found, returns only the first one.
*/
func (tr TrackRepository) FindByName(name string, artistId int, albumId int) (entity domain.Track, err error) {
	var entities domain.Tracks
	_, err = tr.AppContext.DB.Select(&entities, "SELECT * FROM tracks WHERE title = ? AND artist_id = ? AND album_id = ?", name, artistId, albumId)

	if err == nil {
		if len(entities) > 0 {
			entity = entities[0]
		} else {
			err = errors.New("No result found")
		}
	}

	return
}

/**
Fetches tracks having the specified albumId from database ordered by disc number then track number.
*/
func (tr TrackRepository) FindTracksForAlbum(albumId int) (entities domain.Tracks, err error) {
	_, err = tr.AppContext.DB.Select(&entities, "SELECT * FROM tracks WHERE album_id = ? ORDER BY disc, number", albumId)

	return
}

/**
Create or update a track in the Database.
*/
func (tr TrackRepository) Save(entity *domain.Track) (err error) {
	if entity.Id != 0 {
		// Update.
		_, err = tr.AppContext.DB.Update(entity)
		return
	} else {
		// Insert new entity.
		err = tr.AppContext.DB.Insert(entity)
		return
	}

	return nil
}
