#!/bin/bash

GITTER_USER=$1
GITTER_TOKEN=$2
GITTER_ROOM_NAME=$3

: ${GITTER_USER:? required}
: ${GITTER_TOKEN:? required}
: ${DOMAIN:? required}
#: ${GITTER_ROOM_NAME:? required}

debug() {
  [[ "$DEBUG" ]] && echo "-----> $*" 1>&2
}

gitterApi() {
    curl -sX POST \
    -H "Content-Type: application/json" \
    -H "Accept: application/json" \
    -H "Authorization: Bearer ${GITTER_TOKEN}" \
    "$@"
}

joinRoom() {
  userId=$(curl -sH "Authorization: Bearer $GITTER_TOKEN" https://api.gitter.im/v1/user/me | jq .id -r)
  debug joining gitter room: $GITTER_ROOM_NAME userId: $userId

  if [[ $GITTER_ROOM_NAME == "" ]]; then
    debug GITTER_ROOM_NAME unset, no autojoin will happen ...
    return
  fi

  GITTER_ROOM_ID=$( gitterApi "https://api.gitter.im/v1/rooms" \
      -d '{"uri":"'$GITTER_ROOM_NAME'"}' \
      | jq .id -r)
  debug GITTER_ROOM_ID: $GITTER_ROOM_ID

  # join room
  gitterApi "https://api.gitter.im/v1/user/${userId}/rooms" \
      -d '{"id":"'${GITTER_ROOM_ID}'"}' \
      &> /dev/null

  # post welcome message
  gitterApi "https://api.gitter.im/v1/rooms/${GITTER_ROOM_ID}/chatMessages" \
      -d '{"text":"Hi im using session: '${deployment}'"}' \
      &> /dev/null
}

# get assigned or assign new deployment
if ! deployment=$(kubectl get deployments -l ghuser=${GITTER_USER} -o jsonpath='{.items[0].metadata.name}' 2>/dev/null); then
  debug "no session assigned to ${GITTER_USER} ..."

  nextUnassigned=$(kubectl get deployments -l 'user,!ghuser' -o jsonpath='{.items[0].metadata.name}')
  debug nextUnassigned=$nextUnassigned
  deployment="${nextUnassigned}"
  kubectl label deployment $deployment --overwrite ghuser=${GITTER_USER} &>/dev/null
fi

pod=$(kubectl get po -lrun=${deployment} -o jsonpath='{.items[0].metadata.name}')
debug pod=$pod

rndPath=$(kubectl logs ${pod} |sed -n '/HTTP server is listening at/ s/.*:8080//p')
debug rndPath=$rndPath

# deletes "session." prefix from authRedirectUrl's domain
url="http://${deployment}.${DOMAIN#session.}${rndPath}"
debug url=$url

kubectl annotate deployment $deployment --overwrite sessionurl=$url &> /dev/null


externalip=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type == "ExternalIP")].address}') 
nodePort=$(kubectl get svc ${deployment} -o jsonpath="{.spec.ports[0].nodePort}")
sessionUrlNodePort="http://${externalip}:${nodePort}${rndPath}"
kubectl annotate deployments ${deployment} --overwrite sessionurlnp=${sessionUrlNodePort} &>/dev/null

joinRoom

sessionUrl=$(kubectl get deployments $deployment -o jsonpath='{.metadata.annotations.sessionurl}')
cat << EOF
<p/><a href='https://gitter.im/${GITTER_ROOM_NAME}' >Open GITTER chat</a>
<p/><a href="${sessionUrl}">web session - ${deployment} [Ingress]</a>
<p/><a href="${sessionUrlNodePort}">web session - ${deployment} [NodePort]</a>
EOF
