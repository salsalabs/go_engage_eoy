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

//KeyValue returns the value of a key for the YearResult object.
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

//Fill fills in a spreadsheet using data from the years table.
func (y Year) Fill(rt *Runtime, sheet Sheet, row, col int) int {
	var a []YearResult
	rt.DB.Table("years").Select("max(years.id), stats.*").Joins("left join stats on stats.id = years.id").Scan(&a)
	for i, r := range a {
		if i < len(sheet.KeyNames) {
			v := r.KeyValue(i)
			s := sheet.KeyStyles[i]
			rt.Cell(sheet.Name, row, i, v, s)
		} else {
			j := i - len(sheet.KeyNames)
			v := r.Value(j)
			s := r.Style(rt, j)
			rt.Cell(sheet.Name, row, i, v, s)
		}
		row++
	}
	return row
}
