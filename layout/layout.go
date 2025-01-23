package layout

import (
	"fmt"
	"log"
	"print/win32"
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

func CenterElement(pageWidth, elementWidth uint32) uint32 {
	return (pageWidth - elementWidth) / 2
}

func AlignRight(pageWidth, elementWidth uint32) uint32 {
	return pageWidth - elementWidth
}

func AlignLeft() uint32 {
	return 0
}

func AlignRightFrom(startX, elementWidth uint32) uint32 {
	return startX - elementWidth
}

func AlignLeftFrom(startX uint32) uint32 {
	return startX
}

func CenterElementFrom(startX, containerWidth, elementWidth uint32) uint32 {
	return startX + (containerWidth-elementWidth)/2
}

func AlignTopFrom(startY, elementHeight uint32) uint32 {
	return startY
}

func AlignBottomFrom(startY, containerHeight, elementHeight uint32) uint32 {
	return startY + containerHeight - elementHeight
}

func CenterElementVerticallyFrom(startY, containerHeight, elementHeight uint32) uint32 {
	return startY + (containerHeight-elementHeight)/2
}

func truncateString(str string, maxLength int) string {
	if len(str) > maxLength {
		return str[:maxLength-1] + "."
	}
	return str
}

func formatArticles(ticket Ticket) [][]string {
	var formattedArticles [][]string
	for _, article := range ticket.Articles {
		formattedArticle := []string{
			fmt.Sprintf("%-40s", truncateString(fmt.Sprintf("%d x %s", article.Quantite, article.Nom), 40)),        // Première colonne, largeur fixe de 40 caractères
			fmt.Sprintf("%-15s", truncateString(fmt.Sprintf("à %.2f €", article.Prix), 15)),                        // Deuxième colonne, largeur fixe de 10 caractères
			fmt.Sprintf("%10s", truncateString(fmt.Sprintf("%.2f €", article.Prix*float64(article.Quantite)), 10)), // Troisième colonne, largeur fixe de 10 caractères
		}
		formattedArticles = append(formattedArticles, formattedArticle)
	}
	return formattedArticles
}

func DrawArticlesTab(dc win32.HDC, pageWidth, margin, startY uint32, ticket Ticket) {
	formattedArticles := formatArticles(ticket)

	for _, article := range formattedArticles {
		// Dessine la première colonne
		text := article[0]
		x := AlignLeft() + margin
		err := win32.TextOut(dc, x, startY, text, uint32(len(text)))
		if err != nil {
			log.Printf("TextOut failed: %s", err)
		}

		// Dessine la deuxième colonne
		text = article[1]
		x = AlignLeftFrom(4 * pageWidth / 7)
		err = win32.TextOut(dc, x, uint32(startY), text, uint32(len(text)))
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
		x = AlignRight(pageWidth, textWidth)
		err = win32.TextOut(dc, x, uint32(startY), text, uint32(len(text)))
		if err != nil {
			log.Printf("TextOut failed: %s", err)
		}

		startY += textHeight + 20 // Incrémente la position verticale pour le prochain article
	}
}

func DrawHeader(dc win32.HDC, pageWidth, margin, startY uint32, pdv Pdv) uint32 {
	// Dessiner le nom du PDV en gras
	win32.SetBoldFont(dc, true)

	text := pdv.Nom
	textWidth, textHeight, err := win32.GetTextExtentPoint32(syscall.Handle(dc), text)
	if err != nil {
		log.Printf("Get text dimensions failed: %s", err)
	}
	x := CenterElement(pageWidth, textWidth)
	err = win32.TextOut(dc, x, startY, text, uint32(len(text)))
	if err != nil {
		log.Printf("TextOut failed: %s", err)
	}
	startY += textHeight + 50

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
		x = CenterElement(pageWidth, textWidth)
		err = win32.TextOut(dc, x, startY, info, uint32(len(info)))
		if err != nil {
			log.Printf("TextOut failed: %s", err)
		}
		startY += textHeight + 15
	}
	startY += 15

	timestamp := time.Now().Format("02/01/2006 15:04:05")
	textWidth, textHeight, err = win32.GetTextExtentPoint32(syscall.Handle(dc), timestamp)

	if err != nil {
		log.Printf("Get text dimensions failed: %s", err)
	}
	x = CenterElement(pageWidth, textWidth)
	err = win32.TextOut(dc, x, startY, timestamp, uint32(len(timestamp)))
	if err != nil {
		log.Printf("TextOut failed: %s", err)
	}
	startY += textHeight + 30

	// Dessiner la ligne de séparation
	separator := "------------------------------------------------"
	textWidth, textHeight, err = win32.GetTextExtentPoint32(syscall.Handle(dc), separator)
	if err != nil {
		log.Printf("Get text dimensions failed: %s", err)
	}
	x = CenterElement(pageWidth, textWidth)
	err = win32.TextOut(dc, x, startY, separator, uint32(len(separator)))
	if err != nil {
		log.Printf("TextOut failed: %s", err)
	}
	startY += textHeight + 50

	return startY
}
