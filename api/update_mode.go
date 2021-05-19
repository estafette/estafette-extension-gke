package api

type UpdateMode string

const (
	UpdateModeUnknown  UpdateMode = ""
	UpdateModeOff      UpdateMode = "Off"
	UpdateModeInitial  UpdateMode = "Initial"
	UpdateModeRecreate UpdateMode = "Recreate"
	UpdateModeAuto     UpdateMode = "Auto"
)
