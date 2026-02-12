// VDI Calculator - JavaScript Port
// Original Go implementation by Andre Leibovici (myvirtualcloud.net)
// Ported to client-side JavaScript for GitHub Pages

// ============================================================
// Helper functions (port of functions/functions.go)
// ============================================================

function toInt(value) {
  return parseInt(value, 10) || 0;
}

function toFloat(value) {
  return parseFloat(value) || 0;
}

function val(id) {
  return document.getElementById(id).value;
}

function setResult(id, value) {
  document.getElementById(id).textContent = value;
}

// ============================================================
// VM Overhead Calculations (port of vm/vm.go)
// ============================================================

// GetVMDisplayOverhead - calculates display/resolution overhead for memory and storage vswap
// Returns { memory: int, storage: int } in MB
function getVMDisplayOverhead(displayCount, displayResolution, videoRAM) {
  var m = 0, s = 0;
  displayCount = toInt(displayCount);
  displayResolution = String(displayResolution);
  videoRAM = String(videoRAM);

  if (videoRAM === "0") {
    // No 3D graphics - overhead based on display count and resolution
    if (displayCount === 1) {
      switch (displayResolution) {
        case "1": m = 4; s = 107; break;
        case "2": m = 8; s = 111; break;
        case "3": m = 16; s = 203; break;
      }
    } else if (displayCount === 2) {
      switch (displayResolution) {
        case "1": m = 13; s = 163; break;
        case "2": m = 26; s = 190; break;
        case "3": m = 60; s = 203; break;
      }
    } else if (displayCount === 3) {
      switch (displayResolution) {
        case "1": m = 19; s = 207; break;
        case "2": m = 38; s = 248; break;
        case "3": m = 85; s = 461; break;
      }
    } else if (displayCount === 4) {
      switch (displayResolution) {
        case "1": m = 25; s = 252; break;
        case "2": m = 51; s = 306; break;
        case "3": m = 110; s = 589; break;
      }
    }
  } else if (videoRAM === "1") {
    // GPU use case
    m = 96;
    s = 0;
  } else {
    // Software 3D graphics with specific video RAM
    switch (videoRAM) {
      case "64": s = 1076; break;
      case "128": s = 1468; break;
      case "256": s = 1468; break;
      case "512": s = 1916; break;
    }
    m = toInt(videoRAM);
  }

  return { memory: m, storage: s };
}

// GetVMVcpuMemoryOverhead - calculates vm vcpu memory overhead in MB
function getVMVcpuMemoryOverhead(vcpuCount, memorySize) {
  var r = 0;
  var x = toInt(memorySize);
  vcpuCount = String(vcpuCount);

  if (x <= 256) {
    switch (vcpuCount) {
      case "1": r = 21; break;
      case "2": r = 25; break;
      case "4": r = 33; break;
      case "8": r = 49; break;
    }
  } else if (x <= 1024) {
    switch (vcpuCount) {
      case "1": r = 26; break;
      case "2": r = 30; break;
      case "4": r = 38; break;
      case "8": r = 54; break;
    }
  } else if (x <= 4096) {
    switch (vcpuCount) {
      case "1": r = 49; break;
      case "2": r = 53; break;
      case "4": r = 61; break;
      case "8": r = 77; break;
    }
  } else {
    switch (vcpuCount) {
      case "1": r = 140; break;
      case "2": r = 144; break;
      case "4": r = 152; break;
      case "8": r = 169; break;
    }
  }

  return r;
}

// ============================================================
// Host Calculations (port of host/host.go)
// ============================================================

function getHostCoresCount(socketCount, coresPerSocket, coresOverhead) {
  return (toInt(socketCount) * toInt(coresPerSocket)) - toInt(coresOverhead);
}

// GetHostVMCount - number of VMs per host
function getHostVMCount(vmCount, socketCount, coresPerSocket, vmsPerCore, coresOverhead) {
  var capacity = getHostCoresCount(socketCount, coresPerSocket, coresOverhead) * toInt(vmsPerCore);
  return Math.min(toInt(vmCount), capacity);
}

