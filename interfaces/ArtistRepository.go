package interfaces

import (
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

	var artist domain.Artist
	err = ar.AppContext.DB.SelectOne(&artist, "SELECT * FROM artists WHERE id=?", id)
	if err == nil {
		entity = artist

		// Populate albums.
		albumRepo := AlbumRepository{AppContext: ar.AppContext}
		albums, err := albumRepo.FindAlbumsForArtist(id)
		if err == nil {
			entity.Albums = albums
		}
	}

	return
}

/**
Fetches an artist from database based on its name (case insensitive).
*/
func (ar ArtistRepository) FindByName(name string) (entity domain.Artist, err error) {
	var entities domain.Artists
	_, err = ar.AppContext.DB.Select(&entities, "SELECT * FROM artists WHERE name = ?", name)

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
