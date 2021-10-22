#!/usr/bin/env bash

# Exit on error
set -Ee

TOKEN_HEADER="Circle-Token: ${CIRCLECI_TOKEN}"

#Pass parameters to the Circle CI pipeline
PARAMETERS=""
[ "$QKM_NAMESPACE" ] && export PARAMETERS=$PARAMETERS,\"qkm_namespace\":\"$QKM_NAMESPACE\"
[ "$QKM_TAG" ] && export PARAMETERS=$PARAMETERS,\"qkm_tag\":\"$QKM_TAG\"
[ "$QKM_REPOSITORY" ] && export PARAMETERS=$PARAMETERS,\"qkm_repository\":\"$QKM_REPOSITORY\"
[ "$ENVIRONMENT_VALUES" ] && export PARAMETERS=$PARAMETERS,\"environment_values\":\"$ENVIRONMENT_VALUES\"
[ "$B64_MANIFESTS" ] && export PARAMETERS=$PARAMETERS,\"b64_manifests\":\"$B64_MANIFESTS\"
[ "$PARAMETERS" ] && PARAMETERS=${PARAMETERS:1}

echo "Pipeline parameters: $PARAMETERS"

#Create CircleCI pipeline
RESPONSE=$(curl -s --request POST --header "${TOKEN_HEADER}" --header "Content-Type: application/json" --data '{"branch":"'${BRANCH-main}'","parameters":{'${PARAMETERS}'}}' https://circleci.com/api/v2/project/github/ConsenSys/quorum-key-manager-kubernetes/pipeline)
echo $RESPONSE
ID=$(echo $RESPONSE | jq '.id' -r)
NUMBER=$(echo $RESPONSE | jq '.number' -r)

echo "Circle CI pipeline created: $ID"

#Timeout after 4*450 seconds = 30min
SLEEP=4
RETRY=450

for i in $(seq 1 1 $RETRY); do
  sleep $SLEEP

  # Get pipeline status
  STATUS=$(curl -s --request GET --header "${TOKEN_HEADER}" --header "Content-Type: application/json" https://circleci.com/api/v2/pipeline/${ID}/workflow | jq '.items[0].status' -r)
  echo "$i/$RETRY - $STATUS"

  if [[ $STATUS != 'running' ]]; then
    break
  fi

  if [ $i = $RETRY ]; then
    echo "Timeout"
  fi
done

echo "Final status: ${STATUS}"

PIPELINE_ID=$(curl -s --request GET --header "${TOKEN_HEADER}" --header "Content-Type: application/json" https://circleci.com/api/v2/pipeline/${ID}/workflow | jq '.items[0].id' -r)
echo "See the pipeline https://app.circleci.com/pipelines/github/ConsenSys/quorum-key-manager-kubernetes/${NUMBER}/workflows/${PIPELINE_ID}"

if [ "$STATUS" != "success" ]; then
  exit 1
fi
