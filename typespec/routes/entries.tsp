import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

import "../common.tsp";
import "../models/entry.tsp";

using Http;
using OpenAPI;

namespace AdminAPI;

@route("/entries")
@get
@summary("Get latest entries")
op getLatestEntries(
  /**
   * Filter entries by the last edited date
   */
  @query(#{ explode: true }) last_last_edited_at?: utcDateTime,
): Body<GetLatestEntriesRow[]> | GeneratedHelpers.DefaultResponse<
  Description = "An unexpected error occurred",
  Body = ErrorResponse
>;

@route("/entries")
@post
@summary("Create a new entry")
op createEntry(@body body: CreateEntryRequest): (CreatedResponse &
  CreateEntryResponse) | (ConflictResponse &
  ErrorResponse) | GeneratedHelpers.DefaultResponse<
  Description = "An unexpected error occurred",
  Body = ErrorResponse
>;

model CreateEntryRequest {
  /**
   * The title of the new entry
   */
  title?: string;
}

model CreateEntryResponse {
  /**
   * The path of the created entry
   */
  path: string;
}
