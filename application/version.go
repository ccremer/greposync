package application

import (
	"fmt"
	"time"
)

type VersionInfo struct {
	Version string
	Commit  string
	Date    string
}

func (vi VersionInfo) String() string {
	date := ""
	if vi.Date != "" {
		dateLayout := "2006-01-02"
		t, _ := time.Parse(dateLayout, vi.Date)
		date = t.Format(dateLayout)
	}
	commit := vi.Commit
	if len(commit) > 7 {
		commit = vi.Commit[:7]
	}
	return fmt.Sprintf("%s, commit %s, date %s", vi.Version, commit, date)
}
