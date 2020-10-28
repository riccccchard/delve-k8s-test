package utils

import (
	"strconv"
	"strings"
)

// version compare
const (
	GT               = 1        // greater than
	EQ               = 0        // equal to
	LT               = -1       // less than
	INVALID          = -2       // invalid param
	SDK_BASE_VERSION = "1.0.18" //base version 1.0.18
)

// 正常版本号(3位) 1.0.18-snapshot
// hotfix版本号(4位) 1.0.18.1-snapshot
// baseVersion 1.0.18
// 1.0.18 > 1.0.17
func VersionCompare(currentVersion, baseVersion string) int {
	if currentVersion == "" || baseVersion == "" {
		return INVALID
	}
	currentArr := strings.Split(currentVersion, "-")
	current := strings.Split(currentArr[0], ".")
	base := strings.Split(baseVersion, ".")
	minLen := minInt(len(current), len(base))
	currentIsLonger := len(current) > len(base) // 当前版本号是否为4位版本号
	var result int

	// 循环比较版本号
	for i := 0; i < minLen; i++ {
		currentValue, _ := strconv.Atoi(current[i])
		baseValue, _ := strconv.Atoi(base[i])
		result = compareInt(currentValue, baseValue);
		if result == EQ {
			continue
		} else {
			break
		}
	}

	// 四位版本号
	if result == EQ && currentIsLonger {
		return GT
	}
	return result
}

func minInt(n1, n2 int) int {
	if n1 > n2 {
		return n2
	}

	return n1
}

func compareInt(n1, n2 int) int {
	if n1 > n2 {
		return GT
	} else if n1 < n2 {
		return LT
	}

	return EQ
}
