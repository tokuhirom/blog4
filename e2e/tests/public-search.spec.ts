import { test, expect } from '@playwright/test';

// クライアント検索 (search.js) が /search-index.json を読み、ブラウザ側で
// 絞り込むことを検証する。MySQL FULLTEXT / MATCH AGAINST から移行した代替実装の確認。

test('public search: initial query from URL renders matching entry', async ({ page }) => {
    await page.goto('/search?q=Docker');
    await expect(page.locator('.entry-title', { hasText: 'Docker Setup Guide' })).toBeVisible();
});

test('public search: incremental search works for Japanese keywords', async ({ page }) => {
    await page.goto('/search');
    // fetch 完了 (= リスナー登録 + 初期表示) を待つ
    await expect(page.locator('#search-results')).toContainText('Enter keywords');
    await page.locator('.search-input').fill('日本語');
    await expect(
        page.locator('.entry-title', { hasText: '日本語コンテンツのサンプル' }),
    ).toBeVisible();
});

test('public search: AND of multiple terms narrows results', async ({ page }) => {
    await page.goto('/search');
    await expect(page.locator('#search-results')).toContainText('Enter keywords');
    const input = page.locator('.search-input');
    await input.fill('Docker Setup');
    await expect(page.locator('.entry-title', { hasText: 'Docker Setup Guide' })).toBeVisible();
    // 無関係な語を AND で足すと結果が消える
    await input.fill('Docker zzzznotfound');
    await expect(page.locator('#search-results')).toContainText('No results found');
});

test('public search: private entries are not searchable', async ({ page }) => {
    await page.goto('/search');
    await expect(page.locator('#search-results')).toContainText('Enter keywords');
    await page.locator('.search-input').fill('Private Draft Example');
    await page.waitForTimeout(600);
    // private エントリは /search-index.json に含まれないので結果に出ない
    await expect(page.getByText('Private Draft Example')).toHaveCount(0);
});
