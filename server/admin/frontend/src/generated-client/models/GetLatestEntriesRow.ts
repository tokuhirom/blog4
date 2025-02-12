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
 * @interface GetLatestEntriesRow
 */
export interface GetLatestEntriesRow {
    /**
     * 
     * @type {string}
     * @memberof GetLatestEntriesRow
     */
    path?: string;
    /**
     * 
     * @type {string}
     * @memberof GetLatestEntriesRow
     */
    title?: string;
    /**
     * 
     * @type {string}
     * @memberof GetLatestEntriesRow
     */
    body?: string;
    /**
     * 
     * @type {string}
     * @memberof GetLatestEntriesRow
     */
    visibility?: string;
    /**
     * 
     * @type {string}
     * @memberof GetLatestEntriesRow
     */
    format?: string;
    /**
     * 
     * @type {Date}
     * @memberof GetLatestEntriesRow
     */
    publishedAt?: Date | null;
    /**
     * 
     * @type {Date}
     * @memberof GetLatestEntriesRow
     */
    lastEditedAt?: Date | null;
    /**
     * 
     * @type {Date}
     * @memberof GetLatestEntriesRow
     */
    createdAt?: Date | null;
    /**
     * 
     * @type {Date}
     * @memberof GetLatestEntriesRow
     */
    updatedAt?: Date | null;
    /**
     * 
     * @type {string}
     * @memberof GetLatestEntriesRow
     */
    imageUrl?: string | null;
}

/**
 * Check if a given object implements the GetLatestEntriesRow interface.
 */
export function instanceOfGetLatestEntriesRow(value: object): value is GetLatestEntriesRow {
    return true;
}

export function GetLatestEntriesRowFromJSON(json: any): GetLatestEntriesRow {
    return GetLatestEntriesRowFromJSONTyped(json, false);
}

export function GetLatestEntriesRowFromJSONTyped(json: any, ignoreDiscriminator: boolean): GetLatestEntriesRow {
    if (json == null) {
        return json;
    }
    return {
        
        'path': json['Path'] == null ? undefined : json['Path'],
        'title': json['Title'] == null ? undefined : json['Title'],
        'body': json['Body'] == null ? undefined : json['Body'],
        'visibility': json['Visibility'] == null ? undefined : json['Visibility'],
        'format': json['Format'] == null ? undefined : json['Format'],
        'publishedAt': json['PublishedAt'] == null ? undefined : (new Date(json['PublishedAt'])),
        'lastEditedAt': json['LastEditedAt'] == null ? undefined : (new Date(json['LastEditedAt'])),
        'createdAt': json['CreatedAt'] == null ? undefined : (new Date(json['CreatedAt'])),
        'updatedAt': json['UpdatedAt'] == null ? undefined : (new Date(json['UpdatedAt'])),
        'imageUrl': json['ImageUrl'] == null ? undefined : json['ImageUrl'],
    };
}

export function GetLatestEntriesRowToJSON(json: any): GetLatestEntriesRow {
    return GetLatestEntriesRowToJSONTyped(json, false);
}

export function GetLatestEntriesRowToJSONTyped(value?: GetLatestEntriesRow | null, ignoreDiscriminator: boolean = false): any {
    if (value == null) {
        return value;
    }

    return {
        
        'Path': value['path'],
        'Title': value['title'],
        'Body': value['body'],
        'Visibility': value['visibility'],
        'Format': value['format'],
        'PublishedAt': value['publishedAt'] == null ? undefined : ((value['publishedAt'] as any).toISOString()),
        'LastEditedAt': value['lastEditedAt'] == null ? undefined : ((value['lastEditedAt'] as any).toISOString()),
        'CreatedAt': value['createdAt'] == null ? undefined : ((value['createdAt'] as any).toISOString()),
        'UpdatedAt': value['updatedAt'] == null ? undefined : ((value['updatedAt'] as any).toISOString()),
        'ImageUrl': value['imageUrl'],
    };
}

