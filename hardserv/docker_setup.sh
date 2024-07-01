# GET Container IP:
#   docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}'

# Remove Mass Containers:
#   docker rm -v $(docker ps --filter status=created -q)

# Remove all Images:
#   sudo docker rmi -f $(docker images -q)
