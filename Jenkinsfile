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
    sh 'scp -o StrictHostKeyChecking=no bamboo@releases.nanocloud.com:/data/artifacts/community/canary/nanocloud-frontend.tar nanocloud-frontend.tar'
  }

  sh 'docker load -i nanocloud-backend.tar'
  sh 'docker load -i nanocloud-frontend.tar'

  sh 'docker-compose -f modules/docker-compose-build.yml build nanocloud-backend'
  sh 'docker-compose -f modules/docker-compose-build.yml build nanocloud-frontend'
  sh "docker save nanocloud/nanocloud-backend \$(docker history -q nanocloud/nanocloud-backend | tail -n +2 | grep -v \\<missing\\> | tr '\n' ' ') > nanocloud-backend.tar"
  sh "docker save nanocloud/nanocloud-frontend \$(docker history -q nanocloud/nanocloud-frontend | tail -n +2 | grep -v \\<missing\\> | tr '\n' ' ') > nanocloud-frontend.tar"

  sshagent(['releases']) {
    sh 'scp -o StrictHostKeyChecking=no nanocloud-backend.tar bamboo@releases.nanocloud.com:/data/artifacts/community/builds/${BRANCH_NAME}/'
    sh 'scp -o StrictHostKeyChecking=no nanocloud-frontend.tar bamboo@releases.nanocloud.com:/data/artifacts/community/builds/${BRANCH_NAME}/'
  }

  sh 'rm -rf nanocloud-backend.tar'
  sh 'rm -rf nanocloud-frontend.tar'
}
