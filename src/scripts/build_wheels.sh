#!/usr/bin/env bash
# build_wheels.sh — cross-compile the docgen binary for all supported platforms
# and produce a platform-specific Python wheel for each.
#
# Usage (from repo root):
#   bash src/scripts/build_wheels.sh
#
# Output: src/dist/*.whl  (one wheel per platform)
#
# Requirements: Go toolchain, Python 3 with `build` and `wheel` packages.
#   pip install build wheel
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PYTHON_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
REPO_DIR="$(cd "${PYTHON_DIR}/.." && pwd)"
BIN_DIR="${PYTHON_DIR}/cisco_docgen/_bin"
DIST_DIR="${PYTHON_DIR}/dist"

cd "${REPO_DIR}"

# Ensure build tools are available
python -m pip install --quiet build wheel

# Clean previous artifacts
rm -f "${BIN_DIR}"/docgen "${BIN_DIR}"/docgen.exe
rm -rf "${DIST_DIR}"
mkdir -p "${DIST_DIR}"

build_wheel() {
    local goos="$1"
    local goarch="$2"
    local plat_name="$3"
    local bin_name="${4:-docgen}"

    echo "==> Building Go binary: GOOS=${goos} GOARCH=${goarch}"
    rm -f "${BIN_DIR}/docgen" "${BIN_DIR}/docgen.exe"
    GOOS="${goos}" GOARCH="${goarch}" CGO_ENABLED=0 \
        go build -trimpath -o "${BIN_DIR}/${bin_name}" ./cmd/docgen

    echo "==> Building wheel: ${plat_name}"
    (
        cd "${PYTHON_DIR}"
        DOCGEN_WHEEL_PLAT="${plat_name}" python -m build --wheel --outdir "${DIST_DIR}"
    )

    echo "==> Done: $(ls "${DIST_DIR}"/*.whl | tail -1)"
}

build_wheel linux amd64 manylinux_2_17_x86_64.manylinux2014_x86_64
build_wheel linux arm64 manylinux_2_17_aarch64.manylinux2014_aarch64
build_wheel darwin amd64 macosx_10_9_x86_64
build_wheel darwin arm64 macosx_11_0_arm64
build_wheel windows amd64 win_amd64 docgen.exe

# Clean up binary so the source tree is left tidy
rm -f "${BIN_DIR}/docgen" "${BIN_DIR}/docgen.exe"

echo ""
echo "All wheels built:"
ls -1 "${DIST_DIR}"/*.whl
