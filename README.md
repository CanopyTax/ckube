# Concurrent kubectl
![Kelsey Hightower Approved!](https://img.shields.io/badge/Hightower-approved-blue.svg) [![Build Status](https://travis-ci.org/CanopyTax/ckube.svg?branch=master)](https://travis-ci.org/CanopyTax/ckube) [![Go Report Card](https://goreportcard.com/badge/github.com/CanopyTax/ckube)](https://goreportcard.com/report/github.com/CanopyTax/ckube) [![GitHub release](https://img.shields.io/github/release/CanopyTax/ckube.svg)](https://github.com/CanopyTax/ckube/releases/latest) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A cli to simplify working with kubectl for some common workflows

#### Installation
Find the latest binaries [here](https://github.com/canopytax/ckube/releases/) or use `homebrew` on macOS
```
brew tap canopytax/ckube
brew install ckube
```

#### Usage
```$xslt
$ ckube
A CLI to simplify working with kubectl.

Usage:
  ckube [command]

Available Commands:
  deployment  Get information about deployments
  exec        execute a command in a container
  help        Help about any command
  logs        get logs from a service
  ls          Interactive list of pods
  nodes       Lists pods grouped by node
  service     Interactive view of your services
  top         View cpu and memory usage for pods

Flags:
      --context string      the kubernetes context (defaults to value currently used by kubectl)
  -h, --help                help for ckube
      --kubeconfig string   path to kubeconfig file to use for CLI requests (defaults to $KUBECONFIG or $HOME/.kube/config)
  -l, --labels string       Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)
  -n, --namespace string    the kubernetes namespace (defaults to value currently used by kubectl)

Use "ckube [command] --help" for more information about a command.

```

Targeting pods can be burdensome (i.e. in a deployment with replicas). If a k8s cluster has a service called `nginx` you could have several pods which might be named something like this:
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
ckube logs nginx -a

# follow the logs
ckube logs nginx -f
```
![](https://github.com/canopytax/ckube/blob/master/images/logs.gif?raw=true)

Similar concurrent functionality exists for `exec`
![](https://github.com/canopytax/ckube/blob/master/images/exec.gif?raw=true)

![](https://github.com/canopytax/ckube/blob/master/images/complex-exec.gif?raw=true)

You can toggle through your services and see the associated pods. `ckube service`
![](https://github.com/canopytax/ckube/blob/master/images/services.gif?raw=true)

You can use `ckube nodes` to show which pods are running on each node
[![asciicast](https://asciinema.org/a/150564.png)](https://asciinema.org/a/150564)

Get an interactive view of your pods with `ckube ls`. You can further refine the prompt results by searching using `/`. Selecting a pod will print detailed pod information returned by `kubectl describe`
![](https://github.com/canopytax/ckube/blob/master/images/list-interactive.gif?raw=true)

Read the blog post about `ckube` [here](https://devonmoss.com/concurrent-kubectl)

#### Contributing
PR's accepted

If you are looking to build the project or would like to contribute:
```$xslt
go get -d github.com/canopytax/ckube
cd $GOPATH/src/github.com/canopytax/ckube
```
