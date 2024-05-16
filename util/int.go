package util

import "strconv"

// Atoi32 parses the given string and returns an int32 and an error
// from the strconv.Atoi method.
//
// It wraps the returned int from the strconv.Atoi method in the int32() function.
func Atoi32(s string) (int32, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return int32(id), nil
}
