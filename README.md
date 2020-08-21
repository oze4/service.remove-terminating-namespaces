This microservice finds namespaces in terminating state and removes them.

- [Intro](#intro)
- [Why?](#why)
- [Usage](#usage)
- [Details](#Details)
- [Example](/deploy/cronjob.yaml)

## INTRO

This microservice automates the steps outlined [in this article](https://medium.com/@craignewtondev/how-to-fix-kubernetes-namespace-deleting-stuck-in-terminating-state-5ed75792647e)

## WHY?

Sometimes namespaces get stuck in terminating state.  I got sick of following the steps outlined in the [article above](#intro), so I wrote this.

## USAGE

*We assume all scripts are being invoked from the root of this project

 - If you would like to use this microservice in your own K8 environment AND store the container in your own DockerHub:
   - `./scripts/deploy.sh --container_name yourcontainername`
 - To only build container and push to your own DockerHub account:
   - `./scripts/docker-build-push.sh --container_name yourcontainername`
 - To only apply .yaml files:
   - **PREFERRED:** (*should be ran at least once, in order to generae correct `deploy/cronjob.yaml` file*)
     - `./scritps/kubernetes-apply.sh --container_name yourcontainername`
   - **ONLY** apply existing .yaml files
     - `./scritps/kubernetes-apply.sh`
 - To delete containers:
   - `./scritps/kubernetes-delete.sh`

## DETAILS

 - This microservice is designed to run as a `CronJob`
 - The way we interact with Kubernetes means this microservice should run *inside* the cluster (and not external to the cluster like you can do using `kubectl`)
   - This means you will need to create a `ServiceAccount` with proper RBAC
   - **EXAMPLE DEPLOYMENT CAN BE [FOUND HERE](/deploy)**
