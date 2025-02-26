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

import { mapValues } from '../runtime';
/**
 * 
 * @export
 * @interface UploadPost200Response
 */
export interface UploadPost200Response {
    /**
     * 
     * @type {string}
     * @memberof UploadPost200Response
     */
    url: string;
}

/**
 * Check if a given object implements the UploadPost200Response interface.
 */
export function instanceOfUploadPost200Response(value: object): value is UploadPost200Response {
    if (!('url' in value) || value['url'] === undefined) return false;
    return true;
}

export function UploadPost200ResponseFromJSON(json: any): UploadPost200Response {
    return UploadPost200ResponseFromJSONTyped(json, false);
}

export function UploadPost200ResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): UploadPost200Response {
    if (json == null) {
        return json;
    }
    return {
        
        'url': json['url'],
    };
}

export function UploadPost200ResponseToJSON(json: any): UploadPost200Response {
    return UploadPost200ResponseToJSONTyped(json, false);
}

export function UploadPost200ResponseToJSONTyped(value?: UploadPost200Response | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'url': value['url'],
    };
}

