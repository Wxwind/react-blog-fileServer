#!/bin/bash
## author:wxwind

imageName=rb-fileserver
containerName=rbfileServerIns

echo "stop running cantainer and delete image.."
docker stop $containerName &&
    docker rm $containerName &&
    docker rmi $imageName

echo "build image..."
docker build -t $imageName . &&
    docker run -p 7123:7123 -d --name $containerName \
        -v /soft/react_website/react_website_fileServer/assets:/app/assets $imageName