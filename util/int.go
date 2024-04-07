package util

import "strconv"

func Atoi32(s string) (int32, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return int32(id), nil
}
