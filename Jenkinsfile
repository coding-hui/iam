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
      when {
        environment name: 'BUILD', value: 'true'
      }
      steps {
        sh 'mvn clean install package -DskipTests'
      }
    }

    stage('构建镜像') {
      when {
        environment name: 'BUILD', value: 'true'
      }
      steps {
        dir('server') {
          sh "docker build -t ${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_VERSION} ."
        }

      }
    }

    stage('推送镜像') {
      when {
        environment name: 'PUSH', value: 'true'
      }
      steps {
        useCustomStepPlugin(key: 'SYSTEM:artifact_docker_push', version: 'latest', params: [properties:'[]',image:'${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_VERSION}',repo:'docker-repo'])
      }
    }

    stage('部署') {
      when {
        environment name: 'DEPLOY', value: 'true'
      }
      steps {
        cdDeploy(deployType: 'PATCH_IMAGE', application: '${CCI_CURRENT_TEAM}', pipelineName: '${PROJECT_NAME}-${CCI_JOB_NAME}-2959888', image: 'devops-wecoding-docker.pkg.coding.net/wecoding/docker-repo/${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_VERSION}', cloudAccountName: 'wecoding-k8s', namespace: 'wecoding-system', manifestType: 'Deployment', manifestName: 'wecoding-iam', containerName: 'wecoding-iam', credentialId: '0d9cbd34fd9e4112b2712bb8a0ad84a7', personalAccessToken: '${CD_PERSONAL_ACCESS_TOKEN}')
      }
    }

  }
}