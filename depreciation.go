package main

import (
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
)

type DepreciationConfig struct {
	InitialValue int    `toml:"initialValue"`
	ScrapValue   int    `toml:"scrapValue"`
	Period       int    `toml:"period"`
	StartDate    string `toml:"startDate"`
	AccountName  string `toml:"accountName"`
	Category     string `toml:"category"`
	Method       string `toml:"method"`
}

type Config struct {
	Payee          string               `toml:"payee"`
	ExpensesPrefix string               `toml:"expensesPrefix"`
	AssetPrefix    string               `toml:"assetPrefix"`
	Depreciations  []DepreciationConfig `toml:"depreciation"`
}

type Posting struct {
	Date        time.Time
	Description string
	Entries     []SingleEntry
}

type SingleEntry struct {
	AccountName string
	Value       int
}

func buildDepreciation(d DepreciationConfig, c Config) ([]Posting, error) {
	initialDate, err := time.Parse("2006-01-02", d.StartDate)
	if err != nil {
		return nil, err
	}

	postings := []Posting{}
	var deprAmounts []int

	switch d.Method {
	case "sl":
		deprAmounts, err = calculateStraightlineDepreciation(int(d.InitialValue), int(d.ScrapValue), d.Period)
	case "ddb":
		deprAmounts, err = calculateDoubleDecliningBalance(int(d.InitialValue), int(d.ScrapValue), d.Period)
	case "soy":
		deprAmounts, err = calculateSumOfTheYearsDigits(int(d.InitialValue), int(d.ScrapValue), d.Period)
	default:
		deprAmounts, err = calculateStraightlineDepreciation(int(d.InitialValue), int(d.ScrapValue), d.Period)
	}

	if err != nil {
		return nil, err
	}
	for i := 0; i < d.Period; i++ {
		expensesAccount := ""
		assetsAccount := ""
		if d.AccountName != "" {
			expensesAccount = strings.Join([]string{c.ExpensesPrefix, d.Category, d.AccountName}, ":")
			assetsAccount = strings.Join([]string{c.AssetPrefix, d.Category, d.AccountName}, ":")
		} else {
			expensesAccount = strings.Join([]string{c.ExpensesPrefix, d.Category}, ":")
			assetsAccount = strings.Join([]string{c.AssetPrefix, d.Category}, ":")
		}
		entries := []SingleEntry{
			{AccountName: expensesAccount, Value: deprAmounts[i]},
			{AccountName: assetsAccount, Value: -deprAmounts[i]},
		}
		p := Posting{
			Date:        initialDate.AddDate(0, i, 0),
			Description: c.Payee,
			Entries:     entries,
		}

		postings = append(postings, p)
	}

	return postings, nil
}

func Depreciations(cCtx *cli.Context) error {
	filePath := cCtx.String("file")
	config := &Config{}
	_, err := toml.DecodeFile(filePath, config)

	var postings []Posting

	for _, depConfig := range config.Depreciations {
		newPostings, err := buildDepreciation(depConfig, *config)
		postings = append(postings, newPostings...)
		if err != nil {
			return err
		}
	}

	postings = mergePostings(postings)

	if cCtx.Bool("thismonth") {
		currentMonth := time.Now().Month()
		currentYear := time.Now().Year()
		for _, posting := range postings {
			if currentMonth == posting.Date.Month() && currentYear == posting.Date.Year() {
				renderPosting(posting)
			}
		}
	} else {
		for _, posting := range postings {
			renderPosting(posting)
		}
	}

	if err != nil {
		return err
	}

	return nil
}
