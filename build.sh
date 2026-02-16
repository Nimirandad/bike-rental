#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Define ruta del proyecto
project_root="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export PROJECT_ROOT=$project_root
echo ">> Building in $PROJECT_ROOT"

# Validación de argumentos
if [ "$#" -lt 3 ]; then
    echo "Usage: $0 <go_compiler_path> <output_path> <source_dir_or_file> [<goos> <goarch>]"
    exit 1
fi

GO_COMPILER_PATH=$1
OUTPUT_PATH=$2
SOURCE_PATH=$3
GOOS=${4:-$(go env GOOS)}
GOARCH=${5:-$(go env GOARCH)}

echo ">> Compiling for GOOS=$GOOS, GOARCH=$GOARCH"

# Compilación
env GOOS="$GOOS" GOARCH="$GOARCH" "$GO_COMPILER_PATH" build \
  -trimpath \
  -o "$OUTPUT_PATH" "$PROJECT_ROOT/$SOURCE_PATH"

chmod +x "$OUTPUT_PATH"
echo ">> Build complete: $OUTPUT_PATH"
