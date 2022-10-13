pipeline {
  agent any
  stages {
    stage('检出') {
      steps {
        checkout([$class: 'GitSCM',
        branches: [[name: GIT_BUILD_REF]],
        userRemoteConfigs: [[
          url: GIT_REPO_URL,
          credentialsId: CREDENTIALS_ID
        ]]])
      }
    }

    stage('编译') {
      steps {
        sh 'mvn clean install package -DskipTests'
      }
    }

    stage('推送') {
      steps {
        dir('server') {
          sh "docker build -t ${CODING_DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_VERSION} -f ${DOCKERFILE_PATH} ${DOCKER_BUILD_CONTEXT}"
        }

        useCustomStepPlugin(key: 'SYSTEM:artifact_docker_push', version: 'latest', params: [properties:'[]',image:'${CODING_DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_VERSION}',host:'docker.io',project:'wecoding',repo:'iam-server',username:'${PROJECT_TOKEN_GK}',password:'${PROJECT_TOKEN}'])
      }
    }

    stage('部署') {
      steps {
        cdDeploy(deployType: 'PATCH_IMAGE', application: '${CCI_CURRENT_TEAM}', pipelineName: '${PROJECT_NAME}-${CCI_JOB_NAME}-2959888', image: '${CODING_DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_VERSION}', cloudAccountName: 'wecoding-k8s', namespace: 'wecoding-system', manifestType: 'Deployment', manifestName: 'wecoding-iam', containerName: 'wecoding-iam', credentialId: '16c6dc5732f84db1b8c6dfac219dae2b', personalAccessToken: '${CD_PERSONAL_ACCESS_TOKEN}')
      }
    }

  }
  environment {
    CODING_DOCKER_REG_HOST = "${CCI_CURRENT_TEAM}-docker.pkg.${CCI_CURRENT_DOMAIN}"
    CODING_DOCKER_IMAGE_NAME = "${PROJECT_NAME.toLowerCase()}/${DOCKER_REPO_NAME}/${DOCKER_IMAGE_NAME}"
  }
}