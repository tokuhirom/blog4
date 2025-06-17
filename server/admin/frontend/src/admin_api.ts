import * as api from "./generated-client";

export interface DefaultApi {
	getLatestEntries(params?: api.GetLatestEntriesParams): Promise<api.GetLatestEntriesRow[]>;
	createEntry(request: api.CreateEntryRequest): Promise<api.CreateEntryResponse>;
	getAllEntryTitles(): Promise<string[]>;
	deleteEntry(params: { path: string }): Promise<api.EmptyResponse>;
	getEntryByDynamicPath(params: { path: string }): Promise<api.GetLatestEntriesRow>;
	updateEntryBody(params: { path: string }, request: api.UpdateEntryBodyRequest): Promise<api.EmptyResponse>;
	getLinkPallet(params: { path: string }): Promise<api.LinkPalletData>;
	getLinkedEntryPaths(params: { path: string }): Promise<Record<string, string>>;
	regenerateEntryImage(params: { path: string }): Promise<api.EmptyResponse>;
	updateEntryTitle(params: { path: string }, request: api.UpdateEntryTitleRequest): Promise<api.EmptyResponse>;
	updateEntryVisibility(params: { path: string }, request: api.UpdateVisibilityRequest): Promise<api.UpdateVisibilityResponse>;
	uploadFile(request: FormData): Promise<api.UploadFileResponse>;
}

export function createAdminApiClient(): DefaultApi {
	return {
		async getLatestEntries(params) {
			const response = await api.getLatestEntries(params);
			return response.data;
		},
		async createEntry(request) {
			const response = await api.createEntry(request);
			return response.data;
		},
		async getAllEntryTitles() {
			const response = await api.getAllEntryTitles();
			return response.data;
		},
		async deleteEntry(params) {
			const response = await api.deleteEntry(params.path);
			return response.data;
		},
		async getEntryByDynamicPath(params) {
			const response = await api.getEntryByDynamicPath(params.path);
			return response.data;
		},
		async updateEntryBody(params, request) {
			const response = await api.updateEntryBody(params.path, request);
			return response.data;
		},
		async getLinkPallet(params) {
			const response = await api.getLinkPallet(params.path);
			return response.data;
		},
		async getLinkedEntryPaths(params) {
			const response = await api.getLinkedEntryPaths(params.path);
			return response.data;
		},
		async regenerateEntryImage(params) {
			const response = await api.regenerateEntryImage(params.path);
			return response.data;
		},
		async updateEntryTitle(params, request) {
			const response = await api.updateEntryTitle(params.path, request);
			return response.data;
		},
		async updateEntryVisibility(params, request) {
			const response = await api.updateEntryVisibility(params.path, request);
			return response.data;
		},
		async uploadFile(formData) {
			// Extract file from FormData
			const file = formData.get('file') as File;
			if (!file) {
				throw new Error('No file provided');
			}
			const response = await api.uploadFile({ file });
			return response.data;
		},
	};
}
