{
  pkgs,
  pname,
  ...
}:
pkgs.mkShellNoCC {
  inherit pname;
  packages = [
    pkgs.go
    pkgs.gotest
  ];
}
