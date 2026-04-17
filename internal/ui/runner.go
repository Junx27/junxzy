package ui

func RunStep(message string, fn func()) {
	Start(message)
	defer func() {
		if r := recover(); r != nil {
			StopError("Terjadi kesalahan")
			return
		}
	}()
	fn()
	StopSuccess(message + " selesai")
}
