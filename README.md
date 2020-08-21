This microservice finds namespaces in terminating state and removes them.

- [Intro](#intro)
- [Why?](#why?)
- [Install](#installation)
- [Details](#Details)
- [Example](/deploy/cronjob.yaml)

## INTRO

This microservice automates the steps outlined [in this article](https://medium.com/@craignewtondev/how-to-fix-kubernetes-namespace-deleting-stuck-in-terminating-state-5ed75792647e)

## WHY?

Sometimes namespaces get stuck in terminating state.  I got sick of following the steps outlined in the [article above](#intro), so I wrote this.

## INSTALLATION

See [demo deployment](/deploy)

## DETAILS

 - This microservice is designed to run as a `CronJob`
 - The way we interact with Kubernetes means this microservice should run *inside* the cluster (and not external to the cluster like you can do using `kubectl`)
   - This means you will need to create a `ServiceAccount` with proper RBAC
   - **EXAMPLE DEPLOYMENT CAN BE [FOUND HERE](/deploy)**
