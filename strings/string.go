package strings

import (
	"regexp"
	"strings"
	"time"

	ctime "github.com/BigTong/common/time"
	"github.com/BigTong/common/url"
)

func CleanString(src string) string {
	ret := strings.TrimSpace(src)
	ret = strings.Replace(ret, "\n", " ", -1)
	ret = strings.Replace(ret, "\t", " ", -1)
	ret = strings.Replace(ret, "\r", " ", -1)
	ret = strings.Replace(ret, "\u0009", " ", -1)
	ret = strings.Replace(ret, "&nbsp;", " ", -1)
	runes := []rune(ret)
	j := 0
	startSpace := true
	for i := 0; i < len(runes); i++ {
		if runes[i] != rune(' ') {
			startSpace = false
			runes[j] = runes[i]
			j++
			continue
		}

		if startSpace {
			continue
		}
		startSpace = true
		runes[j] = runes[i]
		j++
	}
	return string(runes[0:j])
}

func Segment(str, start, end string) string {
	segs1 := strings.Split(str, start)
	if len(segs1) < 2 {
		return str
	}

	segs := strings.Split(segs1[1], end)
	if len(segs) < 1 {
		return str
	}
	return segs[0]
}

func getDigitNum(str string) string {
	re := regexp.MustCompile("[0-9.,]+")
	digit := re.FindString(str)
	if len(digit) > 0 {
		p1 := strings.Index(str, "-")
		p2 := strings.Index(str, digit)
		if p1 < p2 && p1 >= 0 {
			digit = "-" + digit
		}
	}
	digit = strings.Replace(digit, ",", "", -1)
	if digit == "." || digit == "-" || digit == "" ||
		strings.Contains(digit, "--") {
		return "0"
	}
	return digit
}

type Command struct {
	Name  string
	Param string
}

type StringCmds struct {
	Name string
	Cmds []Command
}

func NewStringCmds(buf string) *StringCmds {
	tks := strings.Split(buf, "::")
	ret := &StringCmds{
		Name: tks[0],
		Cmds: []Command{},
	}
	for _, tk := range tks[1:] {
		p := strings.Index(tk, "(")
		if p < 0 || tk[len(tk)-1] != ')' {
			continue
		}
		cmd := Command{
			Name:  tk[0:p],
			Param: tk[p+1 : len(tk)-1],
		}
		ret.Cmds = append(ret.Cmds, cmd)
	}
	return ret
}

func innerSegment(src, cmd string) string {
	segs := strings.Split(cmd, " ")
	if len(segs) < 2 {
		return src
	}
	return Segment(src, segs[0], segs[1])
}

func fixTime(val string) string {
	re := regexp.MustCompile("[0-9]{1,2}:[0-9]{1,2}")
	if strings.Contains(val, "今天") || len(re.FindString(val)) != 0 {
		return time.Now().Format("2006-01-02")
	}

	if strings.Contains(val, "昨天") {
		duration, _ := time.ParseDuration("-24h")
		return time.Now().Add(duration).Format("2006-01-02")
	}

	return val
}

func regexpReplace(src, cmds string) string {
	re := regexp.MustCompile(cmds)
	src = re.ReplaceAllString(src, "")
	return src
}

func srcReplace(val, cmds string) string {
	cmd := strings.Split(cmds, " ")
	if len(cmd) != 2 {
		return val
	}
	return strings.Replace(cmd[0], cmd[1], val, -1)
}

func replaceKeyword(keyword string, val string) string {
	if strings.Contains(val, "{{$data}}") {
		return strings.Replace(val, "{{$data}}", keyword, -1)
	}
	return val
}

func normalizeCount(src string) string {
	if strings.HasSuffix(src, "万") {
		src = strings.TrimSuffix(src, "万")
		segs := strings.Split(src, ".")
		if len(segs) != 1 {
			if len(segs[1]) < 4 {
				for i := 0; i < 4-len(segs[1]); i++ {
					segs[1] += "0"
				}
				src = segs[0] + segs[1]
			} else {
				src = segs[0] + segs[1][0:3]
			}
		}
	}
	return src
}

func StringExtract(src, cmds string) (string, string) {
	sc := NewStringCmds(cmds)
	val := src
	for _, cmd := range sc.Cmds {
		switch cmd.Name {
		case "trim":
			val = strings.TrimSpace(val)
		case "trim_prefix":
			val = strings.TrimPrefix(val, cmd.Param)
		case "trim_suffix":
			val = strings.TrimSuffix(val, cmd.Param)
		case "digit":
			val = getDigitNum(val)
		case "segment":
			val = innerSegment(val, cmd.Param)
		case "clean":
			val = CleanString(val)
		case "regexp_replace":
			val = regexpReplace(val, cmd.Param)
		case "parse_url":
			val, _ = url.ResolveReference(cmd.Param, val)
		case "fix_time":
			val = fixTime(val)
		case "time_fmt":
			fmtTime, err := ctime.TimeFormat(val)
			val = fmtTime
			if err != nil {
				val = ""
			}
		case "trimlu":
			p1 := strings.Index(val, cmd.Param)
			if p1 >= 0 {
				val = val[p1+len(cmd.Param) : len(val)]
			}
		case "trimru":
			p1 := strings.LastIndex(val, cmd.Param)
			if p1 >= 0 {
				val = val[0:p1]
			}
		case "src_replace":
			val = srcReplace(val, cmd.Param)
		case "regex":
			re := regexp.MustCompile(cmd.Param)
			val = re.FindString(val)
		case "replace_keyword":
			val = replaceKeyword(val, cmd.Param)
		case "add_prefix":
			val = cmd.Param + val
		case "add_suffix":
			val = val + cmd.Param
		case "normal_count":
			val = normalizeCount(val)
		}

	}
	return sc.Name, val
}
