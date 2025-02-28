import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

import "../../../common.tsp";

using Http;
using OpenAPI;

namespace AdminAPI;

@route("/entries/{path}/title")
@put
@summary("Update entry title")
op updateEntryTitle(
  /**
   * The path of the entry to update
   */
  @path path: string,

  @body body: UpdateEntryTitleRequest,
):
  | EmptyResponse
  | (NotFoundResponse & ErrorResponse)
  | (ConflictResponse & ErrorResponse)
  | GeneratedHelpers.DefaultResponse<
      Description = "An unexpected error occurred",
      Body = ErrorResponse
    >;

model UpdateEntryTitleRequest {
  /**
   * The new title for the entry
   */
  title: string;
}