// GetHostCount - number of hosts needed
function getHostCount(vmCount, socketCount, coresPerSocket, vmsPerCore, coresOverhead, clusterHA) {
  var hostVMCount = getHostVMCount(vmCount, socketCount, coresPerSocket, vmsPerCore, coresOverhead);
  if (hostVMCount === 0) return 0;
  var r = toFloat(vmCount) / toFloat(hostVMCount);
  if (String(clusterHA) === "1") {
    r *= 1.125;
  }
  // Go uses FormatFloat with precision 0, which rounds to nearest
  return Number(r.toFixed(0));
}

// GetHostClockUsed - host CPU clock in GHz
function getHostClockUsed(vcpuCount, vcpuMHz, vmCount, socketCount, coresPerSocket, vmsPerCore, coresOverhead) {
  var hostVMCount = getHostVMCount(vmCount, socketCount, coresPerSocket, vmsPerCore, coresOverhead);
  var hostCores = getHostCoresCount(socketCount, coresPerSocket, coresOverhead);
  if (hostCores === 0) return "0.0";
  var r = (toFloat(vcpuCount) * toFloat(vcpuMHz) * toFloat(hostVMCount) / toFloat(hostCores)) / 1000;
  return r.toFixed(1);
}

// GetHostMemory - host memory in GB
function getHostMemory(vmCount, socketCount, coresPerSocket, coresOverhead, vmsPerCore, memorySize, memoryOverhead, displayCount, displayResolution, vcpuCount, videoRAM) {
  var hostVMCount = getHostVMCount(vmCount, socketCount, coresPerSocket, vmsPerCore, coresOverhead);
  var displayOverhead = getVMDisplayOverhead(displayCount, displayResolution, videoRAM);
  var vcpuMemOverhead = getVMVcpuMemoryOverhead(vcpuCount, memorySize);
  var r = Math.floor((hostVMCount * (toInt(memorySize) + displayOverhead.memory + vcpuMemOverhead)) / 1024) + toInt(memoryOverhead);
  return r;
}

// ============================================================
// Storage Calculations (port of storage/storage.go)
// ============================================================

// GetStorageCapacity - total storage capacity in TB
function getStorageCapacity(vmCount, diskSize, capacityOverhead, dedupeRatio, displayCount, displayResolution, videoRAM, memorySize, cloneRefreshRate) {
  var displayOverhead = getVMDisplayOverhead(displayCount, displayResolution, videoRAM);

  var effectiveDiskSize;
  if (String(cloneRefreshRate) !== "0") {
    effectiveDiskSize = toFloat(diskSize) * (toFloat(cloneRefreshRate) / 100);
  } else {
    effectiveDiskSize = toFloat(diskSize);
  }

  // memorySize is for swap, displayOverhead.storage converted MB→GB
  var r = toFloat(vmCount) * (effectiveDiskSize + (toFloat(memorySize) / 1000) + (toFloat(displayOverhead.storage) / 1000));

  if (String(capacityOverhead) !== "0") {
    r += (toFloat(capacityOverhead) / 100) * r;
  }

  if (String(dedupeRatio) !== "0") {
    r -= (toFloat(dedupeRatio) / 100) * r;
  }

  // Convert GB to TB
  return (r / 1000).toFixed(2);
}

// GetStorageDatastoreCount - number of datastores
function getStorageDatastoreCount(vmCount, datastoreVMCount) {
  return Math.ceil(toFloat(vmCount) / toFloat(datastoreVMCount));
}

// GetStorageDatastoreSize - size per datastore in TB
function getStorageDatastoreSize(vmCount, datastoreVMCount, diskSize, capacityOverhead, dedupeRatio, displayCount, displayResolution, videoRAM, memorySize, cloneRefreshRate) {
  var totalCapacity = toFloat(getStorageCapacity(vmCount, diskSize, capacityOverhead, dedupeRatio, displayCount, displayResolution, videoRAM, memorySize, cloneRefreshRate));
  var dsCount = getStorageDatastoreCount(vmCount, datastoreVMCount);
  if (dsCount === 0) return "0.00";
  return (totalCapacity / dsCount).toFixed(2);
}

