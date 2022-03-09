# HOWTO

Disclaimer: Do not do this unless you know how to use kubernetes!

### Run
```shell
kubectl apply -f warship.yaml -n default
kubectl get deployments && kubectl scale deployment ukrainian-warship --replicas 10
```

### Stop
```shell
kubectl get deployments && kubectl delete deployments ukrainian-warship
```

*Thanks too VI D*