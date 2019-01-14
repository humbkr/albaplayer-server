package domain

type Playlist struct {
	Id int               `db:"id"`
	Title string         `db:"title"`
	Tracks []Track 		 `db:"-"`
}

func (p *Playlist) Add(tracks ...Track) {
	p.Tracks = append(p.Tracks, tracks...)
}

func (p *Playlist) Remove(trackIndex int) {
	p.Tracks = append(p.Tracks[:trackIndex], p.Tracks[trackIndex + 1:]...)
}

func (p *Playlist) Clear() {
	p.Tracks = nil
}

type Playlists []Playlist
