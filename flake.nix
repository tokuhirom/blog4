{
  description = "blog4 development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Go
            go
            golangci-lint

            # Tools
            go-task   # Task runner

            # Database
            postgresql_14

            # Node.js (admin frontend)
            nodejs
          ];

          shellHook = ''
            echo "🚀 blog4 development environment"
            echo "Go version: $(go version)"
          '';
        };
      });
}
