package interfaces

import (
	"errors"
	"strings"

	"fmt"
	"git.humbkr.com/jgalletta/alba-player/domain"
)

type AlbumRepository struct {
	AppContext *AppContext
}

/*
Fetches an album from the database.
*/
func (ar AlbumRepository) Find(id int) (entity domain.Album, err error) {
	var album domain.Album
	err = ar.AppContext.DB.SelectOne(&album, "SELECT * FROM albums WHERE id=?", id)
	if err != nil {
		entity = album
		// Populate tracks.
		tracksRepo := TrackRepository{AppContext: ar.AppContext}
		tracks, err := tracksRepo.FindTracksForAlbum(id)
		if err != nil {
			entity.Tracks = tracks
		}
	}

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
Fetches albums having the specified artistId from database ordered by year.
*/
func (tr TrackRepository) FindAlbumsForArtist(artistId int) (entities map[int]domain.Album, err error) {
	var albums domain.Albums

	_, err = tr.AppContext.DB.Select(&albums, "SELECT * FROM albums WHERE artist_id = ? ORDER BY year", artistId)
	if err != nil {
		// Create a map of tracks indexed by the trackId.
		entities = make(map[int]domain.Album)
		for _, album := range albums {
			entities[album.Id] = album
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
