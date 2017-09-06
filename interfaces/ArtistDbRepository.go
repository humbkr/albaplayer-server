package interfaces

import (
	"errors"

	"git.humbkr.com/jgalletta/alba-player/domain"
)

type ArtistDbRepository struct {
	AppContext *AppContext
}

/*
Fetches an artist from the database.
*/
func (ar ArtistDbRepository) Get(id int) (entity domain.Artist, err error) {
	object, err := ar.AppContext.DB.Get(domain.Artist{}, id)
	if err == nil && object != nil {
		entity = object.(domain.Artist)
		ar.populateAlbums(&entity, true)
	}

	return
}

/*
Fetches all artists from the database.

@param hydrate
	If true populate albums tracks.
*/
func (ar ArtistDbRepository) GetAll(hydrate bool) (entities domain.Artists, err error) {
	_, err = ar.AppContext.DB.Select(&entities, "SELECT * FROM artists")
	if hydrate {
		for i := range entities {
			ar.populateAlbums(&entities[i], hydrate)
		}
	}

	return
}

/**
Fetches an artist from database based on its name (case insensitive).
*/
func (ar ArtistDbRepository) GetByName(name string) (entity domain.Artist, err error) {
	var entities domain.Artists
	_, err = ar.AppContext.DB.Select(&entities, "SELECT * FROM artists WHERE name = ?", name)

	if err == nil {
		if len(entities) > 0 {
			entity = entities[0]
		} else {
			err = errors.New("no result found")
		}
	}

	return
}

/**
Create or update an artist in the Database.
*/
func (ar ArtistDbRepository) Save(entity *domain.Artist) (err error) {
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

/**
Delete an artist from the Database.
*/
func (ar ArtistDbRepository) Delete(entity *domain.Artist) (err error) {
	// Delete all albums.
	if len(entity.Albums) == 0 {
		ar.populateAlbums(entity, false)
	}
	albumRepo := AlbumDbRepository{AppContext: ar.AppContext}
	for i := range entity.Albums {
		albumRepo.Delete(&entity.Albums[i])
	}

	// Then delete album.
	_, err = ar.AppContext.DB.Delete(entity)

	return
}

// Check if an artist exists for a given id.
func (ar ArtistDbRepository) Exists(id int) bool {
	entity, err := ar.AppContext.DB.Get(domain.Artist{}, id)
	return err == nil && entity != nil
}

/**
Helper function to populate albums.
 */
func (ar ArtistDbRepository) populateAlbums(artist *domain.Artist, hydrate bool) {
	albumRepo := AlbumDbRepository{AppContext: ar.AppContext}
	if albums, err := albumRepo.GetAlbumsForArtist(artist.Id, hydrate); err == nil {
		artist.Albums = albums
	}
}