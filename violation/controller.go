package violation

import (
	"database/sql"
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
		{1, "ΔΙΑΒΙΒΑΣΤΙΚΟ ΓΙΑ  ΠΡΟΣΤΙΜΑ 20 ΚΑΙ 50 ΕΥΡΩ"},
		{2, "ΔΙΑΒΙΒΑΣΤΙΚΟ ΠΡΟΣΤΙΜΟ 175 ( ΑΦΟΡΑ ΜΕΙΚΤΗ ΥΠΗΡΕΣΙΑ 'Η Α.Τ. Γ.Α.Δ.Α. Κ΄ Γ.Α.Δ.Θ.) ΑΦΑΙΡΕΣΗ ΑΠΟ ΥΠΗΡΕΣΙΑ ΑΠΟΣΤΟΛΗΣ"},
		{3, "ΔΙΑΒΙΒΑΣΤΙΚΟ ΠΡΟΣΤΙΜΟ 175 (ΑΦΟΡΑ ΜΗ ΜΕΙΚΤΗ ΥΠΗΡΕΣΙΑ) ΕΚΔΟΣΗ ΑΠΟΦΑΣΗΣ 1 ΑΠΟ Τ.Τ ΠΕΡΙΟΧΗΣ"},
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

func Init(db *sql.DB, ap string) {
	values := []string{"", "", "", "", "", "", "", "", ""}
	comboIndex := -1
	var button dec.PushButton
	mw := new(MyMainWindow)

	if ap != "" {
		viol := GetByAP(db, ap)
		values = []string{viol.AP, viol.AT, viol.ViolationNumber,
			viol.RegistrationNumber, viol.FirstNameOwner, viol.LastNameOwner, viol.MiddleNameOwner, viol.AddressOwner}
		comboIndex = viol.DocumentType
		button = dec.PushButton{
			Text: "Προβολή",
			OnClicked: func() {
				UpdatePublishDate(db, ap, time.Now().Format("02/01/2006"))
			},
			Font: dec.Font{PointSize: 12},
		}
	} else {
		button = dec.PushButton{
			Text: "Αποθήκευση",
			OnClicked: func() {
				Insert(db, out[0].Text(), out[1].Text(), out[2].Text(), out[3].Text(), out[4].Text(), out[5].Text(),
					out[6].Text(), out[7].Text(), violType.CurrentIndex()+1)
				mw.Close()
			},
			Font: dec.Font{PointSize: 12},
		}
	}

	dec.MainWindow{
		Title:    "Καταχώρηση Παράβασης",
		AssignTo: &mw.MainWindow,
		Bounds:   dec.Rectangle{Width: 1200, Height: 300},
		Layout:   dec.VBox{},
		Children: []dec.Widget{
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
			button,
		},
	}.Create()
	mw.Run()
}
