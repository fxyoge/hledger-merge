{
  pkgs,
  pname,
  perSystem,
  ...
}:
pkgs.mkShellNoCC {
  inherit pname;
  packages = [
    (pkgs.writeShellApplication {
      name = "gomod2nix:update";
      runtimeInputs = [pkgs.git perSystem.gomod2nix.default];
      text = ''
        set -eu
        gomod2nix --outdir nix/packages/hledger-merge
        # shellcheck disable=SC2016
        sed -i '1i # Generated with `nix develop .#renovate -c gomod2nix:update`' nix/packages/hledger-merge/gomod2nix.toml
      '';
    })
    (pkgs.writeShellApplication {
      name = "version:bump";
      text = ''
        set -euo pipefail

        version_file="$1"
        if [ ! -f "$version_file" ]; then
          echo "$version_file does not exist" 1>&2
          exit 1
        fi

        current_version=$(grep 'version =' "$version_file" | awk -F '"' '{print $2}')

        bump_type="patch"
        if git log "$(git describe --tags --abbrev=0)..HEAD" --pretty=format:"%s" | grep -qF '[no bump]'; then
          bump_type="noop"
        elif git log "$(git describe --tags --abbrev=0)..HEAD" --pretty=format:"%s" | grep -q "^BREAKING CHANGE" ||
           git log "$(git describe --tags --abbrev=0)..HEAD" --pretty=format:"%s" | grep -q "^[a-z]*!:"; then
          bump_type="major"
        elif git log "$(git describe --tags --abbrev=0)..HEAD" --pretty=format:"%s" | grep -q "^feat:"; then
          bump_type="minor"
        fi

        IFS='.' read -ra VERSION_PARTS <<< "$current_version"
        major=''${VERSION_PARTS[0]}
        minor=''${VERSION_PARTS[1]}
        patch=''${VERSION_PARTS[2]}

        case $bump_type in
          major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
          minor)
            minor=$((minor + 1))
            patch=0
            ;;
          patch)
            patch=$((patch + 1))
            ;;
        esac

        new_version="$major.$minor.$patch"
        sed -i 's/version = ".*"/version = "'"$new_version"'"/' "$version_file"
        echo "$new_version"
      '';
    })
  ];
}
