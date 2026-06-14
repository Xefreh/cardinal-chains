#!/usr/bin/env bash
# Installs the build dependencies for Cardinal Chains (currently just libyaml).
# Detects the system's package manager and installs the appropriate package.
set -euo pipefail

if [[ "${EUID}" -ne 0 ]]; then
  echo "This script needs to install system packages. Re-running with sudo."
  exec sudo -E "$0" "$@"
fi

install_with() {
  local pm="$1"; shift
  local pkgs=("$@")
  echo "Detected package manager: ${pm}"
  case "${pm}" in
    apt-get) apt-get update && apt-get install -y "${pkgs[@]}" ;;
    dnf)     dnf install -y "${pkgs[@]}" ;;
    yum)     yum install -y "${pkgs[@]}" ;;
    pacman)  pacman -Sy --noconfirm --needed "${pkgs[@]}" ;;
    zypper)  zypper install -y "${pkgs[@]}" ;;
    apk)     apk add --no-cache "${pkgs[@]}" ;;
    brew)    brew install "${pkgs[@]}" ;;
    *) echo "Unsupported package manager: ${pm}" >&2; exit 1 ;;
  esac
}

# Map each package manager to its libyaml dev package name.
declare -A YAML_PKG=(
  [apt-get]=libyaml-dev
  [dnf]=libyaml-devel
  [yum]=libyaml-devel
  [pacman]=libyaml
  [zypper]=libyaml-devel
  [apk]=yaml-dev
  [brew]=libyaml
)

for pm in apt-get dnf yum pacman zypper apk brew; do
  if command -v "${pm}" >/dev/null 2>&1; then
    install_with "${pm}" "${YAML_PKG[${pm}]}"
    echo "Done."
    exit 0
  fi
done

cat >&2 <<'EOF'
Could not detect a supported package manager.
Please install libyaml manually:
  - Debian/Ubuntu:  sudo apt-get install -y libyaml-dev
  - Fedora/RHEL:    sudo dnf install -y libyaml-devel
  - Arch:           sudo pacman -S libyaml
  - openSUSE:       sudo zypper install libyaml-devel
  - Alpine:         sudo apk add yaml-dev
  - macOS (brew):   brew install libyaml
EOF
exit 1
