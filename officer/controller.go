package officer

import (
	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
)

func createFields(labels, names []string, out []*walk.TextEdit) []dec.Widget {
	fields := []dec.Widget{}

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
				},
			},
		}
		fields = append(fields, field)
	}

	return fields
}

func Init(commander bool) {
	out := []*walk.TextEdit{}
	title := "Διοικητής"
	if !commander {
		title = "Αρμόδιος"
	}
	dec.MainWindow{
		Title:  title,
		Bounds: dec.Rectangle{Width: 800, Height: 600},
		Layout: dec.VBox{},
		Children: []dec.Widget{
			dec.HSplitter{
				Children: createFields([]string{"Όνομα", "Επίθετο", "Βαθμός"}, []string{"firstName", "lastName", "rang"}, out),
			},
			dec.PushButton{
				Text: "Αποθήκευση",
				OnClicked: func() {
				},
				Font: dec.Font{PointSize: 12},
			},
		},
	}.Run()
}
