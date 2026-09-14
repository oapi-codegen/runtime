package types

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDate_MarshalJSON(t *testing.T) {
	testDate := time.Date(2019, 4, 1, 0, 0, 0, 0, time.UTC)
	b := struct {
		DateField Date `json:"date"`
	}{
		DateField: Date{testDate},
	}
	jsonBytes, err := json.Marshal(b)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"date":"2019-04-01"}`, string(jsonBytes))
}

func TestDate_UnmarshalJSON(t *testing.T) {
	testDate := time.Date(2019, 4, 1, 0, 0, 0, 0, time.UTC)
	jsonStr := `{"date":"2019-04-01"}`
	b := struct {
		DateField Date `json:"date"`
	}{}
	err := json.Unmarshal([]byte(jsonStr), &b)
	assert.NoError(t, err)
	assert.Equal(t, testDate, b.DateField.Time)
}

func TestDate_Stringer(t *testing.T) {
	t.Run("nil date", func(t *testing.T) {
		var d *Date
		assert.Equal(t, "<nil>", fmt.Sprintf("%v", d))
	})

	t.Run("ptr date", func(t *testing.T) {
		d := &Date{
			Time: time.Date(2019, 4, 1, 0, 0, 0, 0, time.UTC),
		}
		assert.Equal(t, "2019-04-01", fmt.Sprintf("%v", d))
	})

	t.Run("value date", func(t *testing.T) {
		d := Date{
			Time: time.Date(2019, 4, 1, 0, 0, 0, 0, time.UTC),
		}
		assert.Equal(t, "2019-04-01", fmt.Sprintf("%v", d))
	})
}

func TestDate_MarshalText(t *testing.T) {
	date := Date{Time: time.Date(2022, 6, 14, 0, 0, 0, 0, time.UTC)}

	value, err := date.MarshalText()

	assert.NoError(t, err)
	assert.Equal(t, "2022-06-14", string(value))
}

func TestDate_TextRoundTrip(t *testing.T) {
	testDate := time.Date(2022, 6, 14, 0, 0, 0, 0, time.UTC)

	value, err := Date{Time: testDate}.MarshalText()
	assert.NoError(t, err)

	date := Date{}
	err = date.UnmarshalText(value)

	assert.NoError(t, err)
	assert.Equal(t, testDate, date.Time)
}

func TestDate_XMLRoundTrip(t *testing.T) {
	testDate := time.Date(2019, 4, 1, 0, 0, 0, 0, time.UTC)
	type body struct {
		XMLName   xml.Name `xml:"body"`
		DateField Date     `xml:"date"`
	}

	xmlBytes, err := xml.Marshal(body{DateField: Date{testDate}})
	assert.NoError(t, err)
	assert.Equal(t, `<body><date>2019-04-01</date></body>`, string(xmlBytes))

	var b body
	err = xml.Unmarshal(xmlBytes, &b)

	assert.NoError(t, err)
	assert.Equal(t, testDate, b.DateField.Time)
}

func TestDate_UnmarshalText(t *testing.T) {
	testDate := time.Date(2022, 6, 14, 0, 0, 0, 0, time.UTC)
	value := []byte("2022-06-14")

	date := Date{}
	err := date.UnmarshalText(value)

	assert.NoError(t, err)
	assert.Equal(t, testDate, date.Time)
}
