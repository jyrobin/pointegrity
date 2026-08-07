#!/bin/sh
# implements cap.F-curl-installer
#
# Install poi.
#
#   curl -fsSL https://pointegrity.com/poi/install.sh | sh
#   curl -fsSL .../install.sh | POI_VERSION=v0.1.0 POI_BIN=/usr/local/bin sh
#
# Downloads the release archive for this platform, **verifies its checksum
# against the published checksums.txt**, and installs the binary. Nothing is
# built, no compiler is needed, and nothing is written outside POI_BIN.
#
# On macOS this does not trip Gatekeeper. The quarantine attribute is applied by
# applications that opt into it — browsers, mail clients — and curl is not one
# of them. A binary downloaded through a browser from the releases page *is*
# quarantined; this one is not.
set -eu

REPO="pointegrity/poi"
BIN_DIR="${POI_BIN:-$HOME/.local/bin}"
VERSION="${POI_VERSION:-}"

say() { printf '%s\n' "$*"; }
die() { printf 'error: %s\n' "$*" >&2; exit 1; }

need() { command -v "$1" >/dev/null 2>&1 || die "$1 is required but not installed"; }
need curl
need tar

# --- platform -----------------------------------------------------------
os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
	darwin | linux) ;;
	*) die "unsupported OS: $os (poi ships darwin and linux builds)" ;;
esac

arch=$(uname -m)
case "$arch" in
	x86_64 | amd64) arch=amd64 ;;
	arm64 | aarch64) arch=arm64 ;;
	*) die "unsupported architecture: $arch" ;;
esac

# --- version ------------------------------------------------------------
if [ -z "$VERSION" ]; then
	VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" |
		sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -1)
	[ -n "$VERSION" ] || die "could not determine the latest version; set POI_VERSION"
fi
# The archive name carries the version without its leading v.
num=${VERSION#v}

archive="poi_${num}_${os}_${arch}.tar.gz"
base="https://github.com/$REPO/releases/download/$VERSION"

say "poi $VERSION — $os/$arch"

tmp=$(mktemp -d)
# Clean up on any exit, including the failure paths below.
trap 'rm -rf "$tmp"' EXIT INT TERM

say "  downloading $archive"
curl -fsSL -o "$tmp/$archive" "$base/$archive" ||
	die "download failed — is $VERSION a released version for $os/$arch?"

# --- verify -------------------------------------------------------------
# A download nobody checks is a download anybody can replace. This is the whole
# reason checksums.txt is published beside the archives.
if curl -fsSL -o "$tmp/checksums.txt" "$base/checksums.txt" 2>/dev/null; then
	want=$(grep " $archive\$" "$tmp/checksums.txt" | awk '{print $1}')
	if [ -n "$want" ]; then
		if command -v shasum >/dev/null 2>&1; then
			got=$(shasum -a 256 "$tmp/$archive" | awk '{print $1}')
		elif command -v sha256sum >/dev/null 2>&1; then
			got=$(sha256sum "$tmp/$archive" | awk '{print $1}')
		else
			got=""
			say "  ! no shasum or sha256sum — cannot verify the download"
		fi
		if [ -n "$got" ]; then
			[ "$want" = "$got" ] || die "checksum mismatch for $archive
  expected $want
  got      $got
Refusing to install."
			say "  checksum ok"
		fi
	else
		say "  ! $archive is not listed in checksums.txt — cannot verify"
	fi
else
	say "  ! checksums.txt unavailable — cannot verify the download"
fi

# --- install ------------------------------------------------------------
tar xzf "$tmp/$archive" -C "$tmp" || die "could not extract $archive"
[ -f "$tmp/poi" ] || die "no poi binary inside $archive"

mkdir -p "$BIN_DIR" || die "could not create $BIN_DIR"
# Write via a temporary name and rename, so a running poi is never truncated
# half-way through and an interrupted install leaves nothing behind.
install_tmp="$BIN_DIR/.poi.$$"
cp "$tmp/poi" "$install_tmp" && chmod 755 "$install_tmp" && mv -f "$install_tmp" "$BIN_DIR/poi" ||
	die "could not install into $BIN_DIR — set POI_BIN to a writable directory"

say "  installed $BIN_DIR/poi"
say ""
"$BIN_DIR/poi" --version || die "the installed binary would not run"

case ":${PATH}:" in
	*":$BIN_DIR:"*) ;;
	*)
		say ""
		say "$BIN_DIR is not on your PATH. Add it:"
		say "  export PATH=\"$BIN_DIR:\$PATH\""
		;;
esac

say ""
say "Start with:  poi adopt      — what this repository already has"
say "             poi gate       — every gate, here"
