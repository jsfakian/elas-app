package violation

import (
	dec "github.com/lxn/walk/declarative"
)

func createFields(labels, names []string) []dec.Widget {
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
					Name: names[i] + "_input",
					Font: dec.Font{PointSize: 12},
				},
			},
		}
		fields = append(fields, field)
	}

	return fields
}

func Init() {
	dec.MainWindow{
		Title:  "Καταχώρηση Παράβασης",
		Bounds: dec.Rectangle{Width: 800, Height: 600},
		Layout: dec.VBox{},
		Children: []dec.Widget{
			dec.HSplitter{
				Children: []dec.Widget{
					dec.TextLabel{
						Text: "Όνομα",
						Font: dec.Font{PointSize: 12},
					},
					dec.TextEdit{
						CompactHeight: true,
						Font:          dec.Font{PointSize: 12},
					},
				},
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
