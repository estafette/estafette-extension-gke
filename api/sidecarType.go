package api

type SidecarType string

const (
	SidecarTypeOpenresty     SidecarType = "openresty"
	SidecarTypeESP           SidecarType = "esp"
	SidecarTypeESPv2         SidecarType = "espv2"
	SidecarTypeCloudSQLProxy SidecarType = "cloudsqlproxy"
	SidecarTypeIstio         SidecarType = "istio"

	SidecarTypeUnknown SidecarType = ""
)
