package core

var UnknownCD = 4 // FIXME;

func GetDefaultCDByServiceType(serviceType ServiceType) *int {
	if serviceType == VoiceMO || serviceType == SmsMO {
		return &UnknownCD
	}

	return nil
}
