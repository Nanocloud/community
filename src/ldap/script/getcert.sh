#!/bin/bash
mkdir -p /opt/conf
sshpass -p ${PASSWORD} scp -P ${SSH_PORT} -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ${USER}@${SERVER}:/cygdrive/c/users/administrator/ad2012.cer /opt/conf
