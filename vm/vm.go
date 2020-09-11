package vm

/* This package contains external vm functions used across the app.
It must be imported using "vdicalc/vm" */

import (
	f "vdicalc/functions"
)

func main() {

}

// GetVmDisplayMemoryOverhead function
/* This public function calculates the display vm memory overhead */
func GetVmDisplayMemoryOverhead(vmdisplaycount string, vmdisplayresolution string, vmvideoram string) string {
	/* This function is reponsible for establishing the memory overhead based on the number of displays and resolution each vm is using.
	   The values have been obtained from VMware Horizon configuration guide and have been rounded up. The values are in megabytes (MB) per vm.
		1 = 1280x800
		2 = 1920x1200
		3 = 2560x1600
	*/

	var r int

	/* This function analises if 3D graphics video ram is enabled/disabled.
	If disabled it will select the ammount of video ram based of number os displays and display resolution. */
	if vmvideoram == "0" {

		/* This swith case analyses the number of displays in use by the vm */
		switch vmdisplaycount {
		case f.InttoStr(1):
			/* This swith case analyses the resolution in use by the displays */
			switch vmdisplayresolution {
			case "1":
				r = 4
			case "2":
				r = 8
			case "3":
				r = 16
			}

		case f.InttoStr(2):
			switch vmdisplayresolution {
			case "1":
				r = 13
			case "2":
				r = 26
			case "3":
				r = 60
			}

		case f.InttoStr(3):
			switch vmdisplayresolution {
			case "1":
				r = 19
			case "2":
				r = 38
			case "3":
				r = 85
			}
		case f.InttoStr(4):
			switch vmdisplayresolution {
			case "1":
				r = 25
			case "2":
				r = 51
			case "3":
				r = 110
			}
		}

	} else {

		r = f.StrtoInt(vmvideoram)

	}

	return f.InttoStr(r)
}

// GetVmVcpuMemoryOverhead function
/* This public function calculates the vm vcpu memory overhead */
func GetVmVcpuMemoryOverhead(vmvcpucount string, vmmemorysize string) string {
	/* This function is reponsible for establishing the memory overhead based on the number of vcpu each vm is using.
	   The values have been obtained from VMware Docs and have been rounded up. The values are in megabytes (MB) per vm.
	   https://docs.vmware.com/en/VMware-vSphere/6.7/com.vmware.vsphere.resmgmt.doc/GUID-B42C72C1-F8D5-40DC-93D1-FB31849B1114.html */

	var r int
	x := f.StrtoInt(vmmemorysize)

	/* This swith case analyses the ammount of memory in use by the vm */
	switch {
	case x <= 256:
		/* This swith case analyses the number of vcpu in use by the vm */
		switch vmvcpucount {
		case "1":
			r = 21
		case "2":
			r = 25
		case "4":
			r = 33
		case "8":
			r = 49
		}
	case x <= 1024:
		switch vmvcpucount {
		case "1":
			r = 26
		case "2":
			r = 30
		case "4":
			r = 38
		case "8":
			r = 54
		}
	case x <= 4096:
		switch vmvcpucount {
		case "1":
			r = 49
		case "2":
			r = 53
		case "4":
			r = 61
		case "8":
			r = 77
		}
	case x > 4096:
		switch vmvcpucount {
		case "1":
			r = 140
		case "2":
			r = 144
		case "4":
			r = 152
		case "8":
			r = 169
		}
	default:
	}

	return f.InttoStr(r)
}
