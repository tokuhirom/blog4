import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

using Http;
using OpenAPI;

namespace AdminAPI;

@route("/entries/{path}/link-pallet")
@get
@summary("Get linked entry paths")
op getLinkPallet(
  /**
   * The source entry path
   */
  @path path: string,
): LinkPalletData | GeneratedHelpers.DefaultResponse<
  Description = "An unexpected error occurred",
  Body = ErrorResponse
>;

model LinkPalletData {
  /**
   * Array of potential new link titles
   */
  newLinks: string[];

  /**
   * Array of directly linked entries
   */
  links: EntryWithImage[];

  /**
   * Array of two-hop link relationships
   */
  twohops: TwoHopLink[];
}

model TwoHopLink {
  src: EntryWithDestTitle;
  links: EntryWithImage[];
}

model EntryWithDestTitle {
  ...EntryWithImage;
  dstTitle: string;
}

model EntryWithImage {
  path: string;
  title: string;
  body: string;
  visibility: string;
  format: string;
  imageUrl?: string | null;
}
