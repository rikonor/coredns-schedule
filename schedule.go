package schedule

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

type Schedule struct {
	Next plugin.Handler

	ScheduleEntries []ScheduleEntry
}

type ScheduleEntry struct {
	Zone      string
	StartTime DayTime
	EndTime   DayTime
}

type DayTime struct {
	Hour   int
	Minute int
}

// IsDisabled checks if the name asked for is currently disabled
func (e ScheduleEntry) IsDisabled(name string, t time.Time) (bool, error) {
	// Name should be under zone
	if !strings.HasSuffix(name, e.Zone) {
		return false, nil
	}

	startTime, err := TimeOfDay(t, e.StartTime.Hour, e.StartTime.Minute)
	if err != nil {
		return false, err
	}

	endTime, err := TimeOfDay(t, e.EndTime.Hour, e.EndTime.Minute)
	if err != nil {
		return false, err
	}

	// Given time should be within disabled schedule
	if t.Before(startTime) || t.After(endTime) {
		return false, nil
	}

	return true, nil
}

// Name implements the Handler interface.
func (p Schedule) Name() string { return "schedule" }

// ServeDNS checks if the name asked for is currently disabled
func (p Schedule) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	for _, q := range r.Question {
		for _, e := range p.ScheduleEntries {
			isDisabled, err := e.IsDisabled(q.Name, time.Now())
			if err != nil {
				return dns.RcodeServerFailure, plugin.Error("schedule", err)
			}

			if isDisabled {
				return dns.RcodeRefused, plugin.Error("schedule", fmt.Errorf("Currently unavailable"))
			}
		}
	}

	return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
}
