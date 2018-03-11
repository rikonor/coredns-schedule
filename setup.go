package schedule

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("schedule", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	sch := Schedule{
		ScheduleEntries: []ScheduleEntry{},
	}

	for c.Next() {
		// Get the zone
		if !c.NextArg() {
			return plugin.Error("schedule", fmt.Errorf("expected zone"))
		}
		z := c.Val()

		// Get the start time
		if !c.NextArg() {
			return plugin.Error("schedule", fmt.Errorf("expected start time"))
		}
		startTimeStr := c.Val()

		// Get the end time
		if !c.NextArg() {
			return plugin.Error("schedule", fmt.Errorf("expected end time"))
		}
		endTimeStr := c.Val()

		// Expect no additional configuration
		if c.NextArg() {
			return plugin.Error("schedule", c.ArgErr())
		}

		// Parse the start-time
		startTimeHr, startTimeMnt, err := ParseHourMinutePair(startTimeStr)
		if err != nil {
			return plugin.Error("schedule", fmt.Errorf("failed to parse time: %s", err))
		}

		// Parse the end-time
		endTimeHr, endTimeMnt, err := ParseHourMinutePair(endTimeStr)
		if err != nil {
			return plugin.Error("schedule", fmt.Errorf("failed to parse time: %s", err))
		}

		// Append a new schedule entry
		sch.ScheduleEntries = append(sch.ScheduleEntries, ScheduleEntry{
			Zone:      z,
			StartTime: DayTime{Hour: startTimeHr, Minute: startTimeMnt},
			EndTime:   DayTime{Hour: endTimeHr, Minute: endTimeMnt},
		})
	}

	// Add the Plugin to CoreDNS
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		sch.Next = next
		return sch
	})

	fmt.Printf("%+v\n", sch)

	return nil
}

// ParseHourMinutePair parses a string with the format HH:MM into
// an hour and minute integers.
func ParseHourMinutePair(s string) (int, int, error) {
	ps := strings.Split(s, ":")

	if len(ps) != 2 {
		return 0, 0, fmt.Errorf("invalid time. Should be of the format HH:MM")
	}

	hr, err := strconv.Atoi(ps[0])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse hour: %s", err)
	}
	if hr < 0 || hr > 23 {
		return 0, 0, fmt.Errorf("Hour is out of range. 0 <= hr <= 23")
	}

	mnt, err := strconv.Atoi(ps[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse minute: %s", err)
	}
	if mnt < 0 || mnt > 60 {
		return 0, 0, fmt.Errorf("Minute is out of range. 0 <= hr <= 60")
	}

	return hr, mnt, nil
}

// TimeOfDay returns hr:mnt in the given day of t
func TimeOfDay(t time.Time, hr int, mnt int) (time.Time, error) {
	if hr < 0 || hr > 23 {
		return time.Time{}, fmt.Errorf("Hour is out of range. 0 <= hr <= 23")
	}

	if mnt < 0 || mnt > 60 {
		return time.Time{}, fmt.Errorf("Minute is out of range. 0 <= hr <= 60")
	}

	y, m, d := t.Date()
	loc := t.Location()

	return time.Date(y, m, d, hr, mnt, 0, 0, loc), nil
}

// ParseTime parses an HH:MM string and returns
// that time within the day of the given time.
func ParseTime(t time.Time, s string) (time.Time, error) {
	hr, mnt, err := ParseHourMinutePair(s)
	if err != nil {
		return time.Time{}, err
	}

	tod, err := TimeOfDay(t, hr, mnt)
	if err != nil {
		return time.Time{}, err
	}

	return tod, nil
}
