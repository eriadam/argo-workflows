# Multi-cluster

Argo allows you to run workflows where the tasks that make up the workflow run in a different cluster and namespace to
the workflow.

## Kubeconfig

Argo can only create pods in clusters in can connect to - ones it has a `kubeconfig` for. These must be installed in the
system namespace (typcially `argo`):

```bash
kubectl get secret -l workflows.argoproj.io/cluster
```

You can add a new secret with a single `kubeconfig` field:

```bash
kubectl create secret generic other-cluster "--from-literal=kubeconfig=`kubectl config view --context=other --minify --raw -o json`"
kubectl label secret other-cluster workflows.argoproj.io/cluster=other
```

## Labels

It is worthwhile understand how Argo uses labels. Some facts:

* It is not possible to create an ownership reference between resources in different namespaces or clusters.
* It is possible for two different Argos to create pods in the same namespace that belong to different workflows.

So this creates problems:

* How do I make sure pods are deleted if the workflow is deleted?
* How do I know which pod belongs to which workflow?

This is solved using labels:

* `workflows.argoproj.io/cluster` tells you which the cluster of the parent workflow.
* `workflows.argoproj.io/workflow-namespace` tells you which the namespace of the parent workflows.

These labels are only applied if an ownership reference cannot be created, i.e. if if the pod is created in different
cluster or namespace to the workflow.

## Pod Garbage Collection

If a pod is created in another cluster, and the parent workflow is deleted, then Argo must garbage collect it. Normally,
Kubernetes would do this.

⚠️ This garbage collection is done on best effort, and that might be long time after the workflow is deleted. To mitigate
this, use `podGCStrategy`.

