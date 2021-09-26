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
	dateLayout := "2006-01-02"
	t, _ := time.Parse(dateLayout, vi.Date)
	return fmt.Sprintf("%s, commit %s, date %s", vi.Version, vi.Commit[0:7], t.Format(dateLayout))
}
