package domain

import (
	"time"
)

type Track struct {
	Type     string        `db:"-"`
	Id       int           `db:"id"`
	Name     string        `db:"name"`
	AlbumId  int           `db:"album_id"`
	ArtistId int           `db:"artist_id"`
	Cd       string        `db:"cd"`
	Number   string        `db:"number"`
	Duration time.Duration `db:"duration"`
	Path     string        `db:"path"`
}

type Tracks []Track
