package business

import "git.humbkr.com/jgalletta/alba-player/domain"

type LibraryInteractor struct {
	ArtistRepository  domain.ArtistRepository
	AlbumRepository   domain.AlbumRepository
	TrackRepository   domain.TrackRepository
	LibraryRepository LibraryRepository
}

func (interactor LibraryInteractor) UpdateLibrary() {
	interactor.LibraryRepository.Update()
}

func (interactor LibraryInteractor) EraseLibrary() {
	interactor.LibraryRepository.Erase()
}

func (interactor LibraryInteractor) CleanDeadFiles() {
	interactor.LibraryRepository.Clean()
}
