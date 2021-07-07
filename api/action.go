package api

type ActionType string

const (
	ActionDeploySimple  ActionType = "deploy-simple"
	ActionDeployCanary  ActionType = "deploy-canary"
	ActionDeployStable  ActionType = "deploy-stable"
	ActionRestartSimple ActionType = "restart-simple"
	ActionRestartCanary ActionType = "restart-canary"
	ActionRestartStable ActionType = "restart-stable"
	ActionDiffSimple    ActionType = "diff-simple"
	ActionDiffCanary    ActionType = "diff-canary"
	ActionDiffStable    ActionType = "diff-stable"
	ActionDiffDelete    ActionType = "diff-delete"
	ActionDelete        ActionType = "delete"

	ActionRollbackCanary ActionType = "rollback-canary"

	ActionUnknown ActionType = ""
)
