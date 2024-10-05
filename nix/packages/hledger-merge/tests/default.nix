{
  perSystem,
  pkgs,
  ...
}: let
  inherit (perSystem.self) hledger-merge;
in {
  smoke =
    pkgs.runCommand "hledger-merge-smoke-test"
    {
      nativeBuildInputs = [hledger-merge];
    }
    ''
      hledger-merge --help > $out
      if ! grep -qF 'hledger-merge [global options] command [command options]' $out; then
        echo "smoke test failed; help text not found"
        cat $out
        exit 1
      fi
    '';
}
