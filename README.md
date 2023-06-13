# pod-test-k8s
A golang web project to test as pod on k8s

### Creating image
docker build --tag <org>/decode-jwt .
docker image tag <org>/decode-jwt:latest <org>/decode-jwt:0.4
docker push <org>/decode-jwt:0.4
