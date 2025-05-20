#!/bin/sh
TEST_TAG=$(date +"%Y%m%d%H%M")

# build for testing
docker build -t $DOCKER_HUB_HOSTED_HOST/navi-radar-parser-cat240:latest . || { echo 'docker build failed' ; exit 1; }
docker tag $DOCKER_HUB_HOSTED_HOST/navi-radar-parser-cat240:latest $DOCKER_HUB_HOSTED_HOST/navi-radar-parser-cat240:$TEST_TAG || { echo 'docker tag failed' ; exit 1; }
echo "$DOCKER_HUB_HOSTED_PASSWORD" | docker login -u "$DOCKER_HUB_HOSTED_USER" --password-stdin $DOCKER_HUB_HOSTED_HOST || { echo 'docker login failed' ; exit 1; }
docker push $DOCKER_HUB_HOSTED_HOST/navi-radar-parser-cat240:latest || { echo 'docker push failed' ; exit 1; }
docker push $DOCKER_HUB_HOSTED_HOST/navi-radar-parser-cat240:$TEST_TAG || { echo 'docker push failed' ; exit 1; }