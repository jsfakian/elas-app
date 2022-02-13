package violation

import (
	"database/sql"
	"elasapp/docx"
	"elasapp/officer"
	"path"
	"time"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	log "github.com/sirupsen/logrus"
)

var (
	out      [8]*walk.TextEdit
	violType *walk.ComboBox
)

const (
	SampleDir string = "samples/parabasi"
	DocDir    string = "docs"
)

type MyMainWindow struct {
	*walk.MainWindow
}

type DropDownItem struct { // Used in the ComboBox dropdown
	Key  int
	Name string
}

func createFields(labels, names, values []string, comboIndex int) []dec.Widget {
	keys := []*DropDownItem{ // These are the items to populate the drop down list
		{1, "ΔΙΑΒΙΒΑΣΤΙΚΟ_ΓΙΑ_ΠΡΟΣΤΙΜΑ_20_ΚΑΙ_50_ΕΥΡΩ"},
		{2, "ΔΙΑΒΙΒΑΣΤΙΚΟ_ΠΡΟΣΤΙΜΟ_175_(ΑΦΟΡΑ_ΜΕΙΚΤΗ_ΥΠΗΡΕΣΙΑ_'Η_Α.Τ.Γ.Α.Δ.Α._Κ΄_Γ.Α.Δ.Θ.)_ΑΦΑΙΡΕΣΗ_ΑΠΟ_ΥΠΗΡΕΣΙΑ_ΑΠΟΣΤΟΛΗΣ"},
		{3, "ΔΙΑΒΙΒΑΣΤΙΚΟ_ΠΡΟΣΤΙΜΟ_175_(ΑΦΟΡΑ_ΜΗ_ΜΕΙΚΤΗ_ΥΠΗΡΕΣΙΑ)_ΕΚΔΟΣΗ_ΑΠΟΦΑΣΗΣ_1_ΑΠΟ_Τ.Τ_ΠΕΡΙΟΧΗΣ"},
	}

	fields := []dec.Widget{
		dec.HSplitter{
			Children: []dec.Widget{
				dec.TextLabel{
					Text: "Τύπος παράβασης",
					Font: dec.Font{PointSize: 12},
				},
				dec.ComboBox{
					AssignTo:      &violType,
					Value:         nil,    // Initial value if required
					Model:         keys,   // The array of drop down items
					DisplayMember: "Name", // The field to display "DropDownItem.Name"
					BindingMember: "Key",  // The field to bind too, ie the value "DropDownItem.Key"
					Font:          dec.Font{PointSize: 12},
					CurrentIndex:  comboIndex,
					Enabled:       comboIndex == -1,
				},
			},
		},
	}
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
					Text:     values[i],
					AssignTo: &out[i],
					ReadOnly: comboIndex != -1,
				},
			},
		}
		fields = append(fields, field)
	}

	return fields
}

func createDoc(dirName string, db *sql.DB) {
	c := officer.GetByCommander(db, 1)
	o := officer.GetByCommander(db, 0)
	viol := GetByViolationNumber(db, out[2].Text())

	docx.EditDoc(path.Join(SampleDir, violType.Text()+".docx"), path.Join(DocDir, dirName, violType.Text()+".docx"),
		[]string{"armodios", "protokolo", "att", "imniaekdosis", "epithetoidioktiti", "onomaidioktiti",
			"patronimoidioktiti", "dieuthunsiidioktiti", "arithmoskykloforias", "diikitis"},
		[]string{o.FirstName + " " + o.LastName + " " + o.Rank, viol.AP, viol.AT, viol.PublishDate, viol.LastNameOwner,
			viol.FirstNameOwner, viol.MiddleNameOwner, viol.AddressOwner, viol.RegistrationNumber,
			c.FirstName + " " + c.LastName + " " + c.Rank})

	docx.OpenDocx(path.Join(DocDir, dirName, violType.Text()+".docx"))
}

func Init(db *sql.DB, ap string) {
	comboIndex := -1
	viol := GetByAP(db, ap)
	values := []string{
		viol.AP,
		viol.AT,
		viol.ViolationNumber,
		viol.RegistrationNumber,
		viol.FirstNameOwner,
		viol.LastNameOwner,
		viol.MiddleNameOwner,
		viol.AddressOwner,
	}

	comboIndex = viol.DocumentType
	children := []dec.Widget{
		dec.VSplitter{
			Children: createFields(
				[]string{
					"Αριθμός Πρωτοκόλλου",
					"Διεύθυνση Α.Τ.",
					"Αριθμός παράβασης",
					"Αριθμός κυκλοφορίας",
					"Όνομα ιδιοκτήτη",
					"Επώνυμο ιδιοκτήτη",
					"Πατρώνυμο ιδιοκτήτη",
					"Διεύθυνση ιδιοκτήτη",
				},
				[]string{
					"ap",
					"at",
					"violationNumber",
					"registrationNumber",
					"firstNameOwner",
					"lastNameOwner",
					"middleNameOwner",
					"addressOwner",
				},
				values,
				comboIndex,
			),
		},
	}
	mw := new(MyMainWindow)

	if ap != "" {
		button := dec.PushButton{
			Text: "Προβολή",
			OnClicked: func() {
				UpdatePublishDate(db, ap, time.Now().Format("02/01/2006"))
				createDoc(out[2].Text(), db)
			},
			Font: dec.Font{PointSize: 12},
		}
		children = append(children, button)
	} else {
		button := dec.PushButton{
			Text: "Αποθήκευση",
			OnClicked: func() {
				Insert(db, out[0].Text(), out[1].Text(), out[2].Text(), out[3].Text(), out[4].Text(), out[5].Text(),
					out[6].Text(), out[7].Text(), violType.CurrentIndex()+1)
				mw.Close()
			},
			Font: dec.Font{PointSize: 12},
		}
		children = append(children, button)
	}

	dec.MainWindow{
		Title:    "Καταχώρηση Παράβασης",
		AssignTo: &mw.MainWindow,
		Bounds:   dec.Rectangle{Width: 1200, Height: 300},
		Layout:   dec.VBox{},
		Children: children,
	}.Create()
	mw.Run()
}
