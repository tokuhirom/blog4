/* tslint:disable */
/* eslint-disable */
/**
 * Admin API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 1.0.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { mapValues } from '../runtime';
/**
 * 
 * @export
 * @interface UpdateVisibilityResponse
 */
export interface UpdateVisibilityResponse {
    /**
     * The new visibility status for the entry
     * @type {string}
     * @memberof UpdateVisibilityResponse
     */
    visibility: string;
}

/**
 * Check if a given object implements the UpdateVisibilityResponse interface.
 */
export function instanceOfUpdateVisibilityResponse(value: object): value is UpdateVisibilityResponse {
    if (!('visibility' in value) || value['visibility'] === undefined) return false;
    return true;
}

export function UpdateVisibilityResponseFromJSON(json: any): UpdateVisibilityResponse {
    return UpdateVisibilityResponseFromJSONTyped(json, false);
}

export function UpdateVisibilityResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): UpdateVisibilityResponse {
    if (json == null) {
        return json;
    }
    return {
        
        'visibility': json['visibility'],
    };
}

export function UpdateVisibilityResponseToJSON(json: any): UpdateVisibilityResponse {
    return UpdateVisibilityResponseToJSONTyped(json, false);
}

export function UpdateVisibilityResponseToJSONTyped(value?: UpdateVisibilityResponse | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'visibility': value['visibility'],
    };
}

