package interfaces

import (
	"strings"

	"errors"

	"git.humbkr.com/jgalletta/alba-player/domain"
)

type ArtistRepository struct {
	AppContext *AppContext
}

/*
Fetches an artist from the database.
*/
func (ar ArtistRepository) Find(id int) (entity domain.Artist, err error) {
	err = ar.AppContext.DB.SelectOne(&entity, "SELECT * FROM artists WHERE id=?", id)

	return
}

/**
Fetches an artist from database based on its name (case insensitive).
*/
func (ar ArtistRepository) FindByName(name string) (entity domain.Artist, err error) {
	var entities domain.Artists
	_, err = ar.AppContext.DB.Select(&entities, "SELECT * FROM artists WHERE LOWER(name) = ?", strings.ToLower(name))

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
Create or update an artist in the Database.
*/
func (ar ArtistRepository) Save(entity *domain.Artist) (err error) {
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
