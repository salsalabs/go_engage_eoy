package eoy

import (
	"fmt"
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

//count converts an int to a string.
func count(v int32) string {
	return fmt.Sprintf("%d", v)
}

//amount converts a float to a string.
func amount(v float64) string {
	return fmt.Sprintf("%.2f", v)
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
	name := "This year"
	_ = rt.Spreadsheet.NewSheet(name)

	var a []yearResult
	rt.DB.Table("years").Select("max(years.id), giving_stats.*").Joins("left join giving_stats on giving_stats.id = years.id").Scan(&a)
	y := a[0].ID
	header := []string{
		fmt.Sprintf("Performance summary for %v", y),
	}
	rt.Spreadsheet.InsertRow(name, 1)
	rt.Spreadsheet.SetSheetRow(name, "A1", header)

	h := strings.Split(statsHeaders, "\n")
	g := a[0].GivingStat
	for i, t := range h {
		r := []string{t}
		var v string
		switch i {
		case 0:
			v = count(g.AllCount)
		case 1:
			v = amount(g.AllAmount)
		case 2:
			v = count(g.OneTimeCount)
		case 3:
			v = amount(g.OneTimeAmount)
		case 4:
			v = count(g.RecurringCount)
		case 5:
			v = amount(g.RecurringAmount)
		case 6:
			v = count(g.OfflineCount)
		case 7:
			v = amount(g.OfflineAmount)
		case 8:
			v = count(g.RefundsCount)
		case 9:
			v = amount(g.RefundsAmount)
		case 10:
			v = amount(g.Largest)
		case 11:
			v = amount(g.Smallest)
		}
		r = append(r, v)
		axis := fmt.Sprintf("A%d", i+2)
		rt.Spreadsheet.InsertRow(name, i+2)
		rt.Spreadsheet.SetSheetRow(name, axis, &r)
	}
	return err
}

// Months selects data for Months, sorts it, tweaks it, then stores it into
//the spreadsheet.
func Months(rt *Runtime) (err error) {
	name := "Month this year"
	_ = rt.Spreadsheet.NewSheet(name)

	return err
}

// YearOverYear selects data for YearOverYear, sorts it, tweaks it, then stores it into
//the spreadsheet.
func YearOverYear(rt *Runtime) (err error) {
	name := "Year-over-year"
	_ = rt.Spreadsheet.NewSheet(name)
	var a []yearResult
	rt.DB.Table("years").Select("years.id, giving_stats.*").Joins("left join giving_stats on giving_stats.id = years.id").Scan(&a)
	header := []string{
		"Year over year performance",
	}
	//Sheet header
	rt.Spreadsheet.InsertRow(name, 1)
	rt.Spreadsheet.SetSheetRow(name, "A1", header)
	//Column headers
	h := strings.Split(statsHeaders, "\n")
	header = []string{}
	header = append(header, "Year")
	for _, t := range h {
		header = append(header, t)
	}
	rt.Spreadsheet.InsertRow(name, 2)
	rt.Spreadsheet.SetSheetRow(name, "A2", header)

	for rowID, r := range a {
		w := []string{
			count(int32(r.ID)),
		}
		g := r.GivingStat
		for i := range h {
			var v string
			switch i {
			case 0:
				v = count(g.AllCount)
			case 1:
				v = amount(g.AllAmount)
			case 2:
				v = count(g.OneTimeCount)
			case 3:
				v = amount(g.OneTimeAmount)
			case 4:
				v = count(g.RecurringCount)
			case 5:
				v = amount(g.RecurringAmount)
			case 6:
				v = count(g.OfflineCount)
			case 7:
				v = amount(g.OfflineAmount)
			case 8:
				v = count(g.RefundsCount)
			case 9:
				v = amount(g.RefundsAmount)
			case 10:
				v = amount(g.Largest)
			case 11:
				v = amount(g.Smallest)
			}
			w = append(w, v)
		}
		axis := fmt.Sprintf("A%d", rowID+3)
		fmt.Printf("name:%v, axis, %v, values: %v\n", name, axis, w)
		rt.Spreadsheet.InsertRow(name, rowID+3)
		rt.Spreadsheet.SetSheetRow(name, axis, &w)
	}
	return err
}

// MonthOverMonth selects data for MonthOverMonth, sorts it, tweaks it, then stores it into
//the spreadsheet.
func MonthOverMonth(rt *Runtime) (err error) {
	name := "Month-over-month"
	_ = rt.Spreadsheet.NewSheet(name)

	return err
}

// AllDonors selects data for AllDonors, sorts it, tweaks it, then stores it into
//the spreadsheet.
func AllDonors(rt *Runtime) (err error) {
	name := "All donors"
	_ = rt.Spreadsheet.NewSheet(name)

	return err
}

// TopDonors selects data for TopDonors, sorts it, tweaks it, then stores it into
//the spreadsheet.
func TopDonors(rt *Runtime) (err error) {
	name := "Top donors"
	_ = rt.Spreadsheet.NewSheet(name)

	return err
}

// ActivityPages selects data for ActivityPages, sorts it, tweaks it, then stores it into
//the spreadsheet.
func ActivityPages(rt *Runtime) (err error) {
	name := "Activity pages"
	_ = rt.Spreadsheet.NewSheet(name)

	return err
}

// ProjectedRevenue selects data for ProjectedRevenue, sorts it, tweaks it, then stores it into
//the spreadsheet.
func ProjectedRevenue(rt *Runtime) (err error) {
	name := "Projected revenue"
	_ = rt.Spreadsheet.NewSheet(name)

	return err
}

//StoreSpreadsheet saves the spreadsheet to disk.
func (rt *Runtime) StoreSpreadsheet(fn string) (err error) {
	err = rt.Spreadsheet.SaveAs(fn)
	return err
}
