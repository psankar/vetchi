secret_settings(disable_scrub=True)

# Load Kubernetes YAML files
k8s_yaml('api/hermione-tilt.yaml')
k8s_yaml('api/granger-tilt.yaml')
k8s_yaml('sqitch/sqitch-tilt.yaml')

# Define Docker builds with root context to include typespec
docker_build('psankar/granger', '.',
    dockerfile='api/Dockerfile-granger',
)

docker_build('psankar/hermione', '.',
    dockerfile='api/Dockerfile-hermione',
)

docker_build('psankar/vetchi-sqitch', 'sqitch', dockerfile='sqitch/Dockerfile')

# Associate images with Kubernetes resources
k8s_resource('granger', port_forwards='8080:8080')
k8s_resource('hermione', port_forwards='8081:8080')
k8s_resource('sqitch')

