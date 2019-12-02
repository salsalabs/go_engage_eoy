package eoy

import (
	"fmt"
	"time"
)

//Month is used to provide a primary key for storing stats by month.
type Month struct {
	//ID is YYYY-MM
	ID          string
	Month       int
	Year        int
	CreatedDate *time.Time
}

//MonthResult holds a month and a stats record.
type MonthResult struct {
	ID    string
	Month int
	Year  int
	Stat
}

//MOMonth is used to provide a primary key for storing stats by month.
type MOMonth struct {
	Month
}

//MOMonthResult holds a month and a stats record for the "month over month" sheet.
type MOMonthResult struct {
	MonthResult
}

//KeyValue implements KeyValuer by returning the value of a key for the
//MonthResult object.
func (r MonthResult) KeyValue(i int) (key interface{}) {
	switch i {
	case 0:
		key = r.Month
	case 1:
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

//FillKeys implements KeyFiller by filling Excel cells with keys from the
//year table.
func (r MOMonthResult) FillKeys(rt *Runtime, sheet Sheet, row, col int) int {
	m := MonthResult{}
	return m.FillKeys(rt, sheet, row, col)
}

//Fill implements Filler by filling in a spreadsheet using data from the years table.
func (r Month) Fill(rt *Runtime, sheet Sheet, row, col int) int {
	var a []MonthResult
	y := Year{}
	year := y.Largest(rt)
	rt.DB.Order("id, year desc").Where("months.year = ?", year).Table("months").Select("month, year, stats.*").Joins("left join stats on stats.id = months.id").Scan(&a)
	for _, r := range a {
		rt.Spreadsheet.InsertRow(sheet.Name, row+1)
		r.FillKeys(rt, sheet, row, 0)
		r.Stat.Fill(rt, sheet.Name, row, len(sheet.KeyNames))
		row++
	}
	return row
}

//Fill implements Filler by filling in a spreadsheet using data from the years table.
func (r MOMonth) Fill(rt *Runtime, sheet Sheet, row, col int) int {
	var a []MonthResult
	rt.DB.Order("id, year desc").Table("months").Select("month, year, stats.*").Joins("left join stats on stats.id = months.id").Scan(&a)
	for _, r := range a {
		rt.Spreadsheet.InsertRow(sheet.Name, row+1)
		r.FillKeys(rt, sheet, row, 0)
		r.Stat.Fill(rt, sheet.Name, row, len(sheet.KeyNames))
		row++
	}
	return row
}

//NewMonthSheet builds the data used to decorate the "month over month" sheet.
func (rt *Runtime) NewMonthSheet() Sheet {
	filler := Month{}
	result := MonthResult{}
	y := Year{}
	year := y.Largest(rt)
	name := fmt.Sprintf("Months, %d", year)
	sheet := Sheet{
		Titles: []string{
			fmt.Sprintf("Results by month for %d", year),
			"Provided by the Custom Success group At Salsalabs",
		},
		Name:      name,
		KeyNames:  []string{"Month", "Year"},
		KeyStyles: []int{rt.KeyStyle, rt.KeyStyle},
		Filler:    filler,
		KeyFiller: result,
	}
	return sheet
}

//NewMOMonthSheet builds the data used to decorate the "month over month" sheet.
func (rt *Runtime) NewMOMonthSheet() Sheet {
	filler := MOMonth{}
	result := MOMonthResult{}
	sheet := Sheet{
		Titles: []string{
			"Month over Month results",
			"Provided by the Custom Success group At Salsalabs",
		},
		Name:      "Month over month",
		KeyNames:  []string{"Month", "Year"},
		KeyStyles: []int{rt.KeyStyle, rt.KeyStyle},
		Filler:    filler,
		KeyFiller: result,
	}
	return sheet
}
