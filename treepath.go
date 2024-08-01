package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func parseTreeSpec(spec string) ([]int64, error) {
	splitStr := strings.Split(spec, ",")
	treespec := make([]int64, len(splitStr))
	for i, v := range splitStr {
		vi, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("error parsing treespec: %w", err)
		}
		treespec[i] = int64(vi)
	}
	return treespec, nil
}

func numToTreePath(n int64, treespec []int64) (string, error) {
	if n < 0 {
		return "", errors.New("n cannot be less than 0")
	}

	maxN := int64(1)
	for _, v := range treespec {
		maxN *= v
	}
	if n > maxN {
		return "", fmt.Errorf("n cannot be more than %v", maxN)
	}

	nums := make([]int64, len(treespec))
	for i := len(treespec) - 1; i >= 0; i -= 1 {
		nums[i] = n % treespec[i]
		n = n / treespec[i]
	}

	var sb strings.Builder
	for i := 0; i < len(nums)-1; i += 1 {
		sb.WriteString(strconv.Itoa(int(nums[i])))
		sb.WriteString("/")
	}
	sb.WriteString(strconv.Itoa(int(nums[len(nums)-1])))

	return sb.String(), nil
}
