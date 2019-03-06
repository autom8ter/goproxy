package twilio

import (
	"github.com/autom8ter/goproxy"
	"os"
)

var BASE = "https://api.twilio.com/2010-04-01"
var APIKEY = os.Getenv("TWILIO_API_KEY")
var APISECRET = os.Getenv("TWILIO_API_SECRET")
var ACCOUNTSID = os.Getenv("TWILIO_ACCOUNT_SID")
var ASSISTANTSID = os.Getenv("TWILIO_ASSISTANT_SID")
var FLOWSID = os.Getenv("TWILIO_FLOW_SID")

var Messages = &goproxy.ProxyConfig{
	PathPrefix: "/Messages.json",
	TargetUrl:  BASE + "/Accounts/" + ACCOUNTSID,
	Username:   APIKEY,
	Password:   APISECRET,
}

var Voice = &goproxy.ProxyConfig{
	PathPrefix: "/Voice.json",
	TargetUrl:  BASE + "/Accounts/" + ACCOUNTSID,
	Username:   APIKEY,
	Password:   APISECRET,
}

var Accounts = &goproxy.ProxyConfig{
	PathPrefix: "/Accounts.json",
	TargetUrl:  BASE,
	Username:   APIKEY,
	Password:   APISECRET,
}

var AutoPilotAssistants = &goproxy.ProxyConfig{
	PathPrefix: "/Assistants.json",
	TargetUrl:  "https://autopilot.twilio.com/v1",
	Username:   APIKEY,
	Password:   APISECRET,
}

var AutoPilotFieldTypes = &goproxy.ProxyConfig{
	PathPrefix: "/FieldTypes.json",
	TargetUrl:  "https://autopilot.twilio.com/v1/Assistants/" + ASSISTANTSID,
	Username:   APIKEY,
	Password:   APISECRET,
}

var AutoPilotTasks = &goproxy.ProxyConfig{
	PathPrefix: "/Tasks.json",
	TargetUrl:  "https://autopilot.twilio.com/v1/Assistants/" + ASSISTANTSID,
	Username:   APIKEY,
	Password:   APISECRET,
}

var AutoPilotModelBuilds = &goproxy.ProxyConfig{
	PathPrefix: "/ModelBuilds.json",
	TargetUrl:  "https://autopilot.twilio.com/v1/Assistants/" + ASSISTANTSID,
	Username:   APIKEY,
	Password:   APISECRET,
}

var AutoPilotQueries = &goproxy.ProxyConfig{
	PathPrefix: "/Queries.json",
	TargetUrl:  "https://autopilot.twilio.com/v1/Assistants/" + ASSISTANTSID,
	Username:   APIKEY,
	Password:   APISECRET,
}

var AutoPilotDefaults = &goproxy.ProxyConfig{
	PathPrefix: "/Defaults.json",
	TargetUrl:  "https://autopilot.twilio.com/v1/Assistants/" + ASSISTANTSID,
	Username:   APIKEY,
	Password:   APISECRET,
}

var StudioFlows = &goproxy.ProxyConfig{
	PathPrefix: "/Flows.json",
	TargetUrl:  "https://studio.twilio.com/v1",
	Username:   APIKEY,
	Password:   APISECRET,
}
var StudioFlowExecutions = &goproxy.ProxyConfig{
	PathPrefix: "/Executions.json",
	TargetUrl:  "https://studio.twilio.com/v1/Flows/" + FLOWSID,
	Username:   APIKEY,
	Password:   APISECRET,
}

var All = []*goproxy.ProxyConfig{
	Messages,
	Voice,
	Accounts,
	AutoPilotAssistants,
	AutoPilotDefaults,
	AutoPilotFieldTypes,
	AutoPilotModelBuilds,
	AutoPilotQueries,
	AutoPilotTasks,
	StudioFlows,
	StudioFlowExecutions,
}
