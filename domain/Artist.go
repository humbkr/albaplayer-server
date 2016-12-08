package domain

type Artist struct {
    Type       string    `db:"-"`
    Id         int       `db:"id"`
    Name       string    `db:"name"`
}

type Artists []Artist
