node {
    checkout scm
        
    stage('Docker Build') {
        docker.build('dukfaar/levebackend')
    }

    stage('Update Service') {
        sh 'docker service update --force levebackend_levebackend'
    }
}
