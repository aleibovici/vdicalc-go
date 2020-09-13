package host

/* This package contains external 'host' functions used across the app.
It must be imported using "vdicalc/host" */

import (
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
func GetHostCount(vmcount string, hostsocketcount string, hostsocketcorescount string, vmspercore string, hostcoresoverhead string) string {

	r := f.StrtoInt(vmcount) / f.StrtoInt(GetHostVMCount(vmcount, hostsocketcount, hostsocketcorescount, vmspercore, hostcoresoverhead))

	return f.InttoStr(r)
}

// GetHostClockUsed function
/* This public function calculates the host clock rate */
func GetHostClockUsed(vmvcpucount string, vmvcpumhz string, vmcount string, hostsocketcount string, hostsocketcorescount string, vmspercore string, hostcoresoverhead string) string {

	r := (f.StrtoInt(vmvcpucount) * f.StrtoInt(vmvcpumhz) * f.StrtoInt(GetHostVMCount(vmcount, hostsocketcount, hostsocketcorescount, vmspercore, hostcoresoverhead)) / f.StrtoInt(getHostCoresCount(hostsocketcount, hostsocketcorescount, hostcoresoverhead)))

	return f.InttoStr(r)
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
