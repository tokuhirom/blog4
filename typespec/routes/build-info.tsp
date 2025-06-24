import "@typespec/http";
import "@typespec/openapi";
import "../common.tsp";

using Http;
using OpenAPI;

namespace AdminAPI;

@route("/api/build-info")
namespace BuildInfo {
  model BuildInfo {
    buildTime: string;
    gitCommit: string;
    gitShortCommit: string;
    gitBranch: string;
    gitTag?: string;
    githubUrl: string;
  }

  @get
  @operationId("getBuildInfo")
  op get(): {
    @statusCode statusCode: 200;
    @body body: BuildInfo;
  } | {
    @statusCode statusCode: 500;
    @body body: ErrorResponse;
  };
}