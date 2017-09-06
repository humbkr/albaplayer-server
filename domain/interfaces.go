package domain

type ArtistRepository interface {
	// Gets an entity from the datasource.
	//
	// Returns an hydrated entity if entity is found, else an error.
	Get(id int) (entity Artist, err error)

	// Gets all entities from the datasource.
	//
	// If no entities found, returns an empty collection without error.
	GetAll(hydrate bool) (entities Artists, err error)

	// Gets an entity based on its name.
	GetByName(name string) (entity Artist, err error)

	// Saves an entity to a datasource.
	Save(entity *Artist) (err error)

	// Deletes an entity from a datasource.
	//
	// Does not return an error if the entity doesn't exists on the datasource or no entity id is given.
	Delete(entity *Artist) (err error)

	// Tests if an entity exists in datasource.
	Exists(id int) bool
}

type AlbumRepository interface {
	// Gets an entity from a datasource.
	//
	// Returns an hydrated entity if entity is fund, else an error.
	Get(id int) (entity Album, err error)

	// Gets all entities from the datasource.
	//
	// If no entities found, returns an empty collection without error.
	GetAll(hydrate bool) (entities Albums, err error)

	// Gets an entity based on its name.
	GetByName(name string, artistId int) (entity Album, err error)

	// Gets all albums for a given artist.
	//
	// If hydrate == true, hydrate the sub objects. If no album found, returns an empty collection without error.
	GetAlbumsForArtist(artistId int, hydrate bool) (entities Albums, err error)

	// Saves an entity to a datasource.
	Save(entity *Album) (err error)

	// Deletes an entity from a datasource.
	//
	// Does not return an error if the entity doesn't exists on the datasource or no entity id is given.
	Delete(entity *Album) (err error)

	// Tests if an entity exists in datasource.
	Exists(id int) bool
}

type TrackRepository interface {
	// Gets an entity from a datasource.
	//
	// Returns an hydrated entity if entity is fund, else an error.
	Get(id int) (entity Track, err error)

	// Gets all entities from the datasource.
	//
	// If no entities found, returns an empty collection without error.
	GetAll() (entities Tracks, err error)

	// Gets an entity based on its name.
	GetByName(name string, artistId int, albumId int) (entity Track, err error)

	// Gets all tracks for a given album.
	//
	// If hydrate == true, hydrate the sub objects. If no track found, returns an empty collection without error.
	GetTracksForAlbum(albumId int) (entities Tracks, err error)

	// Saves an entity to a datasource.
	Save(entity *Track) (err error)

	// Deletes an entity from a datasource.
	//
	// Does not return an error if the entity doesn't exists on the datasource or no entity id is given.
	Delete(entity *Track) (err error)

	// Tests if an entity exists in datasource.
	Exists(id int) bool
}
