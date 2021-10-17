package main

import (
	"elasapp/officer"
	"elasapp/search"
	"elasapp/violation"
	"strings"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/nguyenthenguyen/docx"
	log "github.com/sirupsen/logrus"
)

func editDoc(filein, fileout string, oldText, newText []string) bool {
	r, err := docx.ReadDocxFile(filein)
	// Or read from memory
	// r, err := docx.ReadDocxFromMemory(data io.ReaderAt, size int64)
	if err != nil {
		log.Error(err)
		return false
	}

	docx1 := r.Editable()

	for i := range oldText {
		log.Info(oldText[i], newText[i])
		docx1.Replace(oldText[i], newText[i], -1)
		//docx1.Replace("old_1_2", "new_1_2", -1)
		//docx1.ReplaceLink("http://example.com/", "https://github.com/nguyenthenguyen/docx")
		//docx1.ReplaceHeader("out with the old", "in with the new")
		//docx1.ReplaceFooter("Change This Footer", "new footer")
	}

	docx1.WriteToFile(fileout)
	r.Close()

	return true
}

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
			violation.Init(GetDB())
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
					dec.TextEdit{AssignTo: &outTE, ReadOnly: true},
				},
			},
		},
	}.Run()

	// Or read from memory
	// r, err := docx.ReadDocxFromMemory(data io.ReaderAt, size int64)
	editDoc("./docs/a1.docx", "new_a1.docx", []string{"ΠΑΠΑΓΕΩΡΓΙΟΥ"}, []string{"ΠΑΠΑΓΕΩΡΓΙΟΥ2"})

	db.Close()

	// Or write to ioWriter
	// docx2.Write(ioWriter io.Writer)

}
