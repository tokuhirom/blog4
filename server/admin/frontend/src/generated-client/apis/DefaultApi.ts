/* tslint:disable */
/* eslint-disable */
/**
 * Entries API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 1.0.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


import * as runtime from '../runtime';
import type {
  CreateEntryRequest,
  CreateEntryResponse,
  ErrorResponse,
  GetLatestEntriesRow,
  UpdateEntryBodyRequest,
  UpdateEntryTitleRequest,
} from '../models/index';
import {
    CreateEntryRequestFromJSON,
    CreateEntryRequestToJSON,
    CreateEntryResponseFromJSON,
    CreateEntryResponseToJSON,
    ErrorResponseFromJSON,
    ErrorResponseToJSON,
    GetLatestEntriesRowFromJSON,
    GetLatestEntriesRowToJSON,
    UpdateEntryBodyRequestFromJSON,
    UpdateEntryBodyRequestToJSON,
    UpdateEntryTitleRequestFromJSON,
    UpdateEntryTitleRequestToJSON,
} from '../models/index';

export interface CreateEntryOperationRequest {
    createEntryRequest: CreateEntryRequest;
}

export interface GetEntryByDynamicPathRequest {
    path: string;
}

export interface GetLatestEntriesRequest {
    lastLastEditedAt?: Date;
}

export interface GetLinkedEntryPathsRequest {
    path: string;
}

export interface UpdateEntryBodyOperationRequest {
    path: string;
    updateEntryBodyRequest: UpdateEntryBodyRequest;
}

export interface UpdateEntryTitleOperationRequest {
    path: string;
    updateEntryTitleRequest: UpdateEntryTitleRequest;
}

/**
 * 
 */
export class DefaultApi extends runtime.BaseAPI {

    /**
     * Create a new entry
     */
    async createEntryRaw(requestParameters: CreateEntryOperationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<CreateEntryResponse>> {
        if (requestParameters['createEntryRequest'] == null) {
            throw new runtime.RequiredError(
                'createEntryRequest',
                'Required parameter "createEntryRequest" was null or undefined when calling createEntry().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/entries`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: CreateEntryRequestToJSON(requestParameters['createEntryRequest']),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => CreateEntryResponseFromJSON(jsonValue));
    }

    /**
     * Create a new entry
     */
    async createEntry(requestParameters: CreateEntryOperationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<CreateEntryResponse> {
        const response = await this.createEntryRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Get all entry titles
     */
    async getAllEntryTitlesRaw(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Array<string>>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/entries/titles`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse<any>(response);
    }

    /**
     * Get all entry titles
     */
    async getAllEntryTitles(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Array<string>> {
        const response = await this.getAllEntryTitlesRaw(initOverrides);
        return await response.value();
    }

    /**
     * Get entry by dynamic path
     */
    async getEntryByDynamicPathRaw(requestParameters: GetEntryByDynamicPathRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<GetLatestEntriesRow>> {
        if (requestParameters['path'] == null) {
            throw new runtime.RequiredError(
                'path',
                'Required parameter "path" was null or undefined when calling getEntryByDynamicPath().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/entries/{path}`.replace(`{${"path"}}`, encodeURIComponent(String(requestParameters['path']))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => GetLatestEntriesRowFromJSON(jsonValue));
    }

    /**
     * Get entry by dynamic path
     */
    async getEntryByDynamicPath(requestParameters: GetEntryByDynamicPathRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<GetLatestEntriesRow> {
        const response = await this.getEntryByDynamicPathRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Get latest entries
     */
    async getLatestEntriesRaw(requestParameters: GetLatestEntriesRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Array<GetLatestEntriesRow>>> {
        const queryParameters: any = {};

        if (requestParameters['lastLastEditedAt'] != null) {
            queryParameters['last_last_edited_at'] = (requestParameters['lastLastEditedAt'] as any).toISOString();
        }

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/entries`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => jsonValue.map(GetLatestEntriesRowFromJSON));
    }

    /**
     * Get latest entries
     */
    async getLatestEntries(requestParameters: GetLatestEntriesRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Array<GetLatestEntriesRow>> {
        const response = await this.getLatestEntriesRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Get linked entry paths
     */
    async getLinkedEntryPathsRaw(requestParameters: GetLinkedEntryPathsRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<{ [key: string]: string | null; }>> {
        if (requestParameters['path'] == null) {
            throw new runtime.RequiredError(
                'path',
                'Required parameter "path" was null or undefined when calling getLinkedEntryPaths().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/entries/{path}/links`.replace(`{${"path"}}`, encodeURIComponent(String(requestParameters['path']))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse<any>(response);
    }

    /**
     * Get linked entry paths
     */
    async getLinkedEntryPaths(requestParameters: GetLinkedEntryPathsRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<{ [key: string]: string | null; }> {
        const response = await this.getLinkedEntryPathsRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Update entry body
     */
    async updateEntryBodyRaw(requestParameters: UpdateEntryBodyOperationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<object>> {
        if (requestParameters['path'] == null) {
            throw new runtime.RequiredError(
                'path',
                'Required parameter "path" was null or undefined when calling updateEntryBody().'
            );
        }

        if (requestParameters['updateEntryBodyRequest'] == null) {
            throw new runtime.RequiredError(
                'updateEntryBodyRequest',
                'Required parameter "updateEntryBodyRequest" was null or undefined when calling updateEntryBody().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/entries/{path}/body`.replace(`{${"path"}}`, encodeURIComponent(String(requestParameters['path']))),
            method: 'PUT',
            headers: headerParameters,
            query: queryParameters,
            body: UpdateEntryBodyRequestToJSON(requestParameters['updateEntryBodyRequest']),
        }, initOverrides);

        return new runtime.JSONApiResponse<any>(response);
    }

    /**
     * Update entry body
     */
    async updateEntryBody(requestParameters: UpdateEntryBodyOperationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<object> {
        const response = await this.updateEntryBodyRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Update entry title
     */
    async updateEntryTitleRaw(requestParameters: UpdateEntryTitleOperationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<object>> {
        if (requestParameters['path'] == null) {
            throw new runtime.RequiredError(
                'path',
                'Required parameter "path" was null or undefined when calling updateEntryTitle().'
            );
        }

        if (requestParameters['updateEntryTitleRequest'] == null) {
            throw new runtime.RequiredError(
                'updateEntryTitleRequest',
                'Required parameter "updateEntryTitleRequest" was null or undefined when calling updateEntryTitle().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/entries/{path}/title`.replace(`{${"path"}}`, encodeURIComponent(String(requestParameters['path']))),
            method: 'PUT',
            headers: headerParameters,
            query: queryParameters,
            body: UpdateEntryTitleRequestToJSON(requestParameters['updateEntryTitleRequest']),
        }, initOverrides);

        return new runtime.JSONApiResponse<any>(response);
    }

    /**
     * Update entry title
     */
    async updateEntryTitle(requestParameters: UpdateEntryTitleOperationRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<object> {
        const response = await this.updateEntryTitleRaw(requestParameters, initOverrides);
        return await response.value();
    }

}
