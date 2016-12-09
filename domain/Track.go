package domain

type Track struct {
	Type     string `db:"-"`
	Id       int    `db:"id"`
	Name     string `db:"name"`
	AlbumId  int    `db:"album_id"`
	ArtistId int    `db:"artist_id"`
	Disc     string `db:"disc"`
	Number   int    `db:"number"`
	Duration int    `db:"duration"` // Duration in seconds.
	Path     string `db:"path"`
}

type Tracks []Track
