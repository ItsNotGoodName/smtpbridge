package app

func closers(closers ...func()) func() {
	return func() {
		for _, v := range closers {
			v()
		}
	}
}
