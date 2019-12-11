package eoy

import (
	"fmt"
	"time"
)

//ActivityForm is used to provide a primary key for storing stats by month.
type ActivityForm struct {
	ID          string
	Name        string
	CreatedDate *time.Time
}

//ActivityFormResult holds a month and a stats record.
type ActivityFormResult struct {
	Name string
	Stat
}

//KeyValue implements KeyValuer by returning the value of a key for the
//ActivityFormResult object.
func (r ActivityFormResult) KeyValue(i int) (key interface{}) {
	switch i {
	case 0:
		key = r.Name
	default:
		fmt.Printf("Error in ActivityFormResult\n%+v\n", r)
		err := fmt.Errorf("Not a valid ActivityFormResult index, %v", i)
		panic(err)
	}
	return key
}

//FillKeys implements KeyFiller by filling Excel cells with keys from the
//year table.
func (r ActivityFormResult) FillKeys(rt *Runtime, sheet Sheet, row, col int) int {
	for j := 0; j < len(sheet.KeyNames); j++ {
		v := r.KeyValue(j)
		s := sheet.KeyStyles[j]
		rt.Cell(sheet.Name, row, col+j, v, s)
	}
	return row
}

//Fill implements Filler by filling in a spreadsheet using data from the years table.
func (r ActivityForm) Fill(rt *Runtime, sheet Sheet, row, col int) int {
	var a []ActivityFormResult
	rt.DB.Order("stats.all_amount desc").Table("activity_forms").Select("name, stats.*").Joins("left join stats on stats.id = activity_forms.id").Scan(&a)
	for _, r := range a {
		rt.Spreadsheet.InsertRow(sheet.Name, row+1)
		r.FillKeys(rt, sheet, row, 0)
		r.Stat.Fill(rt, sheet.Name, row, len(sheet.KeyNames))
		row++
	}
	return row
}

//NewActivityFormSheet builds the data used to decorate the "month over month" sheet.
func (rt *Runtime) NewActivityFormSheet() Sheet {
	filler := ActivityForm{}
	result := ActivityFormResult{}
	sheet := Sheet{
		Titles: []string{
			"Results by Activity Form",
		},
		Name:      "Activity forms",
		KeyNames:  []string{"ActivityForm"},
		KeyStyles: []int{rt.KeyStyle},
		Filler:    filler,
		KeyFiller: result,
	}
	return sheet
}
