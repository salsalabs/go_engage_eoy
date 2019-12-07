package eoy

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/jinzhu/gorm"
	goengage "github.com/salsalabs/goengage/pkg"
)

//Runtime contains the variables that we need to run this application.
type Runtime struct {
	Env           *goengage.Environment
	DB            *gorm.DB
	Log           *log.Logger
	Channels      []chan goengage.Fundraise
	Spreadsheet   *excelize.File
	CountStyle    int
	ValueStyle    int
	KeyStyle      int
	TitleStyle    int
	HeaderStyle   int
	TopDonorLimit int
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

//style contains style declarations, ready to marshall into JSON.
type style struct {
	NumberFormat int  `json:"number_format,omitempty"`
	Font         font `json:"font"`
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
//Note: Excel location is limited to the range of columns for this app!
func Axis(r, c int) string {
	cols := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	s := string(cols[c])
	return fmt.Sprintf("%v%v", s, r+1)
}

//Decorate a sheet by putting it into the spreadsheet as an Excel sheet.
func (rt *Runtime) Decorate(sheet Sheet) (row int) {
	_ = rt.Spreadsheet.NewSheet(sheet.Name)
	//Titles go on separate lines
	for row, t := range sheet.Titles {
		rt.Spreadsheet.InsertRow(sheet.Name, row+1)
		rt.Cell(sheet.Name, row, 0, t, rt.HeaderStyle)
	}
	row = len(sheet.Titles)
	rt.Spreadsheet.InsertRow(sheet.Name, row+1)
	//Key headers are followed by stat headers on a single row.
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
	row = sheet.Filler.Fill(rt, sheet, row, 0)
	return row
}

//NewRuntime creates a runtime object and initializes the rt.
func NewRuntime(e *goengage.Environment, db *gorm.DB, channels []chan goengage.Fundraise) *Runtime {
	w, err := os.Create("eoy.log")
	if err != nil {
		log.Panic(err)
	}

	rt := Runtime{
		Env:         e,
		DB:          db,
		Log:         log.New(w, "EOY: ", log.LstdFlags),
		Channels:    channels,
		Spreadsheet: excelize.NewFile(),
	}

	f := font{Size: 18}
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

	f = font{Size: 21, Bold: True, Color: "darkblue"}
	s = style{
		NumberFormat: 0,
		Font:         f,
	}
	rt.HeaderStyle = rt.StyleInt(s)
	return &rt
}

//Cell stores a value in an Excel cell and sets its style.
func (rt *Runtime) Cell(sheetName string, row, col int, v interface{}, s int) {
	a := Axis(row, col)
	rt.Spreadsheet.SetCellValue(sheetName, a, v)
	rt.Spreadsheet.SetCellStyle(sheetName, a, a, s)
}
