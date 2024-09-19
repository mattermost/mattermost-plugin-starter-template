#!/bin/sh
echo "Clone mattermost"
cd ../../
git clone --depth 1 https://github.com/mattermost/mattermost.git

echo "==> Install mattermost/webapp dependencies"
cd mattermost/webapp
npm i

echo "==> Install mattermost/e2e-tests/playwright dependencies"
cd ../e2e-tests/playwright
npm i

echo "==> Back to plugin folder and install e2e dependencies"
cd ../../../mattermost-plugin-starter-template/e2e
npm i
