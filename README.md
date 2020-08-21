This microservice finds namespaces in terminating state and removes them.

- [How to resolve the root of this issue](#resolve-root-issue)
- [Intro](#intro)
  - [Important Note](#important-note)
- [Why?](#why)
- [Usage](#usage)
  - [Bash](#bash)
  - [Makefile](#makefile)
- [Details](#details)
- [Example](/deploy)

# RESOLVE ROOT ISSUE

[Taken from here](https://medium.com/@cristi.posoiu/this-is-not-the-right-way-especially-in-a-production-environment-190ff670bc62)

>This is not the right way, especially in a production environment.
>
>Today I got into the same problem. By removing the finalizer you’ll end up with leftovers in various states. You should actually find what is keeping the deletion from complete.
>
>See https://github.com/kubernetes/kubernetes/issues/60807#issuecomment-524772920
>
>(also, unfortunately, ‘kubetctl get all’ does not report all things, you need to use similar commands like in the link)
>
>My case — deleting ‘cert-manager’ namespace. In the output of ‘kubectl get apiservice -o yaml’ I found APIService ‘v1beta1.admission.certmanager.k8s.io’ with status=False . This apiservice was part of cert-manager, which I just deleted. So, in 10 seconds after I ‘kubectl delete apiservice v1beta1.admission.certmanager.k8s.io’ , the namespace disappeared.
>
>Hope that helps.

#### How I found the offending ApiService

- CTRL+F in the output .yaml file that was generated
  - Use this search (copy and paste it exactly)
    - `status: "False"`
- Find the name of that ApiService
- Follow the rest of the commands above, but using that ApiService name
- Do this for all offending ApiService (any service that has `status: "False"`

## INTRO

This microservice automates the steps outlined [in this article](https://medium.com/@craignewtondev/how-to-fix-kubernetes-namespace-deleting-stuck-in-terminating-state-5ed75792647e)

### IMPORTANT NOTE

Due to the fact this code was written for a home-lab, we assign `cluster-admin` permissions to the service account used to run the scheduled job pods. If this is not ideal for your scenario, you will need to configure your own ClusterRole ***inside*** of [rbac.yaml](/deploy/rbac.yaml)!

## WHY?

Sometimes namespaces get stuck in terminating state.  I got sick of following the steps outlined in the [article above](#intro), so I wrote this.

## USAGE

*We assume all scripts are being invoked from the root of this project

### Makefile

 - PREFERRED OVER USING BASH
 - See the `Makefile.example` file at the root of this project
   - Change the variable on line 1 to your container name: 
     - `CONTAINER_NAME = yourcontainername`
   - Rename from `Makefile.example` to just `Makefile`
 - You now have access to the following commands:
   - `make deploy`
     - Same as running: `./scripts/deploy.sh --container_name yourcontainername` where `yourcontainername` is what you set the `CONTAINER_NAME` variable to be in the `Makefile`
   - `make docker`
     - Same as running: `./scripts/docker-build-push.sh --container_name yourcontainername` where `yourcontainername` is what you set the `CONTAINER_NAME` variable to be in the `Makefile`
   - `make kubernetes-apply`
     - Same as running: `./scritps/kubernetes-apply.sh --container_name yourcontainername` where `yourcontainername` is what you set the `CONTAINER_NAME` variable to be in the `Makefile`
   - `make kubernetes-delete`
     - Same as running: `./scritps/kubernetes-delete.sh`

### Bash

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
