package calculations

import (
	"net/http"
	"vdicalc/host"
	"vdicalc/storage"
	"vdicalc/validation"
	"vdicalc/virtualization"
)

// Calculate export
/* This public function is reponsible for executing all necessary calculations to obtain results */
func Calculate(fullData map[string]interface{}, r *http.Request) {

	fullData["hostresultscount"] = host.GetHostCount(r.FormValue("vmcount"), r.FormValue("hostsocketcount"), r.FormValue("hostsocketcorescount"), r.FormValue("vmpercorecount"), r.FormValue("hostcoresoverhead"), r.FormValue("virtualizationclusterhostha"))
	fullData["hostresultsclockused"] = host.GetHostClockUsed(r.FormValue("vmvcpucount"), r.FormValue("vmvcpumhz"), r.FormValue("vmcount"), r.FormValue("hostsocketcount"), r.FormValue("hostsocketcorescount"), r.FormValue("vmpercorecount"), r.FormValue("hostcoresoverhead"))
	fullData["hostresultsmemory"] = host.GetHostMemory(r.FormValue("vmcount"), r.FormValue("hostsocketcount"), r.FormValue("hostsocketcorescount"), r.FormValue("hostcoresoverhead"), r.FormValue("vmpercorecount"), r.FormValue("vmmemorysize"), r.FormValue("hostmemoryoverhead"), r.FormValue("vmdisplaycount"), r.FormValue("vmdisplayresolution"), r.FormValue("vmvcpucount"), r.FormValue("vmvideoram"))
	fullData["hostresultsvmcount"] = host.GetHostVMCount(r.FormValue("vmcount"), r.FormValue("hostsocketcount"), r.FormValue("hostsocketcorescount"), r.FormValue("vmpercorecount"), r.FormValue("hostcoresoverhead"))
	fullData["storageresultscapacity"] = storage.GetStorageCapacity(r.FormValue("vmcount"), r.FormValue("vmdisksize"), r.FormValue("storagecapacityoverhead"), r.FormValue("storagededuperatio"), r.FormValue("vmdisplaycount"), r.FormValue("vmdisplayresolution"), r.FormValue("vmvideoram"), r.FormValue("vmmemorysize"), r.FormValue("vmclonesizerefreshrate"))
	fullData["storageresultsdatastorecount"] = storage.GetStorageDatastoreCount(r.FormValue("vmcount"), r.FormValue("storagedatastorevmcount"))
	fullData["storageresultsdatastoresize"] = storage.GetStorageDatastoreSize(r.FormValue("vmcount"), r.FormValue("storagedatastorevmcount"), r.FormValue("vmdisksize"), r.FormValue("storagecapacityoverhead"), r.FormValue("storagededuperatio"), r.FormValue("vmdisplaycount"), r.FormValue("vmdisplayresolution"), r.FormValue("vmvideoram"), r.FormValue("vmmemorysize"), r.FormValue("vmclonesizerefreshrate"))
	fullData["storagedatastorefroentendiops"], fullData["storagedatastorebackendiops"], fullData["storageresultsfrontendiops"], fullData["storageresultsbackendiops"] = storage.GetStorageDatastoreIops(r.FormValue("vmiopscount"), r.FormValue("vmiopsreadratio"), r.FormValue("storagedatastorevmcount"), r.FormValue("storageraidtype"), r.FormValue("vmcount"), r.FormValue("storagedatastorevmcount"))
	fullData["virtualizationresultsclustercount"] = virtualization.GetClusterSize(r.FormValue("vmcount"), r.FormValue("hostsocketcount"), r.FormValue("hostsocketcorescount"), r.FormValue("vmpercorecount"), r.FormValue("hostcoresoverhead"), r.FormValue("virtualizationclusterhostsize"), r.FormValue("virtualizationclusterhostha"))
	fullData["virtualizationresultsmanagementservercount"] = virtualization.GetManagementServerCount(r.FormValue("vmcount"), r.FormValue("virtualizationmanagementservertvmcount"))
	fullData["errorresults"] = validation.ValidateResults(fullData)

}
