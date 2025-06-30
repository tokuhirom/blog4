import { test, expect } from "@playwright/test";

test.describe("Admin Top Page", () => {
	test.beforeEach(async ({ page }) => {
		// Login before each test
		await page.goto("/admin/login");
		await page.getByLabel("Username").fill("admin");
		await page.getByLabel("Password").fill("password");
		await page.getByRole("button", { name: "Login" }).click();

		// Wait for redirect to complete
		await page.waitForLoadState("networkidle");
		await expect(page).toHaveURL("/admin");
	});

	test("should display top page elements correctly", async ({ page }) => {
		// Check header elements
		await expect(page.getByRole("link", { name: "Blog Admin" })).toBeVisible();
		await expect(page.getByRole("button", { name: "New Entry" })).toBeVisible();
		await expect(page.getByRole("button", { name: "Logout" })).toBeVisible();

		// Check search box
		await expect(page.getByPlaceholder("Search entries...")).toBeVisible();

		// Entry cards should be visible (if any exist)
		const entryCards = page.locator('[class*="MuiCard-root"]');
		const count = await entryCards.count();

		if (count > 0) {
			// At least one entry card should be visible
			await expect(entryCards.first()).toBeVisible();
		}
	});

	test("should filter entries using search box", async ({ page }) => {
		// Type in search box
		const searchBox = page.getByPlaceholder("Search entries...");
		await searchBox.fill("Private");

		// Wait a bit for debounce
		await page.waitForTimeout(1300);

		// Check if entries are filtered (if any match)
		const entryCards = page.locator('[class*="MuiCard-root"]');
		const allCards = await entryCards.all();

		// All visible cards should contain 'test' in title or body
		for (const card of allCards) {
			const text = await card.textContent();
			expect(text?.toLowerCase()).toContain("private");
		}
	});

	test("should create new entry using keyboard shortcut", async ({ page }) => {
		// Press 'c' key to create new entry
		await page.keyboard.press("c");

		// Should navigate to new entry page
		await expect(page).toHaveURL(/\/admin\/entry\/.+/);

		// Entry editor should be visible
		await expect(page.getByLabel("Title")).toBeVisible();
		await expect(page.getByText("Body")).toBeVisible();

		// The title should contain "New Entry" with timestamp
		const titleInput = page.getByLabel("Title");
		const titleValue = await titleInput.inputValue();
		expect(titleValue).toContain(new Date().getFullYear().toString());
	});

	test.skip("should NOT create new entry when typing 'c' in search box", async ({
		page,
	}) => {
		// Focus on search box
		const searchBox = page.getByPlaceholder("Search entries...");
		await searchBox.focus();

		// Type 'c' in the search box
		await searchBox.type("c");

		// Wait a moment to ensure no navigation happens
		await page.waitForTimeout(500);

		// Should still be on the admin page
		await expect(page).toHaveURL("/admin");

		// Search box should contain 'c'
		await expect(searchBox).toHaveValue("c");

		// Should not have navigated to a new entry page
		await expect(page).not.toHaveURL(/\/admin\/entry\/.+/);
	});

	test("should create new entry using NEW button", async ({ page }) => {
		// Click New Entry button
		await page.getByRole("button", { name: "New Entry" }).click();

		// Should navigate to new entry page
		await expect(page).toHaveURL(/\/admin\/entry\/.+/);

		// Entry editor should be visible
		await expect(page.getByLabel("Title")).toBeVisible();
		await expect(page.getByText("Body")).toBeVisible();
	});

	test("should navigate to entry when clicking entry card", async ({
		page,
	}) => {
		const entryCards = page.locator('[class*="MuiCard-root"]');
		const count = await entryCards.count();

		if (count > 0) {
			// Get the first entry's title
			const firstCard = entryCards.first();
			const titleElement = firstCard.locator("h2");
			const title = await titleElement.textContent();

			// Click the card
			await firstCard.click();

			// Should navigate to entry edit page
			await expect(page).toHaveURL(/\/admin\/entry\/.+/);

			// The title input should contain the entry title
			const titleInput = page.getByLabel("Title");
			await expect(titleInput).toHaveValue(title || "");
		}
	});
});
