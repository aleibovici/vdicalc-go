# VDI Calculator

Free browser-based tool for sizing Virtual Desktop Infrastructure deployments. It estimates host, storage, virtualization, and Azure instance requirements from your desktop and host inputs.

The current app is a static site in `docs/` (GitHub Pages). Calculations run in the browser — no server or database required. A legacy Go server remains in the repo for reference.

## Features

- Worker profiles for Windows 11 single-session VDI: Task, Office, Knowledge, and Power
- Host CPU, memory, and capacity sizing for modern dual-socket Xeon/EPYC hosts
- Storage capacity, datastore count, and frontend/backend IOps (including RAID write amplification)
- Cluster and management-server counts
- Azure recommendations on the Ds_v5 and NVads_A10_v5 families
- Input validation with warning messages
- Live results as you edit, plus print-friendly output

Profile and Azure defaults follow Microsoft AVD session-host guidance and published Omnissa/Dell Horizon Windows 11 density studies. Treat the output as a starting estimate and validate with a pilot.

## Usage

Open the GitHub Pages site, or open `docs/index.html` locally. Pick a profile or edit the inputs. Results update automatically.

```bash
node docs/js/vdicalc.test.js
```

## Project structure

- `docs/` — GitHub Pages site
  - `index.html` — Calculator UI
  - `js/vdicalc.js` — Calculation engine
  - `js/vdicalc.test.js` — Node regression tests
  - `css/vdicalc.css` — Styles
- `calculations/`, `host/`, `storage/`, `vm/`, `azure/`, `validation/` — Original Go packages (legacy)
- `main.go`, `config/`, `templates/` — Original HTTP server (legacy)

## Author

André Leibovici

## License

This project is licensed under the [MIT License](LICENSE).
