package domain

type Playlist struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
	Tracks   Tracks	`db:"-"`
}

type Playlists []Playlist
