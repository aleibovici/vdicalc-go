package functions

/* This package contains external genericfunctions used across the app.
It must be imported using "vdicalc/functions" */

import (
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
/* This public function convert float64 to string using .00 digits */
func Float64toStr(value float64) string {

	return strconv.FormatFloat(value, 'f', 2, 64)
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
		/* vm are used to display form items for hosts */
		"vmprofilelabel":              class.VM.Vmprofilelabel,
		"vmprofile":                   class.VM.Vmprofile,
		"vmprofileselected":           class.VM.Vmprofileselected,
		"vmcountlabel":                class.VM.Vmcountlabel,
		"vmcount":                     class.VM.Vmcount,
		"vmcountselected":             class.VM.Vmcountselected,
		"vmvcpucountlabel":            class.VM.Vcpucountlabel,
		"vmvcpucount":                 class.VM.Vcpucount,
		"vmvcpucountselected":         class.VM.Vcpucountselected,
		"vmvcpumhzlabel":              class.VM.Vcpumhzlabel,
		"vmvcpumhz":                   class.VM.Vcpumhz,
		"vmvcpumhzselected":           class.VM.Vcpumhzselected,
		"vmpercorecountlabel":         class.VM.Vmpercorecountlabel,
		"vmpercorecount":              class.VM.Vmpercorecount,
		"vmpercorecountselected":      class.VM.Vmpercorecountselected,
		"vmdisplaycountlabel":         class.VM.Displaycountlabel,
		"vmdisplaycount":              class.VM.Displaycount,
		"vmdisplaycountselected":      class.VM.Displaycountselected,
		"vmdisplayresolutionlabel":    class.VM.Displayresolutionlabel,
		"vmdisplayresolution":         class.VM.Displayresolution,
		"vmdisplayresolutionselected": class.VM.Displayresolutionselected,
		"memorysizelabel":             class.VM.Memorysizelabel,
		"memorysize":                  class.VM.Memorysize,
		"memorysizeselected":          class.VM.Memorysizeselected,
		"vmvideoram":                  class.VM.Videoram,
		"vmvideoramlabel":             class.VM.Videoramlabel,
		"vmvideoramselected":          class.VM.Videoramselected,
		"vmdisksizelabel":             class.VM.Disksizelabel,
		"vmdisksize":                  class.VM.Disksize,
		"vmdisksizeselected":          class.VM.Disksizeselected,
		"vmiopscount":                 class.VM.Iopscount,
		"vmiopscountlabel":            class.VM.Iopscountlabel,
		"vmiopscountselected":         class.VM.Iopscountselected,
		"vmiopsreadratio":             class.VM.Iopsreadratio,
		"vmiopsreadratiolabel":        class.VM.Iopsreadratiolabel,
		"vmiopsreadratioselected":     class.VM.Iopsreadratioselected,
		"vmiopswriteratio":            class.VM.Iopswriteratio,
		/* host are used to display form items for hosts */
		"socketcountlabel":           class.Host.Socketcountlabel,
		"socketcount":                class.Host.Socketcount,
		"socketcountselected":        class.Host.Socketcountselected,
		"socketcorescountlabel":      class.Host.Socketcorescountlabel,
		"socketcorescount":           class.Host.Socketcorescount,
		"socketcorescountselected":   class.Host.Socketcorescountselected,
		"hostmemoryoverhead":         class.Host.Memoryoverhead,
		"hostmemoryoverheadselected": class.Host.Memoryoverheadselected,
		"hostmemoryoverheadlabel":    class.Host.Memoryoverheadlabel,
		"hostcoresoverhead":          class.Host.Coresoverhead,
		"hostcoresoverheadselected":  class.Host.Coresoverheadselected,
		"hostcoresoverheadlabel":     class.Host.Coresoverheadlabel,
		/* storage are used to display form items for storage */
		"storagecapacityoverhead":         class.Storage.Capacityoverhead,
		"storagecapacityoverheadlabel":    class.Storage.Capacityoverheadlabel,
		"storagecapacityoverheadselected": class.Storage.Capacityoverheadselected,
		"storagedatastorevmcount":         class.Storage.Datastorevmcount,
		"storagedatastorevmcountlabel":    class.Storage.Datastorevmcountlabel,
		"storagedatastorevmcountselected": class.Storage.Datastorevmcountselected,
		"storagededuperatio":              class.Storage.Deduperatio,
		"storagededuperatiolabel":         class.Storage.Deduperatiolabel,
		"storagededuperatioselected":      class.Storage.Deduperatioselected,
		"storageraidtype":                 class.Storage.Raidtype,
		"storageraidtypelabel":            class.Storage.Raidtypelabel,
		"storageraidtypeselected":         class.Storage.Raidtypeselected,
		/* hostresults are used to display resulting calculation for hosts */
		"hostresultscount":          class.HostResults.Count,
		"hostresultscountlabel":     class.HostResults.Countlabel,
		"hostresultsclockused":      class.HostResults.Hostclockused,
		"hostresultsclockusedlabel": class.HostResults.Hostclockusedlabel,
		"hosresultstmemory":         class.HostResults.Hostmemory,
		"hosresultstmemorylabel":    class.HostResults.Hostmemorylabel,
		"hostresultvmcount":         class.HostResults.Hostvmcount,
		"hostresultvmcountlabel":    class.HostResults.Hostvmcountlabel,
		/* storageresults are used to display resulting calculation for storage */
		"storageresultcapacity":              class.StorageResults.Storagecapacity,
		"storageresultcapacitylabel":         class.StorageResults.Storagecapacitylabel,
		"storageresultdatastorecount":        class.StorageResults.Storagedatastorecount,
		"storageresultdatastorecountlabel":   class.StorageResults.Storagedatastorecountabel,
		"storageresultdatastoresize:":        class.StorageResults.Storagedatastoresize,
		"storageresultdatastoresizelabel":    class.StorageResults.Storagedatastoresizelabel,
		"storagedatastorefroentendiopslabel": class.StorageResults.Storagedatastorefroentendiopslabel,
		"storagedatastorebackendiopslabel":   class.StorageResults.Storagedatastorebackendiopslabel,
	}

	return data

}
