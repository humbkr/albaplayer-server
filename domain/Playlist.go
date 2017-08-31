package domain

type Playlist struct {
	Type string `db:"-"`
	Id   int    `db:"id"`
	Name string `db:"name"`
}

type Playlists []Playlist
