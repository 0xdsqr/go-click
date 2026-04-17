{
  description = "go-click";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShells.default = with pkgs; mkShell {
          packages = [
            # Go compiler and standard tooling.
            go
          ];

          shellHook = ''
            echo "go-click dev shell"
            go version
          '';
        };

        checks.default = pkgs.runCommand "go-click-check" {
          nativeBuildInputs = [ pkgs.go ];
        } ''
          export HOME="$TMPDIR"
          cd ${self}

          test -z "$(gofmt -l .)"
          go vet ./...
          go test ./...
          go build ./...

          touch "$out"
        '';
      });
}
