#!/usr/bin/env bash

PLUGIN_ID=$(build/bin/manifest id)

if [[ -z ${PLUGIN_ID} ]]
then
    echo "Could not find plugin id. Exiting."
    exit 1
fi

PLUGIN_PID=$(ps aux | grep ${PLUGIN_ID} | grep -v "grep" | awk -F " " '{print $2}')

if [[ -z ${PLUGIN_PID} ]]
then
    echo "Could not find plugin PID; the plugin is not running. Exiting."
    exit 1
fi

echo "Located Plugin running with PID: ${PLUGIN_PID}"
dlv attach ${PLUGIN_PID} --listen :2346 --headless=true --api-version=2 --accept-multiclient &
