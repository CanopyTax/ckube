# Concurrent kubectl
A cli to simplify working with kubectl for some common workflows

#### Usage
`ckube` lets you think in terms of [services](https://kubernetes.io/docs/concepts/services-networking/service/) instead of [pods](https://kubernetes.io/docs/concepts/workloads/pods/pod/).

Kubernetes services are often an abstraction over multiple pods, particularly if the replicas are scaled. If a k8s cluster has a service called `nginx` you could have several pods which might be named something like this:
```$xslt
nginx-3528986049-kpd4z
nginx-3528986049-71s10 
nginx-3528986049-f6mwf
nginx-3528986049-ltx6j
nginx-3528986049-m3cmm
nginx-3528986049-h8cnn
nginx-3528986049-6v4c1
```

Getting logs for the nginx service is easy with `ckube`
For a single unspecified nginx pod:
```
# single uspecified nginx pod:
ckube logs nginx

# all nginx pods
ckube logs nginx

# follow the logs
ckube logs nginx -f
```

Similar concurrent functionality exists for `exec`

This readme could really use some gifs of ckube in action