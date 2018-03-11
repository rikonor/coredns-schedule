package schedule

import (
	"testing"
	"time"
)

func TestParseHourMinutePair(t *testing.T) {
	type testCase struct {
		in     string
		outHr  int
		outMnt int
	}

	testCases := []testCase{
		testCase{
			in:     "12:30",
			outHr:  12,
			outMnt: 30,
		},
	}

	for _, tc := range testCases {
		hr, mnt, err := ParseHourMinutePair(tc.in)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if hr != tc.outHr || mnt != tc.outMnt {
			t.Fatalf("wrong output: %d:%d, expected %d:%d", hr, mnt, tc.outHr, tc.outMnt)
		}
	}
}

func TestTimeOfDay(t *testing.T) {
	type testCase struct {
		inT   time.Time
		inHr  int
		inMnt int
		out   time.Time
	}

	testCases := []testCase{
		testCase{
			inT:   time.Date(2018, 1, 1, 5, 20, 14, 11, time.UTC),
			inHr:  3,
			inMnt: 10,
			out:   time.Date(2018, 1, 1, 3, 10, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		out, err := TimeOfDay(tc.inT, tc.inHr, tc.inMnt)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if !out.Equal(tc.out) {
			t.Fatalf("wrong output: %s, expected %s", out, tc.out)
		}
	}
}

func TestParseTime(t *testing.T) {
	type testCase struct {
		in  string
		out time.Time
	}

	tt := time.Date(2018, 1, 1, 3, 10, 0, 0, time.UTC)

	testCases := []testCase{
		testCase{
			in:  "12:30",
			out: time.Date(2018, 1, 1, 12, 30, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		out, err := ParseTime(tt, tc.in)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if !out.Equal(tc.out) {
			t.Fatalf("wrong output: %s, expected %s", out, tc.out)
		}
	}
}
