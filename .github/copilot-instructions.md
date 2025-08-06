# Copilot Instructions for Grafana UpDown.io Plugin

<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

## Project Description

This project is a Grafana plugin (datasource) that allows querying the UpDown.io API to retrieve monitoring data.

## Main Features

1. **Data source configuration**: Allow adding an UpDown.io API key
2. **Service listing**: Retrieve the list of all monitored services
3. **Check status**: Get verification status over a given duration
4. **Time series metrics**: Display data as time series in Grafana

## Technologies Used

- **Frontend**: TypeScript, React, Grafana SDK
- **Backend**: Go, Grafana Backend SDK
- **API Client**: npm package `updown.io`
- **API**: UpDown.io REST API (https://updown.io/api)

## Code Structure

- `src/` : TypeScript/React code for the user interface
- `pkg/` : Go code for the backend
- The plugin uses a backend architecture to support secure authentication

## Development Summary

### Completed Implementation
1. **Full Plugin Architecture**: Complete Grafana datasource plugin with TypeScript frontend and Go backend
2. **API Integration**: Integration with UpDown.io API v3.3.5 using the official npm client
3. **Authentication**: Secure API key management through Grafana's secure JSON data
4. **Query Interface**: React-based query editor with service selection and query type options
5. **Backend Handler**: Go backend with proper datasource implementation for QueryData and CheckHealth
6. **Build System**: Complete build pipeline with webpack (frontend) and mage/go build (backend)
7. **Testing Framework**: Jest for unit tests, Playwright for E2E testing
8. **Documentation**: Comprehensive English documentation (README, CONTRIBUTING, CHANGELOG)
9. **GitHub Integration**: CI/CD workflows, issue templates, and proper repository structure
10. **Plugin Validation**: Grafana plugin validator compliance with proper packaging

### Key Technical Decisions
- **Backend Architecture**: Chosen for secure API key handling and CORS bypass
- **TypeScript Strict Mode**: Enforced type safety throughout the codebase  
- **React Hooks**: Modern React patterns with useCallback for performance optimization
- **Go Modules**: Proper dependency management with github.com/nextmap-io/grafana-updownio-datasource
- **Apache 2.0 License**: Open source licensing for community adoption

### Authorship & Branding
- **Author**: alo-is (GitHub username)
- **Organization**: nextmap-io
- **Repository**: github.com/nextmap-io/grafana-updownio-datasource
- **Plugin ID**: updown-updownio-datasource

### Build & Packaging
- **Frontend Build**: Webpack with TypeScript compilation and asset optimization
- **Backend Build**: Multi-platform Go binaries (linux/darwin/windows, amd64/arm64/arm)
- **Plugin Packaging**: Proper archive structure with single directory named after plugin ID
- **Validation**: Automated plugin validator integration in build pipeline

### Resolved Issues
1. **Go Compilation Errors**: Fixed function signatures, import paths, and type compatibility
2. **Plugin Validator**: Resolved archive structure issues and validator version compatibility
3. **English Localization**: Complete translation from French to English across all files
4. **React Dependencies**: Optimized hook dependencies and resolved deprecation warnings
5. **Module Path**: Updated from placeholder to actual GitHub repository path

## Development Instructions

1. Use Grafana TypeScript types for compatibility
2. Implement secure API key management on the backend side
3. Follow Grafana plugin naming conventions
4. Ensure compatibility with latest Grafana versions
5. Use the `updown.io` API client for all interactions with the UpDown.io API
6. Maintain English language throughout the codebase for public GitHub repository
7. Follow semantic versioning for releases
8. Use the build scripts: `npm run build:all` for complete build, `npm run validate` for plugin validation
