import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

import "../../../common.tsp";

using Http;
using OpenAPI;

namespace AdminAPI;

@route("/entries/{path}/regenerate-image")
@post
@summary("Regenerate entry image")
op regenerateEntryImage(
  /**
   * The path of the entry to regenerate the image for
   */
  @path path: string,
): EmptyResponse | (NotFoundResponse &
  ErrorResponse) | GeneratedHelpers.DefaultResponse<
  Description = "An unexpected error occurred",
  Body = ErrorResponse
>;
