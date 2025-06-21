import { expect, type Page } from '@playwright/test';

export const TEST_USER = 'admin';
export const TEST_PASSWORD = 'admin';

export async function login(page: Page) {
  await page.goto('/admin');
  
  // Check if already logged in
  const loginForm = await page.locator('input[name="username"]').count();
  if (loginForm === 0) {
    // Already logged in
    return;
  }
  
  // Fill login form
  await page.fill('input[name="username"]', TEST_USER);
  await page.fill('input[name="password"]', TEST_PASSWORD);
  await page.click('button[type="submit"]');
  
  // Wait for navigation
  await page.waitForURL('/admin/**');
}

export async function createTestEntry(page: Page, title: string) {
  // Navigate to admin page
  await page.goto('/admin');
  
  // Press 'c' to create new entry
  await page.keyboard.press('c');
  
  // Wait for entry page
  await page.waitForURL(/\/admin\/entry\/.+/);
  
  // Update title
  const titleInput = page.locator('input').filter({ has: page.locator('label:has-text("Title")') });
  await titleInput.fill(title);
  
  // Wait for auto-save
  await page.waitForTimeout(1000);
  
  return page.url();
}

export async function checkNoConsoleErrors(page: Page) {
  // Listen for console errors
  const consoleErrors: string[] = [];
  
  page.on('console', (msg) => {
    if (msg.type() === 'error') {
      consoleErrors.push(msg.text());
    }
  });
  
  return () => {
    // Filter out expected errors (if any)
    const unexpectedErrors = consoleErrors.filter(error => {
      // Add any expected errors here
      return !error.includes('ResizeObserver');
    });
    
    expect(unexpectedErrors).toHaveLength(0);
  };
}