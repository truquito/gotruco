package util

import "strings"

func All(bs ...bool) bool {
	for _, b := range bs {
		if !b {
			return false
		}
	}
	return true
}

func Assert(should bool, callback func()) {
	if !should {
		callback()
	}
}

func Mod(a, b int) int {
	c := a % b
	if c < 0 {
		c += b
	}
	return c
}

func CheckElementsLen(slice []string) bool {
	for _, s := range slice {
		if len(s) == 0 {
			return false
		}
	}
	return true
}

func CheckSliceDuplicates(slice []string) bool {
	seen := make(map[string]bool)
	for _, s := range slice {
		l := strings.ToLower(s)
		if seen[l] {
			return true
		}
		seen[l] = true
	}
	return false
}
