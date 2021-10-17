package officer

import (
	"database/sql"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
)

var out [3]*walk.TextEdit

func createFields(labels, names []string, officer *Officer) []dec.Widget {
	fields := []dec.Widget{}
	values := []string{"", "", ""}

	if officer != nil {
		values[0] = officer.FirstName
		values[1] = officer.LastName
		values[2] = officer.Rank
	}

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
					Text:     values[i],
					Font:     dec.Font{PointSize: 12},
					AssignTo: &out[i],
				},
			},
		}
		fields = append(fields, field)
	}

	return fields
}

func Init(db *sql.DB, commander int) {
	//var out [3]*walk.TextEdit
	title := "Διοικητής"
	if commander == 0 {
		title = "Αρμόδιος"
	}
	officer := GetByCommander(db, commander)
	dec.MainWindow{
		Title:  title,
		Bounds: dec.Rectangle{Width: 800, Height: 200},
		Layout: dec.VBox{},
		Children: []dec.Widget{
			dec.VSplitter{
				Children: createFields(
					[]string{"Όνομα", "Επίθετο", "Βαθμός"},
					[]string{"firstName", "lastName", "rank"},
					officer,
				),
			},
			dec.PushButton{
				Text: "Αποθήκευση",
				OnClicked: func() {
					if officer == nil {
						Insert(db, out[0].Text(), out[1].Text(), out[2].Text(), commander)
						officer = GetByCommander(db, commander)
					} else {
						Update(db, out[0].Text(), out[1].Text(), out[2].Text(), commander)
					}
					//Insert(db, out.Text(), out.Text(), out.Text(), commander)
				},
				Font: dec.Font{PointSize: 12},
			},
		},
	}.Run()
}
