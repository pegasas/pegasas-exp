package icm

func NewMitigateIncident() map[string]interface{} {
	mitigateParameter := make(map[string]interface{})
	mitigateParameter["IsCustomerImpacting"]= "True"
	mitigateParameter["IsNoise"]= "False"
	mitigateParameter["Mitigation"]= "[From ADFExtJobScheduler]: Auto Mitigated."
	mitigateParameter["HowFixed"]= "Fixed By Automation"
	mitigateParameter["MitigateContactAlias"]= "chenghaochen"
	return map[string]interface{} {
        "MitigateParameters" : mitigateParameter,
    }
}