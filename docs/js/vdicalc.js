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

// Non-negative helpers for counts / sizes that must not go below zero
function toNonNegInt(value) {
  return Math.max(0, toInt(value));
}

function toNonNegFloat(value) {
  return Math.max(0, toFloat(value));
}

// Clamp a percentage-like input into [0, 100]
function clampPercent(value) {
  var n = toFloat(value);
  if (n < 0) return 0;
  if (n > 100) return 100;
  return n;
}

function val(id) {
  return document.getElementById(id).value;
}

function setResult(id, value) {
  document.getElementById(id).textContent = value;
}

// ============================================================
// Explicit domain modes (form still stores legacy string sentinels)
// ============================================================

// Graphics regimes encoded historically as #vmvideoram values:
//   "0" = none, "1" = hardware GPU, "64"|"128"|"256"|"512" = software VRAM MB
function parseGraphicsMode(videoRAM) {
  var raw = String(videoRAM);
  if (raw === "0") return { mode: "none", softwareVramMb: 0, raw: raw };
  if (raw === "1") return { mode: "gpu", softwareVramMb: 0, raw: raw };
  return { mode: "software", softwareVramMb: toInt(raw), raw: raw };
}

// Host HA inflate factor (N+1-ish capacity for cluster HA = Yes)
var HA_HOST_FACTOR = 1.125;

// RAID write amplification multipliers (backend write IOps)
var RAID_WRITE_AMPLIFICATION = {
  "5": 4,
  "6": 6,
  "10": 2
};

// Display overhead when graphics mode is none: [displayCount][resolutionCode] → MB
// Resolution codes: "1"=1280×800, "2"=1920×1200, "3"=2560×1600
var DISPLAY_OVERHEAD_NONE = {
  1: { "1": { memory: 4, storage: 107 }, "2": { memory: 8, storage: 111 }, "3": { memory: 16, storage: 203 } },
  2: { "1": { memory: 13, storage: 163 }, "2": { memory: 26, storage: 190 }, "3": { memory: 60, storage: 203 } },
  3: { "1": { memory: 19, storage: 207 }, "2": { memory: 38, storage: 248 }, "3": { memory: 85, storage: 461 } },
  4: { "1": { memory: 25, storage: 252 }, "2": { memory: 51, storage: 306 }, "3": { memory: 110, storage: 589 } }
};

var SOFTWARE_VRAM_STORAGE_OVERHEAD = {
  "64": 1076,
  "128": 1468,
  "256": 1468,
  "512": 1916
};

// ============================================================
// VM Overhead Calculations (port of vm/vm.go)
// ============================================================

// GetVMDisplayOverhead - calculates display/resolution overhead for memory and storage vswap
// Returns { memory: int, storage: int } in MB
function getVMDisplayOverhead(displayCount, displayResolution, videoRAM) {
  var graphics = parseGraphicsMode(videoRAM);
  var count = toInt(displayCount);
  var resolution = String(displayResolution);

  if (graphics.mode === "none") {
    var byRes = DISPLAY_OVERHEAD_NONE[count];
    if (byRes && byRes[resolution]) {
      return { memory: byRes[resolution].memory, storage: byRes[resolution].storage };
    }
    return { memory: 0, storage: 0 };
  }

  if (graphics.mode === "gpu") {
    return { memory: 96, storage: 0 };
  }

  // software 3D graphics with specific video RAM
  var storage = SOFTWARE_VRAM_STORAGE_OVERHEAD[graphics.raw] || 0;
  return { memory: graphics.softwareVramMb, storage: storage };
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
  // Cores overhead must not produce a negative usable core count
  return Math.max(0, (toNonNegInt(socketCount) * toNonNegInt(coresPerSocket)) - toNonNegInt(coresOverhead));
}

// GetHostVMCount - number of VMs per host
function getHostVMCount(vmCount, socketCount, coresPerSocket, vmsPerCore, coresOverhead) {
  var capacity = getHostCoresCount(socketCount, coresPerSocket, coresOverhead) * toNonNegInt(vmsPerCore);
  if (capacity <= 0) return 0;
  return Math.min(toNonNegInt(vmCount), capacity);
}

