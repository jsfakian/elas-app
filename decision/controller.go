package decision

import (
	"database/sql"
	"elasapp/docx"
	"elasapp/officer"
	"elasapp/violation"
	"path"
	"time"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	log "github.com/sirupsen/logrus"
)

var out [10]*walk.TextEdit
var decisionType *walk.ComboBox

const (
	SampleDir string = "samples/apofasi"
	DocDir    string = "docs"
)

type MyMainWindow struct {
	*walk.MainWindow
}

type DropDownItem struct { // Used in the ComboBox dropdown
	Key  int
	Name string
}

func createFields(labels, names, values []string) []dec.Widget {
	keys := []*DropDownItem{ // These are the items to populate the drop down list
		{1, "ΑΙΤΗΣΗ_ΑΠΟΡΡΙΨΗΣ_ΓΙΑ_ΑΟ"},
		{2, "ΑΙΤΗΣΗ_ΔΕΚΤΗ_ΓΙΑ_ΑΟ"},
		{3, "ΑΠΟΡΡΙΨΗ_ΓΙΑ_ΟΛΗ_ΤΗΝ_ΚΛΗΣΗ"},
		{4, "ΑΠΟΦΑΣΗ_ΑΚΥΡ.ΚΛΗΣΗΣ_ΡΑΝΤΑΡ"},
		{5, "ΑΠΟΦΑΣΗ_ΑΠΟΡΡΙΨΗΣ_ΕΝΣΤΑΣΕΩΝ"},
	}
	fields := []dec.Widget{dec.HSplitter{
		Children: []dec.Widget{
			dec.TextLabel{
				Text: "Τύπος παράβασης",
				Font: dec.Font{PointSize: 12},
			},
			dec.ComboBox{
				AssignTo:      &decisionType,
				Value:         nil,    // Initial value if required
				Model:         keys,   // The array of drop down items
				DisplayMember: "Name", // The field to display "DropDownItem.Name"
				BindingMember: "Key",  // The field to bind too, ie the value "DropDownItem.Key"
				Font:          dec.Font{PointSize: 12},
			},
		}}}
	log.Info(names)

	for i := range names {
		field := dec.HSplitter{
			Children: []dec.Widget{
				dec.TextLabel{
					Text: labels[i],
					Name: names[i] + "_label",
					Font: dec.Font{PointSize: 12},
				},
				dec.TextEdit{
					Name:     names[i] + "_input",
					Font:     dec.Font{PointSize: 12},
					AssignTo: &out[i],
					ReadOnly: values[i] != "",
					Text:     values[i],
				},
			},
		}
		fields = append(fields, field)
	}

	return fields
}

func createDoc(dirName string, db *sql.DB) {
	o := officer.GetByCommander(db, 0)
	c := officer.GetByCommander(db, 1)

	viol := violation.GetByViolationNumber(db, out[1].Text())

	docx.EditDoc(path.Join(SampleDir, decisionType.Text()+".docx"), path.Join(DocDir, dirName, decisionType.Text()+".docx"),
		[]string{"Αρχ/κας ΓΙΑΝΝΑΚΑΚΗΣ Αντώνιος", "2515/5/1", "916100093367", "29 Ιουλίου 2021", "2515/5/1/1138-κστ",
			"14/05/2021", "Εμμανουήλ ΤΑΜΠΑΚΑΚΗΣ", "Υπαστυνόμος Β΄"},
		[]string{o.Rank + " " + o.LastName + " " + o.FirstName, out[0].Text(), out[1].Text(), out[2].Text(),
			viol.AP, viol.PublishDate, c.FirstName + " " + c.LastName, c.Rank})

	docx.OpenDocx(path.Join(DocDir, dirName, decisionType.Text()+".docx"))

}

func Init(db *sql.DB, violationNumber string) {
	values := []string{"", violationNumber, time.Now().Format("02/01/2006")}
	mw := new(MyMainWindow)
	dec.MainWindow{
		Title:    "Καταχώρηση Παράβασης",
		AssignTo: &mw.MainWindow,
		Bounds:   dec.Rectangle{Width: 900, Height: 200},
		Layout:   dec.VBox{},
		Children: []dec.Widget{
			dec.VSplitter{
				Children: createFields(
					[]string{
						"Αριθμός πρωτοκόλλου",
						"Αριθμός παράβασης",
						"Ημ/νια ένστασης",
					},
					[]string{
						"ap",
						"violationNumber",
						"ObjectionDate",
					},
					values,
				),
			},
			dec.PushButton{
				Text: "Αποθήκευση",
				OnClicked: func() {
					Insert(db, out[0].Text(), out[1].Text(), out[2].Text())
					createDoc(out[1].Text(), db)
					mw.Close()
				},
				Font: dec.Font{PointSize: 12},
			},
		},
	}.Create()
	mw.Run()
}
