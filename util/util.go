package util

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
