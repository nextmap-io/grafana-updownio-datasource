# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-08-06

### Added
- Initial Grafana plugin for UpDown.io datasource
- UpDown.io API key authentication support
- Four query types:
  - Service List: Display all UpDown.io checks
  - Service Metrics: Uptime, Apdex, response times
  - Downtimes: Downtime history
  - Uptime: Availability percentage as time series
- Configuration interface with API key validation
- Query interface with service selection and grouping options
- Go backend for secure API calls
- TypeScript/React frontend for user interface
- Complete documentation and installation instructions
- Native Grafana visualizations support
- Error handling and detailed logging

### Technical
- Official Grafana plugin SDK usage
- Backend/frontend architecture for security
- Custom HTTP client for UpDown.io API
- Complete TypeScript types for UpDown.io API
- Unit and E2E tests configured
- GitHub Actions workflow for CI/CD
- Apache 2.0 License

[1.0.0]: https://github.com/alo-is/grafana-updownio-datasource/releases/tag/v1.0.0
