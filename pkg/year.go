package eoy

import (
	"fmt"
	"time"
)

//Year is used to provide a primary key for storing stats by year.
type Year struct {
	ID          int
	CreatedDate *time.Time
}

//YearResult holds a year and a stats record.
type YearResult struct {
	ID int
	Stat
}

//KeyValue implements KeyValuer by returning the value of a key for the
//YearResult object.
func (r YearResult) KeyValue(i int) (key interface{}) {
	switch i {
	case 0:
		key = r.ID
	default:
		err := fmt.Errorf("Not a valid YearResult index, %v", i)
		panic(err)
	}
	return key
}

//FillKeys implements KeyFiller by filling Excel cells with keys from the
//year table.
func (r YearResult) FillKeys(rt *Runtime, sheet Sheet, row, col int) int {
	for j := 0; j <= len(sheet.KeyNames); j++ {
		v := r.KeyValue(j)
		s := sheet.KeyStyles[j]
		rt.Cell(sheet.Name, row, 0, v, s)
	}
	return row
}

//Fill implements Filler by filling in a spreadsheet using data from the years table.
func (y Year) Fill(rt *Runtime, sheet Sheet, row, col int) int {
	var a []YearResult
	rt.DB.Table("years").Select("max(years.id), stats.*").Joins("left join stats on stats.id = years.id").Scan(&a)
	for _, r := range a {
		rt.Spreadsheet.InsertRow(sheet.Name, row+1)
		r.FillKeys(rt, sheet, row, len(sheet.KeyNames))
		r.Stat.Fill(rt, sheet.Name, row, len(sheet.KeyNames))
		row++
	}
	return row
}
