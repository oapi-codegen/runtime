package types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDuration(t *testing.T) {
	valid := map[string]Duration{
		"P1Y":              {Years: 1},
		"P3Y6M4DT12H30M5S": {Years: 3, Months: 6, Days: 4, Hours: 12, Minutes: 30, Seconds: 5},
		"P1M":              {Months: 1},
		"PT1M":             {Minutes: 1},
		"P0D":              {},
		"PT0S":             {},
		"PT36H":            {Hours: 36},
		"P2W":              {Weeks: 2},
		"PT0.5S":           {Seconds: 0.5},
		"PT0,5S":           {Seconds: 0.5},
		"PT1H30M2.25S":     {Hours: 1, Minutes: 30, Seconds: 2.25},
		"P10Y10M10D":       {Years: 10, Months: 10, Days: 10},
	}
	for input, want := range valid {
		t.Run(input, func(t *testing.T) {
			got, err := ParseDuration(input)
			require.NoError(t, err)
			assert.Equal(t, want, got)
		})
	}

	invalid := []string{
		"",        // empty
		"P",       // no components
		"PT",      // T without time components
		"1Y",      // missing P
		"P1S",     // seconds without T
		"PT1D",    // date designator in time part
		"P1M1Y",   // out of order
		"PT1S1H",  // out of order
		"P1Y1Y",   // repeated
		"P1W2D",   // week combined with other components
		"P1D2W",   // week combined with other components
		"PT1W",    // week in time part
		"P-1D",    // negative
		"P1.5Y",   // fraction outside seconds
		"PT1.5H",  // fraction outside seconds
		"PT1.S",   // fraction without digits
		"P1YT",    // trailing T
		"P1",      // number without designator
		"PT5X",    // unknown designator
		"p1y",     // lowercase
		" P1Y",    // leading junk
		"P1Y ",    // trailing junk
		"PT1H30M5",
	}
	for _, input := range invalid {
		t.Run("invalid/"+input, func(t *testing.T) {
			_, err := ParseDuration(input)
			require.Error(t, err)
		})
	}
}

func TestDurationString(t *testing.T) {
	// Every parseable value serializes back to its canonical form.
	for _, canonical := range []string{
		"P1Y", "P3Y6M4DT12H30M5S", "PT1M", "PT36H", "P2W", "PT0.5S",
		"PT1H30M2.25S", "P10Y10M10D", "PT0S",
	} {
		d, err := ParseDuration(canonical)
		require.NoError(t, err)
		assert.Equal(t, canonical, d.String())
	}

	// Non-canonical spellings normalize.
	d, err := ParseDuration("PT0,5S")
	require.NoError(t, err)
	assert.Equal(t, "PT0.5S", d.String())
	d, err = ParseDuration("P0D")
	require.NoError(t, err)
	assert.Equal(t, "PT0S", d.String())

	// The zero value is the zero duration.
	assert.Equal(t, "PT0S", Duration{}.String())
}

func TestDurationTimeDuration(t *testing.T) {
	d, err := ParseDuration("P1DT2H30M1.5S")
	require.NoError(t, err)
	td, err := d.TimeDuration()
	require.NoError(t, err)
	assert.Equal(t, 24*time.Hour+2*time.Hour+30*time.Minute+1500*time.Millisecond, td)

	d, err = ParseDuration("P2W")
	require.NoError(t, err)
	td, err = d.TimeDuration()
	require.NoError(t, err)
	assert.Equal(t, 14*24*time.Hour, td)

	// Years and months have no fixed length.
	for _, input := range []string{"P1Y", "P1M"} {
		d, err := ParseDuration(input)
		require.NoError(t, err)
		_, err = d.TimeDuration()
		require.Error(t, err)
	}
}

func TestFromTimeDuration(t *testing.T) {
	d, err := FromTimeDuration(90 * time.Minute)
	require.NoError(t, err)
	assert.Equal(t, "PT1H30M", d.String())

	d, err = FromTimeDuration(25*time.Hour + 500*time.Millisecond)
	require.NoError(t, err)
	assert.Equal(t, "PT25H0.5S", d.String())

	d, err = FromTimeDuration(0)
	require.NoError(t, err)
	assert.Equal(t, "PT0S", d.String())

	_, err = FromTimeDuration(-time.Second)
	require.Error(t, err)
}

func TestDurationJSON(t *testing.T) {
	type payload struct {
		Retry Duration `json:"retry"`
	}

	marshaled, err := json.Marshal(payload{Retry: Duration{Hours: 1, Minutes: 30}})
	require.NoError(t, err)
	assert.JSONEq(t, `{"retry":"PT1H30M"}`, string(marshaled))

	var decoded payload
	require.NoError(t, json.Unmarshal([]byte(`{"retry":"P3Y6M4DT12H30M5S"}`), &decoded))
	assert.Equal(t, Duration{Years: 3, Months: 6, Days: 4, Hours: 12, Minutes: 30, Seconds: 5}, decoded.Retry)

	require.Error(t, json.Unmarshal([]byte(`{"retry":"one hour"}`), &decoded))
	require.Error(t, json.Unmarshal([]byte(`{"retry":42}`), &decoded))
}

func TestDurationTextAndBind(t *testing.T) {
	var d Duration
	require.NoError(t, d.UnmarshalText([]byte("PT15M")))
	assert.Equal(t, Duration{Minutes: 15}, d)

	text, err := d.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "PT15M", string(text))

	var bound Duration
	require.NoError(t, bound.Bind("P1DT12H"))
	assert.Equal(t, Duration{Days: 1, Hours: 12}, bound)

	// Empty input leaves the value untouched, like the other scalar types.
	require.NoError(t, bound.Bind(""))
	assert.Equal(t, Duration{Days: 1, Hours: 12}, bound)

	require.Error(t, bound.Bind("garbage"))
}
