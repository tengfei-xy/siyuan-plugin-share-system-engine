#!/bin/bash
sudo docker-compose down spss_engine  && \
sudo docker images | grep spss_engine | awk '{print $3}' | xargs -n1 sudo docker rmi && \
sudo docker-compose up -d spss_engine