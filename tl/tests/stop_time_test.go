package tests

import (
	"strconv"
	"testing"

	"github.com/interline-io/transitland-lib/tl"
)

type expectStopTime struct {
	ArrivalTime       int
	DepartureTime     int
	ShapeDistTraveled float64
}

type expectTrip struct {
	ExpectError       string
	ArrivalTime       []int
	DepartureTime     []int
	ShapeDistTraveled []float64
}

func expectTripToStopTime(e expectTrip) []tl.StopTime {
	ret := []tl.StopTime{}
	for i := range e.ArrivalTime {
		ret = append(ret, tl.StopTime{
			TripID:            "1",
			StopID:            strconv.Itoa(i),
			StopSequence:      i,
			ArrivalTime:       tl.NewWideTimeFromSeconds(e.ArrivalTime[i]),
			DepartureTime:     tl.NewWideTimeFromSeconds(e.DepartureTime[i]),
			ShapeDistTraveled: tl.NewOFloat(e.ShapeDistTraveled[i]),
		})
	}
	return ret
}

func TestValidateStopTimes(t *testing.T) {
	// base cases
	trips := []expectTrip{
		{"1", []int{10, 20, 30}, []int{10, 20, 30}, []float64{0, 1, 2}}, // all specified
		{"2", []int{10, 0, 30}, []int{10, 0, 30}, []float64{0, 1, 2}},   // ends specified
		{"3", []int{10, 20, 30}, []int{10, 20, 30}, []float64{0, 0, 0}}, // no dist
		{"4", []int{0, 20, 30}, []int{10, 20, 30}, []float64{0, 1, 2}},  // missing first arrival_time
		{"5", []int{10, 20, 30}, []int{10, 20, 0}, []float64{0, 1, 2}},  // missing last departure_time
		{"6", []int{10, 20, 30}, []int{10, 20, 30}, []float64{0, 1, 2}}, // two is OK
	}
	for _, et := range trips {
		t.Run(et.ExpectError, func(t *testing.T) {
			stoptimes := expectTripToStopTime(et)
			if errs := tl.ValidateStopTimes(stoptimes); len(errs) > 0 {
				t.Errorf("got %d errors, expected %d: %s", len(errs), 0, errs)
			}
		})
	}
	// error cases
	errortrips := []expectTrip{
		{"Error:OneStopTime", []int{10}, []int{10}, []float64{0}},
		{"Error:NoFinalArrivalTime", []int{10, 0}, []int{10, 0}, []float64{0, 0}},
		{"SequenceError:departure_time", []int{10, 20, 5}, []int{10, 20, 5}, []float64{0, 1, 2}},
		{"SequenceError:shape_pt_traveled", []int{10, 20, 30}, []int{10, 20, 30}, []float64{1, 2, 1}},
	}
	for _, et := range errortrips {
		t.Run(et.ExpectError, func(t *testing.T) {
			stoptimes := expectTripToStopTime(et)
			if errs := tl.ValidateStopTimes(stoptimes); len(errs) != 1 {
				t.Errorf("expected 1 error, got 0")
			}
		})
	}
	// Check for duplicate IDs
	errorStopSequence := expectTrip{"", []int{10, 20, 30}, []int{10, 20, 30}, []float64{0, 1, 2}}
	t.Run("SequenceError:stop_sequence", func(t *testing.T) {
		stoptimes := expectTripToStopTime(errorStopSequence)
		stoptimes[0].StopSequence = 1
		stoptimes[1].StopSequence = 2
		stoptimes[2].StopSequence = 2
		if errs := tl.ValidateStopTimes(stoptimes); len(errs) != 1 {
			t.Errorf("expected 1 error, got 0")
		}
	})
}

func BenchmarkValidateStopTime(b *testing.B) {
	trip := expectTrip{"1", []int{10, 20, 30, 40, 50, 60}, []int{10, 20, 30, 40, 50, 60}, []float64{0, 1, 2, 3, 4, 5, 6}}
	stoptimes := expectTripToStopTime(trip)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tl.ValidateStopTimes(stoptimes)
	}
}