// GetHostCount - number of hosts needed
function getHostCount(vmCount, socketCount, coresPerSocket, vmsPerCore, coresOverhead, clusterHA) {
  var hostVMCount = getHostVMCount(vmCount, socketCount, coresPerSocket, vmsPerCore, coresOverhead);
  if (hostVMCount <= 0) return 0;
  var r = toNonNegFloat(vmCount) / toFloat(hostVMCount);
  if (String(clusterHA) === "1") {
    r *= HA_HOST_FACTOR;
  }
  // Go uses FormatFloat with precision 0, which rounds to nearest
  return Number(r.toFixed(0));
}

// GetHostClockUsed - host CPU clock in GHz
function getHostClockUsed(vcpuCount, vcpuMHz, vmCount, socketCount, coresPerSocket, vmsPerCore, coresOverhead) {
  var hostVMCount = getHostVMCount(vmCount, socketCount, coresPerSocket, vmsPerCore, coresOverhead);
  var hostCores = getHostCoresCount(socketCount, coresPerSocket, coresOverhead);
  if (hostCores <= 0 || hostVMCount <= 0) return "0.0";
  var r = (toNonNegFloat(vcpuCount) * toNonNegFloat(vcpuMHz) * toFloat(hostVMCount) / toFloat(hostCores)) / 1000;
  return r.toFixed(1);
}

// GetHostMemory - host memory in GB
function getHostMemory(vmCount, socketCount, coresPerSocket, coresOverhead, vmsPerCore, memorySize, memoryOverhead, displayCount, displayResolution, vcpuCount, videoRAM) {
  var hostVMCount = getHostVMCount(vmCount, socketCount, coresPerSocket, vmsPerCore, coresOverhead);
  if (hostVMCount <= 0) return toNonNegInt(memoryOverhead);
  var displayOverhead = getVMDisplayOverhead(displayCount, displayResolution, videoRAM);
  var vcpuMemOverhead = getVMVcpuMemoryOverhead(vcpuCount, memorySize);
  var r = Math.floor((hostVMCount * (toNonNegInt(memorySize) + displayOverhead.memory + vcpuMemOverhead)) / 1024) + toNonNegInt(memoryOverhead);
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

  // memorySize is for swap, displayOverhead.storage converted MB->GB
  var r = toNonNegFloat(vmCount) * (effectiveDiskSize + (toNonNegFloat(memorySize) / 1000) + (toFloat(displayOverhead.storage) / 1000));

  if (String(capacityOverhead) !== "0") {
    r += (clampPercent(capacityOverhead) / 100) * r;
  }

  if (String(dedupeRatio) !== "0") {
    r -= (clampPercent(dedupeRatio) / 100) * r;
  }

  // Convert GB to TB
  return (r / 1000).toFixed(2);
}

// GetStorageDatastoreCount - number of datastores
function getStorageDatastoreCount(vmCount, datastoreVMCount) {
  var perDatastore = toNonNegFloat(datastoreVMCount);
  var vms = toNonNegFloat(vmCount);
  if (perDatastore <= 0 || vms <= 0) return 0;
  return Math.ceil(vms / perDatastore);
}

// GetStorageDatastoreSize - size per datastore in TB
function getStorageDatastoreSize(vmCount, datastoreVMCount, diskSize, capacityOverhead, dedupeRatio, displayCount, displayResolution, videoRAM, memorySize, cloneRefreshRate) {
  var totalCapacity = toFloat(getStorageCapacity(vmCount, diskSize, capacityOverhead, dedupeRatio, displayCount, displayResolution, videoRAM, memorySize, cloneRefreshRate));
  var dsCount = getStorageDatastoreCount(vmCount, datastoreVMCount);
  if (dsCount <= 0) return "0.00";
  return (totalCapacity / dsCount).toFixed(2);
}

