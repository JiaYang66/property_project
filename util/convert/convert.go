package convert

import (
	"fmt"
	"strconv"
)

//StrToFloat64 支持指定精度
func StrToFloat64(str string, len int) (float64, error) {
	lenstr := "%." + strconv.Itoa(len) + "f"
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}
	nstr := fmt.Sprintf(lenstr, value) //指定精度
	val, err := strconv.ParseFloat(nstr, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}
