package domain

type Album struct {
	Id       int    `db:"id"`
	Title    string `db:"title"`
	Image    string `db:"image"`
	Year     string `db:"year"`
	ArtistId int    `db:"artist_id"`
	Tracks   Tracks	`db:"-"`
}

type Albums []Album
