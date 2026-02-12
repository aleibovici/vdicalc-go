# VDI Calculator

## Synopsis

The myvirtualcloud.net VDI Calculator is a free tool for sizing Virtual Desktop Infrastructure deployments. It calculates host, storage, virtualization, and Azure instance requirements based on your VM specifications.

This version runs entirely in the browser as a static site hosted on GitHub Pages — no server or database required.

## Features

- VM sizing with multiple worker profiles (Task, Office, Knowledge, Power)
- Host CPU, memory, and capacity calculations
- Storage capacity, datastore, and IOps calculations
- Virtualization cluster and management server sizing
- Azure instance type recommendations
- Input validation with warning messages
- Print-friendly output

## Usage

Visit the GitHub Pages deployment or open `docs/index.html` in your browser. Configure your VDI parameters and click **Calculate** to see the results.

## Project Structure

- `docs/` — Static GitHub Pages site (client-side JavaScript)
  - `index.html` — Main calculator interface
  - `js/vdicalc.js` — Calculation engine (ported from Go)
  - `css/vdicalc.css` — Styling
- `main.go` — Original Go HTTP server (legacy)
- `config/` — Server-side configuration (legacy)
- `templates/` — Server-side Go HTML templates (legacy)

## Author

André Leibovici

## License

Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements. See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License. You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
