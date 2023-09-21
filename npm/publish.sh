#!/bin/bash
set -euo pipefail 

npm config set //registry.npmjs.org/:_authToken=$NPM_TOKEN
npm ci --ignore-script

PACKAGE_NAME=$(cat package.json | jq -r '.name')
PACKAGE_VERSION=$(cat package.json | jq -r '.version')

echo "Get otp from Optic for ${PACKAGE_NAME}:${PACKAGE_VERSION}"

OTP=$(curl -s \
-d "{ \"packageInfo\": { \"version\": \"$PACKAGE_VERSION\", \"name\": \"$PACKAGE_NAME\" } }" \
-H "Content-Type: application/json" \
-X POST \
-f https://optic-zf3votdk5a-ew.a.run.app/api/generate/$OPTIC_TOKEN)

echo "Publish package ${PACKAGE_NAME}:${PACKAGE_VERSION}"
npm publish --otp ${OTP} --access public --provenance
