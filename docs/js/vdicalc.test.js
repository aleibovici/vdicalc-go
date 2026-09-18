#!/usr/bin/env node
/**
 * Lightweight regression tests for docs/js/vdicalc.js calculation helpers.
 * Run: node docs/js/vdicalc.test.js
 */
const fs = require('fs');
const path = require('path');

let src = fs.readFileSync(path.join(__dirname, 'vdicalc.js'), 'utf8');
// Keep helpers through profiles; drop DOM bootstrap
const cut = src.indexOf('// UI Interactions');
src = cut > 0 ? src.slice(0, cut) : src.slice(0, src.indexOf('// Main Calculate function'));
global.document = {
  getElementById: () => ({ value: '', textContent: '', classList: { add() {}, remove() {}, toggle() {} }, style: {} }),
  querySelectorAll: () => [],
  addEventListener() {}
};
eval(src);

let fails = 0;
function eq(actual, expected, msg) {
  const ok = String(actual) === String(expected);
  if (!ok) {
    fails++;
    console.log('FAIL:', msg, '| expected', expected, 'got', actual);
  } else {
    console.log('PASS:', msg);
  }
}

console.log('=== Happy path (Win11 Task defaults: 2×16 hosts, 4 VMs/core) ===');
eq(getHostVMCount(2000, 2, 16, 4, 0), 128, 'VMs per host');
eq(getHostCount(2000, 2, 16, 4, 0, 0), 16, 'host count');
eq(getHostCount(2000, 2, 16, 4, 0, 1), 18, 'host count with HA');
eq(getHostClockUsed(2, 350, 2000, 2, 16, 4, 0), '2.8', 'host clock');
eq(getHostMemory(2000, 2, 16, 0, 4, 4096, 4, 1, 2, 2, 0), 523, 'host memory');
eq(getStorageCapacity(2000, 128, 5, 0, 1, 2, 0, 4096, 0), '277.63', 'storage capacity');
eq(getStorageDatastoreCount(2000, 100), 20, 'datastore count');
eq(getStorageDatastoreSize(2000, 100, 128, 5, 0, 1, 2, 0, 4096, 0), '13.88', 'datastore size');
const iops = getStorageDatastoreIops(8, 20, 50, 20, 100, 2, 6, 2000, 100);
eq(iops.dsFrontend, 900, 'ds frontend iops');
eq(iops.dsBackend, 4500, 'ds backend iops');
eq(iops.totalFrontend, 18000, 'total frontend iops');
eq(iops.totalBackend, 90000, 'total backend iops');
eq(getClusterSize(2000, 2, 16, 4, 0, 8, 0), 2, 'clusters');
eq(getManagementServerCount(2000, 2000), 1, 'mgmt servers');
eq(getAzureInstanceType(2, 4096, 128, 0), 'D2s_v5 P10', 'azure task');
eq(getAzureInstanceType(2, 8192, 128, 64), 'D2s_v5 P10', 'azure office');
eq(getAzureInstanceType(4, 16384, 128, 128), 'D4s_v5 P10', 'azure knowledge');
eq(getAzureInstanceType(8, 32768, 256, 1), 'NV12ads_A10_v5 P15', 'azure power GPU');

console.log('\n=== Profiles object (Win11 baselines) ===');
eq(profiles['1'].memorysize, '4096', 'task RAM');
eq(profiles['1'].vcpucount, '2', 'task vCPU');
eq(profiles['2'].memorysize, '8192', 'office RAM');
eq(profiles['3'].memorysize, '16384', 'knowledge RAM');
eq(profiles['4'].memorysize, '32768', 'power RAM');
eq(profiles['4'].videoram, '1', 'power GPU');

console.log('\n=== Edge-case guards ===');
eq(getHostCoresCount(1, 2, 6), 0, 'B1: cores overhead > cores → 0');
eq(getHostVMCount(2000, 1, 2, 5, 6), 0, 'B1: VMs/host → 0');
eq(getHostCount(2000, 1, 2, 5, 6, 0), 0, 'B1: host count → 0');
eq(getHostMemory(2000, 1, 2, 6, 5, 1536, 1, 1, 2, 1, 0), 1, 'B1: host memory → overhead only');
eq(getHostClockUsed(1, 500, 2000, 1, 2, 5, 6), '0.0', 'B1: host clock → 0.0');

eq(getStorageDatastoreCount(2000, 0), 0, 'B2: VMs/DS=0 → 0 datastores');
eq(getStorageDatastoreCount(2000, ''), 0, 'B2: VMs/DS empty → 0');
eq(getStorageDatastoreSize(2000, 0, 100, 5, 0, 1, 2, 0, 1536, 0), '0.00', 'B2: DS size → 0.00');
const iopsZeroDs = getStorageDatastoreIops(6, 20, 600, 20, 0, 2, 6, 2000, 0);
eq(iopsZeroDs.totalFrontend, 0, 'B2: total FE IOps not Infinity');
eq(iopsZeroDs.totalBackend, 0, 'B2: total BE IOps not Infinity');

eq(getManagementServerCount(2000, 0), 0, 'B3: max VMs=0 → 0');
eq(getManagementServerCount(2000, ''), 0, 'B3: max VMs empty → 0');
eq(getManagementServerCount(2000, -5), 0, 'B3: max VMs negative → 0');
eq(getClusterSize(2000, 2, 12, 5, 0, 0, 0), 0, 'cluster size 0 → 0');

const iopsHigh = getStorageDatastoreIops(6, 150, 600, 150, 100, 2, 6, 2000, 100);
eq(iopsHigh.dsBackend >= 0, 'true', 'B4: ratio 150% → non-negative backend');
eq(iopsHigh.dsFrontend, 1800, 'B4: frontend unchanged at 150% (clamped reads)');
eq(iopsHigh.dsBackend, 1800, 'B4: ratio clamped to 100% → read-only backend');

const iopsNeg = getStorageDatastoreIops(6, -10, 600, -10, 100, 2, 6, 2000, 100);
eq(iopsNeg.dsBackend, 10800, 'B4: ratio -10% clamped to 0 → all-write backend');

const vClock = validateResults({ vmMemorySize: '4096', hostClockUsed: '6.0', hostVMCount: '120', datastoreCount: '20' });
eq(vClock[0], 'Warning: Host CPU (GHz) above typical turbo limit. (max≈5.5)', 'clock warning uses 5.5 GHz ceiling');

console.log('\nResult:', fails ? fails + ' failure(s)' : 'all passed');
process.exit(fails ? 1 : 0);
