node {
  sshagent(['releases']) {
    sh 'ssh -o StrictHostKeyChecking=no bamboo@releases.nanocloud.com "mkdir /data/artifacts/community/builds/${BRANCH_NAME} || true"'
  }
}

node {

  stage 'Test backend'

  checkout scm

  sh 'docker rmi -f `docker images -aq` || true'

  sshagent(['releases']) {
    sh 'scp -o StrictHostKeyChecking=no bamboo@releases.nanocloud.com:/data/artifacts/community/canary/nanocloud-backend.tar nanocloud-backend.tar'
  }

  sh 'docker load -i nanocloud-backend.tar'

  sh 'docker-compose -f modules/docker-compose-build.yml build nanocloud-backend'
  sh "docker save nanocloud/nanocloud-backend \$(docker history -q nanocloud/nanocloud-backend | tail -n +2 | grep -v \\<missing\\> | tr '\n' ' ') > nanocloud-backend.tar"

  sshagent(['releases']) {
    sh 'scp -o StrictHostKeyChecking=no nanocloud-backend.tar bamboo@releases.nanocloud.com:/data/artifacts/community/builds/${BRANCH_NAME}/'
  }

  sh 'rm -rf nanocloud-backend.tar'
}
