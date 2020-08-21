#!/bin/bash

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
# Use like:                                                           #
# ./this-script.sh --container_name yourcontainername                 #
# Or like this, to write a new cronjob.yaml:                          #
# ./this-script.sh --container_name yourcontainername --write_cronjob #
# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

# This allows us to use named params
container_name=${container_name:-}
write_cronjob=${write_cronjob:-false}
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
    if [ -z "$write_cronjob" ]; then
        printf '\n'
        printf '%s\n' "~~ Writing new CronJob file to: ./deploy/cronjob.yaml ~~"
        # This creates the cronjob yaml based upon the container name
        echo \
"apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: remove-terminating-namespaces-cronjob
spec:
  schedule: \"0 */1 * * *\" # at minute 0 of each hour aka once per hour
  #successfulJobsHistoryLimit: 0
  #failedJobsHistoryLimit: 0
  jobTemplate:
    spec:
      template:
        spec:
          serviceAccountName: svc-remove-terminating-namespaces
          containers:
          - name: remove-terminating-namespaces
            image: $container_name
          restartPolicy: OnFailure" >./deploy/cronjob.yaml
    fi

    kubectl apply -f deploy/rbac.yaml
    kubectl apply -f deploy/cronjob.yaml
fi
