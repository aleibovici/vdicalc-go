package azure

import (
	f "vdicalc/functions"
)

// GetAzureInstanceType Export
func GetAzureInstanceType(vmvcpucount string, vmdisksize string, vmvideoram string) string {

	var result string

	// Cores (vmvcpucount)
	switch f.StrtoInt(vmvcpucount) {
	case 1:
		result = "F1"
	case 2:
		result = "F2"
	case 4:
		result = "F4"
	case 8:
		result = "F8"
	}

	// GPU
	switch f.StrtoInt(vmvideoram) {
	case 1:
		result = "N6"
	}

	// Storage (vmdisksize)
	c := (f.StrtoInt(vmdisksize))
	switch {
	case c <= 32:
		result = result + " | P4"
	case c <= 64:
		result = result + " | P6"
	case c <= 128:
		result = result + " | P10"
	case c <= 256:
		result = result + " | P15"
	case c <= 512:
		result = result + " | P20"
	case c <= 1024:
		result = result + " | P30"
	case c >= 1024:
		result = result + " | P40"
	}

	return result
}
