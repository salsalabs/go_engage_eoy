package eoy

import (
	"fmt"
	"time"
)

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

//KeyValue implements KeyValuer by returning the value of a key for the
//MonthResult object.
func (r MonthResult) KeyValue(i int) (key interface{}) {
	switch i {
	case 0:
		key = r.ID
	case 1:
		key = r.Month
	case 2:
		key = r.Year
	default:
		fmt.Printf("Error in MonthResult\n%+v\n", r)
		err := fmt.Errorf("Not a valid MonthResult index, %v", i)
		panic(err)
	}
	return key
}

//FillKeys implements KeyFiller by filling Excel cells with keys from the
//year table.
func (r MonthResult) FillKeys(rt *Runtime, sheet Sheet, row, col int) int {
	for j := 0; j < len(sheet.KeyNames); j++ {
		v := r.KeyValue(j)
		s := sheet.KeyStyles[j]
		rt.Cell(sheet.Name, row, col+j, v, s)
	}
	return row
}

//Fill implements Filler by filling in a spreadsheet using data from the years table.
func (y Month) Fill(rt *Runtime, sheet Sheet, row, col int) int {
	var a []MonthResult
	rt.DB.Table("months").Select("month, year, stats.*").Joins("left join stats on stats.id = years.id").Scan(&a)
	for _, r := range a {
		rt.Spreadsheet.InsertRow(sheet.Name, row+1)
		r.FillKeys(rt, sheet, row, 0)
		r.Stat.Fill(rt, sheet.Name, row, len(sheet.KeyNames))
		row++
	}
	return row
}
