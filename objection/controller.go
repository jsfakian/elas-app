package objection

import (
	"database/sql"
	"elasapp/docx"
	"elasapp/officer"
	"elasapp/violation"
	"path"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	log "github.com/sirupsen/logrus"
)

var out [10]*walk.TextEdit
var objectionType *walk.ComboBox
var genderDriver *walk.ComboBox

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

func createFields(labels, names, values []string, gender string, documentType int) []dec.Widget {
	keys := []*DropDownItem{ // These are the items to populate the drop down list
		{1, "ΑΠΟΦΑΣΗ_ΑΠΟΡΡΙΨΗ_ΚΛΗΣΗΣ"},
		{2, "ΑΠΟΦΑΣΗ_ΑΠΟΡΡΙΨΗ_ΛΟΓΩ_ΕΚΠΡΟΘΕΣΜΩΝ_ΑΝΤΙΡΡΗΣΕΩΝ"},
		{3, "ΑΠΟΦΑΣΗ_ΑΠΟΡΡΙΨΗΣ_ΓΙΑ_POINT_SYSTEM"},
		{4, "ΑΠΟΦΑΣΗ_ΑΠΟΡΡΙΨΗΣ_ΓΙΑ_ΑΟ"},
		{5, "ΑΠΟΦΑΣΗ_ΑΠΟΡΡΙΨΗΣ_ΓΙΑ_ΙΑΤΡΙΚΟΥΣ_ΛΟΓΟΥΣ"},
		{6, "ΑΠΟΦΑΣΗ_ΔΕΚΤΗ_ΓΙΑ_ΑΟ"},
		{7, "ΑΠΟΦΑΣΗ_ΔΕΚΤΗ_ΓΙΑ_ΚΛΗΣΗ"},
	}
	genderKeys := []*DropDownItem{
		{1, "Αρσενικό"},
		{2, "Θηλυκό"},
	}
	fields := []dec.Widget{dec.HSplitter{
		Children: []dec.Widget{
			dec.TextLabel{
				Text: "Τύπος παράβασης",
				Font: dec.Font{PointSize: 12},
			},
			dec.ComboBox{
				AssignTo:      &objectionType,
				Value:         documentType, // Initial value if required
				Model:         keys,         // The array of drop down items
				DisplayMember: "Name",       // The field to display "DropDownItem.Name"
				BindingMember: "Key",        // The field to bind too, ie the value "DropDownItem.Key"
				Font:          dec.Font{PointSize: 12},
			},
		}}}
	log.Info(names)

	for i := range names {
		if names[i] != "genderDriver" {
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
		} else {
			genderIdx := 1
			if gender == "Θηλυκό" {
				genderIdx = 2
			}
			field := dec.HSplitter{
				Children: []dec.Widget{
					dec.TextLabel{
						Text: "Φύλο",
						Font: dec.Font{PointSize: 12},
					},
					dec.ComboBox{
						AssignTo:      &genderDriver,
						Value:         genderIdx,  // Initial value if required
						Model:         genderKeys, // The array of drop down items
						DisplayMember: "Name",     // The field to display "DropDownItem.Name"
						BindingMember: "Key",      // The field to bind too, ie the value "DropDownItem.Key"
						Font:          dec.Font{PointSize: 12},
					},
				},
			}
			fields = append(fields, field)
		}
	}

	return fields
}

func createDoc(dirName string, db *sql.DB) {
	c := officer.GetByCommander(db, 1)
	viol := violation.GetByViolationNumber(db, out[1].Text())
	toy := "του"
	ston := "στον"
	if genderDriver.Text() == "Θηλυκό" {
		toy = "της"
		ston = "στην"
	}

	docx.EditDoc(path.Join(SampleDir, objectionType.Text()+".docx"), path.Join(DocDir, dirName, objectionType.Text()+".docx"),
		[]string{"<protokolo>", "<imnia_ekdosis>", "<imnia_enstansis>", "<toy_odigos>", "<ston_odigos>",
			"<patronimo_odigou>", "<arithmos_paravasis>", "<diikitis>"},
		[]string{out[0].Text(), out[6].Text(), out[7].Text(), toy + " " + out[3].Text() + " " + out[2].Text(),
			ston + " " + out[3].Text() + " " + out[2].Text(), out[4].Text(), viol.ViolationNumber,
			c.FirstName + " " + c.LastName + " " + c.Rank})

	docx.OpenDocx(path.Join(DocDir, dirName, objectionType.Text()+".docx"))

}

func Init(db *sql.DB, violationNumber string) {
	objection := GetByViolationNumber(db, violationNumber)
	values := []string{
		objection.AP,
		objection.ViolationNumber,
		objection.FirstNameDriver,
		objection.LastNameDriver,
		objection.MiddleNameDriver,
		objection.GenderDriver,
		objection.ObjectionDate,
		objection.PublishDate,
	}
	mw := new(MyMainWindow)
	buttonText := "Αποθήκευση"
	if violationNumber != "" {
		buttonText = "Προβολή"
	}
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
						"Όνομα οδηγού",
						"Επίθετο οδηγού",
						"Πατρόνυμο οδηγού",
						"Φύλο οδηγού",
						"Ημ/νια ένστασης",
						"Ημ/νια έκδοσης",
					},
					[]string{
						"ap",
						"violationNumber",
						"firstNameDriver",
						"lastNameDriver",
						"middleNameDriver",
						"genderDriver",
						"objectionDate",
						"publishDate",
					},
					values,
					objection.GenderDriver,
					objection.DocumentType,
				),
			},
			dec.PushButton{
				Text: buttonText,
				OnClicked: func() {
					createDoc(out[1].Text(), db)
					if buttonText == "Αποθήκευση" {
						Insert(db, out[0].Text(), out[1].Text(), out[2].Text(), out[3].Text(), out[4].Text(),
							genderDriver.Text(), out[6].Text(), out[7].Text(), objectionType.CurrentIndex()+1)
						mw.Close()
					}
				},
				Font: dec.Font{PointSize: 12},
			},
		},
	}.Create()
	mw.Run()
}
