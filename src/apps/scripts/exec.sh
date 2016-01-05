#!/bin/bash

# Load configuration file
source ../src/apps/scripts/configuration.sh

if [ ${#} -lt 1 ]; then
    echo "Not enough arguments"
      exit 1
    fi

    COMMAND=${1}
    SSH=$(which ssh)

    sshpass -p "${PASSWORD}" "${SSH}" -o StrictHostKeyChecking=no -p "${PORT}" "${USER}@${SERVER}" "${COMMAND}"
