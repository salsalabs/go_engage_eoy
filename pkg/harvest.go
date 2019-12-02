package eoy

//harvester declares functions that process data.
type harvester func(rt *Runtime) (err error)

//Harvest retrieves data from the database in various permutations of slicing
//and dicing, then stores them into a spreadsheet.  The spreadsheet is written
//to disk when done.
func (rt *Runtime) Harvest(fn string) (err error) {
	functions := []harvester{
		ThisYear,
		// Months,
		YearOverYear,
		// MonthOverMonth,
		// AllDonors,
		// TopDonors,
		// ActivityPages,
		// ProjectedRevenue,
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

//NewThisYearSheet builds the data used to decorate the "this year" page.
func (rt *Runtime) NewThisYearSheet() Sheet {
	filler := Year{}
	result := YearResult{}
	sheet := Sheet{
		Titles: []string{
			"Results for the year",
			"Provided by the Custom Success group At Salsalabs",
		},
		Name:      "This year",
		KeyNames:  []string{"Year"},
		KeyStyles: []int{rt.KeyStyle},
		Filler:    filler,
		KeyFiller: result,
	}
	return sheet
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

// ThisYear selects data for ThisYear, sorts it, tweaks it, then stores it into
//the spreadsheet.
func ThisYear(rt *Runtime) (err error) {
	sheet := rt.NewThisYearSheet()
	rt.Decorate(sheet)
	return err
}

// YearOverYear selects data for ThisYear, sorts it, tweaks it, then stores it into
//the spreadsheet.
func YearOverYear(rt *Runtime) (err error) {
	sheet := rt.NewYOYearSheet()
	rt.Decorate(sheet)
	return err
}

// // Months selects data for Months, sorts it, tweaks it, then stores it into
// //the spreadsheet.
// func Months(rt *Runtime) (err error) {
// 	sheet := "Month this year"
// 	_ = rt.Spreadsheet.NewSheet(sheet)
//
// 	return err
// }
//
// // YearOverYear selects data for YearOverYear, sorts it, tweaks it, then stores it into
// //the spreadsheet.
// func YearOverYear(rt *Runtime) (err error) {
// 	sheet := "Year-over-year"
// 	_ = rt.Spreadsheet.NewSheet(sheet)
// 	var a []YearResult
// 	rt.DB.Table("years").Select("years.id, stats.*").Joins("left join stats on stats.id = years.id").Order("years.id desc").Scan(&a)
// 	h := "Year over year performance"
// 	//Sheet header
// 	rt.Spreadsheet.InsertRow(sheet, 1)
// 	rt.Spreadsheet.SetCellValue(sheet, "A1", h)
// 	//Column headers
// 	hc := strings.Split(statsHeaders, "\n")
// 	header := []string{}
// 	header = append(header, "Year")
// 	for _, t := range hc {
// 		header = append(header, t)
// 	}
// 	rt.Spreadsheet.InsertRow(sheet, 2)
// 	err = rt.Spreadsheet.SetSheetRow(sheet, "A2", &header)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for rowID, r := range a {
// 		rt.Spreadsheet.InsertRow(sheet, rowID+3)
// 		axis := fmt.Sprintf("A%d", rowID+3)
// 		rt.Spreadsheet.SetCellValue(sheet, axis, r.ID)
// 		g := r.Stat
// 		cols := "BCDEFGHIJKLMNOPQ"
// 		for i := range hc {
// 			c := cellContent(rt, g, i, cols)
// 			axis = fmt.Sprintf("%v%d", c.Column, rowID+3)
// 			rt.Spreadsheet.SetCellValue(sheet, axis, c.Value)
// 			rt.Spreadsheet.SetCellStyle(sheet, axis, axis, c.Style)
// 		}
// 	}
// 	return err
// }
//
// // MonthOverMonth selects data for MonthOverMonth, sorts it, tweaks it, then stores it into
// //the spreadsheet.
// func MonthOverMonth(rt *Runtime) (err error) {
// 	sheet := "Month-over-month"
// 	_ = rt.Spreadsheet.NewSheet(sheet)
//
// 	_ = rt.Spreadsheet.NewSheet(sheet)
// 	var a []MonthResult
// 	rt.DB.Table("months").Select("month, year, stats.*").Joins("left join stats on stats.id = months.id").Order("month,year").Scan(&a)
// 	h := "Month over month performance"
// 	//Sheet header
// 	rt.Spreadsheet.InsertRow(sheet, 1)
// 	rt.Spreadsheet.SetCellValue(sheet, "A1", h)
// 	//Column headers
// 	hc := strings.Split(statsHeaders, "\n")
// 	header := []string{}
// 	header = append(header, "Month")
// 	header = append(header, "Year")
// 	for _, t := range hc {
// 		header = append(header, t)
// 	}
// 	rt.Spreadsheet.InsertRow(sheet, 2)
// 	err = rt.Spreadsheet.SetSheetRow(sheet, "A2", &header)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for rowID, r := range a {
// 		rt.Spreadsheet.InsertRow(sheet, rowID+3)
// 		axis := fmt.Sprintf("A%d", rowID+3)
// 		rt.Spreadsheet.SetCellValue(sheet, axis, r.Month)
// 		axis = fmt.Sprintf("B%d", rowID+3)
// 		rt.Spreadsheet.SetCellValue(sheet, axis, r.Year)
// 		g := r.Stat
// 		cols := "CDEFGHIJKLMNOPQR"
// 		for i := range hc {
// 			c := cellContent(rt, g, i, cols)
// 			axis = fmt.Sprintf("%v%d", c.Column, rowID+3)
// 			rt.Spreadsheet.SetCellValue(sheet, axis, c.Value)
// 			rt.Spreadsheet.SetCellStyle(sheet, axis, axis, c.Style)
// 		}
// 	}
//
// 	return err
// }
//
// // AllDonors selects data for AllDonors, sorts it, tweaks it, then stores it into
// //the spreadsheet.
// func AllDonors(rt *Runtime) (err error) {
// 	sheet := "All donors"
// 	_ = rt.Spreadsheet.NewSheet(sheet)
//
// 	return err
// }
//
// // TopDonors selects data for TopDonors, sorts it, tweaks it, then stores it into
// //the spreadsheet.
// func TopDonors(rt *Runtime) (err error) {
// 	sheet := "Top donors"
// 	_ = rt.Spreadsheet.NewSheet(sheet)
//
// 	return err
// }
//
// // ActivityPages selects data for ActivityPages, sorts it, tweaks it, then stores it into
// //the spreadsheet.
// func ActivityPages(rt *Runtime) (err error) {
// 	sheet := "Activity pages"
// 	_ = rt.Spreadsheet.NewSheet(sheet)
//
// 	return err
// }
//
// // ProjectedRevenue selects data for ProjectedRevenue, sorts it, tweaks it, then stores it into
// //the spreadsheet.
// func ProjectedRevenue(rt *Runtime) (err error) {
// 	sheet := "Projected revenue"
// 	_ = rt.Spreadsheet.NewSheet(sheet)
//
// 	return err
// }

//StoreSpreadsheet saves the spreadsheet to disk.
func (rt *Runtime) StoreSpreadsheet(fn string) (err error) {
	err = rt.Spreadsheet.SaveAs(fn)
	return err
}
