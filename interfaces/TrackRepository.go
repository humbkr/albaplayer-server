package interfaces

import (
	"errors"
	"strings"

	"git.humbkr.com/jgalletta/alba-player/domain"
)

type TrackRepository struct {
	AppContext *AppContext
}

/*
Fetches an album from the database.
*/
func (tr TrackRepository) Find(id int) (entity domain.Track, err error) {
	err = tr.AppContext.DB.SelectOne(&entity, "SELECT * FROM albums WHERE id=?", id)

	return
}

/**
Fetches an album from database.
*/
func (tr TrackRepository) FindByName(name string, artistId int, albumId int) (entity domain.Track, err error) {
	var entities domain.Tracks
	_, err = tr.AppContext.DB.Select(&entities, "SELECT * FROM tracks WHERE LOWER(name) = ? AND artist_id = ? AND album_id = ?", strings.ToLower(name), artistId, albumId)

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
Create or update an album in the Database.
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
