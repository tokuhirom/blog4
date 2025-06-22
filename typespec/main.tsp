import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

import "./common.tsp";

import "./models/entry.tsp";

import "./routes/entries/{path}/body.tsp";
import "./routes/entries/{path}/link-pallet.tsp";
import "./routes/entries/{path}/linked-paths.tsp";
import "./routes/entries/{path}/regenerate-image.tsp";
import "./routes/entries/{path}/title.tsp";
import "./routes/entries/{path}/visibility.tsp";
import "./routes/entries/titles.tsp";
import "./routes/entries/{path}.tsp";
import "./routes/entries.tsp";
import "./routes/upload.tsp";
import "./routes/auth.tsp";

using Http;
using OpenAPI;

@service(#{
  title: "Admin API",
})
@info(#{
  version: "1.0.0",
})
namespace AdminAPI;
