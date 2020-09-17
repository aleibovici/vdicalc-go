package validation

import (
	"vdicalc/config"
	"vdicalc/functions"
)

var (
	hostResults           config.HostResultsConfiguration
	storageResults        config.StorageResultsConfiguration
	virtualizationResults config.VirtualizationResultsConfiguration
	vmResults             config.VMConfigurations
)

func init() {

	hostResults.Hostclockusedlimit = "4.2"                     // Maximum Intel CPU clock rate available
	storageResults.Storagedatastorecountlimit = "500"          // VMs per VMFS Datastore (https://configmax.vmware.com/guest?vmwareproduct=Horizon&release=Horizon%208%202006&categories=50-0)
	virtualizationResults.Managementservercountlimit = "12000" // VMware vCenter maximum (https://configmax.vmware.com/guest?vmwareproduct=Horizon&release=Horizon%208%202006&categories=46-0)
	vmResults.Memorysizelimit = "6128000"                      // RAM per VM (https://configmax.vmware.com/guest?vmwareproduct=vSphere&release=vSphere%207.0&categories=1-0)

}

// ValidateHostResults export
/* This public function validate input field and calucation results, raising errors using ErrorResultsConfiguration */
func ValidateHostResults(hostresultsclockused, virtualizationmanagementservertvmcount, memorysize, storagedatastorevmcount interface{}) config.ErrorResultsConfiguration {

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

	/* Validate Storage */
	if functions.StrtoFloat64(storagedatastorevmcount.(string)) > functions.StrtoFloat64(storageResults.Storagedatastorecountlimit) {
		error := config.ErrorConfiguration{Code: "Warning: ", Description: "Number os datastores above limit."}
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
