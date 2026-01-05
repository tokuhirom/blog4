import { test, expect } from '@playwright/test';

test('admin search functionality', async ({ page }) => {
    // Login to admin
    await page.goto('/admin/login');
    await page.fill('input[name="username"]', 'admin');
    await page.fill('input[name="password"]', 'password');
    await page.click('button[type="submit"]');

    // Wait for redirect to entries page
    await page.waitForURL('/admin/entries/search');

    // Check that we're on the entries page
    await expect(page).toHaveTitle(/Admin - Entry List/);

    // Test search functionality with an existing entry keyword
    const searchInput = page.locator('input[name="q"]');
    await searchInput.fill('Docker');

    // Wait for HTMX to trigger the search (500ms delay + processing time)
    await page.waitForTimeout(1000);

    // Check that search results are displayed
    const entryList = page.locator('#entry-list');
    await expect(entryList).toBeVisible();

    // Verify that at least one entry card is shown (Docker Setup Guide should match)
    const entryCards = page.locator('.entry-card');
    await expect(entryCards.first()).toBeVisible();

    // Verify that the Docker Setup Guide entry is in the results
    await expect(page.locator('.entry-card:has-text("Docker Setup Guide")')).toBeVisible();
});
