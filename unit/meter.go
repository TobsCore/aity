package unit

import "strconv"

type Meter Unit

func (meter Meter) toKM() int64 {
	return int64(meter / 1000)
}

func (meter Meter) String() string {
	return strconv.Itoa(int(meter)) + "m"
}