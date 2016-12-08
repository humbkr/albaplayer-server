package interfaces

import (
	"errors"
	"strings"

	"git.humbkr.com/jgalletta/alba-player/domain"
)

type AlbumRepository struct {
	AppContext *AppContext
}

/*
Fetches an album from the database.
*/
func (ar AlbumRepository) Find(id int) (entity domain.Album, err error) {
	err = ar.AppContext.DB.SelectOne(&entity, "SELECT * FROM albums WHERE id=?", id)

	return
}

/**
Fetches an album from database.
*/
func (ar AlbumRepository) FindByName(name string, artistId int) (entity domain.Album, err error) {
	var entities domain.Albums
	_, err = ar.AppContext.DB.Select(&entities, "SELECT * FROM albums WHERE LOWER(name) = ? AND artist_id = ?", strings.ToLower(name), artistId)

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
func (ar AlbumRepository) Save(entity *domain.Album) (err error) {
	if entity.Id != 0 {
		// Update.
		_, err = ar.AppContext.DB.Update(entity)
		return
	} else {
		// Insert new entity.
		err = ar.AppContext.DB.Insert(entity)
		return
	}

	return nil
}
