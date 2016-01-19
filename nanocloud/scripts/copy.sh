#!/bin/bash

if [ ${#} -lt 1 ]; then
  echo "Not enough arguments"
  exit 1
fi

FILENAME=${1}
SCP=$(which scp)

sshpass -p "${WIN_PASSWORD}" "${SCP}" -o StrictHostKeyChecking=no -P "${WIN_PORT}" "${FILENAME}" "${WIN_USER}@${WIN_SERVER}:~/Desktop/"
if [ $? -eq 0 ]; then
  rm "${FILENAME}"
fi
