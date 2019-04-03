package domain

type Album struct {
	Id       int    `db:"id"`
	Title    string `db:"title"` // Mandatory.
	Year     string `db:"year"`
	ArtistId int    `db:"artist_id"`
	CoverId  int    `db:"cover_id"`
	Tracks   Tracks	`db:"-"`
	AddedAt  int    `db:"added_at"` // Format: YYYYMMDD
}

type Albums []Album
