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

  try {
    sh 'docker-compose -f modules/docker-compose-build.yml run --rm nanocloud-backend make tests '
    withCredentials([[$class: 'StringBinding', credentialsId: 'GITHUB_PASSWORD', variable: 'TOKEN']]) {
      sh 'env GITHUB_PASSWORD=$TOKEN ./scripts/notify-github.sh success "Backend tests" "go tests" "http://jenkins.nanocloud.com:8080/job/Nanocloud%20Community/job/Nanocloud%20Community/job/$BRANCH_NAME/$BUILD_NUMBER/" "Nanocloud/community" $(git rev-list --parents -n 1 $(git rev-parse HEAD) | cut -f 3 -d " ")'
    }
  } catch (all) {
    withCredentials([[$class: 'StringBinding', credentialsId: 'GITHUB_PASSWORD', variable: 'TOKEN']]) {
      sh 'env GITHUB_PASSWORD=$TOKEN ./scripts/notify-github.sh failure "Backend tests" "go tests" "http://jenkins.nanocloud.com:8080/job/Nanocloud%20Community/job/Nanocloud%20Community/job/$BRANCH_NAME/$BUILD_NUMBER/" "Nanocloud/community" $(git rev-list --parents -n 1 $(git rev-parse HEAD) | cut -f 3 -d " ")'
    }

  try {
    sh 'docker-compose -f modules/docker-compose-build.yml run --rm nanocloud-frontend ember test'
    withCredentials([[$class: 'StringBinding', credentialsId: 'GITHUB_PASSWORD', variable: 'TOKEN']]) {
      sh 'env GITHUB_PASSWORD=$TOKEN ./scripts/notify-github.sh success "Frontend tests" "ember tests" "http://jenkins.nanocloud.com:8080/job/Nanocloud%20Community/job/Nanocloud%20Community/job/$BRANCH_NAME/$BUILD_NUMBER/" "Nanocloud/community" $(git rev-list --parents -n 1 $(git rev-parse HEAD) | cut -f 3 -d " ")'
    }
  } catch (all) {
    withCredentials([[$class: 'StringBinding', credentialsId: 'GITHUB_PASSWORD', variable: 'TOKEN']]) {
      sh 'env GITHUB_PASSWORD=$TOKEN ./scripts/notify-github.sh failure "Frontend tests" "ember tests" "http://jenkins.nanocloud.com:8080/job/Nanocloud%20Community/job/Nanocloud%20Community/job/$BRANCH_NAME/$BUILD_NUMBER/" "Nanocloud/community" $(git rev-list --parents -n 1 $(git rev-parse HEAD) | cut -f 3 -d " ")'
    }

  sh "docker save nanocloud/nanocloud-backend \$(docker history -q nanocloud/nanocloud-backend | tail -n +2 | grep -v \\<missing\\> | tr '\n' ' ') > nanocloud-backend.tar"
  sh "docker save nanocloud/nanocloud-frontend \$(docker history -q nanocloud/nanocloud-frontend | tail -n +2 | grep -v \\<missing\\> | tr '\n' ' ') > nanocloud-frontend.tar"

  sshagent(['releases']) {
    sh 'scp -o StrictHostKeyChecking=no nanocloud-backend.tar bamboo@releases.nanocloud.com:/data/artifacts/community/builds/${BRANCH_NAME}/'
    sh 'scp -o StrictHostKeyChecking=no nanocloud-frontend.tar bamboo@releases.nanocloud.com:/data/artifacts/community/builds/${BRANCH_NAME}/'
  }

  sh 'rm -rf nanocloud-backend.tar'
  sh 'rm -rf nanocloud-frontend.tar'
}
