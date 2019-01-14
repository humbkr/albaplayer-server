package interfaces

import (
	"errors"

	"github.com/humbkr/albaplayer-server/internal/alba/domain"
)

// Type for managing table storing the mapping between tracks and playlists.
type playlistTrack struct {
	PlaylistId int	`db:"playlist_id"`
	TrackId int		`db:"track_id"`
	Position int	`db:"position"`
}

type PlaylistDbRepository struct {
	AppContext *AppContext
}

/*
Fetches a playlist from the database.
*/
func (pr PlaylistDbRepository) Get(id int) (entity domain.Playlist, err error) {
	object, err := pr.AppContext.DB.Get(domain.Playlist{}, id)
	if err == nil && object != nil {
		entity = *object.(*domain.Playlist)

		playlistItemsRepo := PlaylistItemDbRepository{AppContext: pr.AppContext}
		playlistItems, errItems := playlistItemsRepo.getItemsForPlaylist(entity.Id)
		if errItems == nil {
			entity.Tracks = playlistItems
		}
	} else {
		err = errors.New("no album found")
	}

	return
}

/*
Fetches all playlists from the database.

@param hydrate
	If true populate playlists tracks. WARNING: can be time consuming
*/
func (pr PlaylistDbRepository) GetAll(hydrate bool) (entities domain.Playlists, err error) {
	_, err = pr.AppContext.DB.Select(&entities, "SELECT * FROM playlists")
	if hydrate {
		for i := range entities {
			playlistItemsRepo := playlistItemDbRepository{AppContext: pr.AppContext}
			playlistItems, errItems := playlistItemsRepo.getItemsForPlaylist(entities[i].Id)
			if errItems == nil {
				for i := range playlistItems {

				}

				entities[i].Tracks = playlistItems
			}
		}
	}

	return
}

/*
Create or update a playlist in the Database.
*/
func (pr PlaylistDbRepository) Save(entity *domain.Playlist) (err error) {
	if entity.Id != 0 {
		// Update.
		_, err = pr.AppContext.DB.Update(entity)
		return
	} else {
		// Insert new entity.
		err = pr.AppContext.DB.Insert(entity)
		return
	}
}

/*
Delete a playlist from the Database.
TODO also delete Playlist items
*/
func (pr PlaylistDbRepository) Delete(entity *domain.Playlist) (err error) {
	_, err = pr.AppContext.DB.Delete(entity)
	return
}

// Check if a playlist exists for a given id.
func (pr PlaylistDbRepository) Exists(id int) bool {
	_, err := pr.Get(id)
	return err == nil
}

// ---------------------------------------------------------------------------------------------------------------------

type playlistItemDbRepository struct {
	AppContext *AppContext
}

/*
Get items for a playlist.
 */
func (pir playlistItemDbRepository) getItemsForPlaylist(playlistId int) (entities domain.PlaylistItems, err error) {
	_, err = pir.AppContext.DB.Select(&entities, "SELECT * FROM playlist_track WHERE playlist_id = ? ORDER BY position", playlistId)
	return
}

/*
Delete items from a playlist.
 */
func (pir playlistItemDbRepository) deleteItemsFromPlaylist(playlistId int) (err error) {
	_, err = pir.AppContext.DB.Exec("DELETE FROM playlist_track WHERE playliet_id = ?", playlistId)
	return
}
