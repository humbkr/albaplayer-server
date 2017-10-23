package interfaces

import (
	"errors"
	"git.humbkr.com/jgalletta/alba-player/domain"
	"git.humbkr.com/jgalletta/alba-player/business"
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
func (ar AlbumDbRepository) GetAll(hydrate bool) (entities []business.AlbumView, err error) {
	if !hydrate {
		query := "SELECT alb.*, COALESCE(art.name, '') ArtistName " +
			     "FROM albums alb LEFT JOIN artists art ON alb.artist_id = art.id"
		_, err = ar.AppContext.DB.Select(&entities, query)

	} else {
		type gorpResult struct {
			business.AlbumView
			domain.Track
			TrackId int
			TrackTitle string
			TrackArtistId int
			TrackCoverId int
		}
		var results []gorpResult

		query := "SELECT alb.*, trk.Id TrackId, trk.Title TrackTitle, trk.artist_id TrackArtistId, trk.cover_id TrackCoverId, trk.disc , trk.number , trk.duration , trk.genre, trk.path " +
			     "FROM albums alb JOIN tracks trk ON alb.id = trk.album_id"

		_, err = ar.AppContext.DB.Select(&results, query)
		if err == nil {
			// Deduplicate stuff.
			var current business.AlbumView
			for _, r := range results {
				track := domain.Track{
					Id: r.TrackId,
					Title: r.TrackTitle,
					AlbumId: r.Id,
					ArtistId: r.TrackArtistId,
					CoverId: r.TrackCoverId,
					Disc: r.Disc,
					Number: r.Number,
					Duration: r.Duration,
					Genre: r.Genre,
					Path: r.Path,
				}

				if current.Id == 0 {
					// Create a new AlbumView.
					current = business.AlbumView{
						Album: domain.Album{Id: r.Id,
							Title: r.Title,
							Year: r.Year,
							ArtistId: r.ArtistId,
						},
					}
				} else if r.Id != current.Id {
					// Put the current AlbumView in the results list.
					entities = append(entities, current)
					// Then change the current AlbumView
					current = business.AlbumView{
						Album: domain.Album{Id: r.Id,
							Title: r.Title,
							Year: r.Year,
							ArtistId: r.ArtistId,
						},
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
