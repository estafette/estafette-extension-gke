package main

type StrategyType string

const (
	StrategyTypeRollingUpdate StrategyType = "RollingUpdate"
	StrategyTypeRecreate      StrategyType = "Recreate"
	StrategyTypeAtomicUpdate  StrategyType = "AtomicUpdate"

	StrategyTypeUnknown StrategyType = ""
)
