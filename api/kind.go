package api

type Kind string

const (
	KindDeployment         Kind = "deployment"
	KindHeadlessDeployment Kind = "headless-deployment"
	KindStatefulset        Kind = "statefulset"
	KindJob                Kind = "job"
	KindCronJob            Kind = "cronjob"
	KindConfig             Kind = "config"
	KindConfigToFile       Kind = "config-to-file"

	KindUnknown Kind = ""
)
