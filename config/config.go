package config

import (
	_ "github.com/go-sql-driver/mysql"
)

/* This package contains external config used across the app.
It must be imported using "vdicalc/config" */

// Configurations exported
type Configurations struct {
	Server         ServerConfigurations
	Variable       VariableConfigurations
	VM             VMConfigurations
	VMProfile01    VMConfigurations
	VMProfile02    VMConfigurations
	VMProfile03    VMConfigurations
	VMProfile04    VMConfigurations
	Host           HostConfigurations
	Storage        StorageConfigurations
	HostResults    HostResultsConfiguration
	StorageResults StorageResultsConfiguration
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Port int
}

// VariableConfigurations exported
type VariableConfigurations struct {
	Title    string
	Titlesub string
	Update   string
	About    string
	Print    string
}

// VMConfigurations exported
type VMConfigurations struct {
	Vmprofilelabel            string
	Vmprofile                 map[string]interface{}
	Vmprofileselected         string
	Vmcountlabel              string
	Vmcount                   int
	Vmcountselected           int
	Vcpucountlabel            string
	Vcpucount                 map[string]interface{}
	Vcpucountselected         string
	Vcpumhzlabel              string
	Vcpumhz                   int
	Vcpumhzselected           int
	Vmpercorecountlabel       string
	Vmpercorecount            map[int]interface{}
	Vmpercorecountselected    string
	Displaycountlabel         string
	Displaycount              map[int]interface{}
	Displaycountselected      string
	Displayresolutionlabel    string
	Displayresolution         map[string]interface{}
	Displayresolutionselected string
	Memorysizelabel           string
	Memorysize                int
	Memorysizeselected        int
	Videoram                  map[string]interface{}
	Videoramselected          string
	Videoramlabel             string
	Disksizelabel             string
	Disksize                  int
	Disksizeselected          int
	Iopscount                 int
	Iopscountselected         string
	Iopscountlabel            string
	Iopsreadratio             int
	Iopsreadratioselected     string
	Iopsreadratiolabel        string
	Iopswriteratio            int
}

// HostConfigurations exported
type HostConfigurations struct {
	Socketcountlabel         string
	Socketcount              map[int]interface{}
	Socketcountselected      string
	Socketcorescountlabel    string
	Socketcorescount         map[int]interface{}
	Socketcorescountselected string
	Memoryoverhead           int
	Memoryoverheadselected   int
	Memoryoverheadlabel      string
	Coresoverhead            int
	Coresoverheadselected    string
	Coresoverheadlabel       string
}

// StorageConfigurations exported
type StorageConfigurations struct {
	Capacityoverhead         map[int]interface{}
	Capacityoverheadlabel    string
	Capacityoverheadselected string
	Datastorevmcount         int
	Datastorevmcountlabel    string
	Datastorevmcountselected string
	Deduperatio              int
	Deduperatiolabel         string
	Deduperatioselected      string
	Raidtype                 map[int]interface{}
	Raidtypelabel            string
	Raidtypeselected         string
}

// HostResultsConfiguration exported
type HostResultsConfiguration struct {
	Countlabel         string
	Count              string
	Hostclockusedlabel string
	Hostclockused      string
	Hostmemory         string
	Hostmemorylabel    string
	Hostvmcount        string
	Hostvmcountlabel   string
}

// StorageResultsConfiguration exported
type StorageResultsConfiguration struct {
	Storagecapacitylabel               string
	Storagecapacity                    string
	Storagedatastorecountabel          string
	Storagedatastorecount              string
	Storagedatastoresize               string
	Storagedatastoresizelabel          string
	Storagedatastorefroentendiops      string
	Storagedatastorefroentendiopslabel string
	Storagedatastorebackendiops        string
	Storagedatastorebackendiopslabel   string
	Storagefrontendiops                string
	Storagefrontendiopslabel           string
	Storagebackendiops                 string
	Storagebackendiopslabel            string
}
