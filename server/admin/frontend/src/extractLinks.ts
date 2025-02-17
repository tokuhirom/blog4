export function extractLinks(markdown: string): string[] {
	const links: string[] = markdown.match(/\[\[(.+?)\]\]/g) || [];
	const seen = new Set();
	const result: string[] = [];
	for (const link of links) {
		if (!seen.has(link.toLowerCase())) {
			seen.add(link.toLowerCase());
			result.push(link.slice(2, -2));
		}
	}
	return result;
}
