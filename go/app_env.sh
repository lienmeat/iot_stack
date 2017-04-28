#!/bin/bash
#This file is used in the build, and deploy scripts so we don't have
#to customize them for every single app

#Package/app to compile
export APP_NAME="iot.ericslien.com"

#Docker hub repo to push to
export HUB_REPO="ericlien"

#Tag used, which should probably mirror code branch
export APP_BRANCH="master"
