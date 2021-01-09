package app

import "time"

var fixture = []Record{
	{
		Start:    time.Date(2020, 12, 31, 0, 0, 0, 0, time.Local),
		Duration: 90 * time.Minute,
		Activity: "Documentation",
		Actors:   []string{"Powerman"},
		Customer: "CompanyA",
		Details:  "update something in notion",
	},
	{
		Start:    time.Date(2020, 12, 30, 23, 0, 0, 0, time.Local),
		Duration: 120 * time.Minute,
		Activity: "Meeting",
		Actors:   []string{"Powerman"},
		Customer: "Company A",
		Details:  "slack",
	},
	{
		Start:    time.Date(2020, 12, 20, 10, 15, 0, 0, time.Local),
		Duration: 60 * time.Minute,
		Activity: "mono",
		Actors:   []string{"Powerman", "Someone"},
		Customer: "Customer B",
	},
	{
		Start:    time.Date(2020, 12, 19, 10, 30, 0, 0, time.Local),
		Duration: 45 * time.Minute,
		Activity: "mono",
		Actors:   []string{"Powerman", "Someone"},
		Customer: "Customer B",
	},
	{
		Start:    time.Date(2020, 11, 15, 18, 30, 0, 0, time.Local),
		Duration: 15 * time.Minute,
		Activity: "mono",
		Actors:   []string{"Powerman"},
		Customer: "Customer B",
		Details:  "update deps",
	},
}

func (a *App) loadFixture() {
	for i := range fixture {
		a.AddRecord(fixture[i])
	}
}
