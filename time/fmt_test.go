package time

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestTimeFmt(t *testing.T) {
	expectedFmt := "2017-04-28 15:02:12"
	dates := []string{
		"2017/04/28 15:02:12",
		"2017/4/28 15:02:12",
		"2017-4-28 15:02:12",
		"2017-04-28 15:02:12",
	}

	for _, d := range dates {
		ret, err := TimeFormat(d)
		assert.Equal(t, expectedFmt, ret)
		assert.Equal(t, nil, err)
	}
}
