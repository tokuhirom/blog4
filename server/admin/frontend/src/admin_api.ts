import { Configuration, DefaultApi } from "./generated-client";

export function createAdminApiClient(): DefaultApi {
	const conf = new Configuration({
		basePath: import.meta.env.VITE_API_BASE_URL,
	});
	return new DefaultApi(conf);
}
