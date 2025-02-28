import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

import "../../../common.tsp";

using Http;
using OpenAPI;

namespace AdminAPI;

@route("/entries/{path}/body")
@put
@summary("Update entry body")
op updateEntryBody(
  /**
   * The path of the entry to update
   */
  @path path: string,

  @body body: UpdateEntryBodyRequest,
): EmptyResponse | (NotFoundResponse &
  ErrorResponse) | GeneratedHelpers.DefaultResponse<
  Description = "An unexpected error occurred",
  Body = ErrorResponse
>;

model UpdateEntryBodyRequest {
  /**
   * The new content of the entry
   */
  body: string;
}

