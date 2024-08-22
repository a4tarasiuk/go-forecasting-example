package core

type ServiceType byte

const (
	VoiceMO ServiceType = iota
	VoiceMT
	SmsMO
	SmsMT
	Data
)
