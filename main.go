package main

import (
	"elasapp/officer"
	"elasapp/search"
	"elasapp/violation"
	"strings"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
)

func mainMenu(buttonNames []string, fNames []func()) dec.VSplitter {
	vSplitter := dec.VSplitter{}
	vSplitter.Children = []dec.Widget{}

	for i := range buttonNames {
		widGet := dec.PushButton{
			Text:      buttonNames[i],
			OnClicked: fNames[i],
			Font:      dec.Font{PointSize: 12},
		}
		vSplitter.Children = append(vSplitter.Children, widGet)
	}

	return vSplitter
}

func main() {

	var inTE, outTE *walk.TextEdit

	buttonNames := []string{"Διοικητής", "Αρμόδιος", "Παράβαση", "Αναζήτηση"}
	fNames := []func(){
		func() {
			officer.Init(GetDB(), 1)
		}, func() {
			officer.Init(GetDB(), 0)
		}, func() {
			violation.Init(GetDB(), "")
		}, func() {
			search.Init(GetDB())
		}, func() {
			outTE.SetText(strings.ToUpper(inTE.Text()))
		},
	}

	dec.MainWindow{
		Title:  "Εφαρμογή ΕΛΑΣ",
		Bounds: dec.Rectangle{Width: 800, Height: 600},
		Layout: dec.VBox{},
		Children: []dec.Widget{
			dec.HSplitter{
				Children: []dec.Widget{
					mainMenu(buttonNames, fNames),
					dec.TextEdit{AssignTo: &inTE},
				},
			},
		},
	}.Run()

	// Or read from memory
	// r, err := docx.ReadDocxFromMemory(data io.ReaderAt, size int64)

	db.Close()

	// Or write to ioWriter
	// docx2.Write(ioWriter io.Writer)

}
