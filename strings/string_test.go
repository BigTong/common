package strings

import (
	"testing"
)

type TestItem struct {
	src      string
	expected string
}

func TestStringClean(t *testing.T) {
	var testCases = []TestItem{
		TestItem{
			src:      "  aa",
			expected: "aa",
		},
		TestItem{
			src:      "aa    bb cc   ",
			expected: "aa bb cc",
		},
		TestItem{
			src:      "  ",
			expected: "",
		},
		TestItem{
			src:      "         aa    bb cc  dd         eee    ",
			expected: "aa bb cc dd eee",
		},
		TestItem{
			src:      "         aa    bb cc  dd         eee    fff",
			expected: "aa bb cc dd eee fff",
		},
	}
	for _, item := range testCases {
		ret := CleanString(item.src)
		if ret != item.expected {
			t.Error("not pass test cases, src:", item.src, " expected:", item.expected)
		}
	}
}
