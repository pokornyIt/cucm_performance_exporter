package main

type AllowedCounters struct {
	allowedCounterName string
	prometheusName     string
	defaultEnabled     bool
}

const (
	PartiallyRegisteredPhone     = "PartiallyRegisteredPhone"
	RegisteredHardwarePhones     = "RegisteredHardwarePhones"
	GatewaysSessionsActive       = "GatewaysSessionsActive"
	GatewaysSessionsFailed       = "GatewaysSessionsFailed"
	PhoneSessionsActive          = "PhoneSessionsActive"
	PhoneSessionsFailed          = "PhoneSessionsFailed"
	AnnunciatorOutOfResources    = "AnnunciatorOutOfResources"
	AnnunciatorResourceActive    = "AnnunciatorResourceActive"
	AnnunciatorResourceAvailable = "AnnunciatorResourceAvailable"
	AnnunciatorResourceTotal     = "AnnunciatorResourceTotal"

	AuthenticatedCallsActive              = "AuthenticatedCallsActive"
	AuthenticatedCallsCompleted           = "AuthenticatedCallsCompleted"
	AuthenticatedPartiallyRegisteredPhone = "AuthenticatedPartiallyRegisteredPhone"
	AuthenticatedRegisteredPhones         = "AuthenticatedRegisteredPhones"

	BRIChannelsActive = "BRIChannelsActive"
	BRISpansInService = "BRISpansInService"

	CallManagerHeartBeat = "CallManagerHeartBeat"
	CallsActive          = "CallsActive"
	CallsAttempted       = "CallsAttempted"
	CallsCompleted       = "CallsCompleted"
	CallsInProgress      = "CallsInProgress"

	CumulativeAllocatedResourceCannotOpenPort = "CumulativeAllocatedResourceCannotOpenPort"

	EncryptedCallsActive               = "EncryptedCallsActive"
	EncryptedCallsCompleted            = "EncryptedCallsCompleted"
	EncryptedPartiallyRegisteredPhones = "EncryptedPartiallyRegisteredPhones"
	EncryptedRegisteredPhones          = "EncryptedRegisteredPhones"

	ExternalCallControlEnabledCallsAttempted          = "ExternalCallControlEnabledCallsAttempted"
	ExternalCallControlEnabledCallsCompleted          = "ExternalCallControlEnabledCallsCompleted"
	ExternalCallControlEnabledFailureTreatmentApplied = "ExternalCallControlEnabledFailureTreatmentApplied"

	FXOPortsActive    = "FXOPortsActive"
	FXOPortsInService = "FXOPortsInService"
	FXSPortsActive    = "FXSPortsActive"
	FXSPortsInService = "FXSPortsInService"

	HWConferenceActive            = "HWConferenceActive"
	HWConferenceCompleted         = "HWConferenceCompleted"
	HWConferenceOutOfResources    = "HWConferenceOutOfResources"
	HWConferenceResourceActive    = "HWConferenceResourceActive"
	HWConferenceResourceAvailable = "HWConferenceResourceAvailable"
	HWConferenceResourceTotal     = "HWConferenceResourceTotal"

	HuntListsInService   = "HuntListsInService"
	IVROutOfResources    = "IVROutOfResources"
	IVRResourceActive    = "IVRResourceActive"
	IVRResourceAvailable = "IVRResourceAvailable"
	IVRResourceTotal     = "IVRResourceTotal"

	InitializationState    = "InitializationState"
	LocationOutOfResources = "LocationOutOfResources"

	MCUConferencesActive          = "MCUConferencesActive"
	MCUConferencesCompleted       = "MCUConferencesCompleted"
	MCUHttpConnectionErrors       = "MCUHttpConnectionErrors"
	MCUHttpNon200OkResponse       = "MCUHttpNon200OkResponse"
	MCUOutOfResources             = "MCUOutOfResources"
	MOHMulticastResourceActive    = "MOHMulticastResourceActive"
	MOHMulticastResourceAvailable = "MOHMulticastResourceAvailable"
	MOHOutOfResources             = "MOHOutOfResources"
	MOHTotalMulticastResources    = "MOHTotalMulticastResources"
	MOHTotalUnicastResources      = "MOHTotalUnicastResources"
	MOHUnicastResourceActive      = "MOHUnicastResourceActive"
	MOHUnicastResourceAvailable   = "MOHUnicastResourceAvailable"

	MTPOutOfResources    = "MTPOutOfResources"
	MTPRequestsThrottled = "MTPRequestsThrottled"
	MTPResourceActive    = "MTPResourceActive"
	MTPResourceAvailable = "MTPResourceAvailable"
	MTPResourceTotal     = "MTPResourceTotal"
)

var (
	AllowedGroupNames   = []string{"Cisco CallManager", "Cisco Recording"}
	AllowedCounterNames = []AllowedCounters{
		{allowedCounterName: CallsActive, prometheusName: "cucm_calls_active", defaultEnabled: true},
		{allowedCounterName: CallsInProgress, prometheusName: "cucm_calls_in_progress", defaultEnabled: true},
		{allowedCounterName: CallsCompleted, prometheusName: "cucm_calls_completed", defaultEnabled: true},
		{allowedCounterName: PartiallyRegisteredPhone, prometheusName: "cucm_partially_registered_phone", defaultEnabled: true},
		{allowedCounterName: RegisteredHardwarePhones, prometheusName: "cucm_registered_hardware_phones", defaultEnabled: true},
		{allowedCounterName: GatewaysSessionsActive, prometheusName: "cucm_gateways_sessions_active", defaultEnabled: true},
		{allowedCounterName: GatewaysSessionsFailed, prometheusName: "cucm_gateways_sessions_failed", defaultEnabled: true},
		{allowedCounterName: PhoneSessionsActive, prometheusName: "cucm_phone_sessions_active", defaultEnabled: true},
		{allowedCounterName: PhoneSessionsFailed, prometheusName: "cucm_phone_sessions_failed", defaultEnabled: true},
	}
)

//func (a *AllowedCounters) inCounter(name string) bool {
//	return name == a.allowedCounterName
//}

func isNameInAllowedCounter(name string) bool {
	for _, v := range AllowedCounterNames {
		if v.allowedCounterName == name {
			return true
		}
	}
	return false
}

//func getPrometheusName(name string) string {
//	for _, v := range AllowedCounterNames {
//		if v.allowedCounterName == name {
//			return v.prometheusName
//		}
//	}
//	return ""
//}
