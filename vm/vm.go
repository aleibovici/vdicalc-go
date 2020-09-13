package vm

/* This package contains external vm functions used across the app.
It must be imported using "vdicalc/vm" */

import (
	f "vdicalc/functions"
)

func main() {

}

// GetVMDisplayOverhead function
/* This public function calculates the display and resolution overhead for memory and storage vswap */
func GetVMDisplayOverhead(vmdisplaycount string, vmdisplayresolution string, vmvideoram string) (string, string) {
	/* This function is reponsible for establishing the memory and storage overhead based on the number of displays and resolution each vm is using. The values have been obtained from VMware Horizon configuration guide and have been rounded up.
	   The values are in megabytes (MB) per vm.
		1 = 1280x800
		2 = 1920x1200
		3 = 2560x1600
	*/

	var m int // memory
	var s int // storage

	/* This function analises if 3D graphics video ram is enabled/disabled.
	If disabled it will select the ammount of video ram based of number os displays and display resolution. */
	if vmvideoram == "0" {

		/* This swith case analyses the number of displays in use by the vm */
		switch vmdisplaycount {
		case f.InttoStr(1):
			/* This swith case analyses the resolution in use by the displays */
			switch vmdisplayresolution {
			case "1":
				m = 4
				s = 107
			case "2":
				m = 8
				s = 111
			case "3":
				m = 16
				s = 203
			}

		case f.InttoStr(2):
			switch vmdisplayresolution {
			case "1":
				m = 13
				s = 163
			case "2":
				m = 26
				s = 190
			case "3":
				m = 60
				s = 203
			}

		case f.InttoStr(3):
			switch vmdisplayresolution {
			case "1":
				m = 19
				s = 207
			case "2":
				m = 38
				s = 248
			case "3":
				m = 85
				s = 461
			}
		case f.InttoStr(4):
			switch vmdisplayresolution {
			case "1":
				m = 25
				s = 252
			case "2":
				m = 51
				s = 306
			case "3":
				m = 110
				s = 589
			}
		}

	} else {

		switch vmvideoram {
		case "64":
			s = 1076
		case "128":
			s = 1468
		case "256":
			s = 1468
		case "512":
			s = 1916
		}

		m = f.StrtoInt(vmvideoram)

	}

	return f.InttoStr(m), f.InttoStr(s)
}

// GetVMVcpuMemoryOverhead function
/* This public function calculates the vm vcpu memory overhead */
func GetVMVcpuMemoryOverhead(vmvcpucount string, vmmemorysize string) string {
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
