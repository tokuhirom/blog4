import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

import "../../../common.tsp";

using Http;
using OpenAPI;

namespace AdminAPI;

@route("/entries/{path}/visibility")
@post
@summary("Update entry visibility")
op updateEntryVisibility(
  /**
   * The path of the entry to update visibility
   */
  @path path: string,

  @body body: UpdateVisibilityRequest,
): UpdateVisibilityResponse | (NotFoundResponse &
  ErrorResponse) | GeneratedHelpers.DefaultResponse<
  Description = "An unexpected error occurred",
  Body = ErrorResponse
>;

model UpdateVisibilityRequest {
  /**
   * The new visibility status for the entry
   */
  visibility: string;
}

model UpdateVisibilityResponse {
  /**
   * The new visibility status for the entry
   */
  visibility: string;
}
