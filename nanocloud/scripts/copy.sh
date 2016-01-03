#!/bin/bash

# Load configuration file
source configuration.sh

if [ ${#} -lt 1 ]; then
    echo "Not enough arguments"
      exit 1
    fi

    FILENAME=${1}
    SCP=$(which scp)

    sshpass -p "${PASSWORD}" "${SCP}" -o StrictHostKeyChecking=no -P "${PORT}" "${FILENAME}" "${USER}@${SERVER}:~/Desktop/"
    if [ $? -eq 0 ]; then
        rm "${FILENAME}"
      fi
