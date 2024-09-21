package main

type Counters struct {
	allowedCounterName string
	prometheusName     string
	defaultEnabled     bool
}

const (
	AnnunciatorOutOfResources    = "AnnunciatorOutOfResources"
	AnnunciatorResourceActive    = "AnnunciatorResourceActive"
	AnnunciatorResourceAvailable = "AnnunciatorResourceAvailable"
	AnnunciatorResourceTotal     = "AnnunciatorResourceTotal"

	PartiallyRegisteredPhone = "PartiallyRegisteredPhone"
	RegisteredHardwarePhones = "RegisteredHardwarePhones"

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

	SWConferenceActive            = "SWConferenceActive"
	SWConferenceCompleted         = "SWConferenceCompleted"
	SWConferenceOutOfResources    = "SWConferenceOutOfResources"
	SWConferenceResourceActive    = "SWConferenceResourceActive"
	SWConferenceResourceAvailable = "SWConferenceResourceAvailable"
	SWConferenceResourceTotal     = "SWConferenceResourceTotal"

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

	RegisteredBOTJabberMRA    = "RegisteredBOTJabberMRA"
	RegisteredBOTJabberNonMRA = "RegisteredBOTJabberNonMRA"
	RegisteredCSFJabberMRA    = "RegisteredCSFJabberMRA"
	RegisteredCSFJabberNonMRA = "RegisteredCSFJabberNonMRA"
	RegisteredTABJabberMRA    = "RegisteredTABJabberMRA"
	RegisteredTABJabberNonMRA = "RegisteredTABJabberNonMRA"
	RegisteredTCTJabberMRA    = "RegisteredTCTJabberMRA"
	RegisteredTCTJabberNonMRA = "RegisteredTCTJabberNonMRA"

	// SIPLineServerAuthorizationChallenges SIP
	SIPLineServerAuthorizationChallenges     = "SIPLineServerAuthorizationChallenges"
	SIPLineServerAuthorizationFailures       = "SIPLineServerAuthorizationFailures"
	SIPTrunkApplicationAuthorizationFailures = "SIPTrunkApplicationAuthorizationFailures"
	SIPTrunkApplicationAuthorizations        = "SIPTrunkApplicationAuthorizations"
	SIPTrunkAuthorizationFailures            = "SIPTrunkAuthorizationFailures"
	SIPTrunkAuthorizations                   = "SIPTrunkAuthorizations"
	SIPTrunkServerAuthenticationChallenges   = "SIPTrunkServerAuthenticationChallenges"
	SIPTrunkServerAuthenticationFailures     = "SIPTrunkServerAuthenticationFailures"
	SystemCallsAttempted                     = "SystemCallsAttempted"

	// TranscoderOutOfResources transcoder
	TranscoderOutOfResources    = "TranscoderOutOfResources"
	TranscoderRequestsThrottled = "TranscoderRequestsThrottled"
	TranscoderResourceActive    = "TranscoderResourceActive"
	TranscoderResourceAvailable = "TranscoderResourceAvailable"
	TranscoderResourceTotal     = "TranscoderResourceTotal"
	UnEncryptedCallFailures     = "UnEncryptedCallFailures"

	// VideoCallsActive video
	VideoCallsActive          = "VideoCallsActive"
	VideoCallsCompleted       = "VideoCallsCompleted"
	VideoOnHoldOutOfResources = "VideoOnHoldOutOfResources"
	VideoOnHoldResourceActive = "VideoOnHoldResourceActive"
	VideoOutOfResources       = "VideoOutOfResources"

	// GatewayRegistrationFailures recording
	GatewayRegistrationFailures = "GatewayRegistrationFailures"
	GatewaysInService           = "GatewaysInService"
	GatewaysOutOfService        = "GatewaysOutOfService"
	GatewaysSessionsActive      = "GatewaysSessionsActive"
	GatewaysSessionsFailed      = "GatewaysSessionsFailed"
	PhoneSessionsActive         = "PhoneSessionsActive"
	PhoneSessionsFailed         = "PhoneSessionsFailed"

	// RegisteredAnalogAccess registered devices
	RegisteredAnalogAccess        = "RegisteredAnalogAccess"
	RegisteredMGCPGateway         = "RegisteredMGCPGateway"
	RegisteredOtherStationDevices = "RegisteredOtherStationDevices"
)

