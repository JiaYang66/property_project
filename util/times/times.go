package times

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
)

//时间 String => Time
func ToTime(str string) (time.Time, error) {
	p := strings.TrimSpace(str)
	if p == "" {
		return time.Time{}, errors.New("%s不能为空")
	}
	t, err := time.ParseInLocation(TimeFormat, str, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

//时间 Time => String
func ToStr(t ...time.Time) string {
	if len(t) == 0 {
		return time.Now().Format(TimeFormat)
	} else {
		return t[0].Format(TimeFormat)
	}
}

//返回当天剩余时间秒数
func GetRemainSecondsOneDay() time.Duration {
	todayLast := time.Now().Format("2006-01-02") + " 23:59:59"
	todayLastTime, _ := time.ParseInLocation(TimeFormat, todayLast, time.Local)
	remainSecond := time.Duration(todayLastTime.Unix()-time.Now().Local().Unix()) * time.Second
	return remainSecond
}

// excel日期字段格式化 yyyy-mm-dd
func ConvertToFormatDay(excelDaysString string) string {
	// 2006-01-02 距离 1900-01-01的天数
	baseDiffDay := 38719 //在网上工具计算的天数需要加2天，什么原因没弄清楚
	curDiffDay := excelDaysString
	b, _ := strconv.Atoi(curDiffDay)
	// 获取excel的日期距离2006-01-02的天数
	realDiffDay := b - baseDiffDay
	//fmt.Println("realDiffDay:",realDiffDay)
	// 距离2006-01-02 秒数
	realDiffSecond := realDiffDay * 24 * 3600
	//fmt.Println("realDiffSecond:",realDiffSecond)
	// 2006-01-02 15:04:05距离1970-01-01 08:00:00的秒数 网上工具可查出
	baseOriginSecond := 1136185445
	resultTime := time.Unix(int64(baseOriginSecond+realDiffSecond), 0).Format("2006-01-02")
	return resultTime
}

type JsonTime struct {
	time.Time
}

func (t *JsonTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	*t = JsonTime{
		now,
	}
	return
}

func (t JsonTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(TimeFormat))
	return []byte(formatted), nil
}

func (t JsonTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *JsonTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JsonTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
