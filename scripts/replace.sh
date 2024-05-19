#!/bin/bash

TARGET=${TARGET:-token.mist}
SCRIPT=${SCRIPT:-deploy.ts}

echo "target: ${TARGET}"
echo "script: ${SCRIPT}"

BYTECODE=$(go run cmd/mist.go < "examples/${TARGET}")

echo "bytecode: ${BYTECODE}"

# sed -i -rz \
#     "s/2,\n\s+data: \"+0x[0-9A-Fa-f]+\"/2,\n        data: \"${BYTECODE}\"/" \
#     playground/scripts/deploy.ts

sed -i -r 's/data: "",/data: "0xabcd",/' "playground/scripts/$SCRIPT"

sed -i -r \
    "s/data: \"0x[0-9A-Fa-f]+\",/data: \"${BYTECODE}\",/" \
    "playground/scripts/$SCRIPT"