var (
	AllowedGroupNames = []string{"Cisco CallManager", "Cisco Recording"}
	SupportedCounters = []Counters{
		// CM - basic
		{allowedCounterName: CallsActive, prometheusName: "cucm_calls_active", defaultEnabled: true},
		{allowedCounterName: CallsAttempted, prometheusName: "cucm_calls_attempted", defaultEnabled: false},
		{allowedCounterName: CallsInProgress, prometheusName: "cucm_calls_in_progress", defaultEnabled: true},
		{allowedCounterName: CallsCompleted, prometheusName: "cucm_calls_completed", defaultEnabled: true},
		{allowedCounterName: PartiallyRegisteredPhone, prometheusName: "cucm_partially_registered_phone", defaultEnabled: true},
		{allowedCounterName: RegisteredHardwarePhones, prometheusName: "cucm_registered_hardware_phones", defaultEnabled: true},
		{allowedCounterName: SystemCallsAttempted, prometheusName: "cucm_system_calls_attempted", defaultEnabled: false},
		{allowedCounterName: UnEncryptedCallFailures, prometheusName: "cucm_un_encrypted_call_failures", defaultEnabled: false},
		// CM - annunciator
		{allowedCounterName: AnnunciatorOutOfResources, prometheusName: "cucm_annunciator_out_of_resources", defaultEnabled: false},
		{allowedCounterName: AnnunciatorResourceActive, prometheusName: "cucm_annunciator_resource_active", defaultEnabled: false},
		{allowedCounterName: AnnunciatorResourceAvailable, prometheusName: "cucm_annunciator_resource_available", defaultEnabled: false},
		{allowedCounterName: AnnunciatorResourceTotal, prometheusName: "cucm_annunciator_resource_total", defaultEnabled: false},
		{allowedCounterName: AuthenticatedCallsActive, prometheusName: "cucm_authenticated_calls_active", defaultEnabled: false},
		{allowedCounterName: AuthenticatedCallsCompleted, prometheusName: "cucm_authenticated_calls_completed", defaultEnabled: false},
		{allowedCounterName: AuthenticatedPartiallyRegisteredPhone, prometheusName: "cucm_authenticated_partially_registeredPhone", defaultEnabled: false},
		{allowedCounterName: AuthenticatedRegisteredPhones, prometheusName: "cucm_authenticated_registered_phones", defaultEnabled: false},
		//{allowedCounterName: BRIChannelsActive, prometheusName: "cucm_bri_channels_active", defaultEnabled: false},
		//{allowedCounterName: BRISpansInService, prometheusName: "cucm_bri_spans_in_service", defaultEnabled: false},
		{allowedCounterName: CallManagerHeartBeat, prometheusName: "cucm_call_manager_heart_beat", defaultEnabled: false},
		{allowedCounterName: CumulativeAllocatedResourceCannotOpenPort, prometheusName: "cucm_cumulative_allocated_resource_cannot_open_port", defaultEnabled: false},
		{allowedCounterName: EncryptedCallsActive, prometheusName: "cucm_encrypted_calls_active", defaultEnabled: false},
		{allowedCounterName: EncryptedCallsCompleted, prometheusName: "cucm_encrypted_calls_completed", defaultEnabled: false},
		{allowedCounterName: EncryptedPartiallyRegisteredPhones, prometheusName: "cucm_encrypted_partially_registered_phones", defaultEnabled: false},
		{allowedCounterName: EncryptedRegisteredPhones, prometheusName: "cucm_encrypted_registered_phones", defaultEnabled: false},
		// mtp
		{allowedCounterName: MTPOutOfResources, prometheusName: "cucm_mtp_out_of_resources", defaultEnabled: false},
		{allowedCounterName: MTPRequestsThrottled, prometheusName: "cucm_mtp_requests_throttled", defaultEnabled: false},
		{allowedCounterName: MTPResourceActive, prometheusName: "cucm_mtp_resource_active", defaultEnabled: false},
		{allowedCounterName: MTPResourceAvailable, prometheusName: "cucm_mtp_resource_available", defaultEnabled: false},
		{allowedCounterName: MTPResourceTotal, prometheusName: "cucm_mtp_resource_total", defaultEnabled: false},
		//sip
		{allowedCounterName: SIPLineServerAuthorizationChallenges, prometheusName: "cucm_sip_line_server_authorization_challenges", defaultEnabled: false},
		{allowedCounterName: SIPLineServerAuthorizationFailures, prometheusName: "cucm_sip_line_server_authorization_failures", defaultEnabled: false},
		{allowedCounterName: SIPTrunkApplicationAuthorizationFailures, prometheusName: "cucm_sip_trunk_application_authorization_failures", defaultEnabled: false},
		{allowedCounterName: SIPTrunkApplicationAuthorizations, prometheusName: "cucm_sip_trunk_application_authorizations", defaultEnabled: false},
		{allowedCounterName: SIPTrunkAuthorizationFailures, prometheusName: "cucm_sip_trunk_authorization_failures", defaultEnabled: false},
		{allowedCounterName: SIPTrunkAuthorizations, prometheusName: "cucm_sip_trunk_authorizations", defaultEnabled: false},
		{allowedCounterName: SIPTrunkServerAuthenticationChallenges, prometheusName: "cucm_sip_trunk_server_authentication_challenges", defaultEnabled: false},
		{allowedCounterName: SIPTrunkServerAuthenticationFailures, prometheusName: "cucm_sip_trunk_server_authentication_failures", defaultEnabled: false},
		//transcoder
		{allowedCounterName: TranscoderOutOfResources, prometheusName: "cucm_transcoder_out_of_resources", defaultEnabled: false},
		{allowedCounterName: TranscoderRequestsThrottled, prometheusName: "cucm_transcoder_requests_throttled", defaultEnabled: false},
		{allowedCounterName: TranscoderResourceActive, prometheusName: "cucm_transcoder_resource_active", defaultEnabled: false},
		{allowedCounterName: TranscoderResourceAvailable, prometheusName: "cucm_transcoder_resource_available", defaultEnabled: false},
		{allowedCounterName: TranscoderResourceTotal, prometheusName: "cucm_transcoder_resource_total", defaultEnabled: false},
		// video
		{allowedCounterName: VideoCallsActive, prometheusName: "cucm_video_calls_active", defaultEnabled: false},
		{allowedCounterName: VideoCallsCompleted, prometheusName: "cucm_video_calls_completed", defaultEnabled: false},
		{allowedCounterName: VideoOnHoldOutOfResources, prometheusName: "cucm_video_on_hold_out_of_resources", defaultEnabled: false},
		{allowedCounterName: VideoOnHoldResourceActive, prometheusName: "cucm_video_on_hold_resource_active", defaultEnabled: false},
		{allowedCounterName: VideoOutOfResources, prometheusName: "cucm_video_out_of_resources", defaultEnabled: false},
		// jabber
		{allowedCounterName: RegisteredBOTJabberMRA, prometheusName: "cucm_registered_bot_jabber_mra", defaultEnabled: false},
		{allowedCounterName: RegisteredBOTJabberNonMRA, prometheusName: "cucm_registered_bot_jabber_non_mra", defaultEnabled: false},
		{allowedCounterName: RegisteredCSFJabberMRA, prometheusName: "cucm_registered_csf_jabber_mra", defaultEnabled: false},
		{allowedCounterName: RegisteredCSFJabberNonMRA, prometheusName: "cucm_registered_csf_jabber_non_mra", defaultEnabled: false},
		{allowedCounterName: RegisteredTABJabberMRA, prometheusName: "cucm_registered_tab_jabber_mra", defaultEnabled: false},
		{allowedCounterName: RegisteredTABJabberNonMRA, prometheusName: "cucm_registered_tab_jabber_non_mra", defaultEnabled: false},
		{allowedCounterName: RegisteredTCTJabberMRA, prometheusName: "cucm_registered_tct_jabber_mra", defaultEnabled: false},
		{allowedCounterName: RegisteredTCTJabberNonMRA, prometheusName: "cucm_registered_tct_jabber_non_mra", defaultEnabled: false},
		// cisco recording
		{allowedCounterName: GatewayRegistrationFailures, prometheusName: "cucm_gateway_registration_failures", defaultEnabled: false},
		{allowedCounterName: GatewaysInService, prometheusName: "cucm_gateways_in_service", defaultEnabled: false},
		{allowedCounterName: GatewaysOutOfService, prometheusName: "cucm_gateways_out_of_service", defaultEnabled: false},
		{allowedCounterName: GatewaysSessionsActive, prometheusName: "cucm_gateways_sessions_active", defaultEnabled: true},
		{allowedCounterName: GatewaysSessionsFailed, prometheusName: "cucm_gateways_sessions_failed", defaultEnabled: true},
		{allowedCounterName: PhoneSessionsActive, prometheusName: "cucm_phone_sessions_active", defaultEnabled: true},
		{allowedCounterName: PhoneSessionsFailed, prometheusName: "cucm_phone_sessions_failed", defaultEnabled: true},
		// HW Conference
		{allowedCounterName: HWConferenceActive, prometheusName: "cucm_hw_conference_active", defaultEnabled: false},
		{allowedCounterName: HWConferenceCompleted, prometheusName: "cucm_hw_conference_completed", defaultEnabled: false},
		{allowedCounterName: HWConferenceOutOfResources, prometheusName: "cucm_hw_conference_out_of_resources", defaultEnabled: false},
		{allowedCounterName: HWConferenceResourceActive, prometheusName: "cucm_hw_conference_resource_active", defaultEnabled: false},
		{allowedCounterName: HWConferenceResourceAvailable, prometheusName: "cucm_hw_conference_resource_available", defaultEnabled: false},
		{allowedCounterName: HWConferenceResourceTotal, prometheusName: "cucm_hw_conference_resource_total", defaultEnabled: false},
		// SW Conference
		{allowedCounterName: SWConferenceActive, prometheusName: "cucm_sw_conference_active", defaultEnabled: false},
		{allowedCounterName: SWConferenceCompleted, prometheusName: "cucm_sw_conference_completed", defaultEnabled: false},
		{allowedCounterName: SWConferenceOutOfResources, prometheusName: "cucm_sw_conference_out_of_resources", defaultEnabled: false},
		{allowedCounterName: SWConferenceResourceActive, prometheusName: "cucm_sw_conference_resource_active", defaultEnabled: false},
		{allowedCounterName: SWConferenceResourceAvailable, prometheusName: "cucm_sw_conference_resource_available", defaultEnabled: false},
		{allowedCounterName: SWConferenceResourceTotal, prometheusName: "cucm_sw_conference_resource_total", defaultEnabled: false},
		// registered info
		{allowedCounterName: RegisteredAnalogAccess, prometheusName: "cucm_registered_analog_access", defaultEnabled: false},
		{allowedCounterName: RegisteredMGCPGateway, prometheusName: "cucm_registered_mgcp_gateway", defaultEnabled: false},
		{allowedCounterName: RegisteredOtherStationDevices, prometheusName: "cucm_registered_other_station_devices", defaultEnabled: false},
	}
)

// isNameInAllowedCounter is name part of supported metrics
func isNameInAllowedCounter(name string) bool {
	for _, v := range SupportedCounters {
		if v.allowedCounterName == name {
			return true
		}
	}
	return false
}
