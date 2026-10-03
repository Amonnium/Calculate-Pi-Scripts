package platform

import "runtime"

type Target struct {
	OS   string
	Arch string
}

func Detect() Target {
	return Target{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}
