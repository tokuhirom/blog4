import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

import "../../common.tsp";

using Http;
using OpenAPI;

namespace AdminAPI;

@route("/entries/titles")
@get
@summary("Get all entry titles")
op getAllEntryTitles(): EntryTitlesResponse | GeneratedHelpers.DefaultResponse<
  Description = "An unexpected error occurred",
  Body = ErrorResponse
>;

/**
 * Array of all entry titles
 */
model EntryTitlesResponse is string[];
