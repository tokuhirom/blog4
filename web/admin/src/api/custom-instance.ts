const baseUrl =
	import.meta.env.VITE_API_BASE_URL || "http://localhost:8181/admin/api";

console.log("Using API base URL:", baseUrl);

export const customInstance = async <T>(
	url: string,
	options: {
		method: "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
		params?: Record<string, unknown>;
		data?: unknown;
		body?: unknown;
		headers?: Record<string, string>;
		signal?: AbortSignal;
	},
): Promise<T> => {
	const { method, params, data, body, headers, signal } = options;

	// Build query string if params exist
	const queryString = params
		? `?${new URLSearchParams(params).toString()}`
		: "";

	const fullUrl = `${baseUrl}${url}${queryString}`;

	// Determine headers and body
	const requestHeaders = { ...headers };
	let requestBody: BodyInit | undefined;

	if (body instanceof FormData) {
		// Don't set Content-Type for FormData, let browser set it with boundary
		requestBody = body;
	} else if (body) {
		// Use body if provided (already stringified by Orval)
		requestBody = body;
	} else if (data) {
		requestHeaders["Content-Type"] = "application/json";
		requestBody = JSON.stringify(data);
	}

	const response = await fetch(fullUrl, {
		method,
		headers: requestHeaders,
		body: requestBody,
		signal,
		credentials: "include", // Include cookies for session management
	});

	if (!response.ok) {
		const errorData = await response.json().catch(() => ({}));
		throw new Error(
			errorData.message || `HTTP error! status: ${response.status}`,
		);
	}

	// Handle empty responses
	if (
		response.status === 204 ||
		response.headers.get("content-length") === "0"
	) {
		return {} as T;
	}

	const jsonData = await response.json();

	// Orval expects response in format { data: T, status: number }
	return {
		data: jsonData,
		status: response.status,
		headers: response.headers,
	} as T;
};
