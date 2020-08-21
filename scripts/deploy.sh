#!/bin/bash

#
# Use like:
# ./this-script.sh --container_name yourcontainername

container_name=${container_name:-}

printf '%s\n' "[NOTE] We assume you have a 'Dockerfile' at the root of this project"

# This allows us to use named params
while [ $# -gt 0 ]; do
    if [[ $1 == *"--"* ]]; then
        param="${1/--/}"
        declare $param="$2"
    fi
    shift
done

if [ -z "$container_name" ]; then
    printf '%s\n' "[ERROR] '--container_name' parameter not found! Please supply a container name like: ./this-script.sh --container_name yourcontainername"
else
    ./docker-build-push.sh --container_name $container_name
    ./kubernetes-apply.sh --container_name $container_name
fi
