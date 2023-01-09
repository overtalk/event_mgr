package utils

import "time"

func UtcTime(hour_offset int64) time.Time {
	utc_time := time.Now().UTC()
	if hour_offset == 0 {
		return utc_time
	}
	return utc_time.Add(time.Duration(hour_offset) * time.Hour)
}

func TsToUtcTime(ts int64, hour_offset int64) time.Time {
	utc_time := time.Unix(ts, 0).UTC()
	if hour_offset == 0 {
		return utc_time
	}
	return utc_time.Add(time.Duration(hour_offset) * time.Hour)
}

func UtcTs() int64 { return UtcTime(0).Unix() }
