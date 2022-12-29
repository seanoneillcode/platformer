package actions

type Actions interface {
	OpenBook(title, text string)
	CloseBook()
}
