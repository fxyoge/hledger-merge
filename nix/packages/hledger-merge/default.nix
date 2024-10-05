{
  pname,
  pkgs,
  flake,
  inputs,
  perSystem,
  ...
}: let
  inherit (pkgs) lib;
in
  perSystem.gomod2nix.buildGoApplication rec {
    inherit pname;
    # there's no good way of tying in the version to a git tag or branch
    # so for simplicity's sake we set the version as the commit revision hash
    # we remove the `-dirty` suffix to avoid a lot of unnecessary rebuilds in local dev
    version = lib.removeSuffix "-dirty" (flake.shortRev or flake.dirtyShortRev);

    # ensure we are using the same version of go to build with
    inherit (pkgs) go;

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
  }
