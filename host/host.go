package host

/* This package contains external 'host' functions used across the app.
It must be imported using "vdicalc/host" */

import (
	"log"
	f "vdicalc/functions"
	vm "vdicalc/vm"
)

func main() {

}

/* This private function calculates the number of cores per host */
func getHostCoresCount(hostsocketcount string, hostsocketcorescount string, hostcoresoverhead string) string {

	r := (f.StrtoInt(hostsocketcount) * f.StrtoInt(hostsocketcorescount)) - f.StrtoInt(hostcoresoverhead)

	return f.InttoStr(r)
}

// GetHostVMCount function
/* This public function calculates the number of vms per host */
func GetHostVMCount(vmcount string, hostsocketcount string, hostsocketcorescount string, vmspercore string, hostcoresoverhead string) string {

	var r int

	if f.StrtoInt(vmcount) < (f.StrtoInt(getHostCoresCount(hostsocketcount, hostsocketcorescount, hostcoresoverhead)) * f.StrtoInt(vmspercore)) {
		r = f.StrtoInt(vmcount)
	} else {
		r = (f.StrtoInt(getHostCoresCount(hostsocketcount, hostsocketcorescount, hostcoresoverhead)) * f.StrtoInt(vmspercore))
	}

	return f.InttoStr(r)
}

// GetHostCount function
/* This public function calculates the number of host */
func GetHostCount(vmcount string, hostsocketcount string, hostsocketcorescount string, vmspercore string, hostcoresoverhead string, clusterhostha string) string {

	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()

	r := f.StrtoFloat64(vmcount) / f.StrtoFloat64(GetHostVMCount(vmcount, hostsocketcount, hostsocketcorescount, vmspercore, hostcoresoverhead))

	/* Add cluster high availability overhead if clusterhostha true (false=0/true=1) */
	if clusterhostha == "1" {
		r *= 1.125
	}

	return f.Float64toStr(r, 0)
}

// GetHostClockUsed function
/* This public function calculates the host clock rate and converts it from MHz to GHz*/
func GetHostClockUsed(vmvcpucount string, vmvcpumhz string, vmcount string, hostsocketcount string, hostsocketcorescount string, vmspercore string, hostcoresoverhead string) string {

	r := (f.StrtoFloat64(vmvcpucount) * f.StrtoFloat64(vmvcpumhz) * f.StrtoFloat64(GetHostVMCount(vmcount, hostsocketcount, hostsocketcorescount, vmspercore, hostcoresoverhead)) / f.StrtoFloat64(getHostCoresCount(hostsocketcount, hostsocketcorescount, hostcoresoverhead))) / 1000

	return f.Float64toStr(r, 1)
}

// GetHostMemory function
/* This public function calculates the host memory */
func GetHostMemory(vmcount string, hostsocketcount string, hostsocketcorescount string, hostcoresoverhead string, vmspercore string, vmmemorysize string, hostmemoryoverhead string, vmdisplaycount string, vmdisplayresolution string, vmvcpucount string, vmvideoram string) string {

	hostvmcount := GetHostVMCount(vmcount, hostsocketcount, hostsocketcorescount, vmspercore, hostcoresoverhead)
	vmdisplaymemoryoverhead, _ := vm.GetVMDisplayOverhead(vmdisplaycount, vmdisplayresolution, vmvideoram)
	vmvcpumemoryoverhead := vm.GetVMVcpuMemoryOverhead(vmvcpucount, vmmemorysize)
	r := ((f.StrtoInt(hostvmcount) * (f.StrtoInt(vmmemorysize) + f.StrtoInt(vmdisplaymemoryoverhead) + f.StrtoInt(vmvcpumemoryoverhead))) / 1024) + f.StrtoInt(hostmemoryoverhead)

	return f.InttoStr(r)
}
