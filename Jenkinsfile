pipeline {
    agent any
    
    environment {
        DOCKER_IMAGE = "user-service"
        COMPOSE_PROJECT_NAME = "user-service"
    }
    
    stages {
        stage('Cleanup') {
            steps {
                script {
                    sh '''
                        docker-compose down || true
                        docker rm -f user-service-app || true
                    '''
                }
            }
        }
        
        stage('Build Docker Image') {
            steps {
                script {
                    sh '''
                        docker build -t ${DOCKER_IMAGE}:latest .
                        docker images | grep ${DOCKER_IMAGE}
                    '''
                }
            }
        }
        
        stage('Start Dependencies') {
            steps {
                script {
                    sh '''
                        docker-compose up -d postgres redis
                        echo "Waiting for databases to start..."
                        sleep 15
                    '''
                }
            }
        }
        
        stage('Run Application') {
            steps {
                script {
                    sh '''
                        docker run -d \
                          --name user-service-app \
                          --network ${COMPOSE_PROJECT_NAME}_default \
                          -p 8081:8081 \
                          -e CONFIG_PATH=/app/config/local-cfg.yaml \
                          ${DOCKER_IMAGE}:latest
                        
                        echo "Waiting for application to start..."
                        sleep 10
                    '''
                }
            }
        }
        
        stage('Health Check') {
            steps {
                script {
                    sh '''
                        echo "Checking application health..."
                        docker ps | grep user-service-app
                        docker logs user-service-app
                    '''
                }
            }
        }
    }
    
    post {
        success {
            echo 'Build completed successfully!'
            echo 'Application is running on http://localhost:8081'
        }
        failure {
            echo 'Build failed!'
            sh 'docker logs user-service-app || true'
        }
        always {
            echo 'Cleaning up...'
        }
    }
}