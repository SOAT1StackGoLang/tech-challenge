#!/bin/bash
# This script is used to update the git repository for each folder starting with msvc
# It will pull the latest changes from the develop branch

# Get the current directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
# Get the list of folders starting with msvc
FOLDERS=$(ls -d $DIR/msvc*/)
# For each folder, pull the latest changes from the develop branch
for folder in $FOLDERS
do
    echo "Updating $folder"
    cd $folder
    git pull origin develop
    cd $DIR
done