package domain

type Track struct {
	Id       int    `db:"id"`
	Title    string `db:"title"`
	AlbumId  int    `db:"album_id"`
	ArtistId int    `db:"artist_id"`
	Disc     string `db:"disc"`
	Number   int    `db:"number"`
	Duration int    `db:"duration"` // Duration in seconds.
	Path     string `db:"path"`
}

type Tracks []Track
