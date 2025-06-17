import { describe, it, expect } from "vitest";
import { extractLinks } from "./extractLinks";

describe("extractLinks", () => {
	it("should extract unique links from markdown", () => {
		const markdown =
			"This is a link to [[Page1]] and another link to [[Page2]]. Also, a repeated link to [[Page1]].";
		const expectedLinks = ["Page1", "Page2"];
		expect(extractLinks(markdown)).toEqual(expectedLinks);
	});

	it("should return an empty array if no links are found", () => {
		const markdown = "This is a text without any links.";
		const expectedLinks: string[] = [];
		expect(extractLinks(markdown)).toEqual(expectedLinks);
	});

	it("should handle case insensitivity", () => {
		const markdown = "Link to [[Page1]] and [[page1]].";
		const expectedLinks = ["Page1"];
		expect(extractLinks(markdown)).toEqual(expectedLinks);
	});
});
