package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

func mergePostings(postings []Posting) []Posting {
	mergedMap := make(map[time.Time]Posting)

	for _, posting := range postings {
		if existing, found := mergedMap[posting.Date]; found {
			existing.Entries = append(existing.Entries, posting.Entries...)
			mergedMap[posting.Date] = existing
		} else {
			mergedMap[posting.Date] = posting
		}
	}

	mergedPostings := make([]Posting, 0, len(mergedMap))
	for _, posting := range mergedMap {
		mergedPostings = append(mergedPostings, posting)
	}

	sort.Slice(mergedPostings, func(i, j int) bool {
		return mergedPostings[i].Date.Before(mergedPostings[j].Date)
	})

	return mergedPostings
}

func renderPosting(p Posting) {
	fmt.Printf("%s %s\n", p.Date.Format("2006-01-02"), p.Description)
	acctLen := 0
	for _, entry := range p.Entries {
		acctLen = max(acctLen, len(entry.AccountName))
	}
	for _, entry := range p.Entries {
		if entry.Value < 0 {
			fmt.Printf("    %s%s%d\n", entry.AccountName, strings.Repeat(" ", max(acctLen-len(entry.AccountName)+4, 4)), entry.Value)
		} else {
			fmt.Printf("    %s%s%d\n", entry.AccountName, strings.Repeat(" ", max(acctLen-len(entry.AccountName)+5, 5)), entry.Value)
		}
	}
}
