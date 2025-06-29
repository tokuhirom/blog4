# E2E Tests for Blog4 Admin

This directory contains end-to-end tests for the Blog4 admin interface using Playwright.

## Running Tests

```bash
# Run all tests
pnpm test:e2e

# Run specific test file
pnpm test:e2e auth.spec.ts

# Run tests with UI mode
pnpm test:e2e:ui

# Run tests in headed browser
pnpm test:e2e:headed

# Run specific test group
pnpm test:e2e --grep "Login Page"
```

## Test Files

- `auth.spec.ts` - Authentication tests organized into:
  - **Login Page**: Form display, error handling
  - **Authentication Flow**: Login/logout, session management
- `top-page.spec.ts` - Tests for the admin top page (requires authentication)

## Test Status

### Working Tests ✅
- Display login form elements
- Show error with empty/incorrect credentials
- Login with correct credentials
- Logout functionality

### Known Issues ⚠️
- Session persistence across page refreshes
- Redirect to login when not authenticated
- Some logout redirect tests

## Notes

- The tests expect the dev server to be running on http://localhost:6173
- Default credentials: `admin:password`
- Some tests are skipped due to current implementation issues