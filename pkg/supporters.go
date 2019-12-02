package eoy

import (
	"fmt"
	"time"
)

//Donor is used to provide a primary key for storing stats by month.
type Donor struct {
	//ID is YYYY-MM
	SupporterID string `gorm:"supporter_id"`
	FirstName   string
	LastName    string
	CreatedDate *time.Time
}

//DonorResult holds a month and a stats record.
type DonorResult struct {
	SupporterID string `gorm:"supporter_id"`
	FirstName   string
	LastName    string
	Stat
}

//TopDonor is used to provide a primary key for storing stats by month.
type TopDonor struct {
	Donor
}

//TopDonorResult holds a month and a stats record for the "month over month" sheet.
type TopDonorResult struct {
	DonorResult
}

//KeyValue implements KeyValuer by returning the value of a key for the
//DonorResult object.
func (r DonorResult) KeyValue(i int) (key interface{}) {
	switch i {
	case 0:
		key = r.FirstName
	case 1:
		key = r.LastName
	default:
		fmt.Printf("Error in DonorResult\n%+v\n", r)
		err := fmt.Errorf("Not a valid DonorResult index, %v", i)
		panic(err)
	}
	return key
}

//FillKeys implements KeyFiller by filling Excel cells with keys from the
//year table.
func (r DonorResult) FillKeys(rt *Runtime, sheet Sheet, row, col int) int {
	for j := 0; j < len(sheet.KeyNames); j++ {
		v := r.KeyValue(j)
		s := sheet.KeyStyles[j]
		rt.Cell(sheet.Name, row, col+j, v, s)
	}
	return row
}

//FillKeys implements KeyFiller by filling Excel cells with keys from the
//year table.
func (r TopDonorResult) FillKeys(rt *Runtime, sheet Sheet, row, col int) int {
	m := DonorResult{}
	return m.FillKeys(rt, sheet, row, col)
}

//Fill implements Filler by filling in a spreadsheet using data from the years table.
func (r Donor) Fill(rt *Runtime, sheet Sheet, row, col int) int {
	var a []DonorResult
	// y := Year{}
	// year := y.Largest(rt)
	//.Where(`stats.created_date LIKE "%?%"`, year)
	rt.DB.Table("supporters").Select("supporters.first_name, supporters.last_name, stats.*").Joins("left join stats on stats.id = supporters.supporter_id").Order("stats.all_amount desc").Scan(&a)
	for _, r := range a {
		rt.Spreadsheet.InsertRow(sheet.Name, row+1)
		r.FillKeys(rt, sheet, row, 0)
		r.Stat.Fill(rt, sheet.Name, row, len(sheet.KeyNames))
		row++
	}
	return row
}

//Fill implements Filler by filling in a spreadsheet using data from the years table.
func (r TopDonor) Fill(rt *Runtime, sheet Sheet, row, col int) int {
	var a []DonorResult
	// y := Year{}
	// year := y.Largest(rt)
	//.Where(`stats.created_date LIKE "%?%"`, year)

	rt.DB.Order("stats.all_amount desc").Table("supporters").Select("first_name, last_name, stats.*").Joins("left join stats on stats.id = supporters.supporter_id").Limit(rt.TopDonorLimit).Scan(&a)
	for _, r := range a {
		rt.Spreadsheet.InsertRow(sheet.Name, row+1)
		r.FillKeys(rt, sheet, row, 0)
		r.Stat.Fill(rt, sheet.Name, row, len(sheet.KeyNames))
		row++
	}
	return row
}

//NewAllDonorsSheet builds the data used to decorate all donors sheet.
func (rt *Runtime) NewAllDonorsSheet() Sheet {
	filler := Donor{}
	result := DonorResult{}
	y := Year{}
	year := y.Largest(rt)
	name := fmt.Sprintf("All donors, %d", year)
	sheet := Sheet{
		Titles: []string{
			fmt.Sprintf("Ranked donors for %d", year),
			"Provided by the Custom Success group At Salsalabs",
		},
		Name:      name,
		KeyNames:  []string{"First Name", "Last Name"},
		KeyStyles: []int{rt.KeyStyle, rt.KeyStyle, rt.KeyStyle},
		Filler:    filler,
		KeyFiller: result,
	}
	return sheet
}

//NewTopDonorsSheet builds the data used to decorate the Top donors for the year sheet.
func (rt *Runtime) NewTopDonorsSheet() Sheet {
	filler := TopDonor{}
	result := TopDonorResult{}
	y := Year{}
	year := y.Largest(rt)
	name := fmt.Sprintf("Top donors for %d", year)
	sheet := Sheet{
		Titles: []string{
			fmt.Sprintf("Top %d donors for %d", rt.TopDonorLimit, year),
			"Provided by the Custom Success group At Salsalabs",
		},
		Name:      name,
		KeyNames:  []string{"First Name", "Last Name"},
		KeyStyles: []int{rt.KeyStyle, rt.KeyStyle, rt.KeyStyle},
		Filler:    filler,
		KeyFiller: result,
	}
	return sheet
}
