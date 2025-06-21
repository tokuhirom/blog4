import { test, expect } from '@playwright/test';
import { login, checkNoConsoleErrors } from './fixtures/test-helpers';

test.describe('Entry CRUD Operations', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('should create a new entry with keyboard shortcut', async ({ page }) => {
    const checkErrors = checkNoConsoleErrors(page);
    
    await page.goto('/admin');
    
    // Press 'c' to create new entry
    await page.keyboard.press('c');
    
    // Should navigate to new entry page
    await expect(page).toHaveURL(/\/admin\/entry\/.+/);
    
    // Should have title input
    const titleInput = page.locator('input[label="Title"]');
    await expect(titleInput).toBeVisible();
    
    // Should have body editor
    const bodyEditor = page.locator('.cm-editor');
    await expect(bodyEditor).toBeVisible();
    
    // Check no console errors
    checkErrors();
  });

  test('should update entry title', async ({ page }) => {
    await page.goto('/admin');
    await page.keyboard.press('c');
    await page.waitForURL(/\/admin\/entry\/.+/);
    
    const newTitle = `Test Entry ${Date.now()}`;
    const titleInput = page.locator('input[label="Title"]');
    
    await titleInput.fill(newTitle);
    
    // Wait for auto-save
    await page.waitForTimeout(1000);
    
    // Verify title was saved
    await expect(titleInput).toHaveValue(newTitle);
    
    // Reload page to verify persistence
    await page.reload();
    await expect(titleInput).toHaveValue(newTitle);
  });

  test('should update entry body', async ({ page }) => {
    await page.goto('/admin');
    await page.keyboard.press('c');
    await page.waitForURL(/\/admin\/entry\/.+/);
    
    const bodyContent = '# Test Content\n\nThis is a test entry.';
    const bodyEditor = page.locator('.cm-content');
    
    await bodyEditor.click();
    await page.keyboard.type(bodyContent);
    
    // Wait for auto-save
    await page.waitForTimeout(1500);
    
    // Reload page to verify persistence
    await page.reload();
    await expect(bodyEditor).toContainText('Test Content');
    await expect(bodyEditor).toContainText('This is a test entry.');
  });

  test('should change entry visibility', async ({ page }) => {
    await page.goto('/admin');
    await page.keyboard.press('c');
    await page.waitForURL(/\/admin\/entry\/.+/);
    
    // Should be private by default
    const privateRadio = page.locator('input[value="private"]');
    await expect(privateRadio).toBeChecked();
    
    // Change to public
    const publicRadio = page.locator('input[value="public"]');
    await publicRadio.click();
    
    // Confirm dialog
    page.on('dialog', dialog => dialog.accept());
    
    // Wait for update
    await page.waitForTimeout(500);
    
    // Should show "Go to User Side Page" link
    await expect(page.locator('text=Go to User Side Page')).toBeVisible();
  });

  test('should delete entry', async ({ page }) => {
    // Create a test entry first
    await page.goto('/admin');
    await page.keyboard.press('c');
    await page.waitForURL(/\/admin\/entry\/.+/);
    
    const entryUrl = page.url();
    
    // Click delete button
    const deleteButton = page.locator('button:has-text("Delete")');
    await deleteButton.click();
    
    // Confirm deletion
    page.on('dialog', dialog => dialog.accept());
    
    // Should redirect to admin page
    await expect(page).toHaveURL('/admin/');
    
    // Entry should no longer be accessible
    await page.goto(entryUrl);
    await expect(page).toHaveURL('/admin/');
  });
});