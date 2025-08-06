# Grafana UpDown.io Datasource Plugin

A Grafana datasource plugin to integrate with the UpDown.io API and visualize your service monitoring metrics.

[![Version](https://img.shields.io/github/v/release/nextmap-io/grafana-updownio-datasource)](https://github.com/nextmap-io/grafana-updownio-datasource/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![CI](https://github.com/nextmap-io/grafana-updownio-datasource/workflows/CI/badge.svg)](https://github.com/nextmap-io/grafana-updownio-datasource/actions)

## Features

- **Service List**: Visualize all your services monitored by UpDown.io
- **Performance Metrics**: Display uptime, Apdex, and response times
- **Downtime History**: Analyze downtimes and their causes
- **Time Series**: Create temporal charts to track your metrics evolution
- **Secure Configuration**: Secure storage of your UpDown.io API key

## Installation

### Prerequisites

- Grafana >= 10.4.0
- UpDown.io account with API key

### Plugin Installation

1. Download the latest version from [GitHub releases](https://github.com/nextmap-io/grafana-updownio-datasource/releases)
2. Extract the archive to your Grafana plugins folder
3. Restart Grafana
4. Configure the UpDown.io datasource in Grafana

#### Install from source

```bash
git clone https://github.com/nextmap-io/grafana-updownio-datasource.git
cd grafana-updownio-datasource
npm install
npm run build
# Copy the dist folder to your Grafana plugins directory
```

## Configuration

1. In Grafana, add a new datasource
2. Select "UpDown.io" from the list
3. Configure the settings:
   - **API URL**: https://updown.io/api (default)
   - **API Key**: Your UpDown.io API key (available in your UpDown.io account)
4. Test the connection
5. Save the configuration

## Usage

### Available Query Types

#### Service List
Retrieves the list of all your monitored services with basic information (URL, alias, uptime, status).

#### Service Metrics
Displays detailed metrics for a specific service:
- Uptime (availability percentage)
- Apdex (performance satisfaction index)
- Request statistics and response times

#### Downtimes
Visualizes service downtime history with:
- Downtime duration
- Interruption causes
- Start and end dates

#### Uptime
Displays service uptime percentage as a time series.

### Creating a Dashboard

1. Create a new dashboard
2. Add a panel
3. Select "UpDown.io" as datasource
4. Choose the query type
5. For specific metrics, select the relevant service
6. Configure the visualization according to your needs

## Development

### Prerequisites

- Node.js >= 22
- Go >= 1.22
- Docker (for testing)
- Git

### Setup

```bash
# Clone the repository
git clone https://github.com/nextmap-io/grafana-updownio-datasource.git
cd grafana-updownio-datasource

# Install dependencies
npm install

# Development mode with watch
npm run dev

# In another terminal, start Grafana
npm run server
```

### Available Scripts

- `npm run build`: Production build
- `npm run build:all`: Build frontend and backend
- `npm run dev`: Development mode with watch
- `npm run test`: Unit tests
- `npm run test:ci`: Tests in CI mode
- `npm run test:all`: Complete test suite (types, lint, tests, build, validation)
- `npm run e2e`: End-to-end tests
- `npm run lint`: Code linting
- `npm run typecheck`: TypeScript type checking
- `npm run validate`: Validate plugin with Grafana plugin-validator
- `npm run server`: Start Grafana for development

### Architecture

#### Frontend (TypeScript/React)
- `src/datasource.ts`: Main datasource class
- `src/types.ts`: TypeScript types for UpDown.io API
- `src/components/`: React components for the interface
- `src/module.ts`: Plugin entry point

#### Backend (Go)
- `pkg/plugin/datasource.go`: Backend implementation
- `pkg/models/`: Data structures
- `pkg/main.go`: Backend entry point

### Plugin Validation

The plugin is validated using Grafana's official plugin validator:

```bash
npm run validate
```

This ensures the plugin meets Grafana's standards and compatibility requirements.

## UpDown.io API

This plugin uses the UpDown.io REST API. For more information:
- [UpDown.io API Documentation](https://updown.io/api)
- [UpDown.io Website](https://updown.io)

### Used Endpoints

- `GET /api/checks` - Service list
- `GET /api/checks/:token` - Service details
- `GET /api/checks/:token/metrics` - Service metrics
- `GET /api/checks/:token/downtimes` - Downtime history

## Contributing

Contributions are welcome! Please:

1. Read the [Contributing Guide](CONTRIBUTING.md)
2. Follow the [Code of Conduct](CODE_OF_CONDUCT.md)
3. Fork the project
4. Create a feature branch
5. Commit your changes
6. Push to the branch
7. Open a Pull Request

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## Support and Questions

- 📖 [UpDown.io API Documentation](https://updown.io/api)
- 🐛 [Report a bug](https://github.com/nextmap-io/grafana-updownio-datasource/issues/new?template=bug_report.md)
- ✨ [Request a feature](https://github.com/nextmap-io/grafana-updownio-datasource/issues/new?template=feature_request.md)
- 💬 [Discussions](https://github.com/nextmap-io/grafana-updownio-datasource/discussions)

## Acknowledgments

- [Grafana](https://grafana.com/) for the excellent plugin framework
- [UpDown.io](https://updown.io/) for their well-documented API
- The open source community for the tools and libraries used
