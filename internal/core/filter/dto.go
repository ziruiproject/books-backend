package filter

import "time"

type Default struct {
	Search    string `json:"search" validate:"max=100"`
	StartDate string `json:"start_date" validate:"omitempty,publication_date"`
	EndDate   string `json:"end_date" validate:"omitempty,publication_date"`
}

type BookFilter struct {
	Default
	Categories []uint `json:"categories"`
	From       string `json:"from" validate:"omitempty,publication_date"`
	To         string `json:"to" validate:"omitempty,publication_date"`
}

func NewDefaultFilter(filter *Default) {
	now := time.Now()
	defaultStart := now.AddDate(-10, 0, 0).Format("2006-01-02")
	defaultEnd := now.AddDate(0, 0, 2).Format("2006-01-02")

	if filter.StartDate == "" {
		filter.StartDate = defaultStart
	}
	if filter.EndDate == "" {
		filter.EndDate = defaultEnd
	}
}

func NewBookFilter(filter *BookFilter) {
	NewDefaultFilter(&filter.Default)
	now := time.Now()
	defaultFrom := now.AddDate(-100, 0, 0).Format("2006-01-02")
	defaultTo := now.Format("2006-01-02")

	if filter.From == "" {
		filter.From = defaultFrom
	}

	if filter.To == "" {
		filter.To = defaultTo
	}
}
