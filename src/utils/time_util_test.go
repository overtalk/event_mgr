package utils

import (
	"fmt"
	"testing"
)

func TestUtcTime(t *testing.T) {
	offset_list := []int64{-1, 8, 0}

	for i := 0; i < len(offset_list); i++ {
		time_offset := offset_list[i]
		time := UtcTime(time_offset)
		fmt.Printf("[TestUtcTime] time_offset:%d time:%v year:%d month:%d day:%d hour:%d minute:%d second:%d\n",
			time_offset, time, time.Year(), int(time.Month()), time.Day(), time.Hour(), time.Minute(), time.Second())
	}
}

func TestUtcTs(t *testing.T) {
	fmt.Println("[TestUtcTs] now ts = ", UtcTs())
}

func TestTsToUtcTime(t *testing.T) {
	fmt.Println("[TestTsToUtcTime] now ts = ", TsToUtcTime(UtcTs(), 8))
}
