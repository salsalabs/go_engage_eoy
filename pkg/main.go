package eoy

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/jinzhu/gorm"
	goengage "github.com/salsalabs/goengage/pkg"
)

//Runtime contains the variables that we need to run this application.
type Runtime struct {
	Env             *goengage.Environment
	DB              *gorm.DB
	Log             *log.Logger
	Channels        []chan goengage.Fundraise
	Spreadsheet     *excelize.File
	CountStyle      int
	ValueStyle      int
	KeyStyle        int
	TitleStyle      int
	HeaderStyle     int
	StatHeaderStyle int
	TopDonorLimit   int
	Year            int
	YearStart       time.Time
	YearEnd         time.Time
	OrgLocation     *time.Location
}

//KeyValuer returns a key value for the specified offset.
type KeyValuer interface {
	KeyValue(i int) interface{}
}

//KeyFiller inserts stats objects into an Excel spreadsheet starting at the
//specified zero-based row.
type KeyFiller interface {
	FillKeys(rt *Runtime, sheet Sheet, row, col int) int
}

//Filler inserts stats objects into an Excel spreadsheet starting at the
//specified zero-based row and column.
type Filler interface {
	Fill(rt *Runtime, sheet Sheet, row, col int) int
}

//Sheet contains the stuff that we need to create and populate a sheet
//in the EOY spreadsheet.
type Sheet struct {
	Name      string
	Titles    []string
	KeyNames  []string
	KeyStyles []int
	Filler    Filler
	KeyFiller KeyFiller
}

//font contains the font definition for a style.
type font struct {
	Bold   bool   `json:"bold,omitempty"`
	Italic bool   `json:"italic,omitempty"`
	Family string `json:"family,omitempty"`
	Size   int    `json:"size,omitempty"`
	Color  string `json:"color,omitempty"`
}

//alignment contains alignment definitions for a style.
type alignment struct {
	Horizontal  string `json:"horizontal,omitempty"`
	ShrinkToFit bool   `json:"shrink_to_fit,omitempty"`
}

//style contains style declarations, ready to marshall into JSON.
type style struct {
	NumberFormat int       `json:"number_format,omitempty"`
	Font         font      `json:"font"`
	Alignment    alignment `json:"alignment"`
}

//StyleInt converts a style for use in an Excelize spreadsheet.
func (rt *Runtime) StyleInt(s style) int {
	b, err := json.Marshal(s)
	if err != nil {
		log.Panic(err)
	}
	x, _ := rt.Spreadsheet.NewStyle(string(b))
	if err != nil {
		log.Panic(err)
	}
	return x
}

//Axis accepts zero-based row and column and returns an Excel location.
//Note: Column offset is limited to the range of columns for this app --
//one alpha digit.
func Axis(r, c int) string {
	cols := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	s := string(cols[c])
	return fmt.Sprintf("%v%v", s, r+1)
}

//Titles inserts titles into a sheet and decorates them.
func (rt *Runtime) Titles(sheet Sheet) (row int) {
	//Titles go on separate lines
	for row, t := range sheet.Titles {
		rt.Spreadsheet.InsertRow(sheet.Name, row+1)
		rt.Cell(sheet.Name, row, 0, t, rt.TitleStyle)
		// Sorry, "2" is a magic number for now...
		w := len(sheet.KeyNames) + int(StatFieldCount) - 2
		left := Axis(row, 0)
		right := Axis(row, w)
		err := rt.Spreadsheet.MergeCell(sheet.Name, left, right)
		if err != nil {
			panic(err)
		}
	}
	row++
	return row
}

//StatHeaders show the topics for stats.  Most of the headers
//will be two columns to cover count and amount.
func (rt *Runtime) StatHeaders(sheet Sheet, row int) int {
	s := Stat{}
	s.Headers(rt, sheet.Name, row, len(sheet.KeyNames))
	row++
	return row
}

