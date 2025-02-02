package rt

import (
	"testing"

	"github.com/interline-io/transitland-lib/internal/testutil"
	"github.com/interline-io/transitland-lib/tlcsv"
)

func newTestValidator() *Validator {
	r, err := tlcsv.NewReader(testutil.RelPath("test/data/rt/bart-rt.zip"))
	if err != nil {
		panic(err)
	}
	fi, err := NewValidatorFromReader(r)
	if err != nil {
		panic(err)
	}
	return fi
}

func TestValidateHeader(t *testing.T) {
	fi := newTestValidator()
	msg, err := ReadFile(testutil.RelPath("test/data/rt/bart-trip-updates.pb"))
	if err != nil {
		t.Error(err)
	}
	header := msg.GetHeader()
	errs := fi.ValidateHeader(header, msg)
	for _, err := range errs {
		_ = err
		// fmt.Println(err)
	}
}

func TestValidateTripUpdate(t *testing.T) {
	fi := newTestValidator()
	msg, err := ReadFile(testutil.RelPath("test/data/rt/bart-trip-updates.pb"))
	if err != nil {
		t.Error(err)
	}
	ents := msg.GetEntity()
	if len(ents) == 0 {
		t.Error("no entities")
	}
	trip := ents[0].TripUpdate
	if trip == nil {
		t.Error("expected TripUpdate")
	}
	errs := fi.ValidateTripUpdate(trip, msg)
	for _, err := range errs {
		_ = err
		// fmt.Println(err)
	}
}

func TestValidateVehiclePosition(t *testing.T) {

}

func TestValidateAlert(t *testing.T) {

}
