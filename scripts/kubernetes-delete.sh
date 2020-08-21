#!/bin/bash

kubectl delete -f deploy/rbac.yaml
kubectl delete -f deploy/cronjob.yaml
