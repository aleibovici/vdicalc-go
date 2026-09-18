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
    r *= 1.125;
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
    default:
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
  var vram = toInt(videoRAM);

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

  // GPU-accelerated NVads A10 v5 (AVD graphics / power workloads)
  if (vram === 1) {
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

function validateResults(data) {
  var errors = [];

  // VM memory size limit: 6,128,000 MB (vSphere max)
  if (toFloat(data.vmMemorySize) > 6128000) {
    errors.push("Warning: VM memory size above limit.");
    return errors;
  }

  // Host CPU clock limit: ~5.5 GHz modern turbo ceiling (soft warning)
  if (toFloat(data.hostClockUsed) > 5.5) {
    errors.push("Warning: Host CPU (GHz) above typical turbo limit. (max≈5.5)");
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
  var errorBanner = document.getElementById("errorBanner");
  var errorText = document.getElementById("errorText");
  if (errorBanner) errorBanner.classList.remove("visible");
  if (errorText) errorText.textContent = "";

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

  if (errors.length > 0 && errorBanner && errorText) {
    errorText.textContent = errors[0];
    errorBanner.classList.add("visible");
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

function loadProfileById(profileId) {
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

  calculate();
}

// ============================================================
// UI Interactions (vanilla JS - no jQuery)
// ============================================================

// Section collapse/expand
function toggleSection(sectionId) {
  var section = document.getElementById(sectionId);
  if (section) {
    section.classList.toggle("collapsed");
  }
}

// About modal
function showAbout() {
  document.getElementById("aboutModal").classList.add("visible");
  document.body.style.overflow = "hidden";
}

function hideAbout() {
  document.getElementById("aboutModal").classList.remove("visible");
  document.body.style.overflow = "";
}

// ============================================================
// Initialize on DOM ready
// ============================================================

document.addEventListener("DOMContentLoaded", function () {
  // Profile tab click handlers
  var tabs = document.querySelectorAll(".chip[data-profile], .profile-tab[data-profile]");
  tabs.forEach(function (tab) {
    tab.addEventListener("click", function () {
      tabs.forEach(function (t) {
        t.classList.remove("active");
        t.setAttribute("aria-selected", "false");
      });
      tab.classList.add("active");
      tab.setAttribute("aria-selected", "true");
      loadProfileById(tab.getAttribute("data-profile"));
    });
  });

  // Auto-calculate on any form input change
  var form = document.getElementById("vdiform");
  if (form) {
    form.addEventListener("input", function () {
      calculate();
    });
    form.addEventListener("change", function () {
      calculate();
    });
  }

  // Close modal on overlay click
  var modalOverlay = document.getElementById("aboutModal");
  if (modalOverlay) {
    modalOverlay.addEventListener("click", function (e) {
      if (e.target === modalOverlay) {
        hideAbout();
      }
    });
  }

  // Close modal on Escape key
  document.addEventListener("keydown", function (e) {
    if (e.key === "Escape") {
      hideAbout();
    }
  });
});
