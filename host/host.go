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
func getHostCoresCount(socketcount string, socketcorescount string, hostcoresoverhead string) string {

	r := (f.StrtoInt(socketcount) * f.StrtoInt(socketcorescount)) - f.StrtoInt(hostcoresoverhead)

	return f.InttoStr(r)
}

// GetHostVmCount function
/* This public function calculates the number of vms per host */
func GetHostVmCount(vmcount string, hostsocketcount string, hostsocketcorescount string, vmspercore string, hostcoresoverhead string) string {

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
func GetHostCount(vmcount string, socketcount string, socketcorescount string, vmspercore string, hostcoresoverhead string) string {

	r := f.StrtoInt(vmcount) / f.StrtoInt(GetHostVmCount(vmcount, socketcount, socketcorescount, vmspercore, hostcoresoverhead))

	return f.InttoStr(r)
}

// GetHostClockUsed function
/* This public function calculates the host clock rate */
func GetHostClockUsed(vmvcpucount string, vmvcpumhz string, vmcount string, socketcount string, socketcorescount string, vmspercore string, hostcoresoverhead string) string {

	r := (f.StrtoInt(vmvcpucount) * f.StrtoInt(vmvcpumhz) * f.StrtoInt(GetHostVmCount(vmcount, socketcount, socketcorescount, vmspercore, hostcoresoverhead)) / f.StrtoInt(getHostCoresCount(socketcount, socketcorescount, hostcoresoverhead)))

	return f.InttoStr(r)
}

// GetHostMemory function
/* This public function calculates the host memory */
func GetHostMemory(vmcount string, socketcount string, socketcorescount string, hostcoresoverhead string, vmspercore string, vmmemorysize string, hostmemoryoverhead string, vmdisplaycount string, vmdisplayresolution string, vmvcpucount string, vmvideoram string) string {

	hostvmcount := GetHostVmCount(vmcount, socketcount, socketcorescount, vmspercore, hostcoresoverhead)
	vmdisplaymemoryoverhead := vm.GetVmDisplayMemoryOverhead(vmdisplaycount, vmdisplayresolution, vmvideoram)
	vmvcpumemoryoverhead := vm.GetVmVcpuMemoryOverhead(vmvcpucount, vmmemorysize)
	r := ((f.StrtoInt(hostvmcount) * (f.StrtoInt(vmmemorysize) + f.StrtoInt(vmdisplaymemoryoverhead) + f.StrtoInt(vmvcpumemoryoverhead))) / 1024) + f.StrtoInt(hostmemoryoverhead)

	return f.InttoStr(r)
}
