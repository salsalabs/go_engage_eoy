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

//YOYear is used to provide a primary key for storing stats by year.
type YOYear struct {
	Year
}

//YOYearResult holds a year and a stats record.
type YOYearResult struct {
	YearResult
}

//KeyValue implements KeyValuer by returning the value of a key for the
//YearResult object.
func (r YearResult) KeyValue(i int) (key interface{}) {
	switch i {
	case 0:
		key = r.ID
	default:
		fmt.Printf("Error in YearResult\n%+v\n", r)
		err := fmt.Errorf("Not a valid YearResult index, %v", i)
		panic(err)
	}
	return key
}

//FillKeys implements KeyFiller by filling Excel cells with keys from the
//year table.
func (r YearResult) FillKeys(rt *Runtime, sheet Sheet, row, col int) int {
	for j := 0; j < len(sheet.KeyNames); j++ {
		v := r.KeyValue(j)
		s := sheet.KeyStyles[j]
		rt.Cell(sheet.Name, row, col+j, v, s)
	}
	return row
}

//Fill implements Filler by filling in a spreadsheet using data from the years table.
func (y Year) Fill(rt *Runtime, sheet Sheet, row, col int) int {
	var a []YearResult
	rt.DB.Table("years").Select("max(years.id), stats.*").Joins("left join stats on stats.id = years.id").Scan(&a)
	for _, r := range a {
		rt.Spreadsheet.InsertRow(sheet.Name, row+1)
		r.FillKeys(rt, sheet, row, 0)
		r.Stat.Fill(rt, sheet.Name, row, len(sheet.KeyNames))
		row++
	}
	return row
}

//NewThisYearSheet builds the data used to decorate the "this year" page.
func (rt *Runtime) NewThisYearSheet() Sheet {
	filler := Year{}
	result := YearResult{}
	y := Year{}
	year := y.Largest(rt)
	name := fmt.Sprintf("%v Summary", year)
	sheet := Sheet{
		Titles: []string{
			fmt.Sprintf("Results for %v", year),
			"Provided by the Custom Success group At Salsalabs",
		},
		Name:      name,
		KeyNames:  []string{"Year"},
		KeyStyles: []int{rt.KeyStyle},
		Filler:    filler,
		KeyFiller: result,
	}
	return sheet
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

//Largest returns the most recent year in the years database.
func (y Year) Largest(rt *Runtime) int {
	var x int
	row := rt.DB.Table("years").Select("MAX(id)").Row()
	row.Scan(&x)
	return x
}