// GetStorageDatastoreIops - IOps calculations
// Returns { dsFrontend, dsBackend, totalFrontend, totalBackend }
function getStorageDatastoreIops(iopsCount, iopsReadRatio, iopsBootCount, iopsBootReadRatio, datastoreVMCount, concurrentBootVMs, raidType, vmCount, datastoreVMCountForDS) {
  // Boot
  var dsFrontendBootIops = toInt(iopsBootCount) * toInt(concurrentBootVMs);
  var dsBackendBootReadIops = Math.floor(((toFloat(iopsBootReadRatio) / 100) * toFloat(iopsBootCount)) * toFloat(concurrentBootVMs));
  var dsBackendBootWriteIops = Math.floor(((1 - (toFloat(iopsBootReadRatio) / 100)) * toFloat(iopsBootCount)) * toFloat(concurrentBootVMs));

  // Steady state
  var dsFrontendIops = toInt(iopsCount) * toInt(datastoreVMCount);
  var dsBackendReadIops = Math.floor(((toFloat(iopsReadRatio) / 100) * toFloat(iopsCount)) * toFloat(datastoreVMCount));
  var dsBackendWriteIops = Math.floor(((1 - (toFloat(iopsReadRatio) / 100)) * toFloat(iopsCount)) * toFloat(datastoreVMCount));

  // RAID write amplification
  switch (String(raidType)) {
    case "5":
      dsBackendBootWriteIops *= 4;
      dsBackendWriteIops *= 4;
      break;
    case "6":
      dsBackendBootWriteIops *= 6;
      dsBackendWriteIops *= 6;
      break;
    case "10":
      dsBackendBootWriteIops *= 2;
      dsBackendWriteIops *= 2;
      break;
  }

  var dsBackendBootIops = dsBackendBootReadIops + dsBackendBootWriteIops;
  var dsBackendIops = dsBackendReadIops + dsBackendWriteIops;

  var dsCount = getStorageDatastoreCount(vmCount, datastoreVMCountForDS);
  var totalFrontendIops = (dsFrontendBootIops + dsFrontendIops) * dsCount;
  var totalBackendIops = (dsBackendBootIops + dsBackendIops) * dsCount;

  return {
    dsFrontend: dsFrontendBootIops + dsFrontendIops,
    dsBackend: dsBackendBootIops + dsBackendIops,
    totalFrontend: totalFrontendIops,
    totalBackend: totalBackendIops
  };
}

// ============================================================
// Virtualization Calculations (port of virtualization/virtualization.go)
// ============================================================

function getClusterSize(vmCount, socketCount, coresPerSocket, vmsPerCore, coresOverhead, clusterHostSize, clusterHA) {
  var hostCount = getHostCount(vmCount, socketCount, coresPerSocket, vmsPerCore, coresOverhead, clusterHA);
  return Math.ceil(toFloat(hostCount) / toFloat(clusterHostSize));
}

function getManagementServerCount(vmCount, maxVMsPerServer) {
  return Math.ceil(toFloat(vmCount) / toFloat(maxVMsPerServer));
}

// ============================================================
// Azure Instance Type (port of azure/azure.go)
// ============================================================

function getAzureInstanceType(vcpuCount, memorySize, diskSize, videoRAM) {
  var result = "";
  var memory = toFloat(memorySize) / 1024; // MB to GiB
  var vcpu = toInt(vcpuCount);
  var vram = toInt(videoRAM);

  // VM instance type based on cores and memory
  switch (vcpu) {
    case 1:
      if (memory <= 2.1) result = "F1";
      else if (memory <= 4.1) result = "F2";
      else if (memory <= 8.1) result = "F4";
      else if (memory <= 16.1) result = "F8";
      else result = "F16";
      break;
    case 2:
      if (memory <= 4.1) result = "F2";
      else if (memory <= 8.1) result = "F4";
      else if (memory <= 16.1) result = "F8";
      else result = "F16";
      break;
    case 4:
      if (memory <= 8.1) result = "F4";
      else if (memory <= 16.1) result = "F8";
      else result = "F16";
      break;
    case 8:
      if (memory <= 16) result = "F8";
      else result = "F16";
      break;
  }

  // GPU instance types override
  if (vram === 1) {
    switch (vcpu) {
      case 1:
      case 2:
      case 4:
        if (memory <= 14.1) result = "NV4as";
        else if (memory <= 28.1) result = "NV8as";
        break;
      case 8:
        if (memory <= 28.1) result = "NV8as";
        else if (memory <= 56.1) result = "NV16as";
        else result = "NV32as";
        break;
    }
  }

  // Disk instance type
  var c = toInt(diskSize);
  if (c <= 32) result += " P4";
  else if (c <= 64) result += " P6";
  else if (c <= 128) result += " P10";
  else if (c <= 256) result += " P15";
  else if (c <= 512) result += " P20";
  else if (c <= 1024) result += " P30";
  else result += " P40";

  return result;
}

