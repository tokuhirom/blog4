import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

import "../../common.tsp";
import "../../models/entry.tsp";

using Http;
using OpenAPI;

namespace AdminAPI;

@route("/entries/{path}")
@delete
@summary("Delete an entry")
op deleteEntry(
  /**
   * The path of the entry to delete
   */
  @path path: string,
): EmptyResponse | (NotFoundResponse &
  ErrorResponse) | GeneratedHelpers.DefaultResponse<
  Description = "An unexpected error occurred",
  Body = ErrorResponse
>;

@route("/entries/{path}")
@get
@summary("Get entry by dynamic path")
op getEntryByDynamicPath(
  /**
   * The path of the entry
   */
  @path path: string,
): GetLatestEntriesRow | GeneratedHelpers.DefaultResponse<
  Description = "An unexpected error occurred",
  Body = ErrorResponse
>;
