package interfaces

import (
	"errors"
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
	if err == nil {
		entity = album
		ar.populateTracks(&entity)
	}

	return
}

/**
Fetches an album from database.
*/
func (ar AlbumRepository) FindByName(name string, artistId int) (entity domain.Album, err error) {
	var entities domain.Albums
	_, err = ar.AppContext.DB.Select(&entities, "SELECT * FROM albums WHERE title = ? AND artist_id = ?", name, artistId)

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
func (ar AlbumRepository) FindAlbumsForArtist(artistId int) (entities domain.Albums, err error) {
	_, err = ar.AppContext.DB.Select(&entities, "SELECT * FROM albums WHERE artist_id = ? ORDER BY year", artistId)
	if err == nil {
		for i := range entities {
			ar.populateTracks(&entities[i])
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

/**
Helper function to populate tracks.
 */
func (ar AlbumRepository) populateTracks(album *domain.Album) {
	tracksRepo := TrackRepository{AppContext: ar.AppContext}
	if tracks, err := tracksRepo.FindTracksForAlbum(album.Id); err == nil {
		album.Tracks = tracks
	}
}