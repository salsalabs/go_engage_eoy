package eoy

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/jinzhu/gorm"
	goengage "github.com/salsalabs/goengage/pkg"
)

//Runtime contains the variables that we need to run this application.
type Runtime struct {
	Env         *goengage.Environment
	DB          *gorm.DB
	Log         *log.Logger
	Channels    []chan goengage.Fundraise
	Spreadsheet *excelize.File
	CountStyle  int
	ValueStyle  int
	KeyStyle    int
	TitleStyle  int
	HeaderStyle int
}

//NewRuntime creates a runtime object and initializes the rt.
func NewRuntime(e *goengage.Environment, db *gorm.DB, channels []chan goengage.Fundraise) *Runtime {
	w, err := os.Create("eoy.log")
	if err != nil {
		log.Panic(err)
	}
	s := excelize.NewFile()
	countStyle, _ := s.NewStyle(`{"number_format": 3}`)
	valueStyle, _ := s.NewStyle(`{"number_format": 3}`)
	keyStyle, _ := s.NewStyle(`{"number_format": 0}`)
	headerStyle, _ := s.NewStyle(`{"number_format": 0}`)

	rt := Runtime{
		Env:         e,
		DB:          db,
		Log:         log.New(w, "EOY: ", log.LstdFlags),
		Channels:    channels,
		Spreadsheet: s,
		CountStyle:  countStyle,
		ValueStyle:  valueStyle,
		KeyStyle:    keyStyle,
		HeaderStyle: headerStyle,
	}

	return &rt
}

//KeyValuer returns a key value for the specified offset.
type KeyValuer interface {
	KeyValue(i int) interface{}
}

//Filler inserts stats objects into an Excel spreadsheet starting at the
//specified zero-based row.
type Filler interface {
	Fill(rt *Runtime, sheet Sheet, row, col int) int
}

//ActivityForm contains a basic set of values for an activity page.
type ActivityForm struct {
	ID          string
	Name        string
	CreatedDate *time.Time
}

//ActivityFormResult holds a month and a stats record.
type ActivityFormResult struct {
	ID   string
	Name int
	Stat
}

//Month is used to provide a primary key for storing stats by month.
type Month struct {
	//ID is YYYY-MM
	ID          string
	Year        int
	Month       int
	CreatedDate *time.Time
}

//MonthResult holds a month and a stats record.
type MonthResult struct {
	ID    string
	Year  int
	Month int
	Stat
}

//Cell stores a value in an Excel cell and sets its style.
func (rt *Runtime) Cell(sheetName string, row, col int, v interface{}, s int) {
	a := Axis(row, col)
	rt.Spreadsheet.SetCellValue(sheetName, a, v)
	rt.Spreadsheet.SetCellStyle(sheetName, a, a, s)
}

//Sheet contains the stuff that we need to create and populate a sheet
//in the EOY spreadsheet.
type Sheet struct {
	Name      string
	Titles    []string
	KeyNames  []string
	KeyStyles []int
	Filler    Filler
}

//Axis accepts zero-based row and column and returns an Excel location.
//Note: Excel location is limited to the range of columns for this app!
func Axis(r, c int) string {
	cols := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	s := string(cols[c])
	return fmt.Sprintf("%v%v", s, r+1)
}

//NewThisYearSheet builds the data used to decorate the "this year" page.
func (rt *Runtime) NewThisYearSheet() Sheet {
	filler := Year{}
	sheet := Sheet{
		Titles: []string{
			"Results for the year",
			"Provided by the Custom Success group At Salsalabs",
		},
		Name:      "This year",
		KeyNames:  []string{"Year"},
		KeyStyles: []int{rt.KeyStyle},
		Filler:    filler,
	}
	return sheet
}

//Decorate a sheet by putting it into the spreadsheet as an Excel sheet.
func (rt *Runtime) Decorate(sheet Sheet) (row int) {
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
	for i := int(ID); i < int(StatFieldCount); i++ {
		col := len(sheet.KeyNames) + i
		h := stat.Header(i)
		s := stat.Style(rt, i)
		rt.Cell(sheet.Name, row, col, h, s)
	}
	row++
	row = sheet.Filler.Fill(rt, sheet, row, 0)
	return row
}
