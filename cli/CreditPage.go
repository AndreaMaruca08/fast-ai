package cli

import "fmt"

func CreditPage(page *Page) *Page {
	ClearTerminal()

	page = NewPage("______        _           ___  _ \n"+
		"|  ___|      | |         / _ \\(_)\n"+
		"| |_ __ _ ___| |_ ______/ /_\\ \\_ \n"+
		"|  _/ _` / __| __|______|  _  | |\n"+
		"| || (_| \\__ \\ |_       | | | | |\n"+
		"\\_| \\__,_|___/\\__|      \\_| |_/_|\n"+
		"                                 \n"+
		"                                 ",
		"\nCrediti: https://github.com/AndreaMaruca08", true,
	)
	currentPage = page
	fmt.Print("\033[34m")
	page.Update()
	fmt.Print("\033[0m")

	return page
}
