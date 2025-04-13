package utils

import (
	"errors"
	"strconv"
	"strings"
)

// ConvertOffset converts an offset to a part number (int32), the offset
// is expected to be in the following format <begin>-<end> and the
// offset is assumed to be consistent.
func ConvertOffset(offset string) (int32, error) {
	ranges := strings.Split(offset, "-")
	if len(ranges) != 2 {
		return -1, errors.New("invalid offset")
	}

	start, err := strconv.Atoi(ranges[0])
	if err != nil {
		return -1, err
	}

	end, err := strconv.Atoi(ranges[1])
	if err != nil {
		return -1, err
	}

	chunk := end - start

	partNumber := int32(start / chunk)

	return partNumber, nil
}
