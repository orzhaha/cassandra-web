node("slave") {

    stage('prepare') {
        checkout scm
    }

    stage('docker build') {
        sh 'make docker-build'
    }

    if ("${env.BRANCH_NAME}" == "master") {
        stage('docker push') {
            sh 'make docker-push'
        }
    }
}