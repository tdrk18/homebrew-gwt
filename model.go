package main

type Model struct {
	Contexts     []Context
	Cursor       int
	SelectedPath string

	InputMode   InputMode
	InputText   string
	DeleteIndex int

	Error string
}

func newModel(contexts []Context) Model {
	return Model{
		Contexts: contexts,
		Cursor:   initialCursor(contexts),
	}
}
