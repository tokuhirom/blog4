import { test, expect } from '@playwright/test';

async function login(page) {
    await page.goto('/admin/login');
    await page.fill('input[name="username"]', 'admin');
    await page.fill('input[name="password"]', 'password');
    await page.click('button[type="submit"]');
    await page.waitForURL('/admin/entries/search');
}

// admin の一覧・検索は全件を /admin/api/entries から取得し、クライアント側で絞り込む。

test('admin search: client-side filtering by keyword', async ({ page }) => {
    await login(page);
    await expect(page).toHaveTitle(/Admin - Entry List/);
    // 全件ロード完了 (グリッド表示) を待つ
    await expect(page.locator('.entry-grid')).toBeVisible();

    await page.locator('input[placeholder="Search entries..."]').fill('Docker');
    await expect(page.locator('.entry-card:has-text("Docker Setup Guide")')).toBeVisible();
});

test('admin search: includes private entries', async ({ page }) => {
    await login(page);
    await expect(page.locator('.entry-grid')).toBeVisible();

    // public 検索と異なり、admin は private エントリも対象になる
    await page.locator('input[placeholder="Search entries..."]').fill('Private Draft');
    await expect(page.locator('.entry-card:has-text("Private Draft Example")')).toBeVisible();
});

test('admin search: Japanese keyword', async ({ page }) => {
    await login(page);
    await expect(page.locator('.entry-grid')).toBeVisible();

    await page.locator('input[placeholder="Search entries..."]').fill('日本語');
    await expect(page.locator('.entry-card:has-text("日本語コンテンツのサンプル")')).toBeVisible();
});
