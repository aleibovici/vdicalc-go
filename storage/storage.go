package storage

/* This package contains external 'storage' functions used across the app.
It must be imported using "vdicalc/host" */

import (
	"math"
	f "vdicalc/functions"
)

func main() {

}

// GetStorageCapacity function
/* This public function calculates the storage capacity.
It returns result in terabytes. */
func GetStorageCapacity(vmcount string, vmdisksize string, storagecapacityoverhead string) string {

	r := (f.StrtoFloat64(vmcount) * f.StrtoFloat64(vmdisksize)) / 1000
	r *= (1 + (f.StrtoFloat64(storagecapacityoverhead) / 100))

	return f.InttoStr(int(r))
}

// GetStorageDatastoreCount function
/* This public function calculates the number of datastores required based on the maximum number of VMs per datastore provided the user */
func GetStorageDatastoreCount(vmcount string, datastorevmcount string) string {

	r := int(math.Ceil(f.StrtoFloat64(vmcount) / f.StrtoFloat64(datastorevmcount)))

	return f.InttoStr(r)
}

// GetStorageDatastoreSize function
/* This public function calculates the size of the datastores based on total capacity required and the number of datastores determined */
func GetStorageDatastoreSize(vmcount string, datastorevmcount string, vmdisksize string, storagecapacityoverhead string) string {

	r := f.StrtoFloat64(GetStorageCapacity(vmcount, vmdisksize, storagecapacityoverhead)) / f.StrtoFloat64(GetStorageDatastoreCount(vmcount, datastorevmcount))

	return f.Float64toStr(r)
}
