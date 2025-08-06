# Contributing Guide

Thank you for your interest in contributing to the Grafana UpDown.io plugin! This guide will help you understand how to participate in the project development.

## Code of Conduct

By participating in this project, you agree to abide by the [Code of Conduct](CODE_OF_CONDUCT.md).

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with:
- A clear description of the problem
- Steps to reproduce the bug
- The Grafana version used
- The plugin version
- Error logs if available

### Proposing Features

To propose a new feature:
1. Check that it doesn't already exist in the issues
2. Create an issue describing:
   - The problem it would solve
   - The proposed solution
   - Alternatives considered

   - Alternatives considered

### Development

#### Prerequisites

- Node.js >= 18
- Go >= 1.21
- Docker (for testing)
- Git

#### Environment Setup

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR-USERNAME/grafana-updownio-datasource.git
   cd grafana-updownio-datasource
   ```

3. Install dependencies:
   ```bash
   npm install
   ```

4. Configure Git:
   ```bash
   git remote add upstream https://github.com/nextmap-io/grafana-updownio-datasource.git
   ```

#### Development Workflow

1. Create a branch for your feature:
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. Develop in watch mode:
   ```bash
   npm run dev
   ```

3. Start Grafana in development:
   ```bash
   npm run server
   ```

4. Test your changes:
   ```bash
   npm test
   npm run e2e
   ```

5. Check linting:
   ```bash
   npm run lint
   npm run typecheck
   ```

#### Code Structure

- `src/`: Frontend TypeScript/React code
- `pkg/`: Backend Go code
- `tests/`: E2E tests
- `.config/`: Webpack and tools configuration

#### Code Conventions

- **TypeScript**: Follow configured ESLint rules
- **Go**: Use `gofmt` and follow Go conventions
- **Commits**: Descriptive commit messages in English
- **Tests**: Add tests for new features

### Pull Requests

1. Ensure all tests pass
2. Update documentation if necessary
3. Add an entry in CHANGELOG.md
4. Create the Pull Request to the `main` branch

#### Pull Request Template

```markdown
## Description

Brief description of changes

## Type of Change

- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation

## Testing

- [ ] Unit tests pass
- [ ] E2E tests pass
- [ ] Manually tested

## Checklist

- [ ] Code self-reviewed
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] No conflicts
```

## Plugin Architecture

### Frontend (TypeScript/React)

- `src/datasource.ts`: Main data source class
- `src/types.ts`: TypeScript types for UpDown.io API
- `src/components/`: React components for the interface
- `src/module.ts`: Plugin entry point

### Backend (Go)

- `pkg/plugin/datasource.go`: Backend implementation
- `pkg/models/`: Data structures
- `pkg/main.go`: Backend entry point

### UpDown.io API

The plugin uses these endpoints:
- `GET /api/checks`: List of services
- `GET /api/checks/{token}`: Service details
- `GET /api/checks/{token}/metrics`: Metrics
- `GET /api/checks/{token}/downtimes`: Downtime

## Questions and Support

- Create an issue for technical questions
- Check the [UpDown.io documentation](https://updown.io/api)
- Check existing issues before creating a new one

## License

By contributing, you agree that your contributions will be licensed under Apache 2.0.

## Plugin Architecture

### Frontend (TypeScript/React)

- `src/datasource.ts`: Main data source class
- `src/types.ts`: TypeScript types for UpDown.io API
- `src/components/`: React components for the interface
- `src/module.ts`: Plugin entry point

### Backend (Go)

- `pkg/plugin/datasource.go`: Backend implementation
- `pkg/models/`: Data structures
- `pkg/main.go`: Backend entry point

### UpDown.io API

The plugin uses these endpoints:
- `GET /api/checks`: List of services
- `GET /api/checks/{token}`: Service details
- `GET /api/checks/{token}/metrics`: Metrics
- `GET /api/checks/{token}/downtimes`: Downtime

## Questions and Support

- Create an issue for technical questions
- Check the [UpDown.io documentation](https://updown.io/api)
- Check existing issues before creating a new one

## License

By contributing, you agree that your contributions will be licensed under Apache 2.0.
