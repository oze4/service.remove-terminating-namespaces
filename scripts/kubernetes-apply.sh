#!/bin/bash

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
# Use like:                                                           #
# ./this-script.sh                                                    #
# Or like this, to write a new cronjob.yaml:                          #
# ./this-script.sh --container_name yourcontainername                 #
# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

# This allows us to use named params
container_name=${container_name:-}
while [ $# -gt 0 ]; do
    if [[ $1 == *"--"* ]]; then
        param="${1/--/}"
        declare $param="$2"
    fi
    shift
done

# If --container_name param was NOT used, just apply existing files
if [ -z "$container_name" ]; then 
    # Do nothing, we apply at the very end
    return
else # If --container_name param WAS used, write new cronjob.yaml then apply
    printf '\n'
    printf '%s\n' "~~ Writing new CronJob file to: ./deploy/cronjob.yaml ~~"
    # This creates the cronjob yaml based upon the container name
    echo \
"apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: remove-terminating-namespaces-cronjob
spec:
  schedule: \"0 */1 * * *\" # at minute 0 of each hour aka once per hour at the top of the hour (h:00 where h is the current hour)
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

# Apply yaml files
kubectl apply -f deploy/rbac.yaml
kubectl apply -f deploy/cronjob.yaml
