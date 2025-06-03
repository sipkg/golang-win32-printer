// Simple demo that can be cross compiled from Linux to
// Windows and executed via Wine with a printer named `PDF`
//
// ```go
// GOOS=windows GOARCH=amd64 go build cmd
// ```
package main

import (
	"log"

	"github.com/sipkg/golang-win32-printer/ticket"
	"github.com/sipkg/golang-win32-printer/win32"
)

var t = ticket.Ticket{
	PdvID: 1,
	Articles: []ticket.Article{
		{Nom: "Article 1", Quantite: 2, Prix: 150.0},
		{Nom: "Article 2", Quantite: 1, Prix: 5.0},
		{Nom: "Article 3", Quantite: 3, Prix: 1.0},
		{Nom: "Article 4", Quantite: 2, Prix: 2.0},
		{Nom: "Article 5", Quantite: 2, Prix: 3.0},
		{Nom: "Article 6", Quantite: 4, Prix: 5.0},
		{Nom: "Article avec un nom vraiment très long histoire de voir ce que ça fait", Quantite: 1, Prix: 2.0},
	},
	Total: 338,
}

var pdv = ticket.Pdv{
	ID:      1,
	Nom:     "Mon Magasin",
	Adresse: "123 Rue Exemple, 75000 Paris",
	Tel:     "01 23 45 67 89",
	Mail:    "contact@monmagasin.fr",
}

const (
	textHeight = 64
	margin     = 350
)

func main() {
	// printName := "Microsoft Print to PDF"
	printName := "PDF"
	dc, err := win32.CreateDC(printName)
	if err != nil {
		log.Fatalf("CreateDC failed: %s", err)
	}
	err = win32.StartDCPrinter(dc, "gdiDoc")
	if err != nil {
		log.Fatalf("StartDCPrinter failed: %s", err)
	}
	err = win32.StartPage(dc)
	if err != nil {
		log.Fatalf("StartPage failed: %s", err)
	}

	win32.SetFont(dc, win32.CourierNew)
	width, err := win32.GetDeviceCaps(dc, win32.HORZRES)
	if err != nil {
		log.Fatalf("Retreiving page width failed: %s", err)
	}
	_, err = win32.GetDeviceCaps(dc, win32.VERTRES)
	if err != nil {
		log.Fatalf("Retreiving page height failed: %s", err)
	}

	oldcol, err := win32.SetTextColor(dc, win32.RGB(0, 0, 0))
	if err != nil {
		log.Printf("SetTextColor failed: %s", err)
	}
	log.Printf("Before color was %v", oldcol)

	_, err = win32.SetTextSize(dc, textHeight)
	if err != nil {
		log.Printf("SetTextSize failed: %s", err)
	}

	headerStopsY := ticket.DrawHeader(dc, width, margin, pdv)

	tabStopsY, totalArticles, aPayer := ticket.DrawArticlesTab(dc, width, margin, headerStopsY, t)

	ticket.DrawFooter(dc, width, margin, tabStopsY, totalArticles, aPayer)

	err = win32.EndPage(dc)
	if err != nil {
		log.Printf("EndPage failed: %s", err)
	}
	err = win32.EndDoc(dc)
	if err != nil {
		log.Printf("EndDoc failed: %s", err)
	}
	err = win32.DeleteDC(dc)
	if err != nil {
		log.Printf("DeleteDC failed: %s", err)
	}
}
