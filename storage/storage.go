package storage

/* This package contains external 'storage' functions used across the app.
It must be imported using "vdicalc/host" */

import (
	"math"
	f "vdicalc/functions"
	vm "vdicalc/vm"
)

func main() {

}

// GetStorageCapacity function
/* This public function calculates the storage capacity.
It returns result in terabytes. */
func GetStorageCapacity(vmcount string, vmdisksize string, storagecapacityoverhead string, storagededuperatio string, vmdisplaycount string, vmdisplayresolution string, vmvideoram string, vmmemorysize string, vmclonesizerefreshrate string) string {

	_, vmdisplaystorageoverhead := vm.GetVMDisplayOverhead(vmdisplaycount, vmdisplayresolution, vmvideoram)

	/* 	vmclonesizerefreshrate specifies linked clone type VM and this function defines is it has been selected.
	   	It changes storage capacity based on when the VM is refresh */
	var _vmdisksize float64
	if vmclonesizerefreshrate != "0" {
		_vmdisksize = f.StrtoFloat64(vmdisksize) * (f.StrtoFloat64(vmclonesizerefreshrate) / 100)
	} else {
		_vmdisksize = f.StrtoFloat64(vmdisksize)
	}

	/* vmmemorysize is used for VM swap calculation
	vmdisplaystorageoverhead is converted from MB to GB to match vmdisksize
	vmmemorysize is converted from MB to GB to match vmdisksize */
	r := (f.StrtoFloat64(vmcount) * (_vmdisksize + (f.StrtoFloat64(vmmemorysize) / 1000) + (f.StrtoFloat64(vmdisplaystorageoverhead) / 1000)))

	if storagecapacityoverhead != "0" {
		r += (f.StrtoFloat64(storagecapacityoverhead) / 100) * r
	}

	if storagededuperatio != "0" {
		r -= (f.StrtoFloat64(storagededuperatio) / 100) * r
	}

	/* Results are converted from GB to TB */
	return f.Float64toStr((r / 1000), 2)
}

// GetStorageDatastoreCount function
/* This public function calculates the number of datastores required based on the maximum number of VMs per datastore provided the user */
func GetStorageDatastoreCount(vmcount string, datastorevmcount string) string {

	r := int(math.Ceil(f.StrtoFloat64(vmcount) / f.StrtoFloat64(datastorevmcount)))

	return f.InttoStr(r)
}

// GetStorageDatastoreSize function
/* This public function calculates the size of the datastores based on total capacity required and the number of datastores determined */
func GetStorageDatastoreSize(vmcount string, datastorevmcount string, vmdisksize string, storagecapacityoverhead string, storagededuperatio string, vmdisplaycount string, vmdisplayresolution string, vmvideoram string, vmmemorysize string, vmclonesizerefreshrate string) string {

	r := f.StrtoFloat64(GetStorageCapacity(vmcount, vmdisksize, storagecapacityoverhead, storagededuperatio, vmdisplaycount, vmdisplayresolution, vmvideoram, vmmemorysize, vmclonesizerefreshrate)) / f.StrtoFloat64(GetStorageDatastoreCount(vmcount, datastorevmcount))

	return f.Float64toStr(r, 2)
}

// GetStorageDatastoreIops function
/* This public function calculates the amount of frontend and backend IOPs per datastore and for the full storage.
For write IOps the function calculate the write amplification based on raid levels. */
func GetStorageDatastoreIops(vmiopscount string, vmiopsreadratio string, storagedatastorevmcount string, storageraidtype string, vmcount string, datastorevmcount string) (string, string, string, string) {

	datastoreFrontendIops := f.StrtoInt(vmiopscount) * f.StrtoInt(storagedatastorevmcount)
	datastoreBackendReadIops := int(((f.StrtoFloat64(vmiopsreadratio) / 100) * f.StrtoFloat64(vmiopscount)) * f.StrtoFloat64(storagedatastorevmcount))
	datastoreBackendWriteIops := int(((1 - (f.StrtoFloat64(vmiopsreadratio) / 100)) * f.StrtoFloat64(vmiopscount)) * f.StrtoFloat64(storagedatastorevmcount))

	switch storageraidtype {
	case "5":
		datastoreBackendWriteIops *= 4
	case "6":
		datastoreBackendWriteIops *= 6
	case "10":
		datastoreBackendWriteIops *= 2
	}

	datastoreBackendIops := datastoreBackendReadIops + datastoreBackendWriteIops

	datastoreCount := GetStorageDatastoreCount(vmcount, datastorevmcount)
	storageFrontendIops := datastoreFrontendIops * f.StrtoInt(datastoreCount)
	storageBackendIops := datastoreBackendIops * f.StrtoInt(datastoreCount)

	return f.InttoStr(datastoreFrontendIops), f.InttoStr(datastoreBackendIops), f.InttoStr(storageFrontendIops), f.InttoStr(storageBackendIops)
}
