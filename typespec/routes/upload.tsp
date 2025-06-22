import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

using Http;
using OpenAPI;

namespace AdminAPI;

model UploadFileResponse {
  url: url;
}

model UploadFileRequest {
  file: HttpPart<bytes>;
}

@route("/upload") @post op uploadFile(
  @header contentType: "multipart/form-data",
  @multipartBody body: UploadFileRequest,
): UploadFileResponse | BadRequestResponse;
