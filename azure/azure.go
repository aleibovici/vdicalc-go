package azure

import (
	f "vdicalc/functions"
)

// GetAzureInstanceType Export
/* This public function is reponsible to recommend the correct Azure VM instance type */
/* The source for definitions: https://azure.microsoft.com/en-us/pricing/details/virtual-machines/windows/#n-series */
func GetAzureInstanceType(vmvcpucount string, vmmemorysize string, vmdisksize string, vmvideoram string) string {

	/* 	F1 / 1 core / 2 GiB
	   	F2 / 2 core / 4 GiB
	   	F4 / 4 core / 8 GiB
	   	F8 / 8 core / 16 GiB
	   	F16 / 16 core / 32 GiB */

	var result string
	memory := (f.StrtoFloat64(vmmemorysize)) / 1024

	/* This function select the VM instance type based on number of cores and memory */
	/* We use float and .1 decimal to ensure that a smaller VM instance type is selected if memory is close enough */
	switch f.StrtoInt(vmvcpucount) {
	case 1:
		switch {
		case memory <= 2.1:
			result = "F1"
		case memory <= 4.1:
			result = "F2"
		case memory <= 8.1:
			result = "F4"
		case memory <= 16.1:
			result = "F8"
		case memory > 16.1:
			result = "F16"
		}
	case 2:
		switch {
		case memory <= 4.1:
			result = "F2"
		case memory <= 8.1:
			result = "F4"
		case memory <= 16.1:
			result = "F8"
		case memory > 16.1:
			result = "F16"
		}
	case 4:
		switch {
		case memory <= 8.1:
			result = "F4"
		case memory <= 16.1:
			result = "F8"
		case memory > 16.1:
			result = "F16"
		}
	case 8:
		switch {
		case memory <= 16:
			result = "F8"
		case memory > 16:
			result = "F16"
		}
	}

	/* This function select the VM instance type based on GPU requirement */
	/* The calculator only support calculations up to 8 cores due to vSphere calculations. N6 has 6 cores */
	switch f.StrtoInt(vmvideoram) {
	case 1:
		result = "N6"
	}

	/* This function select the disk instance */
	c := (f.StrtoInt(vmdisksize))
	switch {
	case c <= 32:
		result += " | P4"
	case c <= 64:
		result += " | P6"
	case c <= 128:
		result += " | P10"
	case c <= 256:
		result += " | P15"
	case c <= 512:
		result += " | P20"
	case c <= 1024:
		result += " | P30"
	case c >= 1024:
		result += " | P40"
	}

	return result
}
