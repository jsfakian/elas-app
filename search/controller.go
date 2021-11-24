package search

import (
	"database/sql"
	"elasapp/decision"
	"elasapp/violation"
	"sort"
	"strings"

	"github.com/lxn/walk"
	dec "github.com/lxn/walk/declarative"
	log "github.com/sirupsen/logrus"
)

var mydb *sql.DB

type searchViolation struct {
	index      int
	viol       *violation.Violation
	objection  bool
	objection2 bool
	checked    bool
}

type violationArray struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*searchViolation
}

func newViolationArray(db *sql.DB) *violationArray {
	m := new(violationArray)
	m.CreateRows(db)
	return m
}

// Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *violationArray) RowCount() int {
	return len(m.items)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *violationArray) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.index
	case 1:
		return item.viol.AP
	case 2:
		return item.viol.ViolationNumber
	case 3:
		return item.viol.RegistrationNumber
	case 4:
		if item.objection {
			return "Ναι"
		} else {
			return "Όχι"
		}
	case 5:
		if item.objection2 {
			return "Ναι"
		} else {
			return "Όχι"
		}
	}

	log.Info("Col:", col)

	panic("unexpected col")
}

// Called by the TableView to retrieve if a given row is checked.
func (m *violationArray) Checked(row int) bool {
	return m.items[row].checked
}

// Called by the TableView when the user toggled the check box of a given row.
func (m *violationArray) SetChecked(row int, checked bool) error {
	if m.items[row].objection && m.items[row].objection2 {
		m.items[row].checked = true
		return nil
	}
	m.items[row].checked = false
	decision.Init(mydb)
	return nil
}

// Called by the TableView to sort the model.
func (m *violationArray) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}
			return !ls
		}

		switch m.sortColumn {
		case 0:
		case 4:
		case 5:
			return c(a.index < b.index)
		case 1:
			return c(strings.Compare(a.viol.AP, b.viol.AP) == -1)
		case 2:
			return c(strings.Compare(a.viol.ViolationNumber, b.viol.ViolationNumber) == -1)
		case 3:
			return c(strings.Compare(a.viol.RegistrationNumber, b.viol.RegistrationNumber) == -1)
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

func (m *violationArray) CreateRows(db *sql.DB) {
	// Create some random data.
	m.items = []*searchViolation{}

	violations := violation.GetByAll(db)

	for i, viol := range violations {
		sv := new(searchViolation)
		sv.index = i + 1
		sv.viol = viol
		sv.objection = false
		sv.objection2 = false
		sv.checked = false
		m.items = append(m.items, sv)
	}

	// Notify TableView and other interested parties about the reset.
	m.PublishRowsReset()
	m.Sort(m.sortColumn, m.sortOrder)
}

func Init(db *sql.DB) {
	var tv *walk.TableView
	mydb = db

	APBitmap, err := walk.NewBitmap(walk.Size{100, 1})
	if err != nil {
		panic(err)
	}
	defer APBitmap.Dispose()

	canvas, err := walk.NewCanvasFromImage(APBitmap)
	if err != nil {
		panic(err)
	}
	defer APBitmap.Dispose()

	canvas.GradientFillRectangle(walk.RGB(255, 0, 0), walk.RGB(0, 255, 0), walk.Horizontal, walk.Rectangle{0, 0, 100, 1})

	canvas.Dispose()

	model := newViolationArray(db)

	dec.MainWindow{
		Title:  "Αναζήτηση Παράβασης",
		Bounds: dec.Rectangle{Width: 950, Height: 600},
		Layout: dec.VBox{},
		Font:   dec.Font{PointSize: 12},
		Children: []dec.Widget{
			dec.VSplitter{
				Children: []dec.Widget{
					dec.PushButton{
						Text: "Απόφαση για έφεση",
						OnClicked: func() {
						},
					},
					dec.TableView{
						AssignTo:         &tv,
						AlternatingRowBG: true,
						CheckBoxes:       true,
						ColumnsOrderable: true,
						MultiSelection:   true,
						Columns: []dec.TableViewColumn{
							{Title: "#"},
							{Title: "Πρωτόκολλο", Width: 160},
							{Title: "Παράβαση", Width: 130},
							{Title: "Αριθμός κυκλοφορίας", Width: 170},
							{Title: "Απόφαση Ένστασης", Width: 160},
							{Title: "Απόφαση Ένστασης (2)", Width: 180},
						},
						StyleCell: func(style *walk.CellStyle) {
							item := model.items[style.Row()]

							if item.checked {
								if style.Row()%2 == 0 {
									style.BackgroundColor = walk.RGB(159, 215, 255)
								} else {
									style.BackgroundColor = walk.RGB(143, 199, 239)
								}
							}
						},
						Model: model,
						OnItemActivated: func() {
							violation.Init(db, model.items[tv.CurrentIndex()].viol.AP)
						},
					},
				},
			},
		},
	}.Run()
}
