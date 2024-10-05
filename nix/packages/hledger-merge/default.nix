{
  pname,
  pkgs,
  flake,
  inputs,
  perSystem,
  ...
} @ args:
  perSystem.gomod2nix.buildGoApplication rec {
    inherit pname;
    # ensure we are using the same version of go to build with
    inherit (pkgs) go;

    version = "0.0.2";

    src = let
      filter = inputs.nix-filter.lib;
    in
      filter {
        root = ../../../.;
        exclude = [
          "nix/"
          "docs/"
          ".github/"
          "README.md"
          "default.nix"
          "shell.nix"
          ".env"
          ".envrc"
        ];
      };

    modules = ./gomod2nix.toml;

    CGO_ENABLED = 0;

    ldflags = [
      "-s"
      "-w"
      "-X github.com/fxyoge/hledger-merge/build.Name=${pname}"
      "-X github.com/fxyoge/hledger-merge/build.Version=v${version}"
    ];

    preCheck = ''
      XDG_CACHE_HOME=$(mktemp -d)
      export XDG_CACHE_HOME
    '';

    meta = with lib; {
      description = "hledger-merge merges hledger files";
      homepage = "https://github.com/fxyoge/hledger-merge";
      license = licenses.mit;
      mainProgram = "hledger-merge";
    };

    passthru.tests = (import ./tests) args;
  }
