package ticket

import (
	"fmt"
	"log"
	"github.com/clementuu/golang-win32-printer/layout"
	"github.com/clementuu/golang-win32-printer/win32"
	"syscall"
	"time"
)

type Article struct {
	Nom      string
	Prix     float64
	Quantite int
}

type Ticket struct {
	PdvID    int
	Articles []Article
	Total    float64
}

type Pdv struct {
	ID      int
	Nom     string
	Adresse string
	Tel     string
	Mail    string
}

func formatArticles(ticket Ticket) ([][]string, int) {
	var formattedArticles [][]string
	var totalArticles = 0
	for _, article := range ticket.Articles {
		totalArticles += article.Quantite
		formattedArticle := []string{
			fmt.Sprintf("%-35s", layout.TruncateString(fmt.Sprintf("%d x %s", article.Quantite, article.Nom), 35)),        // Première colonne, largeur fixe de 35 caractères
			fmt.Sprintf("%-15s", layout.TruncateString(fmt.Sprintf("à %.2f €", article.Prix), 15)),                        // Deuxième colonne, largeur fixe de 15 caractères
			fmt.Sprintf("%10s", layout.TruncateString(fmt.Sprintf("%.2f €", article.Prix*float64(article.Quantite)), 10)), // Troisième colonne, largeur fixe de 10 caractères
		}
		formattedArticles = append(formattedArticles, formattedArticle)
	}
	return formattedArticles, totalArticles
}

func DrawSeparator(dc win32.HDC, pageWidth, startY uint32) uint32 {
	// Dessiner la ligne de séparation
	separator := "--------------------------------------------------------"
	textWidth, textHeight, err := win32.GetTextExtentPoint32(syscall.Handle(dc), separator)
	if err != nil {
		log.Printf("Get text dimensions failed: %s", err)
	}
	x := layout.CenterElement(pageWidth, textWidth)
	err = win32.TextOut(dc, x, startY, separator, uint32(len(separator)))
	if err != nil {
		log.Printf("TextOut failed: %s", err)
	}

	startY += textHeight

	return startY + 100
}

func DrawArticlesTab(dc win32.HDC, pageWidth, margin, startY uint32, ticket Ticket) (uint32, int, float64) {
	formattedArticles, totalArticles := formatArticles(ticket)

	for _, article := range formattedArticles {
		// Dessine la première colonne
		text := article[0]
		x := layout.AlignLeft() + margin
		err := win32.TextOut(dc, x, startY, text, uint32(len(text)))
		if err != nil {
			log.Printf("TextOut failed: %s", err)
		}

		// Dessine la deuxième colonne
		text = article[1]
		x = layout.AlignLeftFrom(4 * pageWidth / 7)
		err = win32.TextOut(dc, x, startY, text, uint32(len(text)))
		if err != nil {
			log.Printf("TextOut failed: %s", err)
		}

		// Dessine la troisième colonne
		text = article[2]
		textWidth, textHeight, err := win32.GetTextExtentPoint32(syscall.Handle(dc), text)
		if err != nil {
			log.Printf("Get text dimensions failed: %s", err)
			continue
		}
		x = layout.AlignRight(pageWidth-margin, textWidth)
		err = win32.TextOut(dc, x, startY, text, uint32(len(text)))
		if err != nil {
			log.Printf("TextOut failed: %s", err)
		}

		startY += textHeight + 20 // Incrémente la position verticale pour le prochain article
	}

	startY = DrawSeparator(dc, pageWidth, startY)

	return startY, totalArticles, ticket.Total
}

