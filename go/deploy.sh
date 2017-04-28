#!/bin/bash

. app_env.sh

image_name="$HUB_REPO/$APP_NAME:$APP_BRANCH"
echo "Deploying $image_name"

docker push $image_name

ssh lienmeat@ericslien.com "cd ~/iot/go && ./run.sh"
