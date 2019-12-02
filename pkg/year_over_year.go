package eoy

import (
	"fmt"
	"time"
)

//YOYear is used to provide a primary key for storing stats by year.
type YOYear struct {
	ID          int
	CreatedDate *time.Time
}

//YOYearResult holds a year and a stats record.
type YOYearResult struct {
	ID int
	Stat
}

//KeyValue implements KeyValuer by returning the value of a key for the
//YOYearResult object.
func (r YOYearResult) KeyValue(i int) (key interface{}) {
	switch i {
	case 0:
		key = r.ID
	default:
		fmt.Printf("Error in YOYearResult\n%+v\n", r)
		err := fmt.Errorf("Not a valid YOYearResult index, %v", i)
		panic(err)
	}
	return key
}

//FillKeys implements KeyFiller by filling Excel cells with keys from the
//year table.
func (r YOYearResult) FillKeys(rt *Runtime, sheet Sheet, row, col int) int {
	for j := 0; j < len(sheet.KeyNames); j++ {
		v := r.KeyValue(j)
		s := sheet.KeyStyles[j]
		rt.Cell(sheet.Name, row, col+j, v, s)
	}
	return row
}

//Fill implements Filler by filling in a spreadsheet using data from the years table.
func (y YOYear) Fill(rt *Runtime, sheet Sheet, row, col int) int {
	var a []YOYearResult
	rt.DB.Order("years.id desc").Table("years").Select("years.id, stats.*").Joins("left join stats on stats.id = years.id").Scan(&a)
	for _, r := range a {
		rt.Spreadsheet.InsertRow(sheet.Name, row+1)
		r.FillKeys(rt, sheet, row, 0)
		r.Stat.Fill(rt, sheet.Name, row, len(sheet.KeyNames))
		row++
	}
	return row
}

//NewYOYearSheet builds the data used to decorate the "this year" page.
func (rt *Runtime) NewYOYearSheet() Sheet {
	filler := YOYear{}
	result := YOYearResult{}
	sheet := Sheet{
		Titles: []string{
			"Year over Year results",
			"Provided by the Custom Success group At Salsalabs",
		},
		Name:      "Year over year",
		KeyNames:  []string{"Year"},
		KeyStyles: []int{rt.KeyStyle},
		Filler:    filler,
		KeyFiller: result,
	}
	return sheet
}
