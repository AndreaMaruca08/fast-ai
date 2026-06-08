package cli

func HomePage() *Page {
	ClearTerminal()
	page := NewPage(
		"Home",
		"Questo è un client gratuito per interagire con AI senza consumi di ram e cpu alti\n"+
			"Nessun lag dovuti all'app\n"+
			"ESC per uscire\n"+
			"1 - Home\n"+
			"2 - Chat\n"+
			"3 - Chat coding LOW\n"+
			"4 - Chat coding HIGH\n"+
			"5 - Usage and info\n", false)
	page.Update()

	return page
}
