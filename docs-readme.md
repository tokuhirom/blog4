# Blog4 API Documentation

The API documentation for Blog4 is automatically generated and deployed to GitHub Pages.

## Viewing the Documentation

The API documentation is available at: https://tokuhirom.github.io/blog4/

## Local Development

To generate and preview the API documentation locally:

```bash
# Install dependencies
pnpm install --frozen-lockfile

# Generate OpenAPI spec from TypeSpec
pnpm run tsp

# Build static HTML documentation
pnpm run docs:build

# Preview documentation locally
pnpm run docs:preview
```

## Automatic Deployment

The documentation is automatically built and deployed to GitHub Pages whenever changes are pushed to the `main` branch. This is handled by the GitHub Actions workflow in `.github/workflows/deploy-api-docs.yml`.

## Technology Stack

- **TypeSpec**: API specification language
- **Redoc**: API documentation generator
- **GitHub Pages**: Documentation hosting
