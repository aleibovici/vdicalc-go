package virtualization

import (
	"math"
	"vdicalc/functions"
	"vdicalc/host"
)

// GetClusterSize function
/* This public function retrieve the virtualization cluster size */
func GetClusterSize(vmcount string, hostsocketcount string, hostsocketcorescount string, vmpercorecount string, hostcoresoverhead string, virtualizationclusterhostsize string, clusterhostha string) string {

	r := int(math.Ceil(functions.StrtoFloat64(host.GetHostCount(vmcount, hostsocketcount, hostsocketcorescount, vmpercorecount, hostcoresoverhead, clusterhostha)) / functions.StrtoFloat64(virtualizationclusterhostsize)))

	return functions.InttoStr(r)

}

// GetManagementServerCount function
/* This public function calculates the number of virtualization management servers based on VM limit per management server */
func GetManagementServerCount(vmcount string, virtualizationmanagementservertvmcount string) string {

	r := int(math.Ceil(functions.StrtoFloat64(vmcount) / functions.StrtoFloat64(virtualizationmanagementservertvmcount)))

	return functions.InttoStr(r)
}
