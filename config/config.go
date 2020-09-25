package config

/* This package contains external config used across the app.
It must be imported using "vdicalc/config" */

// Configurations exported
type Configurations struct {
	Server                ServerConfigurations
	Variable              VariableConfigurations
	VM                    VMConfigurations
	VMProfile01           VMConfigurations
	VMProfile02           VMConfigurations
	VMProfile03           VMConfigurations
	VMProfile04           VMConfigurations
	Host                  HostConfigurations
	Storage               StorageConfigurations
	Virtualization        VirtualizationConfigurations
	HostResults           HostResultsConfiguration
	StorageResults        StorageResultsConfiguration
	VirtualizationResults VirtualizationResultsConfiguration
	ErrorResults          ErrorResultsConfiguration
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Port        int
	AuthAddress string
}

// VariableConfigurations exported
type VariableConfigurations struct {
	Title     string
	Titlesub  string
	Update    string
	About     string
	Print     string
	Guide     string
	Load      string
	Save      string
	Usersaves map[string]interface{}
}

// VMConfigurations exported
type VMConfigurations struct {
	Vmprofilelabel               string
	Vmprofile                    map[string]interface{}
	Vmprofileselected            string
	Vmcountlabel                 string
	Vmcount                      int
	Vmcountselected              int
	Vcpucountlabel               string
	Vcpucount                    map[string]interface{}
	Vcpucountselected            string
	Vcpumhzlabel                 string
	Vcpumhz                      int
	Vcpumhzselected              int
	Vcpumhztooltip               string
	Vmpercorecountlabel          string
	Vmpercorecount               map[int]interface{}
	Vmpercorecountselected       string
	Displaycountlabel            string
	Displaycount                 map[int]interface{}
	Displaycountselected         string
	Displayresolutionlabel       string
	Displayresolution            map[string]interface{}
	Displayresolutionselected    string
	Memorysizelabel              string
	Memorysize                   string
	Memorysizelimit              string
	Memorysizeselected           string
	Videoram                     map[string]interface{}
	Videoramselected             string
	Videoramlabel                string
	Disksizelabel                string
	Disksize                     int
	Disksizeselected             int
	Iopscount                    int
	Iopscountselected            string
	Iopscountlabel               string
	Iopscounttooltip             string
	Iopsreadratio                int
	Iopsreadratioselected        string
	Iopsreadratiolabel           string
	Iopsreadratiotooltip         string
	Iopswriteratio               int
	Clonesizerefreshrate         map[string]interface{}
	Clonesizerefreshratelabel    string
	Clonesizerefreshrateselected string
	Clonesizerefreshratetooltip  string
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
	Memoryoverheadtooltip    string
	Coresoverhead            int
	Coresoverheadselected    string
	Coresoverheadlabel       string
	Coresoverheadtooltip     string
}

// StorageConfigurations exported
type StorageConfigurations struct {
	Capacityoverhead         map[int]interface{}
	Capacityoverheadlabel    string
	Capacityoverheadselected string
	Capacityoverheadtooltip  string
	Datastorevmcount         int
	Datastorevmcountlabel    string
	Datastorevmcountselected string
	Datastorevmcounttooltip  string
	Deduperatio              int
	Deduperatiolabel         string
	Deduperatioselected      string
	Deduperatiotooltip       string
	Raidtype                 map[int]interface{}
	Raidtypelabel            string
	Raidtypeselected         string
}

// VirtualizationConfigurations exported
type VirtualizationConfigurations struct {
	Clusterhostsize                 map[int]interface{}
	Clusterhostsizelabel            string
	Clusterhostsizeselected         string
	Clusterhostsizetooltip          string
	Managementservervmcount         int
	Managementservervmcountlabel    string
	Managementservervmcountselected string
	Managementservervmcounttooltip  string
}

// HostResultsConfiguration exported
type HostResultsConfiguration struct {
	Countlabel         string
	Count              string
	Hostclockusedlabel string
	Hostclockused      string
	Hostclockusedlimit string
	Hostmemory         string
	Hostmemorylimit    string
	Hostmemorylabel    string
	Hostvmcount        string
	Hostvmcountlimit   string
	Hostvmcountlabel   string
}

// StorageResultsConfiguration exported
type StorageResultsConfiguration struct {
	Storagecapacitylabel               string
	Storagecapacity                    string
	Storagedatastorecountabel          string
	Storagedatastorecount              string
	Storagedatastorecountlimit         string
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

// VirtualizationResultsConfiguration exported
type VirtualizationResultsConfiguration struct {
	Clustercount               string
	Clustercountlabel          string
	Managementservercount      string
	Managementservercountlabel string
	Managementservercountlimit string
}

// ErrorResultsConfiguration exported
type ErrorResultsConfiguration struct {
	Error []ErrorConfiguration
}

// ErrorConfiguration exported
type ErrorConfiguration struct {
	Code        string
	Description string
}