// ============================================================
// Validation (port of validation/validation.go)
// ============================================================

function validateResults(data) {
  var errors = [];

  // VM memory size limit: 6,128,000 MB (vSphere max)
  if (toFloat(data.vmMemorySize) > 6128000) {
    errors.push("Warning: VM memory size above limit.");
    return errors;
  }

  // Host CPU clock limit: 4.2 GHz (Intel max)
  if (toFloat(data.hostClockUsed) > 4.2) {
    errors.push("Warning: Host CPU (GHz) above limit. (max=4.2)");
    return errors;
  }

  // VMs per host limit: 200 (Horizon 8)
  if (toFloat(data.hostVMCount) > 200) {
    errors.push("Warning: Number of VMs per host above limit. (max=200)");
    return errors;
  }

  // Datastore count limit: 500 (VMFS)
  if (toFloat(data.datastoreCount) > 500) {
    errors.push("Warning: Number of datastores above limit (max=500).");
    return errors;
  }

  return errors;
}

// ============================================================
// Main Calculate function (port of calculations/calculations.go)
// ============================================================

function calculate() {
  // Clear previous errors
  setResult("errorresults", "");

  // Read all form values
  var vmCount = val("vmcount");
  var vcpuCount = val("vmvcpucount");
  var vcpuMHz = val("vmvcpumhz");
  var vmsPerCore = val("vmpercorecount");
  var memorySize = val("vmmemorysize");
  var displayCount = val("vmdisplaycount");
  var displayResolution = val("vmdisplayresolution");
  var videoRAM = val("vmvideoram");
  var diskSize = val("vmdisksize");
  var iopsCount = val("vmiopscount");
  var iopsReadRatio = val("vmiopsreadratio");
  var iopsBootCount = val("vmiopsbootcount");
  var iopsBootReadRatio = val("vmiopsbootreadratio");
  var cloneRefreshRate = val("vmclonesizerefreshrate");

  var socketCount = val("hostsocketcount");
  var coresPerSocket = val("hostsocketcorescount");
  var memoryOverhead = val("hostmemoryoverhead");
  var coresOverhead = val("hostcoresoverhead");

  var capacityOverhead = val("storagecapacityoverhead");
  var datastoreVMCount = val("storagedatastorevmcount");
  var dedupeRatio = val("storagededuperatio");
  var raidType = val("storageraidtype");
  var concurrentBootVMs = val("storageconcurrentbootvmcount");

  var clusterHostSize = val("virtualizationclusterhostsize");
  var clusterHA = val("virtualizationclusterhostha");
  var mgmtServerVMMax = val("virtualizationmanagementservertvmcount");

  // Host calculations
  var hostCount = getHostCount(vmCount, socketCount, coresPerSocket, vmsPerCore, coresOverhead, clusterHA);
  var hostClockUsed = getHostClockUsed(vcpuCount, vcpuMHz, vmCount, socketCount, coresPerSocket, vmsPerCore, coresOverhead);
  var hostMemory = getHostMemory(vmCount, socketCount, coresPerSocket, coresOverhead, vmsPerCore, memorySize, memoryOverhead, displayCount, displayResolution, vcpuCount, videoRAM);
  var hostVMCount = getHostVMCount(vmCount, socketCount, coresPerSocket, vmsPerCore, coresOverhead);

  // Storage calculations
  var storageCapacity = getStorageCapacity(vmCount, diskSize, capacityOverhead, dedupeRatio, displayCount, displayResolution, videoRAM, memorySize, cloneRefreshRate);
  var datastoreCount = getStorageDatastoreCount(vmCount, datastoreVMCount);
  var datastoreSize = getStorageDatastoreSize(vmCount, datastoreVMCount, diskSize, capacityOverhead, dedupeRatio, displayCount, displayResolution, videoRAM, memorySize, cloneRefreshRate);
  var iops = getStorageDatastoreIops(iopsCount, iopsReadRatio, iopsBootCount, iopsBootReadRatio, datastoreVMCount, concurrentBootVMs, raidType, vmCount, datastoreVMCount);

  // Virtualization calculations
  var clusterCount = getClusterSize(vmCount, socketCount, coresPerSocket, vmsPerCore, coresOverhead, clusterHostSize, clusterHA);
  var mgmtServerCount = getManagementServerCount(vmCount, mgmtServerVMMax);

  // Azure recommendation
  var azureInstance = getAzureInstanceType(vcpuCount, memorySize, diskSize, videoRAM);

  // Validation
  var errors = validateResults({
    vmMemorySize: memorySize,
    hostClockUsed: hostClockUsed,
    hostVMCount: hostVMCount,
    datastoreCount: datastoreCount
  });

  if (errors.length > 0) {
    setResult("errorresults", errors[0]);
  }

  // Display results
  setResult("hostresultscount", hostCount);
  setResult("hostresultsclockused", hostClockUsed);
  setResult("hostresultsvmcount", hostVMCount);
  setResult("hostresultsmemory", hostMemory);
  setResult("storageresultscapacity", storageCapacity);
  setResult("storageresultsdatastorecount", datastoreCount);
  setResult("storageresultsdatastoresize", datastoreSize);
  setResult("storagedatastorefroentendiops", iops.dsFrontend);
  setResult("storagedatastorebackendiops", iops.dsBackend);
  setResult("storageresultsfrontendiops", iops.totalFrontend);
  setResult("storageresultsbackendiops", iops.totalBackend);
  setResult("virtualizationresultsclustercount", clusterCount);
  setResult("virtualizationresultsmanagementservercount", mgmtServerCount);
  setResult("azureresultsinstancetype", azureInstance);
}

