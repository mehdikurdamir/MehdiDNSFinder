#!/usr/bin/env bash
# Cross-compiles the Windows app and builds the installer on Linux.
# Requires: golang, nsis, binutils-mingw-w64-x86-64
#   sudo apt-get install -y golang-go nsis binutils-mingw-w64-x86-64
set -euo pipefail
cd "$(dirname "$0")"

# windres needs a C preprocessor by default; a plain `cat` stub is enough for an icon-only .rc
PP="$(mktemp)"
printf '#!/bin/sh\nfor a; do f="$a"; done\ncat "$f"\n' > "$PP"; chmod +x "$PP"
x86_64-w64-mingw32-windres --preprocessor="$PP" -O coff -i resource.rc -o rsrc_windows_amd64.syso
rm -f "$PP"

GOOS=windows GOARCH=amd64 go build -ldflags "-H windowsgui -s -w" -o MehdiDNSFinder.exe .
makensis -V2 installer.nsi
mkdir -p dist && mv MehdiDNSFinder-Setup.exe dist/ && mv MehdiDNSFinder.exe dist/
echo "Done: dist/MehdiDNSFinder-Setup.exe"
