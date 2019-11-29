package eoy

import (
	"fmt"
	"log"
	"strings"
)

//statsHeaders lists the stats fields.  Split by a newline to get a list.
const statsHeaders = `All Count
All Amount
OneTime Count
OneTime Amount
Recurring Count
Recurring Amount
Offline Count
Offline Amount
Refunds Count
Refunds Amount
Largest
Smallest`

//harvester declares functions that process data.
type harvester func(rt *Runtime) (err error)

//yearResult holds a year and a giving_stats record.
type yearResult struct {
	ID int
	GivingStat
}

//month result holds a month and a giving_stats record.
type monthResult struct {
	ID    string
	Year  int
	Month int
	GivingStat
}

//content describes the content of a cell.
type content struct {
	Value  interface{}
	Style  int
	Header string
	Column string
}

//cellContent returns a content for one of the stats values.  The
//index value is for the list of headers.
func cellContent(rt *Runtime, g GivingStat, i int, cols string) content {
	hc := strings.Split(statsHeaders, "\n")
	countStyle, _ := rt.Spreadsheet.NewStyle(`{"number_format": 3}`)
	valueStyle, _ := rt.Spreadsheet.NewStyle(`{"number_format": 3}`)
	var v interface{}
	var style int
	switch i {
	case 0:
		v = g.AllCount
		style = countStyle
	case 1:
		v = g.AllAmount
		style = valueStyle
	case 2:
		v = g.OneTimeCount
		style = countStyle
	case 3:
		v = g.OneTimeAmount
		style = valueStyle
	case 4:
		v = g.RecurringCount
		style = countStyle
	case 5:
		v = g.RecurringAmount
		style = valueStyle
	case 6:
		v = g.OfflineCount
		style = countStyle
	case 7:
		v = g.OfflineAmount
		style = valueStyle
	case 8:
		v = g.RefundsCount
		style = countStyle
	case 9:
		v = g.RefundsAmount
		style = valueStyle
	case 10:
		v = g.Largest
		style = valueStyle
	case 11:
		v = g.Smallest
		style = valueStyle
	}
	col := string(cols[i])
	c := content{
		Value:  v,
		Style:  style,
		Header: hc[i],
		Column: col,
	}
	return c
}

//Harvest retrieves data from the database in various permutations of slicing
//and dicing, then stores them into a spreadsheet.  The spreadsheet is written
//to disk when done.
func (rt *Runtime) Harvest(fn string) (err error) {
	functions := []harvester{
		ThisYear,
		Months,
		YearOverYear,
		MonthOverMonth,
		AllDonors,
		TopDonors,
		ActivityPages,
		ProjectedRevenue,
	}

	for _, r := range functions {
		err := r(rt)
		if err != nil {
			return err
		}
	}
	err = rt.StoreSpreadsheet(fn)
	return err
}

// ThisYear selects data for ThisYear, sorts it, tweaks it, then stores it into
//the spreadsheet.
func ThisYear(rt *Runtime) (err error) {
	sheet := "This year"
	_ = rt.Spreadsheet.NewSheet(sheet)

	var a []yearResult
	rt.DB.Table("years").Select("max(years.id), giving_stats.*").Joins("left join giving_stats on giving_stats.id = years.id").Scan(&a)
	y := a[0].ID
	header := []string{
		fmt.Sprintf("Performance summary for %v", y),
	}
	rt.Spreadsheet.InsertRow(sheet, 1)
	err = rt.Spreadsheet.SetSheetRow(sheet, "A1", &header)
	if err != nil {
		log.Fatal(err)
	}

	h := strings.Split(statsHeaders, "\n")
	g := a[0].GivingStat
	for i := range h {
		c := cellContent(rt, g, i, "BCDEFGHIJKLMNOPQ")
		rt.Spreadsheet.InsertRow(sheet, i+2)
		axis := fmt.Sprintf("A%d", i+2)
		rt.Spreadsheet.SetCellValue(sheet, axis, c.Header)
		axis = fmt.Sprintf("B%d", i+2)
		rt.Spreadsheet.SetCellValue(sheet, axis, c.Value)
		rt.Spreadsheet.SetCellStyle(sheet, axis, axis, c.Style)
	}
	return err
}

// Months selects data for Months, sorts it, tweaks it, then stores it into
//the spreadsheet.
func Months(rt *Runtime) (err error) {
	sheet := "Month this year"
	_ = rt.Spreadsheet.NewSheet(sheet)

	return err
}

