package icm

import (
	"encoding/xml"
	"k8s-delegater/pkg/common"
	"strings"
)

func NewCreateIncident() map[string]interface{} {
	dictionary := make(map[string]interface{})
	dictionary["CommitDate"] = nil
	dictionary["Component"] = nil
	dictionary["CorrelationId"] = nil
	dictionary["CustomFields"] = nil
	dictionary["CustomerName"] = nil
	dictionary["Description"] = nil
	dictionary["DescriptionEntries"] = []interface{}{NewDescriptionEntry()}
	dictionary["ExtendedData"] = nil
	dictionary["HowFixed"] = nil
	dictionary["ImpactStartDate"] = nil
	dictionary["ImpactedServices"] = nil
	dictionary["ImpactedTeams"] = nil
	dictionary["IncidentSubType"] = nil
	dictionary["IncidentType"] = nil
	dictionary["IsCustomerImpacting"] = nil
	dictionary["IsNoise"] = nil
	dictionary["IsSecurityRisk"] = nil
	dictionary["Keywords"] = nil
	dictionary["MitigatedDate"] = nil
	dictionary["Mitigation"] = nil
	dictionary["MonitorId"] = "NONE://Default"
	dictionary["OccurringLocation"] = NewOccurringLocation()
	dictionary["OwningAlias"] = nil
	dictionary["OwningContactFullName"] = nil
	dictionary["RaisingLocation"] = NewRaisingLocation()
	dictionary["ReproSteps"] = nil
	dictionary["ResolutionDate"] = nil
	dictionary["RoutingId"] = "NONE://Default"
	dictionary["ServiceResponsible"] = nil
	dictionary["Severity"] = nil
	dictionary["Source"] = NewSource()
	dictionary["Status"] = "Active"
	dictionary["SubscriptionId"] = nil
	dictionary["SupportTicketId"] = nil
	dictionary["Title"] = nil
	dictionary["TrackingTeams"] = nil
	dictionary["TsgId"] = nil
	dictionary["TsgOutput"] = nil
	dictionary["ValueSpecifiedFields"] = "None"
	return dictionary
}

func NewDescriptionEntry() map[string]interface{} {
	dict := make(map[string]interface{})
	dict["Cause"] = "Other"
	dict["ChangedBy"] = nil
	dict["Date"] = "$now$"
	dict["DescriptionEntryId"] = 0
	dict["RenderType"] = "Html"
	dict["SubmitDate"] = "$now$"
	dict["SubmittedBy"] = "ADFExtJobScheduler"
	dict["Text"] = "Required"

	descriptionEntry := make(map[string]interface{})
	descriptionEntry["DescriptionEntry"] = dict

	return descriptionEntry
}

func NewOccurringLocation() map[string]interface{} {
	dict := make(map[string]interface{})
	dict["DataCenter"] = nil
	dict["DeviceGroup"] = nil
	dict["DeviceName"] = nil
	dict["Environment"] = strings.ToUpper(common.GetEnv())
	dict["ServiceInstanceId"] = nil
	return dict
}

func NewRaisingLocation() map[string]interface{} {
	dict := make(map[string]interface{})
	dict["DataCenter"] = nil
	dict["DeviceGroup"] = nil
	dict["DeviceName"] = nil
	dict["Environment"] = nil
	dict["ServiceInstanceId"] = nil
	return dict
}

func NewSource() map[string]interface{} {
	dict := make(map[string]interface{})
	dict["CreateDate"] = "$now$"
	dict["CreatedBy"] = "Monitor"
	dict["IncidentId"] = "$uuid$"
	dict["ModifiedDate"] = "$now$"
	dict["Origin"] = "Monitor"
	dict["Revision"] = nil
	dict["SourceId"] = "00000000-0000-0000-0000-000000000000"
	return dict
}

type CreateIncidentResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    Body     `xml:"Body"`
}

type Body struct {
	AddOrUpdateIncident2Response AddOrUpdateIncident2Response `xml:"AddOrUpdateIncident2Response"`
}

type AddOrUpdateIncident2Response struct {
	AddOrUpdateIncident2Result AddOrUpdateIncident2Result `xml:"AddOrUpdateIncident2Result"`
}

type AddOrUpdateIncident2Result struct {
	IncidentId int    `xml:"IncidentId"`
	Status     string `xml:"Status"`
}
