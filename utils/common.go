package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/satori/go.uuid"
)

// 获取当前系统的真实时间戳
func GetRealTime() int64 {
	return time.Now().Unix()
}

func NewUUIDV4() string {
	u4 := uuid.NewV4()
	return u4.String()
}

func ParseTimeKey(key string) (time.Time, error) {
	// 2019-01-01 00:00:00
	return time.Parse("2006-01-02 15:04:05", key)
}

func GetTimeDeltaSeconds(startTime, endTime string) (int, error) {
	start, err := ParseTimeKey(startTime)
	if err != nil {
		return 0, fmt.Errorf("parse start time failed: %v", err)
	}

	end, err := ParseTimeKey(endTime)
	if err != nil {
		return 0, fmt.Errorf("parse end time failed: %v", err)
	}

	return int(end.Sub(start).Seconds()), nil
}

// [1,2,3]int to "1,2,3"
func IntSliceToString(intSlice []int) string {
	result := ""
	for _, v := range intSlice {
		result += strconv.Itoa(v)
		result += ","
	}
	return result[:len(result)-1]
}
