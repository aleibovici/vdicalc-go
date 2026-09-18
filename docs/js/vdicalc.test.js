#!/usr/bin/env node
/**
 * Lightweight regression tests for docs/js/vdicalc.js calculation helpers.
 * Run: node docs/js/vdicalc.test.js
 */
const fs = require('fs');
const path = require('path');

let src = fs.readFileSync(path.join(__dirname, 'vdicalc.js'), 'utf8');
src = src.slice(0, src.indexOf('// Main Calculate function'));
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

console.log('=== Happy path (Task defaults) ===');
eq(getHostVMCount(2000, 2, 12, 5, 0), 120, 'VMs per host');
eq(getHostCount(2000, 2, 12, 5, 0, 0), 17, 'host count');
eq(getHostCount(2000, 2, 12, 5, 0, 1), 19, 'host count with HA');
eq(getHostClockUsed(1, 500, 2000, 2, 12, 5, 0), '2.5', 'host clock');
eq(getHostMemory(2000, 2, 12, 0, 5, 1536, 1, 1, 2, 1, 0), 187, 'host memory');
eq(getStorageCapacity(2000, 100, 5, 0, 1, 2, 0, 1536, 0), '213.46', 'storage capacity');
eq(getStorageDatastoreCount(2000, 100), 20, 'datastore count');
eq(getStorageDatastoreSize(2000, 100, 100, 5, 0, 1, 2, 0, 1536, 0), '10.67', 'datastore size');
const iops = getStorageDatastoreIops(6, 20, 600, 20, 100, 2, 6, 2000, 100);
eq(iops.dsFrontend, 1800, 'ds frontend iops');
eq(iops.dsBackend, 9000, 'ds backend iops');
eq(iops.totalFrontend, 36000, 'total frontend iops');
eq(iops.totalBackend, 180000, 'total backend iops');
eq(getClusterSize(2000, 2, 12, 5, 0, 8, 0), 3, 'clusters');
eq(getManagementServerCount(2000, 2000), 1, 'mgmt servers');
eq(getAzureInstanceType(1, 1536, 100, 0), 'F1 P10', 'azure');

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
// Clamped to 100% read → write amp 0; BE = reads only
// boot read: 600*2=1200; steady read: 6*100=600; BE=1800
eq(iopsHigh.dsBackend, 1800, 'B4: ratio clamped to 100% → read-only backend');

const iopsNeg = getStorageDatastoreIops(6, -10, 600, -10, 100, 2, 6, 2000, 100);
// Clamped to 0% read → all writes * RAID6: boot 1200*6 + steady 600*6 = 10800
eq(iopsNeg.dsBackend, 10800, 'B4: ratio -10% clamped to 0 → all-write backend');

console.log('\nResult:', fails ? fails + ' failure(s)' : 'all passed');
process.exit(fails ? 1 : 0);
