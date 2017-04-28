package time

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var datePatterns = []string{
	time.RFC3339,
	"2006年1月2日",
	"2006年01月",
	"2006年01月02日",
	"2006/1/2",
	"2006/01/02",
	"2006-1-2",
	"2006-01-02",
	"2006/01/02 15:04:05",
	"2006/1/02 15:04:05",
	"2006/01/2 15:04:05",
	"2006/01/02 15:04:05",
	"2006/1/2 15:04:05",
	"2006-01-02 15:04",
	"2006-01-02 15:04:05",
	"2006-1-2 15:04:05",
	"20060102",
	"Mon Jan 02 15:04:05 CST 2006",
	"Mon, 2 Jan 2006 15:04:05 -0700 (MST)",
	"Mon, 2 Jan 2006 15:04:05 -0700 (MST+07:00)",
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"2 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 -0700 (MST)",
	"Mon, 02 Jan 2006 15:04:05 -0700 (MST+07:00)",
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"02 Jan 2006 15:04:05 -0700",
}

func ToTimeStamp(src string) (int64, bool) {
	if (len(src) == 10 ||
		len(src) == 13) &&
		strings.HasPrefix(src, "1") {
		timeStamp, err := strconv.ParseInt(src, 10, 64)
		if err != nil {
			return 0, false
		}

		if timeStamp > 10000000000 {
			timeStamp /= 1000
		}
		return timeStamp, true
	}
	return 0, false
}

func GuessDate(buf string) (time.Time, error) {
	buf = strings.TrimSpace(buf)
	stamp, ok := ToTimeStamp(buf)
	if ok {
		return time.Unix(stamp, 0), nil
	}

	for _, p := range datePatterns {
		if tm, err := time.Parse(p, buf); err == nil {
			return tm, nil
		}
	}
	return time.Now(), errors.New("fail to parse date")
}

func TimeFormat(date string) (string, error) {
	if len(date) == 0 {
		return "", errors.New("get empty date string")
	}
	date = strings.Split(date, ".")[0]
	tm, err := GuessDate(date)
	return tm.Format("2006-01-02 15:04:05"), err
}
