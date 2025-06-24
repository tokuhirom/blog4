import { test, expect } from "@playwright/test";

test.describe("Admin Authentication", () => {
	test.describe("Login Page", () => {
		test("should display login form elements", async ({ page }) => {
			await page.goto("/admin/login");

			// Login form should be visible
			await expect(
				page.getByRole("heading", { name: "Admin Login" }),
			).toBeVisible();
			await expect(page.getByLabel("Username")).toBeVisible();
			await expect(page.getByLabel("Password")).toBeVisible();
			await expect(page.getByRole("button", { name: "Login" })).toBeVisible();
		});

		test("should show error with empty credentials", async ({ page }) => {
			await page.goto("/admin/login");

			// Click login without filling fields
			await page.getByRole("button", { name: "Login" }).click();

			// Should stay on login page
			await expect(page).toHaveURL(/\/admin\/login/);
		});

		test("should show error with incorrect credentials", async ({ page }) => {
			await page.goto("/admin/login");

			// Try to login with incorrect credentials
			await page.getByLabel("Username").fill("wronguser");
			await page.getByLabel("Password").fill("wrongpassword");
			await page.getByRole("button", { name: "Login" }).click();

			// Should show error message
			await expect(
				page.getByText("Invalid username or password"),
			).toBeVisible();

			// Should still be on login page
			await expect(page).toHaveURL(/\/admin\/login/);
		});
	});

	test.describe("Authentication Flow", () => {
		test("should login successfully with correct credentials", async ({
			page,
		}) => {
			await page.goto("/admin/login");

			// Login with correct credentials
			await page.getByLabel("Username").fill("admin");
			await page.getByLabel("Password").fill("password");
			await page.getByRole("button", { name: "Login" }).click();

			// Wait for navigation
			await page.waitForLoadState("networkidle");

			// Should redirect to admin top page
			await expect(page).toHaveURL("/admin");

			// Should see the admin header buttons
			await expect(page.getByRole("button", { name: "New Entry" })).toBeVisible();
			await expect(page.getByRole("button", { name: "Logout" })).toBeVisible();
		});

		test("should logout successfully", async ({ page }) => {
			// First login
			await page.goto("/admin/login");
			await page.getByLabel("Username").fill("admin");
			await page.getByLabel("Password").fill("password");
			await page.getByRole("button", { name: "Login" }).click();

			// Wait for navigation
			await page.waitForLoadState("networkidle");

			await expect(page).toHaveURL("/admin");

			// Click logout button
			await page.getByRole("button", { name: "Logout" }).click();

			// Wait for navigation to complete
			await page.waitForLoadState("networkidle");

			// Should redirect to login page
			await expect(page).toHaveURL(/\/admin\/login/, { timeout: 10000 });
		});
	});
});