package violation

import (
	"database/sql"
	"elasapp/officer"
	"os"
	"path"
	"time"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	"github.com/nguyenthenguyen/docx"
	log "github.com/sirupsen/logrus"
)

var (
	out              [10]*walk.TextEdit
	violType         *walk.ComboBox
	selectedDocument int
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
		{1, "e1_40_ΜΟΝΟ_ΠΡΟΣΤΙΜΟ"},
		{2, "e3_ΠΡΟΣΤΙΜΟ_ΚΑΙ_ΑΦΑΙΡ.ΣΕ_ΜΙΚΤΑ"},
		{3, "e5 ΑΠΟΦΑΣΗ ΑΦΑΙΡ. ΣΕ Τ.Τ"},
	}

	fields := []dec.Widget{dec.HSplitter{
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
					Text:     values[i],
					AssignTo: &out[i],
				},
			},
		}
		fields = append(fields, field)
	}

	return fields
}

func editDoc(filein, fileout string, oldText, newText []string) bool {
	log.Info(filein)
	log.Info(fileout)
	r, err := docx.ReadDocxFile(filein)
	// Or read from memory
	// r, err := docx.ReadDocxFromMemory(data io.ReaderAt, size int64)
	if err != nil {
		log.Error(err)
		return false
	}

	docx1 := r.Editable()

	for i := range oldText {
		log.Info(oldText[i], " || ", newText[i])
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

func createDoc(dirName string, db *sql.DB) {
	err := os.Mkdir(path.Join(DocDir, dirName), 0755)
	if err != nil {
		openDocx(path.Join(DocDir, dirName, violType.Text()+".docx"))
		return
	}

	o := officer.GetByCommander(db, 0)
	c := officer.GetByCommander(db, 1)

	editDoc(path.Join(SampleDir, violType.Text()+".docx"), path.Join(DocDir, dirName, violType.Text()+".docx"),
		[]string{"Αρχ/κας ΓΙΑΝΝΑΚΑΚΗΣ Αντώνιος", "2515/5/1/1227-κθ", `Α'Α.Τ.ΗΡΑΚΛΕΙΟΥ
		Αγίου Αρτεμίου 1 ΗΡΑΚΛΕΙΟ
		Τ.Κ.: 71601`, "916100093367", "άγνωστο οδηγό", "ΗΚΚ-\n2892", "ΗΚΚ-2892", "ΛΕΝΤΕΡΗΣ ΓΕΩΡΓΙΟΣ πατρ. ΑΠ",
			"ΚΑΤΑΥΛΙΣΜΟΣ ΑΛΙΚΑΡΝΑΣΣΟΥ ΗΡΑΚΛΕΙΟ ΤΚ: 71601", "Εμμανουήλ ΤΑΜΠΑΚΑΚΗΣ", "Υπαστυνόμος Β΄"},
		[]string{o.Rank + " " + o.LastName + " " + o.FirstName,
			out[0].Text(), out[1].Text(), out[2].Text(), out[3].Text() + " " + out[4].Text(), out[5].Text(),
			out[5].Text(), out[6].Text() + " " + out[7].Text() + " πατρ. " + out[8].Text(), out[9].Text(),
			c.FirstName + " " + c.LastName, c.Rank})

	openDocx(path.Join(DocDir, dirName, violType.Text()+".docx"))
}

func Init(db *sql.DB, ap string) {
	values := []string{"", "", "", "", "", "", "", "", "", "", ""}
	comboIndex := -1
	var button dec.PushButton
	mw := new(MyMainWindow)

	if ap != "" {
		viol := GetByAP(db, ap)
		values = []string{viol.AP, viol.AT, viol.ViolationNumber, viol.FirstNameDriver, viol.LastNameDriver,
			viol.RegistrationNumber, viol.FirstNameOwner, viol.LastNameOwner, viol.MiddleNameOwner, viol.AddressOwner}
		comboIndex = viol.DocumentType
		button = dec.PushButton{
			Text: "Προβολή",
			OnClicked: func() {
				UpdatePublishDate(db, ap, time.Now().String())
				createDoc(out[2].Text(), db)
			},
			Font: dec.Font{PointSize: 12},
		}
	} else {
		button = dec.PushButton{
			Text: "Αποθήκευση",
			OnClicked: func() {
				Insert(db, out[0].Text(), out[1].Text(), out[2].Text(), out[3].Text(), out[4].Text(), out[5].Text(),
					out[6].Text(), out[7].Text(), out[8].Text(), out[9].Text(), violType.CurrentIndex()+1)
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
						"Όνομα οδηγού",
						"Επώνυμο οδηγού",
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
						"firstNameDriver",
						"lastNameDriver",
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