func DrawHeader(dc win32.HDC, pageWidth, startY uint32, pdv Pdv) uint32 {
	// Dessiner le nom du PDV en gras
	win32.SetBoldFont(dc, true)

	text := pdv.Nom
	textWidth, textHeight, err := win32.GetTextExtentPoint32(syscall.Handle(dc), text)
	if err != nil {
		log.Printf("Get text dimensions failed: %s", err)
	}
	originalTextHeight, err := win32.SetTextSize(dc, int32(5*textHeight/3))
	if err != nil {
		log.Printf("SetTextSize failed: %s", err)
	}
	x := layout.CenterElement(pageWidth, 5*textWidth/3)
	err = win32.TextOut(dc, x, startY, text, uint32(len(text)))
	if err != nil {
		log.Printf("TextOut failed: %s", err)
	}
	startY += 5*textHeight/3 + 50

	_, err = win32.SetTextSize(dc, originalTextHeight)
	if err != nil {
		log.Printf("SetTextSize failed: %s", err)
	}

	// Dessiner les autres informations du PDV
	win32.SetBoldFont(dc, false)

	infos := []string{
		pdv.Adresse,
		pdv.Tel,
		pdv.Mail,
	}

	for _, info := range infos {
		textWidth, textHeight, err := win32.GetTextExtentPoint32(syscall.Handle(dc), info)
		if err != nil {
			log.Printf("Get text dimensions failed: %s", err)
		}
		x = layout.CenterElement(pageWidth, textWidth)
		err = win32.TextOut(dc, x, startY, info, uint32(len(info)))
		if err != nil {
			log.Printf("TextOut failed: %s", err)
		}
		startY += textHeight + 20
	}
	startY += 20

	timestamp := time.Now().Format("02/01/2006 15:04:05")
	textWidth, textHeight, err = win32.GetTextExtentPoint32(syscall.Handle(dc), timestamp)

	if err != nil {
		log.Printf("Get text dimensions failed: %s", err)
	}
	x = layout.CenterElement(pageWidth, textWidth)
	err = win32.TextOut(dc, x, startY, timestamp, uint32(len(timestamp)))
	if err != nil {
		log.Printf("TextOut failed: %s", err)
	}
	startY += textHeight + 30

	startY = DrawSeparator(dc, pageWidth, startY)

	return startY
}

func DrawFooter(dc win32.HDC, pageWidth, margin, startY uint32, totalArticles int, total float64) {
	totalLine := []string{
		"Total à payer:",
		fmt.Sprintf("%10s", fmt.Sprintf("%2.f €", total)),
	}

	win32.SetBoldFont(dc, true)

	textWidth, textHeight, err := win32.GetTextExtentPoint32(syscall.Handle(dc), totalLine[0])
	if err != nil {
		log.Printf("Get text dimensions failed: %s", err)
	}
	originalTextHeight, err := win32.SetTextSize(dc, int32(6*textHeight/5))
	if err != nil {
		log.Printf("SetTextSize failed: %s", err)
	}
	x := layout.CenterElement(pageWidth/2, textWidth)
	err = win32.TextOut(dc, x, startY, totalLine[0], uint32(len(totalLine[0])))
	if err != nil {
		log.Printf("TextOut failed: %s", err)
	}

	textWidth, textHeight, err = win32.GetTextExtentPoint32(syscall.Handle(dc), totalLine[1])
	if err != nil {
		log.Printf("Get text dimensions failed: %s", err)
	}
	x = layout.AlignRight(pageWidth-margin, textWidth)

	err = win32.TextOut(dc, x, startY, totalLine[1], uint32(len(totalLine[1])))
	if err != nil {
		log.Printf("TextOut failed: %s", err)
	}

	startY += textHeight + 100

	win32.SetBoldFont(dc, false)
	win32.SetTextSize(dc, int32(originalTextHeight))

	startY = DrawSeparator(dc, pageWidth, startY)

	text := "Nous vous remercions de votre visite !"
	textWidth, textHeight, err = win32.GetTextExtentPoint32(syscall.Handle(dc), text)
	if err != nil {
		log.Printf("Get text dimensions failed: %s", err)
	}
	x = layout.CenterElement(pageWidth, textWidth)
	err = win32.TextOut(dc, x, startY, text, uint32(len(text)))
	if err != nil {
		log.Printf("TextOut failed: %s", err)
	}

	startY += textHeight + 20

	ticketNum := 123456789
	text = fmt.Sprintf("Nombre d'articles: %d, Ticket n° %d\n", totalArticles, ticketNum)
	textWidth, _, err = win32.GetTextExtentPoint32(syscall.Handle(dc), text)
	if err != nil {
		log.Printf("Get text dimensions failed: %s", err)
	}
	x = layout.CenterElement(pageWidth, textWidth)
	err = win32.TextOut(dc, x, startY, text, uint32(len(text)))
	if err != nil {
		log.Printf("TextOut failed: %s", err)
	}
}

func PrintA4(printName string, margin uint32, textHeight int32, pdv Pdv, ticket Ticket) {
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

	headerStopsY := DrawHeader(dc, width, margin, pdv)

	tabStopsY, totalArticles, aPayer := DrawArticlesTab(dc, width, margin, headerStopsY, ticket)

	DrawFooter(dc, width, margin, tabStopsY, totalArticles, aPayer)

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
