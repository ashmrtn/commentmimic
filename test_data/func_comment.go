package test_data

import (
	"errors"
)

const (
	// someConst is responsible for nothing.
	someConst = 24
	// OtherConst is responsible for everything.
	otherConst = 42
)

var (
	// x is.
	x = 7
	// y is not.
	y = 17
)

// someFunction is a test function that just has a comment to practice parsing.
func someFunction(a, b int, c string) error {
	if a < b {
		return errors.New("a less than b")
	}

	// This is some check to ensure validity later on.
	if len(c) == 0 {
		return errors.New("c is empty")
	}

	return nil
}

// This comment does not start with the function name.
func otherFunction(a, b int, c string) error {
	if a < b {
		return errors.New("a less than b")
	}

	// This is some check to ensure validity later on.
	if len(c) == 0 {
		return errors.New("c is empty")
	}

	return nil
}
