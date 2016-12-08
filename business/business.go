package business

import "git.humbkr.com/jgalletta/alba-player/domain"

type CollectionInteractor struct {
	ArtistRepository  domain.ArtistRepository
	AlbumRepository   domain.AlbumRepository
	TrackRepository   domain.TrackRepository
	LibraryRepository domain.LibraryRepository
}

func (interactor CollectionInteractor) UpdateLibrary() {
	interactor.LibraryRepository.Update()
}

func (interactor CollectionInteractor) EraseLibrary() {
	interactor.LibraryRepository.Erase()
}
