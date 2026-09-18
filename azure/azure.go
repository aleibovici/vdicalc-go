package azure

import (
	f "vdicalc/functions"
)

// GetAzureInstanceType Export
/* Recommends Azure VM + Premium SSD for AVD single-session (Ds_v5 / NVads_A10_v5).
   Source: https://learn.microsoft.com/en-us/windows-server/remote/remote-desktop-services/session-host-virtual-machine-sizing-guidelines */
func GetAzureInstanceType(vmvcpucount string, vmmemorysize string, vmdisksize string, vmvideoram string) string {

	var result string
	memory := (f.StrtoFloat64(vmmemorysize)) / 1024

	switch f.StrtoInt(vmvcpucount) {
	case 1:
		switch {
		case memory <= 8.1:
			result = "D2s_v5"
		case memory <= 16.1:
			result = "D4s_v5"
		case memory <= 32.1:
			result = "D8s_v5"
		default:
			result = "D16s_v5"
		}
	case 2:
		switch {
		case memory <= 8.1:
			result = "D2s_v5"
		case memory <= 16.1:
			result = "D4s_v5"
		case memory <= 32.1:
			result = "D8s_v5"
		default:
			result = "D16s_v5"
		}
	case 4:
		switch {
		case memory <= 16.1:
			result = "D4s_v5"
		case memory <= 32.1:
			result = "D8s_v5"
		case memory <= 64.1:
			result = "D16s_v5"
		default:
			result = "D32s_v5"
		}
	case 8:
		switch {
		case memory <= 32.1:
			result = "D8s_v5"
		case memory <= 64.1:
			result = "D16s_v5"
		default:
			result = "D32s_v5"
		}
	default:
		result = "D4s_v5"
	}

	switch f.StrtoInt(vmvideoram) {
	case 1:
		switch f.StrtoInt(vmvcpucount) {
		case 1, 2, 4:
			switch {
			case memory <= 55.1:
				result = "NV6ads_A10_v5"
			case memory <= 110.1:
				result = "NV12ads_A10_v5"
			default:
				result = "NV18ads_A10_v5"
			}
		case 8:
			switch {
			case memory <= 110.1:
				result = "NV12ads_A10_v5"
			case memory <= 220.1:
				result = "NV18ads_A10_v5"
			default:
				result = "NV36ads_A10_v5"
			}
		}
	}

	c := (f.StrtoInt(vmdisksize))
	switch {
	case c <= 32:
		result += " P4"
	case c <= 64:
		result += " P6"
	case c <= 128:
		result += " P10"
	case c <= 256:
		result += " P15"
	case c <= 512:
		result += " P20"
	case c <= 1024:
		result += " P30"
	default:
		result += " P40"
	}

	return result
}
