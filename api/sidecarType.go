package api

type SidecarType string

const (
	SidecarTypeOpenresty     SidecarType = "openresty"
	SidecarTypeESP           SidecarType = "esp"
	SidecarTypeCloudSQLProxy SidecarType = "cloudsqlproxy"
	SidecarTypeIstio         SidecarType = "istio"

	SidecarTypeUnknown SidecarType = ""
)
