import { test, expect } from '@playwright/test';
import { login, checkNoConsoleErrors } from './fixtures/test-helpers';

test.describe('Navigation and Link Pallet', () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('should navigate from "New Links to Create"', async ({ page }) => {
    const checkErrors = checkNoConsoleErrors(page);
    
    // Create an entry with wiki links
    await page.goto('/admin');
    await page.keyboard.press('c');
    await page.waitForURL(/\/admin\/entry\/.+/);
    
    // Add content with wiki links
    const bodyEditor = page.locator('.cm-content');
    await bodyEditor.click();
    await page.keyboard.type('This links to [[NonExistentPage]]');
    
    // Wait for link pallet to update
    await page.waitForTimeout(1500);
    
    // Check that "New Links to Create" section appears
    await expect(page.locator('text=New Links to Create')).toBeVisible();
    
    // Click on the new link card
    const newLinkCard = page.locator('text=NonExistentPage').first();
    await newLinkCard.click();
    
    // Should navigate to new entry
    await expect(page).toHaveURL(/\/admin\/entry\/.+/);
    
    // New entry should have the correct title
    const titleInput = page.locator('input[label="Title"]');
    await expect(titleInput).toHaveValue('NonExistentPage');
    
    checkErrors();
  });

  test('should show two-hop links with preview', async ({ page }) => {
    const checkErrors = checkNoConsoleErrors(page);
    
    // This test assumes there are entries with two-hop relationships
    // Navigate to an entry that has two-hop links
    await page.goto('/admin');
    
    // Find and click on first entry (if any)
    const firstEntry = page.locator('[class*="AdminEntryCardItem"]').first();
    const entryExists = await firstEntry.count() > 0;
    
    if (entryExists) {
      await firstEntry.click();
      await page.waitForURL(/\/admin\/entry\/.+/);
      
      // Check if two-hop links section exists
      const twoHopSection = page.locator('text=Two-hop Links');
      if (await twoHopSection.count() > 0) {
        await expect(twoHopSection).toBeVisible();
        
        // Two-hop link cards should show title and body preview
        const twoHopCards = page.locator('[class*="EntryCardItem"]');
        const cardCount = await twoHopCards.count();
        
        if (cardCount > 0) {
          // Check first card has content
          const firstCard = twoHopCards.first();
          await expect(firstCard).toBeVisible();
          
          // Card should contain some text (title or body preview)
          const cardText = await firstCard.textContent();
          expect(cardText).toBeTruthy();
          expect(cardText?.length).toBeGreaterThan(0);
        }
      }
    }
    
    checkErrors();
  });

  test('should navigate between entries using link pallet', async ({ page }) => {
    await page.goto('/admin');
    
    // Find and click on first entry
    const firstEntry = page.locator('[class*="AdminEntryCardItem"]').first();
    if (await firstEntry.count() > 0) {
      await firstEntry.click();
      await page.waitForURL(/\/admin\/entry\/.+/);
      
      // Check if Direct Links section exists
      const directLinksSection = page.locator('text=Direct Links');
      if (await directLinksSection.count() > 0) {
        // Click on a linked entry
        const linkedEntry = page.locator('[class*="AdminEntryCardItem"]').nth(1);
        if (await linkedEntry.count() > 0) {
          const linkedEntryUrl = await linkedEntry.getAttribute('href');
          await linkedEntry.click();
          
          // Should navigate to the linked entry
          await expect(page).toHaveURL(/\/admin\/entry\/.+/);
          
          // Should be a different URL
          expect(page.url()).not.toBe(linkedEntryUrl);
        }
      }
    }
  });
});