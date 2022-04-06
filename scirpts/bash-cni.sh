#!/bin/bash

log () {
    >&2 echo "$1"
}

log "Command: $CNI_COMMAND"
log "Container Id: $CNI_CONTAINERID"
log "Path to Netowork Namespace: $CNI_NETNS"
log "Networker interface: $CNI_IFNAME"
log "CNI PATH: $CNI_PATH"

case $CNI_COMMAND in
ADD)
    echo "{}"
    ;;
DEL)
    echo "{}"
    ;;
VERSION)
    echo '{
    "cniVersion": "0.3.1", 
    "supportedVersions": [ "0.3.0", "0.3.1", "0.4.0" ] 
}'
    ;;
*)
    echo "Unknow Command: $CNI_COMMAND"
    exit 1
    ;;
esac
