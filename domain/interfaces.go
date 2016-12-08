package domain

type ArtistRepository interface {
	Find(id int) (entity Artist, err error)
	FindByName(name string) (entity Artist, err error)
	Save(entity *Artist) (err error)
}

type AlbumRepository interface {
	Find(id int) (entity Album, err error)
	FindByName(name string, artistId int) (entity Album, err error)
	Save(entity *Album) (err error)
}

type TrackRepository interface {
	Find(id int) (entity Track, err error)
	FindByName(name string, artistId int, albumId int) (entity Track, err error)
	Save(entity *Track) (err error)
}

type LibraryRepository interface {
	Update()
	Erase()
}
