import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

using Http;
using OpenAPI;

namespace AdminAPI;

model UploadFileResponse {
  url: url;
}

model UploadFileBody {
  file: HttpPart<bytes>;
}

@route("/upload") @post op uploadFile(
  @header contentType: "multipart/form-data",
  @multipartBody body: UploadFileBody,
): UploadFileResponse | BadRequestResponse;
