package interfaces

import (
	"errors"
	"git.humbkr.com/jgalletta/alba-player/internal/alba/domain"
)

type AlbumDbRepository struct {
	AppContext *AppContext
}

/*
Fetches an album from the database.
*/
func (ar AlbumDbRepository) Get(id int) (entity domain.Album, err error) {
	object, err := ar.AppContext.DB.Get(domain.Album{}, id)
	if err == nil && object != nil {
		entity = *object.(*domain.Album)
		ar.populateTracks(&entity)
	} else {
		err = errors.New("no album found")
	}

	return
}

/*
Fetches all albums from the database.

@param hydrate
	If true populate albums tracks. WARNING: VERY time consuming
*/
func (ar AlbumDbRepository) GetAll(hydrate bool) (entities domain.Albums, err error) {
	if !hydrate {
		query := "SELECT * FROM albums"
		_, err = ar.AppContext.DB.Select(&entities, query)

	} else {
		type gorpResult struct {
			albumId int
			albumTitle string
			albumYear string
			albumArtistId int
			domain.Track
		}
		var results []gorpResult

		query := "SELECT alb.Id albumId, alb.Title albumTitle, alb.Year albumYeat, alb.ArtistId albumArtistId, trk.* " +
			     "FROM albums alb JOIN tracks trk ON alb.id = trk.album_id"

		_, err = ar.AppContext.DB.Select(&results, query)
		if err == nil {
			// Deduplicate stuff.
			var current domain.Album
			for _, r := range results {
				track := domain.Track{
					Id: r.Id,
					Title: r.Title,
					AlbumId: r.albumId,
					ArtistId: r.ArtistId,
					CoverId: r.CoverId,
					Disc: r.Disc,
					Number: r.Number,
					Duration: r.Duration,
					Genre: r.Genre,
					Path: r.Path,
				}

				if current.Id == 0 {
					current = domain.Album{
						Id: r.albumId,
						Title: r.albumTitle,
						Year: r.albumYear,
						ArtistId: r.albumArtistId,
					}
				} else if r.Id != current.Id {
					entities = append(entities, current)
					// Then change the current album
					current = domain.Album{
						Id: r.albumId,
						Title: r.albumTitle,
						Year: r.albumYear,
						ArtistId: r.albumArtistId,
					}
				}
				current.Tracks = append(current.Tracks, track)
			}
		}
	}

	return
}

/**
Fetches an album from database.
*/
func (ar AlbumDbRepository) GetByName(name string, artistId int) (entity domain.Album, err error) {
	var entities domain.Albums
	_, err = ar.AppContext.DB.Select(&entities, "SELECT * FROM albums WHERE title = ? AND artist_id = ?", name, artistId)

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
Fetches albums having the specified artistId from database ordered by year.
*/
func (ar AlbumDbRepository) GetAlbumsForArtist(artistId int, hydrate bool) (entities domain.Albums, err error) {
	_, err = ar.AppContext.DB.Select(&entities, "SELECT * FROM albums WHERE artist_id = ? ORDER BY year", artistId)
	if err == nil && hydrate {
		for i := range entities {
			ar.populateTracks(&entities[i])
		}
	}

	return
}

/**
Create or update an album in the Database.
*/
func (ar AlbumDbRepository) Save(entity *domain.Album) (err error) {
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
Delete an album from the Database.
*/
func (ar AlbumDbRepository) Delete(entity *domain.Album) (err error) {
	// Delete all tracks.
	if len(entity.Tracks) == 0 {
		ar.populateTracks(entity)
	}
	tracksRepo := TrackDbRepository{AppContext: ar.AppContext}
	for i := range entity.Tracks {
		tracksRepo.Delete(&entity.Tracks[i])
	}

	// Then delete album.
	_, err = ar.AppContext.DB.Delete(entity)
	return
}

// Check if an album exists for a given id.
func (ar AlbumDbRepository) Exists(id int) bool {
	_, err := ar.Get(id)
	return err == nil
}

/**
Helper function to populate tracks.
 */
func (ar AlbumDbRepository) populateTracks(album *domain.Album) {
	tracksRepo := TrackDbRepository{AppContext: ar.AppContext}
	if tracks, err := tracksRepo.GetTracksForAlbum(album.Id); err == nil {
		album.Tracks = tracks
	}
}
