package core

type DistributionModel byte

const (
	RetrospectiveDistributionModel DistributionModel = iota
	MovingAverageDistributionModel
	YTDDistributionModel
	ManualVolumeDistributionModel
)
