package conparser

import (
	"time"
	"strings"
	"strconv"
)

var data map[string][]int
var err error

func init() {
	data = make(map[string][]int, 5)
	data["day"] = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	data["week"] = []int{0, 1, 2, 3, 4, 5, 6}
	data["hour"] = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}
	data["month"] = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	data["minute"] = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59}
}

func IsRun(last_ran_time int64) bool {
	timestamp := time.Now().Unix()
	if timestamp - last_ran_time < 59 {
		return false
	}
	return true
}

func Parse(str string) bool {
	bits := strings.Split(str, " ")
	if len(bits) != 5 {
		return false
	}

	return _parse(bits)
}

func _parse(bits []string) bool {
	minute := bits[0]
	hour := bits[1]
	day := bits[2]
	month := bits[3]
	week := bits[4]

	if !_analysis(month, "month") {
		return false
	}

	if !_analysis(week, "week") {
		return false
	}

	if !_analysis(day, "day") {
		return false
	}

	if !_analysis(hour, "hour") {
		return false
	}

	if !_analysis(minute, "minute") {
		return false
	}
	return true
}

func _analysis(str, date_type string) bool {
	step := 0
	unit := make([]int, 0)
	result_unit := make([]int, 0)
	pos := strings.Index(str, "/")
	if pos != -1 {
		sc := strings.Split(str, "/")
		if len(sc) != 2 {
			return false
		}
		if sc[0] != "*" {
			return false
		}

		step, err = strconv.Atoi(sc[1])
		if err != nil {
			return false
		}

		str = sc[0]
	}

	if str == "*" {
		if _, ok := data[date_type]; !ok {
			return false
		}
		unit = data[date_type]
	} else {
		sections := make([]string, 0)
		pos = strings.Index(str, ",")
		if pos != -1 {
			sections = strings.Split(str, ",")
		} else {
			sections = append(sections, str)
		}

		if len(sections) > 0 {
			for _, v := range sections {
				pos = strings.Index(v, "-")
				if pos != -1 {
					se := strings.Split(v, "-")
					if len(se) != 2 {
						return false
					}
					start, err := strconv.Atoi(se[0])
					if err != nil {
						return false
					}
					end, err := strconv.Atoi(se[1])
					if err != nil {
						return false
					}

					if start > end {
						end, start = start, end
					}

					for i := start; i<=end; i++ {
						unit = append(unit, i)
					}
				} else {
					vint, err := strconv.Atoi(v)
					if err != nil {
						return false
					}
					unit = append(unit, vint)
				}
			}
		}
	}

	if step > 0 {
		if len(unit) > 0 {
			i := 0
			for _, v := range unit {
				if i % step == 0 {
					result_unit = append(result_unit, v)
				}
				i++
			}
		}
	} else {
		result_unit = unit
	}

	now_date := _get_now_number(date_type)
	retval := false
	for _, v := range result_unit {
		if v == now_date {
			retval = true
			break
		}
	}
	return retval
}

func _get_now_number(date_type string) int {
	str := time.Now().Format("2006/01/02 15/04/05")
	sc := strings.Split(str, " ")
	ymd := strings.Split(sc[0], "/")
	hms := strings.Split(sc[1], "/")

	switch(date_type) {
	case "month":
		mint, err := strconv.Atoi(ymd[1])
		if err != nil {
			return -1
		}
		return mint
	case "week":
		weeks := make(map[string]int, 7)
		weeks["Sunday"] = 0
		weeks["Monday"] = 1
		weeks["Tuesday"] = 2
		weeks["Wednesday"] = 3
		weeks["Thursday"] = 4
		weeks["Friday"] = 5
		weeks["Saturday"] = 6
		if _, ok := weeks[time.Now().Weekday().String()]; !ok {
			return -1
		}
		return weeks[time.Now().Weekday().String()]
	case "day":
		dint, err := strconv.Atoi(ymd[2])
		if err != nil {
			return -1
		}
		return dint
	case "hour":
		hint, err := strconv.Atoi(hms[0])
		if err != nil {
			return -1
		}
		return hint
	case "minute":
		mint, err := strconv.Atoi(hms[1])
		if err != nil {
			return -1
		}
		return mint
	}
	return -1
}

