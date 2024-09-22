#!/bin/bash

# Disable cgo checks
# gomdb uses cgo incorrectly, so we need to disable cgo checks
export GODEBUG=cgocheck=0

# Resolve the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Execute tukkid-bin from the script's directory
"${SCRIPT_DIR}/tukkid-bin" "$@"