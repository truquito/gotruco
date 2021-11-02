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