// YearOverYear selects data for YearOverYear, sorts it, tweaks it, then stores it into
//the spreadsheet.
func YearOverYear(rt *Runtime) (err error) {
	sheet := "Year-over-year"
	_ = rt.Spreadsheet.NewSheet(sheet)
	var a []yearResult
	rt.DB.Table("years").Select("years.id, giving_stats.*").Joins("left join giving_stats on giving_stats.id = years.id").Order("years.id desc").Scan(&a)
	h := "Year over year performance"
	//Sheet header
	rt.Spreadsheet.InsertRow(sheet, 1)
	rt.Spreadsheet.SetCellValue(sheet, "A1", h)
	//Column headers
	hc := strings.Split(statsHeaders, "\n")
	header := []string{}
	header = append(header, "Year")
	for _, t := range hc {
		header = append(header, t)
	}
	rt.Spreadsheet.InsertRow(sheet, 2)
	err = rt.Spreadsheet.SetSheetRow(sheet, "A2", &header)
	if err != nil {
		log.Fatal(err)
	}
	for rowID, r := range a {
		rt.Spreadsheet.InsertRow(sheet, rowID+3)
		axis := fmt.Sprintf("A%d", rowID+3)
		rt.Spreadsheet.SetCellValue(sheet, axis, r.ID)
		g := r.GivingStat
		cols := "BCDEFGHIJKLMNOPQ"
		for i := range hc {
			c := cellContent(rt, g, i, cols)
			axis = fmt.Sprintf("%v%d", c.Column, rowID+3)
			rt.Spreadsheet.SetCellValue(sheet, axis, c.Value)
			rt.Spreadsheet.SetCellStyle(sheet, axis, axis, c.Style)
		}
	}
	return err
}

// MonthOverMonth selects data for MonthOverMonth, sorts it, tweaks it, then stores it into
//the spreadsheet.
func MonthOverMonth(rt *Runtime) (err error) {
	sheet := "Month-over-month"
	_ = rt.Spreadsheet.NewSheet(sheet)

	_ = rt.Spreadsheet.NewSheet(sheet)
	var a []monthResult
	rt.DB.Table("months").Select("month, year, giving_stats.*").Joins("left join giving_stats on giving_stats.id = months.id").Order("month,year").Scan(&a)
	h := "Month over month performance"
	//Sheet header
	rt.Spreadsheet.InsertRow(sheet, 1)
	rt.Spreadsheet.SetCellValue(sheet, "A1", h)
	//Column headers
	hc := strings.Split(statsHeaders, "\n")
	header := []string{}
	header = append(header, "Month")
	header = append(header, "Year")
	for _, t := range hc {
		header = append(header, t)
	}
	rt.Spreadsheet.InsertRow(sheet, 2)
	err = rt.Spreadsheet.SetSheetRow(sheet, "A2", &header)
	if err != nil {
		log.Fatal(err)
	}
	for rowID, r := range a {
		rt.Spreadsheet.InsertRow(sheet, rowID+3)
		axis := fmt.Sprintf("A%d", rowID+3)
		rt.Spreadsheet.SetCellValue(sheet, axis, r.Month)
		axis = fmt.Sprintf("B%d", rowID+3)
		rt.Spreadsheet.SetCellValue(sheet, axis, r.Year)
		g := r.GivingStat
		cols := "CDEFGHIJKLMNOPQR"
		for i := range hc {
			c := cellContent(rt, g, i, cols)
			axis = fmt.Sprintf("%v%d", c.Column, rowID+3)
			rt.Spreadsheet.SetCellValue(sheet, axis, c.Value)
			rt.Spreadsheet.SetCellStyle(sheet, axis, axis, c.Style)
		}
	}

	return err
}

// AllDonors selects data for AllDonors, sorts it, tweaks it, then stores it into
//the spreadsheet.
func AllDonors(rt *Runtime) (err error) {
	sheet := "All donors"
	_ = rt.Spreadsheet.NewSheet(sheet)

	return err
}

// TopDonors selects data for TopDonors, sorts it, tweaks it, then stores it into
//the spreadsheet.
func TopDonors(rt *Runtime) (err error) {
	sheet := "Top donors"
	_ = rt.Spreadsheet.NewSheet(sheet)

	return err
}

// ActivityPages selects data for ActivityPages, sorts it, tweaks it, then stores it into
//the spreadsheet.
func ActivityPages(rt *Runtime) (err error) {
	sheet := "Activity pages"
	_ = rt.Spreadsheet.NewSheet(sheet)

	return err
}

// ProjectedRevenue selects data for ProjectedRevenue, sorts it, tweaks it, then stores it into
//the spreadsheet.
func ProjectedRevenue(rt *Runtime) (err error) {
	sheet := "Projected revenue"
	_ = rt.Spreadsheet.NewSheet(sheet)

	return err
}

//StoreSpreadsheet saves the spreadsheet to disk.
func (rt *Runtime) StoreSpreadsheet(fn string) (err error) {
	err = rt.Spreadsheet.SaveAs(fn)
	return err
}