// GetStorageDatastoreIops - IOps calculations
// Returns { dsFrontend, dsBackend, totalFrontend, totalBackend }
function getStorageDatastoreIops(iopsCount, iopsReadRatio, iopsBootCount, iopsBootReadRatio, datastoreVMCount, concurrentBootVMs, raidType, vmCount, datastoreVMCountForDS) {
  var readRatio = clampPercent(iopsReadRatio);
  var bootReadRatio = clampPercent(iopsBootReadRatio);
  var dsVMs = toNonNegInt(datastoreVMCount);
  var bootVMs = toNonNegInt(concurrentBootVMs);
  var steadyIops = toNonNegInt(iopsCount);
  var bootIops = toNonNegInt(iopsBootCount);

  // Boot
  var dsFrontendBootIops = bootIops * bootVMs;
  var dsBackendBootReadIops = Math.floor(((bootReadRatio / 100) * bootIops) * bootVMs);
  var dsBackendBootWriteIops = Math.floor(((1 - (bootReadRatio / 100)) * bootIops) * bootVMs);

  // Steady state
  var dsFrontendIops = steadyIops * dsVMs;
  var dsBackendReadIops = Math.floor(((readRatio / 100) * steadyIops) * dsVMs);
  var dsBackendWriteIops = Math.floor(((1 - (readRatio / 100)) * steadyIops) * dsVMs);

  // RAID write amplification (explicit map; unknown RAID → no multiplier)
  var raidAmp = RAID_WRITE_AMPLIFICATION[String(raidType)];
  if (raidAmp) {
    dsBackendBootWriteIops *= raidAmp;
    dsBackendWriteIops *= raidAmp;
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
  var hostsPerCluster = toNonNegFloat(clusterHostSize);
  if (hostsPerCluster <= 0 || hostCount <= 0) return 0;
  return Math.ceil(toFloat(hostCount) / hostsPerCluster);
}

function getManagementServerCount(vmCount, maxVMsPerServer) {
  var maxVMs = toNonNegFloat(maxVMsPerServer);
  var vms = toNonNegFloat(vmCount);
  if (maxVMs <= 0 || vms <= 0) return 0;
  return Math.ceil(vms / maxVMs);
}

// ============================================================
// Azure Instance Type
// Updated for AVD single-session baselines (D*s_v5 / NVads_A10_v5)
// Source: Microsoft session-host sizing guidelines (2024+)
// ============================================================

function getAzureInstanceType(vcpuCount, memorySize, diskSize, videoRAM) {
  var result = "";
  var memory = toFloat(memorySize) / 1024; // MB to GiB
  var vcpu = toInt(vcpuCount);

  // General-purpose Ds_v5 family (replaces legacy F-series defaults)
  // Windows 11 single-session typically starts at 2 vCPUs.
  switch (vcpu) {
    case 1:
      // Legacy 1-vCPU profiles: recommend minimum modern 2-vCPU SKU
      if (memory <= 8.1) result = "D2s_v5";
      else if (memory <= 16.1) result = "D4s_v5";
      else if (memory <= 32.1) result = "D8s_v5";
      else result = "D16s_v5";
      break;
    case 2:
      if (memory <= 8.1) result = "D2s_v5";
      else if (memory <= 16.1) result = "D4s_v5";
      else if (memory <= 32.1) result = "D8s_v5";
      else result = "D16s_v5";
      break;
    case 4:
      if (memory <= 16.1) result = "D4s_v5";
      else if (memory <= 32.1) result = "D8s_v5";
      else if (memory <= 64.1) result = "D16s_v5";
      else result = "D32s_v5";
      break;
    case 8:
      if (memory <= 32.1) result = "D8s_v5";
      else if (memory <= 64.1) result = "D16s_v5";
      else result = "D32s_v5";
      break;
    default:
      result = "D4s_v5";
      break;
  }

  // GPU-accelerated NVads A10 v5 only when graphics mode is hardware GPU
  // (software VRAM MB values intentionally stay on Ds_v5 — same as legacy sentinel "1")
  var graphics = parseGraphicsMode(videoRAM);
  if (graphics.mode === "gpu") {
    switch (vcpu) {
      case 1:
      case 2:
      case 4:
        if (memory <= 55.1) result = "NV6ads_A10_v5";
        else if (memory <= 110.1) result = "NV12ads_A10_v5";
        else result = "NV18ads_A10_v5";
        break;
      case 8:
        if (memory <= 110.1) result = "NV12ads_A10_v5";
        else if (memory <= 220.1) result = "NV18ads_A10_v5";
        else result = "NV36ads_A10_v5";
        break;
      default:
        result = "NV6ads_A10_v5";
        break;
    }
  }

  // Premium SSD managed disk tier
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

// Ordered soft-limit rules. Returns explicit ValidationState:
//   { status: 'ok'|'warn', issues: string[] }
// Collects every matching issue (legacy early-return hid later limits).
function validateResults(data) {
  var issues = [];

  // VM memory size limit: 6,128,000 MB (vSphere max)
  if (toFloat(data.vmMemorySize) > 6128000) {
    issues.push("Warning: VM memory size above limit.");
  }

  // Host CPU clock limit: ~5.5 GHz modern turbo ceiling (soft warning)
  if (toFloat(data.hostClockUsed) > 5.5) {
    issues.push("Warning: Host CPU (GHz) above typical turbo limit. (max≈5.5)");
  }

  // VMs per host limit: 200 (Horizon 8)
  if (toFloat(data.hostVMCount) > 200) {
    issues.push("Warning: Number of VMs per host above limit. (max=200)");
  }

  // Datastore count limit: 500 (VMFS)
  if (toFloat(data.datastoreCount) > 500) {
    issues.push("Warning: Number of datastores above limit (max=500).");
  }

  return {
    status: issues.length ? "warn" : "ok",
    issues: issues
  };
}

function applyValidationBanner(validation) {
  var errorBanner = document.getElementById("errorBanner");
  var errorText = document.getElementById("errorText");
  if (!errorBanner || !errorText) return;

  if (!validation || validation.status === "ok" || !validation.issues.length) {
    errorBanner.classList.remove("visible");
    errorBanner.setAttribute("role", "status");
    errorText.textContent = "";
    return;
  }

  errorText.textContent = validation.issues.join(" · ");
  errorBanner.classList.add("visible");
  errorBanner.setAttribute("role", "alert");
}

// ============================================================
// Main Calculate function (port of calculations/calculations.go)
// ============================================================

function calculate() {
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

  // Azure recommendation (always-on parallel path — not an Azure-vs-on-prem UI mode)
  var azureInstance = getAzureInstanceType(vcpuCount, memorySize, diskSize, videoRAM);

  // Validation (explicit status; does not alter sizing numbers)
  var validation = validateResults({
    vmMemorySize: memorySize,
    hostClockUsed: hostClockUsed,
    hostVMCount: hostVMCount,
    datastoreCount: datastoreCount
  });
  applyValidationBanner(validation);

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

// Windows 11 single-session baselines (Microsoft AVD + Omnissa/Dell Horizon guidance).
// Density (~vCPU per core) aligned with modern Xeon/EPYC dual-socket hosts.
var profiles = {
  "1": { // Task Worker — light / data entry
    vcpucount: "2", vcpumhz: "350", vmpercorecount: "4", memorysize: "4096",
    displaycount: "1", displayresolution: "2", videoram: "0", disksize: "128",
    iopscount: "8", iopsreadratio: "20", iopsbootcount: "50",
    iopsbootreadratio: "20", clonesizerefreshrate: "0"
  },
  "2": { // Office Worker — Outlook / Office / browser / light Teams
    vcpucount: "2", vcpumhz: "400", vmpercorecount: "3", memorysize: "8192",
    displaycount: "1", displayresolution: "2", videoram: "64", disksize: "128",
    iopscount: "12", iopsreadratio: "20", iopsbootcount: "60",
    iopsbootreadratio: "20", clonesizerefreshrate: "0"
  },
  "3": { // Knowledge Worker — heavier Office + web / light creative
    vcpucount: "4", vcpumhz: "400", vmpercorecount: "2", memorysize: "16384",
    displaycount: "1", displayresolution: "2", videoram: "128", disksize: "128",
    iopscount: "15", iopsreadratio: "20", iopsbootcount: "80",
    iopsbootreadratio: "20", clonesizerefreshrate: "0"
  },
  "4": { // Power User — CAD / media / GPU path
    vcpucount: "8", vcpumhz: "500", vmpercorecount: "1", memorysize: "32768",
    displaycount: "2", displayresolution: "3", videoram: "1", disksize: "256",
    iopscount: "25", iopsreadratio: "20", iopsbootcount: "100",
    iopsbootreadratio: "20", clonesizerefreshrate: "0"
  }
};

// Profile field keys → form control IDs (VM section only; host/storage/virt stay shared)
var PROFILE_FIELD_MAP = {
  vcpucount: "vmvcpucount",
  vcpumhz: "vmvcpumhz",
  vmpercorecount: "vmpercorecount",
  memorysize: "vmmemorysize",
  displaycount: "vmdisplaycount",
  displayresolution: "vmdisplayresolution",
  videoram: "vmvideoram",
  disksize: "vmdisksize",
  iopscount: "vmiopscount",
  iopsreadratio: "vmiopsreadratio",
  iopsbootcount: "vmiopsbootcount",
  iopsbootreadratio: "vmiopsbootreadratio",
  clonesizerefreshrate: "vmclonesizerefreshrate"
};

function profileMatchesForm(profileId) {
  var p = profiles[profileId];
  if (!p) return false;
  var key;
  for (key in PROFILE_FIELD_MAP) {
    if (!Object.prototype.hasOwnProperty.call(PROFILE_FIELD_MAP, key)) continue;
    var el = document.getElementById(PROFILE_FIELD_MAP[key]);
    if (!el || String(el.value) !== String(p[key])) return false;
  }
  return true;
}

function loadProfileById(profileId) {
  var p = profiles[profileId];
  if (!p) return;

  var key;
  for (key in PROFILE_FIELD_MAP) {
    if (!Object.prototype.hasOwnProperty.call(PROFILE_FIELD_MAP, key)) continue;
    document.getElementById(PROFILE_FIELD_MAP[key]).value = p[key];
  }

  setProfileSelection(profileId);
  calculate();
}

// ============================================================
// Explicit UI state machine (profile + about modal)
// ============================================================

// profileId: '1'|'2'|'3'|'4' for presets, null for custom (form drifted from preset)
// aboutOpen: boolean driving modal visibility / body scroll / focus
var uiState = {
  profileId: "1",
  aboutOpen: false
};

var aboutLastFocus = null;

function setProfileSelection(profileId) {
  uiState.profileId = profileId || null;

  var tabs = document.querySelectorAll(".chip[data-profile], .profile-tab[data-profile]");
  tabs.forEach(function (tab) {
    var id = tab.getAttribute("data-profile");
    var selected = uiState.profileId !== null && id === uiState.profileId;
    tab.classList.toggle("active", selected);
    tab.setAttribute("aria-selected", selected ? "true" : "false");
    tab.tabIndex = selected ? 0 : -1;
  });

  var customBadge = document.getElementById("profileCustomBadge");
  if (customBadge) {
    var isCustom = uiState.profileId === null;
    customBadge.hidden = !isCustom;
    customBadge.setAttribute("aria-hidden", isCustom ? "false" : "true");
  }

  var tablist = document.getElementById("profileTabs");
  if (tablist) {
    tablist.setAttribute(
      "aria-label",
      uiState.profileId === null ? "Workload profile (custom values)" : "Workload profile"
    );
  }
}

function syncProfileSelectionFromForm() {
  if (uiState.profileId !== null && profileMatchesForm(uiState.profileId)) {
    return;
  }
  // Prefer rematching another preset if the user landed exactly on one
  var id;
  for (id in profiles) {
    if (!Object.prototype.hasOwnProperty.call(profiles, id)) continue;
    if (profileMatchesForm(id)) {
      setProfileSelection(id);
      return;
    }
  }
  setProfileSelection(null);
}

// ============================================================
// UI Interactions (vanilla JS - no jQuery)
// ============================================================

// Section collapse/expand — keep aria-expanded in sync with .collapsed
function toggleSection(sectionId) {
  var section = document.getElementById(sectionId);
  if (!section) return;
  var collapsed = section.classList.toggle("collapsed");
  var head = section.querySelector(".panel__head");
  if (head) head.setAttribute("aria-expanded", collapsed ? "false" : "true");
}

// About modal — explicit open/closed state
function showAbout() {
  if (uiState.aboutOpen) return;
  uiState.aboutOpen = true;
  aboutLastFocus = document.activeElement;

  var modal = document.getElementById("aboutModal");
  if (!modal) return;
  modal.classList.add("visible");
  modal.setAttribute("aria-hidden", "false");
  document.body.style.overflow = "hidden";

  var closeBtn = modal.querySelector(".modal-content__close");
  if (closeBtn && typeof closeBtn.focus === "function") closeBtn.focus();
}

function hideAbout() {
  if (!uiState.aboutOpen) return;
  uiState.aboutOpen = false;

  var modal = document.getElementById("aboutModal");
  if (modal) {
    modal.classList.remove("visible");
    modal.setAttribute("aria-hidden", "true");
  }
  document.body.style.overflow = "";

  if (aboutLastFocus && typeof aboutLastFocus.focus === "function") {
    aboutLastFocus.focus();
  }
  aboutLastFocus = null;
}

// ============================================================
// Initialize on DOM ready
// ============================================================

document.addEventListener("DOMContentLoaded", function () {
  // Align JS profile state with the chip marked active in HTML (Task by default)
  var initialActive = document.querySelector(".chip.active[data-profile], .profile-tab.active[data-profile]");
  setProfileSelection(initialActive ? initialActive.getAttribute("data-profile") : "1");

  // Ensure panel heads expose expand state
  document.querySelectorAll(".panel").forEach(function (panel) {
    var head = panel.querySelector(".panel__head");
    if (head && !head.hasAttribute("aria-expanded")) {
      head.setAttribute("aria-expanded", panel.classList.contains("collapsed") ? "false" : "true");
    }
  });

  // Profile tab click handlers
  var tabs = document.querySelectorAll(".chip[data-profile], .profile-tab[data-profile]");
  tabs.forEach(function (tab) {
    tab.addEventListener("click", function () {
      loadProfileById(tab.getAttribute("data-profile"));
    });
  });

  // Arrow-key roving tabindex within the profile tablist
  var tablist = document.getElementById("profileTabs");
  if (tablist) {
    tablist.addEventListener("keydown", function (e) {
      var keys = ["ArrowLeft", "ArrowRight", "Home", "End"];
      if (keys.indexOf(e.key) === -1) return;
      var items = Array.prototype.slice.call(tabs);
      if (!items.length) return;
      var currentIndex = items.indexOf(document.activeElement);
      if (currentIndex < 0) currentIndex = items.findIndex(function (t) { return t.classList.contains("active"); });
      if (currentIndex < 0) currentIndex = 0;
      var next = currentIndex;
      if (e.key === "ArrowRight") next = (currentIndex + 1) % items.length;
      else if (e.key === "ArrowLeft") next = (currentIndex - 1 + items.length) % items.length;
      else if (e.key === "Home") next = 0;
      else if (e.key === "End") next = items.length - 1;
      e.preventDefault();
      items[next].focus();
      items[next].click();
    });
  }

  // Auto-calculate on any form input change; detect profile drift → custom
  var form = document.getElementById("vdiform");
  if (form) {
    var onFormEdit = function () {
      syncProfileSelectionFromForm();
      calculate();
    };
    form.addEventListener("input", onFormEdit);
    form.addEventListener("change", onFormEdit);
  }

  // Close modal on overlay click
  var modalOverlay = document.getElementById("aboutModal");
  if (modalOverlay) {
    modalOverlay.setAttribute("aria-hidden", "true");
    modalOverlay.addEventListener("click", function (e) {
      if (e.target === modalOverlay) {
        hideAbout();
      }
    });
  }

  // Close modal on Escape key only when open
  document.addEventListener("keydown", function (e) {
    if (e.key === "Escape" && uiState.aboutOpen) {
      hideAbout();
    }
  });
});
