package domain

type ArtistRepository interface {
	Find(id int) (entity Artist, err error)
	FindAll(hydrate bool) (entities Artists, err error)
	FindByName(name string) (entity Artist, err error)
	Save(entity *Artist) (err error)
	Delete(artistId int) (err error)
}

type AlbumRepository interface {
	Find(id int) (entity Album, err error)
	FindAll(hydrate bool) (entities Albums, err error)
	FindByName(name string, artistId int) (entity Album, err error)
	Save(entity *Album) (err error)
	Delete(albumId int) (err error)
}

type TrackRepository interface {
	Find(id int) (entity Track, err error)
	FindAll() (entities Tracks, err error)
	FindByName(name string, artistId int, albumId int) (entity Track, err error)
	Save(entity *Track) (err error)
	Delete(trackId int) (err error)
}
