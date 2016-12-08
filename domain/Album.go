package domain

type Album struct {
    Type       string    `db:"-"`
    Id         int       `db:"id"`
    Name       string    `db:"name"`
    Image      string    `db:"image"`
    Year       string    `db:"year"`
    ArtistId   int       `db:"artist_id"`
}

type Albums []Album
