package domain

type Artist struct {
	Id   	int     `db:"id"`
	Name 	string  `db:"name"`
	Albums 	Albums  `db:"-"`
}

type Artists []Artist
