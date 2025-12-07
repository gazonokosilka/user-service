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
                        docker-compose down -v || true
                        docker rm -f user-service-app || true
                        docker volume rm user-service_user_service_postgres_data || true
                        docker volume rm user-service_redis_data || true
                        docker network prune -f || true
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
                        sleep 20
                        
                        # Проверка что PostgreSQL готов
                        docker exec user-service-db pg_isready -U postgres || echo "PostgreSQL not ready yet"
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
                          --network user-service_default \
                          -p 8081:8081 \
                          -e POSTGRES_HOST=user-service-db \
                          -e POSTGRES_PORT=5432 \
                          -e POSTGRES_USER=postgres \
                          -e POSTGRES_PASSWORD=postgres \
                          -e POSTGRES_DB=user_service \
                          -e REDIS_ADDRESS=user-service_redis:6379 \
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
                        
                        echo "Application logs:"
                        docker logs user-service-app
                        
                        echo "Testing database connection from app container:"
                        docker exec user-service-app nc -zv user-service-db 5432 || true
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
            sh '''
                echo "Application logs:"
                docker logs user-service-app || true
                
                echo "Database logs:"
                docker logs user-service-db || true
                
                echo "Redis logs:"
                docker logs user-service_redis || true
                
                echo "Network info:"
                docker network inspect user-service_default || true
            '''
        }
        always {
            echo 'Cleaning up...'
        }
    }
}