//Headers decorates the data headers for a spreadsheet.
func (rt *Runtime) Headers(sheet Sheet, row int) int {
	//Key headers are followed by stat headers on a single row.
	rt.Spreadsheet.InsertRow(sheet.Name, row+1)
	for i, t := range sheet.KeyNames {
		s := sheet.KeyStyles[i]
		rt.Cell(sheet.Name, row, i, t, s)
	}
	stat := Stat{}
	for i := int(AllCount); i < int(StatFieldCount); i++ {
		//"-1" because we are skipping ID
		col := len(sheet.KeyNames) + i - 1
		h := stat.Header(i)
		s := stat.Style(rt, i)
		rt.Cell(sheet.Name, row, col, h, s)
	}
	row++
	return row
}

//Widths sets the widths in a sheet.
func (rt *Runtime) Widths(sheet Sheet, row int) int {
	left := Axis(row, 0)
	left = left[0:1]
	// Sorry, "2" is a magic number for now...
	w := len(sheet.KeyNames) + int(StatFieldCount) - 2
	right := Axis(row, w)
	right = right[0:1]
	err := rt.Spreadsheet.SetColWidth(sheet.Name, left, right, 16.0)
	if err != nil {
		panic(err)
	}
	//Kludge!  Activity form names are a lot longer than supporter names.
	if strings.Contains(sheet.Name, "Activity") {
		w, _ := rt.Spreadsheet.GetColWidth(sheet.Name, "A")
		_ = rt.Spreadsheet.SetColWidth(sheet.Name, "A", "A", w*4.0)
	}
	return row
}

//Decorate a sheet by putting it into the spreadsheet as an Excel sheet.
func (rt *Runtime) Decorate(sheet Sheet) (row int) {
	_ = rt.Spreadsheet.NewSheet(sheet.Name)
	row = rt.Titles(sheet)
	row = rt.StatHeaders(sheet, row)
	row = rt.Headers(sheet, row)
	row = sheet.Filler.Fill(rt, sheet, row, 0)
	row = rt.Widths(sheet, row)
	return row
}

//NewRuntime creates a runtime object and initializes the rt.
func NewRuntime(e *goengage.Environment, db *gorm.DB, channels []chan goengage.Fundraise, year int, topLimit int, loc string) *Runtime {
	w, err := os.Create("eoy.log")
	if err != nil {
		log.Panic(err)
	}

	yearStart := time.Date(year, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year, time.Month(12), 31, 23, 59, 59, 999, time.UTC)
	orgLocation, err := time.LoadLocation(loc)

	rt := Runtime{
		Env:           e,
		DB:            db,
		Log:           log.New(w, "EOY: ", log.LstdFlags),
		Channels:      channels,
		Spreadsheet:   excelize.NewFile(),
		Year:          year,
		TopDonorLimit: topLimit,
		YearStart:     yearStart,
		YearEnd:       yearEnd,
		OrgLocation:   orgLocation,
	}

	f := font{Size: 12}
	s := style{
		NumberFormat: 3,
		Font:         f,
	}
	rt.CountStyle = rt.StyleInt(s)

	s = style{
		NumberFormat: 3,
		Font:         f,
	}
	rt.ValueStyle = rt.StyleInt(s)

	s = style{
		NumberFormat: 0,
		Font:         f,
	}
	rt.KeyStyle = rt.StyleInt(s)
	a := alignment{Horizontal: "center"}
	s = style{
		NumberFormat: 0,
		Font:         f,
		Alignment:    a,
	}
	rt.HeaderStyle = rt.StyleInt(s)
	f2 := font{Size: 16, Bold: true}
	s = style{
		NumberFormat: 0,
		Font:         f2,
		Alignment:    a,
	}
	rt.TitleStyle = rt.StyleInt(s)
	return &rt
}

//Cell stores a value in an Excel cell and sets its style.
func (rt *Runtime) Cell(sheetName string, row, col int, v interface{}, s int) {
	a := Axis(row, col)
	rt.Spreadsheet.SetCellValue(sheetName, a, v)
	rt.Spreadsheet.SetCellStyle(sheetName, a, a, s)
}

//GoodYear returns true if the specified time is between Runtime.YearStart and Runtime.YearEnd.
func (rt *Runtime) GoodYear(t *time.Time) bool {
	if t.Before(rt.YearStart) || t.After(rt.YearEnd) {
		return false
	}
	return true
}
