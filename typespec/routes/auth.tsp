import "@typespec/http";
import "../common.tsp";

using Http;

namespace AdminAPI;

model LoginRequest {
  username: string;
  password: string;
}

model LoginResponse {
  success: boolean;
  message?: string;
}

model AuthCheckResponse {
  authenticated: boolean;
  username?: string;
}

@route("/auth")
interface Auth {
  @post
  @route("/login")
  @tag("Auth")
  @summary("Login with username and password")
  login(@body request: LoginRequest): {
    @statusCode statusCode: 200;
    @body body: LoginResponse;
  } | {
    @statusCode statusCode: 401;
    @body body: ErrorResponse;
  };

  @post
  @route("/logout")
  @tag("Auth")
  @summary("Logout and invalidate session")
  logout(): {
    @statusCode statusCode: 200;
    @body body: EmptyResponse;
  } | {
    @statusCode statusCode: 401;
    @body body: ErrorResponse;
  };

  @get
  @route("/check")
  @tag("Auth")
  @summary("Check if user is authenticated")
  check(): {
    @statusCode statusCode: 200;
    @body body: AuthCheckResponse;
  } | {
    @statusCode statusCode: 401;
    @body body: ErrorResponse;
  };
}