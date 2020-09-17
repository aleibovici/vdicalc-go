package functions

/* This package contains external genericfunctions used across the app.
It must be imported using "vdicalc/functions" */

import (
	"log"
	"net/http"
	"os"
	"strconv"
	c "vdicalc/config"
)

// StrtoInt function
/* This public function convert string to int */
func StrtoInt(value string) int {
	r, err := strconv.Atoi(value)
	if err != nil {
	}
	return r
}

// InttoStr function
/* This public function convert int to string */
func InttoStr(value int) string {
	r := strconv.Itoa(value)
	return r
}

// StrtoFloat64 function
/* This public function convert string to float64 */
func StrtoFloat64(value string) float64 {
	r, _ := strconv.ParseFloat(value, 8)
	return r
}

// Float64toStr function
/* This public function convert float64 to string with variable precision */
func Float64toStr(value float64, prec int) string {

	return strconv.FormatFloat(value, 'f', prec, 64)
}

// MustGetenv is a helper function for getting environment variables.
// Displays a warning if the environment variable is not set.
func MustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.\n", k)
	}
	return v
}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request) string {

	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}

	return r.RemoteAddr
}

// DataLoad fuction
/* This function builds created the map interface with all items for the HTML content.
It expect type 'Configurations' as defined in 'config.go' */
func DataLoad(class c.Configurations) map[string]interface{} {

	data := map[string]interface{}{
		"title":    class.Variable.Title,
		"titlesub": class.Variable.Titlesub,
		"update":   class.Variable.Update,
		"about":    class.Variable.About,
		"print":    class.Variable.Print,
		"guide":    class.Variable.Guide,
		/* vm are used to display form items for hosts */
		"vmprofilelabel":                 class.VM.Vmprofilelabel,
		"vmprofile":                      class.VM.Vmprofile,
		"vmprofileselected":              class.VM.Vmprofileselected,
		"vmcountlabel":                   class.VM.Vmcountlabel,
		"vmcount":                        class.VM.Vmcount,
		"vmcountselected":                class.VM.Vmcountselected,
		"vmvcpucountlabel":               class.VM.Vcpucountlabel,
		"vmvcpucount":                    class.VM.Vcpucount,
		"vmvcpucountselected":            class.VM.Vcpucountselected,
		"vmvcpumhzlabel":                 class.VM.Vcpumhzlabel,
		"vmvcpumhz":                      class.VM.Vcpumhz,
		"vmvcpumhzselected":              class.VM.Vcpumhzselected,
		"vmvcpumhztooltip":               class.VM.Vcpumhztooltip,
		"vmpercorecountlabel":            class.VM.Vmpercorecountlabel,
		"vmpercorecount":                 class.VM.Vmpercorecount,
		"vmpercorecountselected":         class.VM.Vmpercorecountselected,
		"vmdisplaycountlabel":            class.VM.Displaycountlabel,
		"vmdisplaycount":                 class.VM.Displaycount,
		"vmdisplaycountselected":         class.VM.Displaycountselected,
		"vmdisplayresolutionlabel":       class.VM.Displayresolutionlabel,
		"vmdisplayresolution":            class.VM.Displayresolution,
		"vmdisplayresolutionselected":    class.VM.Displayresolutionselected,
		"vmmemorysizelabel":              class.VM.Memorysizelabel,
		"vmmemorysize":                   class.VM.Memorysize,
		"vmmemorysizeselected":           class.VM.Memorysizeselected,
		"vmvideoram":                     class.VM.Videoram,
		"vmvideoramlabel":                class.VM.Videoramlabel,
		"vmvideoramselected":             class.VM.Videoramselected,
		"vmdisksizelabel":                class.VM.Disksizelabel,
		"vmdisksize":                     class.VM.Disksize,
		"vmdisksizeselected":             class.VM.Disksizeselected,
		"vmiopscount":                    class.VM.Iopscount,
		"vmiopscountlabel":               class.VM.Iopscountlabel,
		"vmiopscountselected":            class.VM.Iopscountselected,
		"vmiopscounttooltip":             class.VM.Iopscounttooltip,
		"vmiopsreadratio":                class.VM.Iopsreadratio,
		"vmiopsreadratiolabel":           class.VM.Iopsreadratiolabel,
		"vmiopsreadratioselected":        class.VM.Iopsreadratioselected,
		"vmiopsreadratiotooltip":         class.VM.Iopsreadratiotooltip,
		"vmiopswriteratio":               class.VM.Iopswriteratio,
		"vmclonesizerefreshrate":         class.VM.Clonesizerefreshrate,
		"vmclonesizerefreshratetooltip":  class.VM.Clonesizerefreshratetooltip,
		"vmclonesizerefreshratelabel":    class.VM.Clonesizerefreshratelabel,
		"vmclonesizerefreshrateselected": class.VM.Clonesizerefreshrateselected,
		/* host are used to display form items for hosts */
		"hostsocketcountlabel":         class.Host.Socketcountlabel,
		"hostsocketcount":              class.Host.Socketcount,
		"hostsocketcountselected":      class.Host.Socketcountselected,
		"hostsocketcorescountlabel":    class.Host.Socketcorescountlabel,
		"hostsocketcorescount":         class.Host.Socketcorescount,
		"hostsocketcorescountselected": class.Host.Socketcorescountselected,
		"hostmemoryoverhead":           class.Host.Memoryoverhead,
		"hostmemoryoverheadselected":   class.Host.Memoryoverheadselected,
		"hostmemoryoverheadlabel":      class.Host.Memoryoverheadlabel,
		"hostmemoryoverheadtooltip":    class.Host.Memoryoverheadtooltip,
		"hostcoresoverhead":            class.Host.Coresoverhead,
		"hostcoresoverheadselected":    class.Host.Coresoverheadselected,
		"hostcoresoverheadlabel":       class.Host.Coresoverheadlabel,
		"hostcoresoverheadtooltip":     class.Host.Coresoverheadtooltip,
		/* storage are used to display form items for storage */
		"storagecapacityoverhead":         class.Storage.Capacityoverhead,
		"storagecapacityoverheadlabel":    class.Storage.Capacityoverheadlabel,
		"storagecapacityoverheadselected": class.Storage.Capacityoverheadselected,
		"storagecapacityoverheadtooltip":  class.Storage.Capacityoverheadtooltip,
		"storagedatastorevmcount":         class.Storage.Datastorevmcount,
		"storagedatastorevmcountlabel":    class.Storage.Datastorevmcountlabel,
		"storagedatastorevmcountselected": class.Storage.Datastorevmcountselected,
		"storagedatastorevmcounttooltip":  class.Storage.Datastorevmcounttooltip,
		"storagededuperatio":              class.Storage.Deduperatio,
		"storagededuperatiolabel":         class.Storage.Deduperatiolabel,
		"storagededuperatioselected":      class.Storage.Deduperatioselected,
		"storagededuperatiotooltip":       class.Storage.Deduperatiotooltip,
		"storageraidtype":                 class.Storage.Raidtype,
		"storageraidtypelabel":            class.Storage.Raidtypelabel,
		"storageraidtypeselected":         class.Storage.Raidtypeselected,
		/* virtualization are used to display form items for storage */
		"virtualizationclusterhostsize":                  class.Virtualization.Clusterhostsize,
		"virtualizationclusterhostsizelabel":             class.Virtualization.Clusterhostsizelabel,
		"virtualizationclusterhostsizeselected":          class.Virtualization.Clusterhostsizeselected,
		"virtualizationclusterhostsizetooltip":           class.Virtualization.Clusterhostsizetooltip,
		"virtualizationmanagementservertvmcount":         class.Virtualization.Managementservervmcount,
		"virtualizationmanagementservertvmcountlabel":    class.Virtualization.Managementservervmcountlabel,
		"virtualizationmanagementservertvmcountselected": class.Virtualization.Managementservervmcountselected,
		"virtualizationmanagementservertvmcounttooltip":  class.Virtualization.Managementservervmcounttooltip,
		/* hostresults are used to display resulting calculation for hosts */
		"hostresultscount":          class.HostResults.Count,
		"hostresultscountlabel":     class.HostResults.Countlabel,
		"hostresultsclockused":      class.HostResults.Hostclockused,
		"hostresultsclockusedlabel": class.HostResults.Hostclockusedlabel,
		"hostresultsmemory":         class.HostResults.Hostmemory,
		"hostresultsmemorylabel":    class.HostResults.Hostmemorylabel,
		"hostresultsvmcount":        class.HostResults.Hostvmcount,
		"hostresultsvmcountlabel":   class.HostResults.Hostvmcountlabel,
		/* storageresults are used to display resulting calculation for storage */
		"storageresultscapacity":                    class.StorageResults.Storagecapacity,
		"storageresultscapacitylabel":               class.StorageResults.Storagecapacitylabel,
		"storageresultsdatastorecount":              class.StorageResults.Storagedatastorecount,
		"storageresultsdatastorecountlabel":         class.StorageResults.Storagedatastorecountabel,
		"storageresultsdatastoresize:":              class.StorageResults.Storagedatastoresize,
		"storageresultsdatastoresizelabel":          class.StorageResults.Storagedatastoresizelabel,
		"storageresultsdatastorefroentendiopslabel": class.StorageResults.Storagedatastorefroentendiopslabel,
		"storageresultsdatastorebackendiopslabel":   class.StorageResults.Storagedatastorebackendiopslabel,
		"storageresultsfrontendiops":                class.StorageResults.Storagefrontendiops,
		"storageresultsfrontendiopslabel":           class.StorageResults.Storagefrontendiopslabel,
		"storageresultsbackendiops":                 class.StorageResults.Storagebackendiops,
		"storageresultsbackendiopslabel":            class.StorageResults.Storagebackendiopslabel,
		/* virtualizationresults are used to display resulting calculation for virtualization */
		"virtualizationresultsclustercount":               class.VirtualizationResults.Clustercount,
		"virtualizationresultsclustercountlabel":          class.VirtualizationResults.Clustercountlabel,
		"virtualizationresultsmanagementservercount":      class.VirtualizationResults.Managementservercount,
		"virtualizationresultsmanagementservercountlabel": class.VirtualizationResults.Managementservercountlabel,
		/* errorresults are used to display error and limits */
		"errorresults": class.ErrorResults,
	}

	return data

}
