package gotransit

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"sort"
	"time"
)

// FeedVersion represents a single GTFS data source.
type FeedVersion struct {
	FeedID               int        `db:"feed_id" gorm:"index"`
	FeedType             string     `db:"feed_type" gorm:"index"`
	SHA1                 string     `db:"sha1" gorm:"index"`
	File                 string     `db:"file"`
	URL                  string     `db:"url"`
	EarliestCalendarDate *time.Time `db:"earliest_calendar_date" gorm:"index;not null"`
	LatestCalendarDate   *time.Time `db:"latest_calendar_date" gorm:"index;not null"`
	FetchedAt            *time.Time `db:"fetched_at"`
	ID                   int        `db:"id"`
	CreatedAt            *time.Time `db:"created_at"`
	UpdatedAt            *time.Time `db:"updated_at"`
}

// TableName sets the table name prefix.
func (ent *FeedVersion) TableName() string {
	return "feed_versions"
}

// NewFeedVersion returns a FeedVersion from a Reader.
func NewFeedVersion(reader Reader) (FeedVersion, error) {
	fv := FeedVersion{}
	// TODO
	// If a zip file, get SHA1
	// if v, ok := reader.(*Reader); ok {
	// 	if a, ok := v.Adapter.(*ZipAdapter); ok {
	// 		if s, err := fileSHA1(a.Path()); err == nil {
	// 			fv.SHA1 = s
	// 		}
	// 	}
	// }
	// Get Calendar Dates
	times := []time.Time{}
	for c := range reader.Calendars() {
		times = append(times, c.StartDate)
		times = append(times, c.EndDate)
	}
	for c := range reader.CalendarDates() {
		if c.ExceptionType == 1 {
			times = append(times, c.Date)
		}
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i].Before(times[j])
	})
	if len(times) > 0 {
		if times[0].Before(times[len(times)-1]) {
			fv.EarliestCalendarDate = &times[0]
			fv.LatestCalendarDate = &times[len(times)-1]
		}
	}
	return fv, nil
}

// fileSHA1 returns the SHA1 hash of the zip file
func fileSHA1(path string) (string, error) {
	var sha1string string
	file, err := os.Open(path)
	if err != nil {
		return sha1string, err
	}
	defer file.Close()
	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return sha1string, err
	}
	sha1bytes := hash.Sum(nil)[:20]
	sha1string = hex.EncodeToString(sha1bytes)
	return sha1string, nil
}
