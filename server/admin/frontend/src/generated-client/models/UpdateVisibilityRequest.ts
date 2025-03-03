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
 * @interface UpdateVisibilityRequest
 */
export interface UpdateVisibilityRequest {
    /**
     * The new visibility status for the entry
     * @type {string}
     * @memberof UpdateVisibilityRequest
     */
    visibility: string;
}

/**
 * Check if a given object implements the UpdateVisibilityRequest interface.
 */
export function instanceOfUpdateVisibilityRequest(value: object): value is UpdateVisibilityRequest {
    if (!('visibility' in value) || value['visibility'] === undefined) return false;
    return true;
}

export function UpdateVisibilityRequestFromJSON(json: any): UpdateVisibilityRequest {
    return UpdateVisibilityRequestFromJSONTyped(json, false);
}

export function UpdateVisibilityRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): UpdateVisibilityRequest {
    if (json == null) {
        return json;
    }
    return {
        
        'visibility': json['visibility'],
    };
}

export function UpdateVisibilityRequestToJSON(json: any): UpdateVisibilityRequest {
    return UpdateVisibilityRequestToJSONTyped(json, false);
}

export function UpdateVisibilityRequestToJSONTyped(value?: UpdateVisibilityRequest | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'visibility': value['visibility'],
    };
}

