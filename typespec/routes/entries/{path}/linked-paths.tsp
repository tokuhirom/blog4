import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

import "../../../common.tsp";

using Http;
using OpenAPI;

namespace AdminAPI;

@route("/entries/{path}/linked-paths")
@get
@summary("Get linked entry paths")
op getLinkedEntryPaths(
  /**
   * The source entry path
   */
  @path path: string,
): LinkedEntryPathsResponse | (NotFoundResponse &
  ErrorResponse) | GeneratedHelpers.DefaultResponse<
  Description = "An unexpected error occurred",
  Body = ErrorResponse
>;

/**
 * Object where keys are lowercase destination entry titles and values are their paths (null if entry doesn't exist)..
 */
model LinkedEntryPathsResponse {
  ...Record<string>;
}
