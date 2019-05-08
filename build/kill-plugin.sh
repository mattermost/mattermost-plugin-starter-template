#!/usr/bin/env bash

set -euf -o pipefail

PLUGIN_ID=$(build/bin/manifest id)

if [[ -z ${PLUGIN_ID} ]]
then
    echo "Could not find plugin id. Exiting."
    exit 1
fi

PLUGIN_PID=$(ps aux | grep "plugins/${PLUGIN_ID}" | grep -v "grep" | awk -F " " '{print $2}')

echo "Located Plugin running with PID: ${PLUGIN_PID}. Killing."
kill -9 ${PLUGIN_PID}
