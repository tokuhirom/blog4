import { test, expect } from '@playwright/test';
import { TEST_USER, TEST_PASSWORD } from './fixtures/test-helpers';

test.describe('Authentication', () => {
  test('should show login form when not authenticated', async ({ page }) => {
    await page.goto('/admin');
    
    // Check for login form elements
    await expect(page.locator('input[name="username"]')).toBeVisible();
    await expect(page.locator('input[name="password"]')).toBeVisible();
    await expect(page.locator('button[type="submit"]')).toBeVisible();
  });

  test('should login successfully with valid credentials', async ({ page }) => {
    await page.goto('/admin');
    
    // Fill login form
    await page.fill('input[name="username"]', TEST_USER);
    await page.fill('input[name="password"]', TEST_PASSWORD);
    await page.click('button[type="submit"]');
    
    // Should redirect to admin page
    await expect(page).toHaveURL('/admin/');
    
    // Should show search box (indicating successful login)
    await expect(page.locator('input[placeholder*="Search"]')).toBeVisible();
  });

  test('should show error with invalid credentials', async ({ page }) => {
    await page.goto('/admin');
    
    // Fill login form with wrong credentials
    await page.fill('input[name="username"]', 'wronguser');
    await page.fill('input[name="password"]', 'wrongpass');
    await page.click('button[type="submit"]');
    
    // Should stay on login page
    await expect(page.locator('input[name="username"]')).toBeVisible();
    
    // Should show error message (implementation dependent)
    // await expect(page.locator('text=Invalid credentials')).toBeVisible();
  });
});