#!/bin/bash 

image_name="quant_api"

docker stop ${image_name} 
docker rm ${image_name} 
docker run --name ${image_name} -d ${image_name}
