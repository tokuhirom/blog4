import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

using Http;
using OpenAPI;

namespace AdminAPI;

model UploadFileResponse {
  url: url;
}

@route("/upload") @post op uploadFile(
  @header contentType: "multipart/form-data",
  @body body: {
    file: bytes;
  },
): UploadFileResponse | BadRequestResponse;
