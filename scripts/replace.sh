#!/bin/bash

FILENAME=token.mist
SCRIPT=deploy.ts

BYTECODE=$(go run cmd/mist.go <examples/${FILENAME})

echo "bytecode: ${BYTECODE}"

# sed -i -rz \
#     "s/2,\n\s+data: \"+0x[0-9A-Fa-f]+\"/2,\n        data: \"${BYTECODE}\"/" \
#     playground/scripts/deploy.ts

sed -i -r 's/data: "",/data: "0xabcd",/' "playground/scripts/$SCRIPT"

sed -i -r \
    "s/data: \"0x[0-9A-Fa-f]+\",/data: \"${BYTECODE}\",/" \
    "playground/scripts/$SCRIPT"
