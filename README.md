# Concurrent kubectl
A cli to simplify working with kubectl for some common workflows

#### Usage
`ckube` lets you think in terms of [services](https://kubernetes.io/docs/concepts/services-networking/service/) instead of [pods](https://kubernetes.io/docs/concepts/workloads/pods/pod/) (mostly).


```$xslt
$ ckube

A CLI to simplify working with kubectl.

Usage:
  ckube [command]

Available Commands:
  exec        execute a command in a container
  help        Help about any command
  logs        get logs from a service
  ls          list pods in kubernetes
  top         View cpu and memory usage for pods

Flags:
      --context string     the kubernetes context (defaults to value currently used by kubectl)
  -h, --help               help for ckube
  -n, --namespace string   the kubernetes namespace (defaults to value currently used by kubectl)

Use "ckube [command] --help" for more information about a command.

```

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