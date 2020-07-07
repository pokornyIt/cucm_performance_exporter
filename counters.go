package main

type AllowedCounters struct {
	allowedCounterName string
	prometheusName     string
}

const (
	CallsActive              = "CallsActive"
	CallsInProgress          = "CallsInProgress"
	CallsCompleted           = "CallsCompleted"
	PartiallyRegisteredPhone = "PartiallyRegisteredPhone"
	RegisteredHardwarePhones = "RegisteredHardwarePhones"
	GatewaysSessionsActive   = "GatewaysSessionsActive"
	GatewaysSessionsFailed   = "GatewaysSessionsFailed"
	PhoneSessionsActive      = "PhoneSessionsActive"
	PhoneSessionsFailed      = "PhoneSessionsFailed"
)

var (
	AllowedGroupNames   = []string{"Cisco CallManager", "Cisco Recording"}
	AllowedCounterNames = []AllowedCounters{
		{allowedCounterName: CallsActive, prometheusName: "cucm_calls_active"},
		{allowedCounterName: CallsInProgress, prometheusName: "cucm_calls_in_progress"},
		{allowedCounterName: CallsCompleted, prometheusName: "cucm_calls_completed"},
		{allowedCounterName: PartiallyRegisteredPhone, prometheusName: "cucm_partially_registered_phone"},
		{allowedCounterName: RegisteredHardwarePhones, prometheusName: "cucm_registered_hardware_phones"},
		{allowedCounterName: GatewaysSessionsActive, prometheusName: "cucm_gateways_sessions_active"},
		{allowedCounterName: GatewaysSessionsFailed, prometheusName: "cucm_gateways_sessions_failed"},
		{allowedCounterName: PhoneSessionsActive, prometheusName: "cucm_phone_sessions_active"},
		{allowedCounterName: PhoneSessionsFailed, prometheusName: "cucm_phone_sessions_failed"},
	}
)

func (a *AllowedCounters) inCounter(name string) bool {
	return name == a.allowedCounterName
}

func isNameInAllowedCounter(name string) bool {
	for _, v := range AllowedCounterNames {
		if v.allowedCounterName == name {
			return true
		}
	}
	return false
}

func getPrometheusName(name string) string {
	for _, v := range AllowedCounterNames {
		if v.allowedCounterName == name {
			return v.prometheusName
		}
	}
	return ""

}
