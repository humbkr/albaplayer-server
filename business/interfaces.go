package business
// TODO: not sure this is used anywhere.
type LibraryRepository interface {
	Erase()
	Update()
	Clean()
}