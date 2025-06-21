import { test, expect } from '@playwright/test';
import { login } from './fixtures/test-helpers';

test.describe('MUI Warnings Check', () => {
  test('should not have MUI Grid warnings', async ({ page }) => {
    const consoleWarnings: string[] = [];
    
    // Capture console warnings
    page.on('console', (msg) => {
      if (msg.type() === 'warning') {
        consoleWarnings.push(msg.text());
      }
    });
    
    await login(page);
    
    // Navigate to pages that use Grid
    await page.goto('/admin');
    await page.waitForLoadState('networkidle');
    
    // Navigate to entry page (uses Grid)
    const firstEntry = page.locator('[class*="AdminEntryCardItem"]').first();
    if (await firstEntry.count() > 0) {
      await firstEntry.click();
      await page.waitForURL(/\/admin\/entry\/.+/);
      await page.waitForLoadState('networkidle');
    }
    
    // Check for MUI Grid-related warnings
    const muiGridWarnings = consoleWarnings.filter(warning => 
      warning.includes('Grid') || 
      warning.includes('MUI') ||
      warning.includes('size prop') ||
      warning.includes('item prop')
    );
    
    // Should have no MUI Grid warnings
    expect(muiGridWarnings).toHaveLength(0);
  });

  test('should not have any console errors on main pages', async ({ page }) => {
    const consoleErrors: string[] = [];
    
    page.on('console', (msg) => {
      if (msg.type() === 'error') {
        consoleErrors.push(msg.text());
      }
    });
    
    await login(page);
    
    // Test main admin page
    await page.goto('/admin');
    await page.waitForLoadState('networkidle');
    
    // Test entry creation
    await page.keyboard.press('c');
    await page.waitForURL(/\/admin\/entry\/.+/);
    await page.waitForLoadState('networkidle');
    
    // Filter out expected/acceptable errors
    const unexpectedErrors = consoleErrors.filter(error => {
      return !error.includes('ResizeObserver') && // Common browser warning
             !error.includes('401') && // Auth errors before login
             !error.includes('favicon'); // Missing favicon
    });
    
    expect(unexpectedErrors).toHaveLength(0);
  });

  test('should use correct Grid2 component', async ({ page }) => {
    await login(page);
    await page.goto('/admin');
    
    // Check that Grid2 is being used (has size prop instead of item prop)
    const gridElements = page.locator('[class*="MuiGrid2"]');
    const gridCount = await gridElements.count();
    
    // Should have Grid2 elements
    expect(gridCount).toBeGreaterThan(0);
    
    // Old Grid elements should not exist
    const oldGridElements = page.locator('[class*="MuiGrid-item"]');
    const oldGridCount = await oldGridElements.count();
    expect(oldGridCount).toBe(0);
  });
});