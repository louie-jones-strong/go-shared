package elements

import (
	"testing"
	"time"

	"github.com/louie-jones-strong/go-shared/dataframe/apptype"
	"github.com/stretchr/testify/assert"
)

func TestUnit_DateEncoding(t *testing.T) {
	dateTime := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	elem := NewDateTimeElement(dateTime)

	str := elem.ToString()

	res, err := apptype.StringToTime(str)
	assert.NoError(t, err)

	assert.Equal(t, dateTime, res)
}
