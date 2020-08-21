Apply in this order:

**NOTE** since this is being used for a homelab, we assign `cluster-admin` permissions to the service account used to run the scheduled job pods.

If this is not ideal for your scenario, you will need to configure your own ClusterRole ***inside*** `rbac.yaml`!

1. rbac.yaml
2. cronjob.yaml
