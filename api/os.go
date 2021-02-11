package api

type OperatingSystem string

const (
	OperatingSystemLinux   OperatingSystem = "linux"
	OperatingSystemWindows OperatingSystem = "windows"

	OperatingSystemUnknown OperatingSystem = ""
)
