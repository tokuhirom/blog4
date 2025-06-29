import { expect, test } from "@playwright/experimental-ct-react";
import Footer from "./Footer";

test.describe("Footer Component", () => {
	test("should not render when build info is not available", async ({
		mount,
		page,
	}) => {
		// Mock the getBuildInfo API to return an error
		await page.route("**/admin/api/build-info", (route) => {
			route.fulfill({
				status: 500,
				contentType: "application/json",
				body: JSON.stringify({ message: "Internal Server Error" }),
			});
		});

		const component = await mount(<Footer />);

		// Wait for the component to mount and API call to complete
		await page.waitForTimeout(100);

		// Footer should not be visible when build info fails to load
		const footerElements = await page.locator("footer").count();
		expect(footerElements).toBe(0);
	});

	test("should render build info when available", async ({ mount, page }) => {
		// Mock the getBuildInfo API
		const mockBuildInfo = {
			buildTime: "2024-01-01T12:00:00Z",
			gitCommit: "abcdef1234567890",
			gitShortCommit: "abcdef1",
			gitBranch: "main",
			gitTag: "v1.0.0",
			githubUrl: "https://github.com/example/repo/commit/abcdef1234567890",
		};

		await page.route("**/admin/api/build-info", (route) => {
			route.fulfill({
				status: 200,
				contentType: "application/json",
				body: JSON.stringify(mockBuildInfo),
			});
		});

		const component = await mount(<Footer />);

		// Wait for the footer to be visible
		const footer = page.locator("footer");
		await expect(footer).toBeVisible();

		// Check that build time is displayed
		const buildTimeText = new Date(mockBuildInfo.buildTime).toLocaleString();
		await expect(footer).toContainText(buildTimeText);

		// Check that branch and tag are displayed
		await expect(footer).toContainText(
			`${mockBuildInfo.gitBranch}@${mockBuildInfo.gitTag}`,
		);

		// Check that the commit link is present
		const commitLink = footer.locator(`a[href="${mockBuildInfo.githubUrl}"]`);
		await expect(commitLink).toBeVisible();
		await expect(commitLink).toHaveText(mockBuildInfo.gitShortCommit);
	});

	test("should render without git tag when not available", async ({
		mount,
		page,
	}) => {
		// Mock the getBuildInfo API without gitTag
		const mockBuildInfo = {
			buildTime: "2024-01-01T12:00:00Z",
			gitCommit: "abcdef1234567890",
			gitShortCommit: "abcdef1",
			gitBranch: "feature/test",
			githubUrl: "https://github.com/example/repo/commit/abcdef1234567890",
		};

		await page.route("**/admin/api/build-info", (route) => {
			route.fulfill({
				status: 200,
				contentType: "application/json",
				body: JSON.stringify(mockBuildInfo),
			});
		});

		const component = await mount(<Footer />);

		// Wait for the footer to be visible
		const footer = page.locator("footer");
		await expect(footer).toBeVisible();

		// Should show branch without @ when no tag
		await expect(footer).toContainText(mockBuildInfo.gitBranch);
		await expect(footer).not.toContainText("@");
	});

	test("should handle API failure gracefully", async ({ mount, page }) => {
		// Use page.on to listen for console messages
		let consoleErrorLogged = false;
		page.on("console", (msg) => {
			if (
				msg.type() === "error" &&
				msg.text().includes("Failed to fetch build info")
			) {
				consoleErrorLogged = true;
			}
		});

		await page.route("**/admin/api/build-info", (route) => {
			route.abort("failed");
		});

		const component = await mount(<Footer />);

		// Wait for the component to mount and API call to complete
		await page.waitForTimeout(100);

		// Footer should not be visible when API fails
		const footerElements = await page.locator("footer").count();
		expect(footerElements).toBe(0);

		// Check that error was logged
		expect(consoleErrorLogged).toBe(true);
	});
});
