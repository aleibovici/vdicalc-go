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

	vmResults.Memorysizelimit = "6128000"                      // RAM per VM (https://configmax.vmware.com/guest?vmwareproduct=vSphere&release=vSphere%207.0&categories=1-0)
	hostResults.Hostclockusedlimit = "4.2"                     // Maximum Intel CPU clock rate available
	hostResults.Hostvmcountlimit = "200"                       // Virtual machines per host (https://configmax.vmware.com/guest?vmwareproduct=Horizon&release=Horizon%208%202006&categories=2-0)
	storageResults.Storagedatastorecountlimit = "500"          // VMs per VMFS Datastore (https://configmax.vmware.com/guest?vmwareproduct=Horizon&release=Horizon%208%202006&categories=50-0)
	virtualizationResults.Managementservercountlimit = "12000" // VMware vCenter maximum (https://configmax.vmware.com/guest?vmwareproduct=Horizon&release=Horizon%208%202006&categories=46-0)

}

// ValidateResults export
/* This public function validate inputs and results, raising errors using ErrorResultsConfiguration */
func ValidateResults(fullData map[string]interface{}) config.ErrorResultsConfiguration {

	var errorList config.ErrorResultsConfiguration

	for key, values := range fullData {

		switch key {
		case "vmmemorysizeselected":

			if functions.StrtoFloat64(values.(string)) > functions.StrtoFloat64(vmResults.Memorysizelimit) {
				error := config.ErrorConfiguration{Code: "Warning: ", Description: "VM memory size above limit."}
				errorList.Error = append(errorList.Error, error)
				return errorList
			}

		case "hostresultsclockused":

			if functions.StrtoFloat64(values.(string)) > functions.StrtoFloat64(hostResults.Hostclockusedlimit) {
				error := config.ErrorConfiguration{Code: "Warning: ", Description: "Host CPU (GHz) above limit. (max=4.2)"}
				errorList.Error = append(errorList.Error, error)
				return errorList
			}

		case "hostresultsvmcount":

			if functions.StrtoFloat64(values.(string)) > functions.StrtoFloat64(hostResults.Hostvmcountlimit) {
				error := config.ErrorConfiguration{Code: "Warning: ", Description: "Number os VMs per host avobe limit. (max=200)"}
				errorList.Error = append(errorList.Error, error)
				return errorList
			}

		case "storageresultsdatastorecount":

			if functions.StrtoFloat64(values.(string)) > functions.StrtoFloat64(storageResults.Storagedatastorecountlimit) {
				error := config.ErrorConfiguration{Code: "Warning: ", Description: "Number os datastores above limit (max=500)."}
				errorList.Error = append(errorList.Error, error)
				return errorList
			}
		}

	}

	return errorList

}
