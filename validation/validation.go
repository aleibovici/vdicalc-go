package validation

import (
	"vdicalc/config"
	"vdicalc/functions"
)

var (
	hostResults           config.HostResultsConfiguration
	virtualizationResults config.VirtualizationResultsConfiguration
	vmResults             config.VMConfigurations
)

func init() {

	hostResults.Hostclockusedlimit = "4.2"
	hostResults.Hostmemorylimit = "222"
	virtualizationResults.Managementservercountlimit = "10000"
	vmResults.Memorysizelimit = "1048576"

}

// ValidateHostResults export
func ValidateHostResults(hostresultsclockused, virtualizationmanagementservertvmcount, memorysize interface{}) config.ErrorResultsConfiguration {

	var errorList config.ErrorResultsConfiguration

	/* Validate VM */
	if functions.StrtoFloat64(memorysize.(string)) > functions.StrtoFloat64(vmResults.Memorysizelimit) {
		error := config.ErrorConfiguration{Code: "Warning: ", Description: "VM memory size above limit."}
		errorList.Error = append(errorList.Error, error)
		return errorList
	}

	/* Validate Host */
	if functions.StrtoFloat64(hostresultsclockused.(string)) > functions.StrtoFloat64(hostResults.Hostclockusedlimit) {
		error := config.ErrorConfiguration{Code: "Warning: ", Description: "Host CPU (GHz) above limit."}
		errorList.Error = append(errorList.Error, error)
		return errorList
	}

	/* Validate Virtualization */
	if functions.StrtoFloat64(virtualizationmanagementservertvmcount.(string)) > functions.StrtoFloat64(virtualizationResults.Managementservercountlimit) {
		error := config.ErrorConfiguration{Code: "Warning: ", Description: "Number of VMs per management server above limit."}
		errorList.Error = append(errorList.Error, error)
		return errorList
	}

	return errorList

}
