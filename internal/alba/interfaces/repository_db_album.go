package interfaces

import (
	"errors"
	"reflect"
	"strconv"
	"time"

	"github.com/humbkr/albaplayer-server/internal/alba/business"
	"github.com/humbkr/albaplayer-server/internal/alba/domain"
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

// GetMultiple fetches albums from db based on filters.
func (ar AlbumDbRepository) GetMultiple(filter business.EntityFilter) (entities domain.Albums, err error) {
	// Build the query from the given filter options.
	selectClause := "SELECT id, title, year, artist_id, cover_id, added_at FROM albums"
	if filter.Hydrate {
		selectClause = "SELECT alb.Id AlbumId, alb.Title AlbumTitle, alb.Year AlbumYear, alb.artist_id AlbumArtistId, alb.cover_id AlbumCoverId, alb.added_at AlbumAddedAt, trk.* " +
			"FROM albums alb, tracks trk WHERE alb.id = trk.album_id"
	}

	limitClause := ""
	if filter.Limit != 0 {
		limitClause = " LIMIT " + strconv.Itoa(filter.Limit)
	}

	orderClause := ""
	if filter.Random {
		orderClause = " ORDER BY RANDOM()"
	} else if filter.Sort != "" {
		// Get the column name of an album field.
		field, found := reflect.TypeOf(domain.Album{}).FieldByName(filter.Sort)
		if found {
			tableFieldName := field.Tag.Get("db")
			orderClause = " ORDER BY " + tableFieldName

			if filter.SortOrder == "ASC" || filter.SortOrder == "DESC" {
				orderClause += " " + filter.SortOrder
			}
		}
	}

	query := selectClause + orderClause + limitClause

	// Query the database.
	if !filter.Hydrate {
		// It's simple.
		_, err = ar.AppContext.DB.Select(&entities, query)
	} else {
		// It becomes complex.
		type gorpResult struct {
			AlbumId int
			AlbumTitle string
			AlbumYear string
			AlbumArtistId int
			AlbumCoverId int
			AlbumAddedAt int
			domain.Track
			// Cannot select domain.album.ArtistId or domain.track.AlbumId because of a Gorp error...
			// So we have to join on trk.album_id, but then Gorp cannot do the mapping with gorpResult, so we have
			// to add this property in the struct. TODO get rid of gorp.
			Album_id int
		}
		var results []gorpResult

		_, err = ar.AppContext.DB.Select(&results, query)
		if err == nil {
			// Deduplicate stuff.
			var currentAlbum domain.Album
			for _, r := range results {
				// If row is not about the same album as the current one, add the current
				// one to the results.
				if currentAlbum.Id != 0  && r.Id != currentAlbum.Id {
					entities = append(entities, currentAlbum)
				}

				// Create album object if we have none or the one from the currentAlbum row
				// is different from the one from the previous row.
				if currentAlbum.Id == 0  || r.Id != currentAlbum.Id {
					currentAlbum = domain.Album{
						Id: r.AlbumId,
						Title: r.AlbumTitle,
						Year: r.AlbumYear,
						ArtistId: r.AlbumArtistId,
						CoverId: r.AlbumCoverId,
						AddedAt: r.AlbumAddedAt,
					}
				}

				// Create track object.
				track := domain.Track{
					Id: r.Id,
					Title: r.Title,
					AlbumId: r.AlbumId,
					ArtistId: r.ArtistId,
					CoverId: r.CoverId,
					Disc: r.Disc,
					Number: r.Number,
					Duration: r.Duration,
					Genre: r.Genre,
					Path: r.Path,
				}

				// Add the track from the current row to the current album.
				currentAlbum.Tracks = append(currentAlbum.Tracks, track)
			}

			// Add the last currentAlbum to the results.
			entities = append(entities, currentAlbum)
		}
	}

	return
}

/*
Fetches all albums from the database.

@param hydrate
	If true populate albums tracks. WARNING: VERY time consuming
*/
func (ar AlbumDbRepository) GetAll(hydrate bool) (entities domain.Albums, err error) {
	return ar.GetMultiple(business.EntityFilter{ Hydrate: hydrate })
}

/*
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

/*
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

// GetRandom returns X random albums.
func (ar AlbumDbRepository) GetRandom(number int, hydrate bool) (entities domain.Albums, err error) {
	_, err = ar.AppContext.DB.Select(&entities, "SELECT * FROM albums ORDER BY RANDOM() LIMIT ?", number)
	if err == nil && hydrate {
		for i := range entities {
			ar.populateTracks(&entities[i])
		}
	}

	return
}

// GetLastEntries returns the X albums that were added last.
func (ar AlbumDbRepository) GetLastEntries(number int, hydrate bool) (entities domain.Albums, err error) {
	_, err = ar.AppContext.DB.Select(&entities, "SELECT * FROM albums ORDER BY added_at LIMIT ?", number)
	if err == nil && hydrate {
		for i := range entities {
			ar.populateTracks(&entities[i])
		}
	}

	return
}

/*
Create or update an album in the Database.
*/
func (ar AlbumDbRepository) Save(entity *domain.Album) (err error) {
	if entity.Id != 0 {
		// Update.
		_, err = ar.AppContext.DB.Update(entity)
		return
	} else {
		// Insert new entity.
		currentTime := time.Now()
		entity.AddedAt, _ = strconv.Atoi(currentTime.Format(domain.DATE_FORMAT))

		err = ar.AppContext.DB.Insert(entity)
		return
	}
}

/*
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
