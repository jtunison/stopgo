package stopgo

import (
	"code.google.com/p/go-qrcode"
	"code.google.com/p/gofpdf"
	"fmt"
	"os"
)

type pdfWriter struct {
	pdf *gofpdf.Fpdf
	fl  *os.File
}

func (pw *pdfWriter) Write(p []byte) (n int, err error) {
	if pw.pdf.Ok() {
		return pw.fl.Write(p)
	}
	return
}

func (pw *pdfWriter) Close() (err error) {
	if pw.fl != nil {
		pw.fl.Close()
		pw.fl = nil
	}
	if pw.pdf.Ok() {
		fmt.Printf("Successfully generated resume.pdf\n")
	} else {
		fmt.Printf("%s\n", pw.pdf.Error())
	}
	return
}

func docWriter(pdf *gofpdf.Fpdf, fileStr string) *pdfWriter {
	pw := new(pdfWriter)
	pw.pdf = pdf
	if pdf.Ok() {
		var err error
		pw.fl, err = os.Create(fileStr)
		if err != nil {
			pdf.SetErrorf("Error opening output file %s", fileStr)
		}
	}
	return pw
}

func Write(outputName string, resume *Resume) {
	var y0 float64
	pdf := gofpdf.New("P", "mm", "letter", "")
	const (
		pageWd                        = 216.0 // letter 216 mm x 279 mm
		pageHeight                    = 279.0
		margin                        = 10.0
		gutter                        = 4
		colNum                        = 2
		mainColumnWidth       float64 = (pageWd - 2*margin - gutter) * 3 / 4
		supplementColumnWidth float64 = (pageWd - 2*margin - gutter) * 1 / 4
		colWd                         = (pageWd - 2*margin - (colNum-1)*gutter) / colNum
		fontSize                      = 9.75
		fontFamily                    = "Roboto"
	)
	lineHeight := pdf.PointConvert(fontSize) * 1.25
	columnWidths := [2]float64{mainColumnWidth, supplementColumnWidth}

	// set up font
	gofpdf.MakeFont("font/Roboto/Roboto-Light.ttf", "font/cp1252.map", "font", nil, true)
	gofpdf.MakeFont("font/Roboto/Roboto-Regular.ttf", "font/cp1252.map", "font", nil, true)
	gofpdf.MakeFont("font/Droid_Serif/DroidSerif.ttf", "font/cp1252.map", "font", nil, true)
	gofpdf.MakeFont("font/Droid_Serif/DroidSerif-Bold.ttf", "font/cp1252.map", "font", nil, true)
	gofpdf.MakeFont("font/Playfair_Display/PlayfairDisplay-Regular.ttf", "font/cp1252.map", "font", nil, true)
	gofpdf.MakeFont("font/Playfair_Display/PlayfairDisplay-Bold.ttf", "font/cp1252.map", "font", nil, true)
	gofpdf.MakeFont("font/glyphicons-halflings-regular.ttf", "font/glyphicons.map", "font", nil, true)
	pdf.SetFontLocation("font")
	pdf.SetTitle("", true)
	pdf.SetAuthor("John Tunison", true)
	pdf.SetSubject("Resume", true)
	pdf.SetCreator("John Tunison", true)
	pdf.SetKeywords("rockstar", true)
	pdf.AddFont("Roboto", "", "Roboto-Light.json")
	pdf.AddFont("Roboto", "B", "Roboto-Regular.json")
	pdf.AddFont("halflings", "", "glyphicons-halflings-regular.json")
	pdf.AddFont("DroidSerif", "", "DroidSerif.json")
	pdf.AddFont("DroidSerif", "B", "DroidSerif-Bold.json")
	pdf.AddFont("Playfair", "", "PlayfairDisplay-Regular.json")
	pdf.AddFont("Playfair", "B", "PlayfairDisplay-Bold.json")

	setCol := func(col int) {
		x := margin
		for j := 0; j < col; j++ {
			x += columnWidths[j] + gutter
		}

		// log.Printf("setCol(%d) -> x = %f (%s)", col, x, columnWidths)
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
	}

	bullet := func(column int, text string) {
		// see http://www.fpdf.org/~~V/en/script/script38.php
		// see http://www.ascii-code.com for bullet character list
		bulletString := "\x95"
		bulletWidth := pdf.GetStringWidth(bulletString) + gutter/2
		columnWidth := columnWidths[column]
		pdf.Cell(bulletWidth, lineHeight, bulletString)
		pdf.MultiCell(columnWidth-bulletWidth, lineHeight, text, "", "L", false)
	}
	mainExperience := func(role string, start int, end int) {
		x := pdf.GetX()
		y := pdf.GetY()
		pdf.SetFont(fontFamily, "B", fontSize)
		pdf.MultiCell(mainColumnWidth, lineHeight, role, "", "L", false)
		pdf.SetXY(x, y)
		pdf.MultiCell(mainColumnWidth, lineHeight, fmt.Sprintf("%d - %d", start, end), "", "R", false)
		pdf.SetFont(fontFamily, "", fontSize)

	}
	horizontalRule := func(width, thickness float64) {
		x := pdf.GetX()
		y := pdf.GetY()
		pdf.SetLineWidth(thickness)
		pdf.SetDrawColor(191, 191, 191)
		pdf.Line(x, y, x+width, y)
		pdf.Ln(2)
	}
	heading := func(column int, text string) {

		x := pdf.GetX()
		y := pdf.GetY()
		height := pdf.PointConvert(fontSize * 1.25)
		columnWidth := columnWidths[column]

		// draw line first, then text (so text overlays line)
		switch column {
		case 0:
			pdf.SetXY(margin, y+height)
		case 1:
			pdf.SetXY(margin+mainColumnWidth+gutter, y+height)
		}
		horizontalRule(columnWidth, 0.2)

		// now heading text
		pdf.SetXY(x, y)
		pdf.SetFont("DroidSerif", "", fontSize*1.25)
		pdf.MultiCell(columnWidth, height, text, "", "L", false)
		pdf.Ln(2)
	}

	pdf.SetHeaderFunc(func() {
		titleStr := resume.Name

		x := margin
		y := pdf.GetY()
		lineHeight := pdf.PointConvert(fontSize * 3)

		// then the qr code
		// var png []byte
		// png, err := qrcode.Encode("https://example.org", qrcode.Medium, 256)
		qrCodeSize := lineHeight * 1.25
		qrcode.WriteFile(resume.Links.Website, qrcode.Medium, 256, "tmp/qr.png")
		pdf.Image("tmp/qr.png", pageWd-margin-qrCodeSize+2, y-2, qrCodeSize, qrCodeSize, false, "", 0, resume.Links.Website)

		// write horizontal rule first
		pdf.SetXY(margin, y+lineHeight)
		horizontalRule(mainColumnWidth+gutter+supplementColumnWidth, 0.4)

		// then write the name
		pdf.SetFont("DroidSerif", "B", fontSize*3)
		pdf.SetXY(x, y)
		pdf.SetTextColor(0, 0, 0)
		pdf.Write(lineHeight, titleStr)
		pdf.Ln(-1)

		// then the location
		pdf.Ln(1)
		lineHeight = pdf.PointConvert(fontSize) * 1.25
		pdf.SetFont(fontFamily, "", fontSize)
		pdf.Write(lineHeight, fmt.Sprintf("%s  \x95  %s  \x95  %s", resume.Location, resume.Links.Website, resume.Email))

		pdf.Ln(10)
		y0 = pdf.GetY()
	})

	pdf.SetFooterFunc(func() {
		footerStr := fmt.Sprintf("Resume generated by http://github.com/jtunison/stopgo.")
		pdf.SetFont(fontFamily, "", fontSize*3/4)
		width := pdf.GetStringWidth(footerStr)
		lineHeight = pdf.PointConvert(fontSize*3/4)
		x := (pageWd - width) / 2
		y := pageHeight - lineHeight - margin
		pdf.SetXY(x, y)
		pdf.SetTextColor(128, 128, 160)
		pdf.Write(lineHeight, footerStr)
	})

	pdf.AddPage()
	pdf.SetFont(fontFamily, "", fontSize)
	setCol(0)

	// Summary
	heading(0, "Summary")
	pdf.SetFont(fontFamily, "", fontSize)
	pdf.MultiCell(mainColumnWidth, lineHeight, resume.Summary, "", "L", false)
	pdf.Ln(-1)

	// Work History
	heading(0, "Work History")
	for _, experience := range resume.History {

		mainExperience(experience.Role, experience.StartYear, experience.EndYear)

		//Put a hyperlink
		pdf.SetTextColor(80, 139, 200)
		pdf.WriteLinkString(lineHeight, experience.Company, experience.CompanyUrl)
		pdf.SetTextColor(0, 0, 0)
		pdf.Ln(-1)
		pdf.Ln(1)

		for _, bulletContent := range experience.Bullets {
			bullet(0, bulletContent)
		}
		pdf.Ln(-1)
	}

	// Education
	heading(0, "Education")
	mainExperience(resume.Education[0].Institution, resume.Education[0].StartYear, resume.Education[0].EndYear)
	pdf.MultiCell(mainColumnWidth, lineHeight, resume.Education[0].Degree, "", "L", false)

	// right hand side
	pdf.SetY(y0)
	setCol(1)
	lineHeight = pdf.PointConvert(fontSize) * 1.4
	for _, supplement := range resume.Supplements {
		heading(1, supplement.Heading)

		for _, bulletContent := range supplement.Bullets {
			pdf.SetFont(fontFamily, "", fontSize)
			bullet(1, bulletContent)
		}

		pdf.Ln(-1)
	}

	pdf.OutputAndClose(docWriter(pdf, outputName))
}
