const baseUrl = import.meta.env.VITE_API_BASE_URL || "";

export const customInstance = async <T>({
	url,
	method,
	params,
	data,
	headers,
	signal,
	body,
}: {
	url: string;
	method: "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
	params?: Record<string, any>;
	data?: any;
	body?: any;
	headers?: Record<string, string>;
	signal?: AbortSignal;
}): Promise<T> => {
	// Build query string if params exist
	const queryString = params
		? `?${new URLSearchParams(params).toString()}`
		: "";

	const fullUrl = `${baseUrl}${url}${queryString}`;

	// Determine headers and body
	let requestHeaders = { ...headers };
	let requestBody: any;

	if (body instanceof FormData) {
		// Don't set Content-Type for FormData, let browser set it with boundary
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

	return response.json();
};