// ============================================================
// VM Profiles (port of config/config.yml profiles)
// ============================================================

var profiles = {
  "1": { // Task Worker
    vcpucount: "1", vcpumhz: "500", vmpercorecount: "5", memorysize: "1536",
    displaycount: "1", displayresolution: "2", videoram: "0", disksize: "100",
    iopscount: "6", iopsreadratio: "20", iopsbootcount: "600",
    iopsbootreadratio: "20", clonesizerefreshrate: "0"
  },
  "2": { // Office Worker
    vcpucount: "1", vcpumhz: "500", vmpercorecount: "5", memorysize: "2048",
    displaycount: "1", displayresolution: "2", videoram: "64", disksize: "100",
    iopscount: "8", iopsreadratio: "20", iopsbootcount: "600",
    iopsbootreadratio: "20", clonesizerefreshrate: "0"
  },
  "3": { // Knowledge Worker
    vcpucount: "2", vcpumhz: "315", vmpercorecount: "4", memorysize: "2048",
    displaycount: "1", displayresolution: "2", videoram: "64", disksize: "100",
    iopscount: "9", iopsreadratio: "20", iopsbootcount: "600",
    iopsbootreadratio: "20", clonesizerefreshrate: "0"
  },
  "4": { // Power User
    vcpucount: "2", vcpumhz: "625", vmpercorecount: "2", memorysize: "4096",
    displaycount: "2", displayresolution: "3", videoram: "128", disksize: "100",
    iopscount: "11", iopsreadratio: "20", iopsbootcount: "600",
    iopsbootreadratio: "20", clonesizerefreshrate: "0"
  }
};

function loadProfile() {
  var profileId = val("vmprofile");
  var p = profiles[profileId];
  if (!p) return;

  document.getElementById("vmvcpucount").value = p.vcpucount;
  document.getElementById("vmvcpumhz").value = p.vcpumhz;
  document.getElementById("vmpercorecount").value = p.vmpercorecount;
  document.getElementById("vmmemorysize").value = p.memorysize;
  document.getElementById("vmdisplaycount").value = p.displaycount;
  document.getElementById("vmdisplayresolution").value = p.displayresolution;
  document.getElementById("vmvideoram").value = p.videoram;
  document.getElementById("vmdisksize").value = p.disksize;
  document.getElementById("vmiopscount").value = p.iopscount;
  document.getElementById("vmiopsreadratio").value = p.iopsreadratio;
  document.getElementById("vmiopsbootcount").value = p.iopsbootcount;
  document.getElementById("vmiopsbootreadratio").value = p.iopsbootreadratio;
  document.getElementById("vmclonesizerefreshrate").value = p.clonesizerefreshrate;
}

// ============================================================
// About Modal
// ============================================================

function showAbout() {
  $('#aboutModal').modal('show');
}
