---
title: "Getting Down with the Kubernetes Krew"
panelTitle: "Getting Down with the Kubernetes Krew"
date: 2020-02-16T18:00:00-01:00
author: "@alexdmoss"
description: "I've picked out some kubectl plugins I've found handy, which I discovered through the Kubernetes plugin manager krew"
banner: "/images/boat-crew.jpg"
tags: [ "Kubernetes", "krew" ]
categories: [ "Kubernetes", "Tips" ]
draft: true
---

/// TODO scrub mw --> sandpit

/// doctor                          Scans your cluster and reports anomalies.
/// eksporter                       Export resources and removes a pre-defined set ...
/// evict-pod                       Evicts the given pod
/// kudo                            Declaratively build, install, and run operators...
/// match-name                      Match names of pods and other API objects

In this post, after a brief intro on what *krew* actually **is**, I run through a dozen or so plugins installed using it that I find to be pretty handy time-savers. Hope you find at least one of them useful!

{{< figure src="/images/boat-crew.jpg?width=600px&classes=shadow" attr="Photo by Quino Al on Unsplash" attrlink="https://unsplash.com/photos/5gdEeq_sWCU" >}}

---

So, long time no post. I was resolved to posting more frequently in 2020 ... and we're halfway through February ... so not so successfully then! Time to try to turn the tide ...

On this occasion, I'll be talking about [krew](https://github.com/kubernetes-sigs/krew). The focus of the post is some of the plugins that work for me, but there is a brief intro to what `krew` actually is if you've not heard about it before.

> Notice how I managed to get this far *without* any mention of "top 10 krew plugins to help you devop harder". That is a very, very deliberate avoidance of a clickbait-y title. Will be interesting to see whether the article still gets picked up on Medium!

## Intro to Krew

In short, it's a plugin manager for `kubectl`. Think `brew` or `apt`, but for your kubectl commands. Plugins have been a thing in Kubernetes since v1.12, and if you want to know a bit more about how they work then [Ahmet Alp Balkan](https://twitter.com/ahmetb) (creator of krew, as well as indispensible tools like [kubectx & kubens](https://github.com/ahmetb/kubectx/)) covers it really well [here](https://ahmet.im/blog/kubectl-plugins/).

Installation instructions can be found [here](https://github.com/kubernetes-sigs/krew). Basic usage is as you'd expect (I'm lazy, so I alias `krew='kubectl krew'`):

- `krew list` - list your plugins installed through krew. If you only use krew to do this, this should be the same as `kubectl plugin list`
- `krew search <string>` - find plugins. Without `<string>`, all are listed
- `krew info <plugin>` - more info about what a plugin does. I always check this to make sure a plugin isn't going to do something dubious - it invariably contains a link to the source code. And lets be honest, you're basically running someone else's bash on your PC with these plugins
- `krew install <plugin>` - installs it - it can then be used through `kubectl <plugin>`
- `krew upgrade` - upgrade your plugins

> At time of writing, there are now ~70 plugins available through it. I feel like when I first started looking at krew in mid-2020, nthere were only a couple of dozen, so its use is clearly accelerating!

---

Below are some of the ones I've found useful time-savers. Hope you find something in here that appeals! They are loosely grouped by ...

- [Intro to Krew](#intro-to-krew)
- [Viewing Resources](#viewing-resources)
- [Resource Usage](#resource-usage)
- [Housekeeping](#housekeeping)
- [Security](#security)
- [Debugging](#debugging)
- [And finally ...](#and-finally)

---

## Viewing Resources

First up in this section we have `kubectl tail` (usually - `kail`). It's a handy plugin for tailing logs - it's just nice and simple. There are a number of alternatives out there (I know some folks like [`stern`](https://github.com/wercker/stern), which I've never got round to trying). The main feature I like is being able to target a deployment/service/ingress (e.g. `kubectl tail --ing=mosstech`) and be able to watch all the logs from the pods behind it, without necessarily needing to know how the resource has been labelled.

`kubectl get-all` [does exactly what it says on the tin](https://en.wikipedia.org/wiki/Does_exactly_what_it_says_on_the_tin). They are not kidding - this *really* gets everything, as the plugin info page says. I've found this to be increasingly useful with the uptake of `CustomResourceDefinition` - i.e. where you may not be able to remember every type of object that exists in a namespace any more (usually because I want to delete stuff ...). This plugin helps a lot with that:

```sh
[~ (⎈ |sandpit-prod:prometheus)]$ kubectl get-all -n=prometheus
NAME                                                                                        NAMESPACE   AGE
configmap/prometheus-sandpit-rulefiles-0                                                    prometheus  14d
endpoints/prometheus-operated                                                               prometheus  20d
endpoints/prometheus-operator                                                               prometheus  20d
persistentvolumeclaim/prometheus-storage-prometheus-sandpit-0                               prometheus  20d
# [...] snip for brevity
prometheus.monitoring.coreos.com/sandpit                                                    prometheus  20d
prometheusrule.monitoring.coreos.com/prometheus-k8s-rules                                   prometheus  14d
servicemonitor.monitoring.coreos.com/prometheus                                             prometheus  16d
servicemonitor.monitoring.coreos.com/prometheus-operator                                    prometheus  16d
poddisruptionbudget.policy/prometheus                                                       prometheus  20d
```

Almost the opposite is `kubectl neat` which, as the name sort-of-implies, neatens up verbose output. If you find your mind glazing over the system-injected annotations and such when get you get a pod's details, this one might be for you - it strips the gumpf from `kubectl get` - for example:

![kubectl neat](/images/krew-neat.png)

Lastly for this section, are a couple of plugins that help visualise the relationship between a resource and other things in the cluster. `kubectl pod-dive` is good if you want to know what surrounds a particular pod, e.g.

```sh
[~ (⎈ |sandpit-prod:prometheus)]$ k pod-dive prometheus-sandpit-0
[node]      gke-sandpit-prod-np-0-947e5a45-84zk [ready]
[namespace]  ├─┬ prometheus
[type]       │ └─┬ statefulset
[workload]   │   └─┬ prometheus-sandpit [2 replicas]
[pod]        │     └─┬ prometheus-sandpit-0 [running]
[containers] │       ├── prometheus [228 restarts]
             │       ├── prometheus-config-reloader [0 restarts]
             │       └── rules-configmap-reloader [0 restarts]
            ...
[siblings]   ├── grafana-0
             ├── nginx-ingress-controller-85744dcf89-nswr9
             ├── fluentd-gcp-v3.1.1-bzh54
             ├── heapster-v1.6.1-5b6bf6cc74-cnnlh
             ├── kube-proxy-gke-sandpit-prod-np-0-947e5a45-84zk
             ├── prometheus-to-sd-cjkwx
             ├── node-exporter-7ml5r

WAITING:
    prometheus crashloopbackoff (Back-off 5m0s restarting failed container=prometheus pod=prometheus-sandpit-0_prometheus(4de6fa0b-5007-11ea-ac4f-42010af0028a))
TERMINATION:
    prometheus error (code 1)
```

Somewhat similarly, `kubectl tree` shows the hierarchy for a particular resource - especially useful when working with `CustomResourceDefinition`s. For example, here's a view from my failing Prometheus pods (installed using the `prometheus` CoreOS operator):

```sh
[~ (⎈ |sandpit-prod:prometheus)]$ kubectl tree prometheus sandpit
NAMESPACE   NAME                                                  READY  REASON              AGE
prometheus  Prometheus/sandpit                                    -                          20d
prometheus  ├─ConfigMap/prometheus-sandpit-rulefiles-0            -                          14d
prometheus  ├─Secret/prometheus-sandpit                           -                          20d
prometheus  ├─Service/prometheus-operated                         -                          20d
prometheus  └─StatefulSet/prometheus-sandpit                      -                          20d
prometheus    ├─ControllerRevision/prometheus-sandpit-6f6fcd787b  -                          20d
prometheus    ├─Pod/prometheus-sandpit-0                          False  ContainersNotReady  29h
prometheus    └─Pod/prometheus-sandpit-1                          False  ContainersNotReady  3h
```

---

## Resource Usage

There are a whole bunch of plugins in this area - I've picked out a couple that I think do a bit more than is easily achieved with some `kubectl`-foo. Some of these were particularly helpful when trying to [squeeze more out of my tiny 'home' GKE cluster](/2019/09/28/squeezing-gke-system-resources-in-small-clusters/).

First, `kubectl resource-capacity` offers a nice and simple view for an overall cluster usage:

```sh
NODE                                  CPU REQUESTS   CPU LIMITS     MEMORY REQUESTS   MEMORY LIMITS
*                                     1136m (60%)    5795m (308%)   3147Mi (59%)      10983Mi (208%)
gke-sandpit-prod-np-0-3be8d635-6s5n   694m (73%)     3406m (362%)   1812Mi (68%)      5774Mi (218%)
gke-sandpit-prod-np-0-947e5a45-84zk   442m (47%)     2389m (254%)   1335Mi (50%)      5209Mi (197%)
```

Occasionally, info about disk utilisation is useful, and not super-obvious from other sources (like `kubectl get`) ... enter `kubectl df-pv`:

```sh
PVC                                       NAMESPACE    POD                    SIZE          USED        AVAILABLE     PERCENTUSED   IUSED   IFREE     PERCENTIUSED
prometheus-storage-prometheus-sandpit-1   prometheus   prometheus-sandpit-1   21003583488   541843456   20444962816   2.58          3729    1306991   0.28
grafana-data-grafana-1                    grafana      grafana-1              1023303680    3493888     1003032576    0.34          23      65513     0.04
prometheus-storage-prometheus-sandpit-0   prometheus   prometheus-sandpit-0   21003583488   533467136   20453339136   2.54          3744    1306976   0.29
grafana-data-grafana-0                    grafana      grafana-0              1023303680    80232448    926294016     7.84          24      65512     0.04
```

And finally, one that is handy for spotting pods without resource requests/limits (it's not obvious from the capture below, but it highlights these in red):

//// picture instead, and need to show the plugin name! Maybe do asciinema

```sh
 Resource                                                          Requested  %Requested      Limit  %Limit  Allocatable  Free
  attachable-volumes-gce-pd                                                0          0%          0      0%          254   254
  ├─ gke-sandpit-prod-np-0-3be8d635-6s5n                                   0          0%          0      0%          127   127
  └─ gke-sandpit-prod-np-0-947e5a45-84zk                                   0          0%          0      0%          127   127
  cpu                                                                  1036m         55%      4795m    255%            1     0
  ├─ gke-sandpit-prod-np-0-3be8d635-6s5n                                692m         74%      3049m    324%         940m     0
  │  ├─ cert-manager-5c5f4b9b49-rj7zw                                    10m                    20m
  │  ├─ default-backend-66b68fd9d8-n5xwl                                 10m                    20m
  │  ├─ fluentd-gcp-v3.1.1-b8hks                                         10m                   100m
  │  ├─ grafana-1                                                        10m                   100m
  │  ├─ heapster-7f9cb9f8d5-x5zp2                                        63m                    63m
  │  ├─ kube-dns-795f7b8488-2scvh                                        30m                   200m
  │  ├─ kube-dns-795f7b8488-fth52                                        30m                   200m
  │  ├─ kube-ops-view-667db9cf57-sdqlx                                   50m                   200m
  │  ├─ kube-ops-view-redis-74c9895647-ztv76                             50m                   100m
  │  ├─ kube-proxy-gke-sandpit-prod-np-0-3be8d635-6s5n                  100m                      0
  │  ├─ kube-state-metrics-569465d5cb-55mb4                              10m                   100m
  │  ├─ l7-default-backend-fd59995cd-l85xz                               10m                    10m
  │  ├─ metrics-server-v0.3.1-57c75779f-lqkqv                            48m                   143m
  │  ├─ nginx-ingress-controller-85744dcf89-thkb4                        20m                   200m
  │  ├─ node-exporter-rqb4c                                              10m                    50m
  │  ├─ prometheus-sandpit-1                                             40m                    40m
  │  ├─ prometheus-operator-6977fd7f46-hhftc                             10m                   200m
  │  ├─ prometheus-to-sd-skmd9                                            1m                     3m
  │  ├─ right-sizer-789b55558-7lkpg                                      10m                    80m
  │  ├─ stackdriver-metadata-agent-cluster-level-69884d6fd6-l5vwf        40m                      0
  └─ gke-sandpit-prod-np-0-947e5a45-84zk                                344m         37%      1746m    186%         940m     0
     ├─ fluentd-gcp-v3.1.1-bzh54                                         10m                   100m
     ├─ grafana-0                                                        10m                   100m
     ├─ heapster-v1.6.1-5b6bf6cc74-52zqd                                 23m                    33m
     ├─ kube-proxy-gke-sandpit-prod-np-0-947e5a45-84zk                  100m                      0
     ├─ nginx-ingress-controller-85744dcf89-nswr9                        20m                   200m
     ├─ node-exporter-7ml5r                                              10m                    50m
     ├─ prometheus-sandpit-0                                             40m                    40m
     ├─ prometheus-to-sd-cjkwx                                            1m                     3m
  ephemeral-storage                                                        0          0%          0      0%          16G   16G
  ├─ gke-sandpit-prod-np-0-3be8d635-6s5n                                   0          0%          0      0%           8G    8G
  └─ gke-sandpit-prod-np-0-947e5a45-84zk                                   0          0%          0      0%           8G    8G
  memory                                                           2711128Ki         50%  9150040Ki    169%          5Gi     0
  ├─ gke-sandpit-prod-np-0-3be8d635-6s5n                           1707608Ki         63%  5228120Ki    194%          2Gi     0
  │  ├─ contact-handler-api-f9475bd5-78kjw                              50Mi                  100Mi
  │  ├─ default-backend-66b68fd9d8-n5xwl                                 5Mi                   10Mi
  │  ├─ fluentd-gcp-v3.1.1-b8hks                                        50Mi                  100Mi
  │  ├─ grafana-1                                                       50Mi                  250Mi
  │  ├─ heapster-7f9cb9f8d5-x5zp2                                   215640Ki               215640Ki
  │  ├─ kube-dns-795f7b8488-2scvh                                       90Mi                  150Mi
  │  ├─ kube-dns-795f7b8488-fth52                                       90Mi                  150Mi
  │  ├─ kube-ops-view-667db9cf57-sdqlx                                  50Mi                  100Mi
  │  ├─ kube-ops-view-redis-74c9895647-ztv76                            50Mi                  100Mi
  │  ├─ kube-state-metrics-569465d5cb-55mb4                             50Mi                  100Mi
  │  ├─ l7-default-backend-fd59995cd-l85xz                              20Mi                   20Mi
  │  ├─ metrics-server-v0.3.1-57c75779f-lqkqv                          105Mi                  355Mi
  │  ├─ nginx-ingress-controller-85744dcf89-thkb4                       10Mi                  100Mi
  │  ├─ node-exporter-rqb4c                                             10Mi                   40Mi
  │  ├─ prometheus-sandpit-1                                            50Mi                   50Mi
  │  ├─ prometheus-operator-6977fd7f46-hhftc                            50Mi                  100Mi
  │  ├─ prometheus-to-sd-skmd9                                          20Mi                   20Mi
  │  ├─ right-sizer-789b55558-7lkpg                                     25Mi                   50Mi
  │  ├─ stackdriver-metadata-agent-cluster-level-69884d6fd6-l5vwf       50Mi                      0
  └─ gke-sandpit-prod-np-0-947e5a45-84zk                                980Mi         37%     3830Mi    145%          2Gi     0
     ├─ fluentd-gcp-v3.1.1-bzh54                                        50Mi                  100Mi
     ├─ grafana-0                                                       50Mi                  250Mi
     ├─ heapster-v1.6.1-5b6bf6cc74-52zqd                               140Mi                  170Mi
     ├─ nginx-ingress-controller-85744dcf89-nswr9                       10Mi                  100Mi
     ├─ node-exporter-7ml5r                                             10Mi                   40Mi
     ├─ prometheus-sandpit-0                                            50Mi                   50Mi
     ├─ prometheus-to-sd-cjkwx                                          20Mi                   20Mi
  pods                                                                     0          0%          0      0%          220   220
  ├─ gke-sandpit-prod-np-0-3be8d635-6s5n                                   0          0%          0      0%          110   110
  └─ gke-sandpit-prod-np-0-947e5a45-84zk                                   0          0%          0      0%          110   110
```

---

## Housekeeping

If you're running a cluster that's been around for a while, there are a few plugins that help you track down old or unused object:

- `kubectl deprecations` - as the name suggests, warns you of objects in your cluster using older versions of objects. On my home setup, it was pretty good at telling me about all the APIs I was using that I'd specified using `v1beta1` for example :smile
- `kubectl outdated` is a very neat idea - it finds images that running older versions than what's out there, but I need to figure out how to allow it read my private GCR. The [asciinema](https://github.com/replicatedhq/outdated) shows this better than a copy/paste of the output, as it shows the colour highlighting
- `kubectl prune-unused configmaps --dry-run` is a handy one - it wipes out (... in other words you should definitely run with `--dry-run` first ...) unused `Secret`s or `ConfigMap`s. Great for clusters I use for experimentation - which often leads to trashing things that don't work out

//// TODO outdated with private GCR

---

## Security

I'll be the first to admit I didn't "get" the RBAC setup in Kubernetes straight away, especially in terms of how it ties through to GCP stuff. Even though I feel like I now do, a couple of plugins to visualise permissions I find really useful.

The first is `kubectl access-matrix`, which in an RBAC-enabled cluster (which you really, really should have by now) has two modes:

- `kubectl access-matrix [-n=namespace]` displays a table showing what the current user can do against each resource type. With no namespace, it looks at the cluster scope
- `kubectl access-matrix for <resource>` displays, for a given resource (e.g. `pod`), which users/groups/serviceaccounts can perform which roles in which namespaces if applicable.

It is the latter I find particularly helpful in our `PodSecurityPolicy`-enabled, multi-tenant (i.e. permissions per namespace) cluster at work.

Related, is `kubectl who-can`, which as the name suggests, shows which subjects can perform what actions on what objects, for example, the following tests who can create a `Service` in the sandpit namespace:

```sh
[~ (⎈ |sandpit-prod:default)]$ kubectl who-can create services -n=sandpit

No subjects found with permissions to create services assigned through RoleBindings

CLUSTERROLEBINDING                                    SUBJECT                                       TYPE            SA-NAMESPACE
cert-manager                                          cert-manager                                  ServiceAccount  kube-system
cluster-admin                                         system:masters                                Group
elastic-cluster-admin-binding                         some-sa@some-project.iam.gserviceaccount.com  User
prometheus-operator                                   prometheus-operator                           ServiceAccount  prometheus
system:controller:clusterrole-aggregation-controller  clusterrole-aggregation-controller            ServiceAccount  kube-system
system:controller:persistent-volume-binder            persistent-volume-binder                      ServiceAccount  kube-system
```

In a better policed cluster than my sandpit (:grin) you'd expect to see something a bit tighter than this!

If you're worried about the hardening of your pods and willing to trust the folks at [kubesec.io](https://kubesec.io), then `kubectl kubesec-scan` is quite interesting. Here's some output from an old deployment I know still has some dodgy stuff in it :grin

```sh
[~ $ kubectl kubesec-scan deployment dodgy-app -n=dodgy
scanning deployment dodgy-app in namespace dodgy
kubesec.io score: -56
-----------------
Critical
1. containers[] .securityContext .capabilities .add | index("SYS_ADMIN")
CAP_SYS_ADMIN is the most privileged capability and should always be avoided
2. containers[] .securityContext .privileged == true
Privileged containers can allow almost completely unrestricted host access
-----------------
Advise1. containers[] .securityContext .runAsNonRoot == true
Force the running image to run as a non-root user to ensure least privilege
2. containers[] .securityContext .capabilities .drop
Reducing kernel capabilities available to a container limits its attack surface
3. containers[] .securityContext .readOnlyRootFilesystem == true
An immutable root filesystem can prevent malicious binaries being added to PATH and increase attack cost
4. containers[] .securityContext .runAsUser > 10000
Run as a high-UID user to avoid conflicts with the host's user table
5. containers[] .securityContext .capabilities .drop | index("ALL")
Drop all capabilities and add only those required to reduce syscall attack surface
```

Sounds like some good stuff to be doing ... :eye-open

And a quick one for the lazy / keyboard efficient - `kubectl view-secret [<secret>] [<key>] [-all]` - no need to base64 decode things - tedious!

---

## Debugging

In the last section, I cover a few plugins that can help with debugging issues.

First up, another one for those who don't like typing - `kubectl iexec` - simplifies the `kubectl exec -it <pod> /bin/sh` wrapper by offering an interactive menu to pick the pod (+ container) you want to exec onto.

> If this breaks for some reason, `kubectl pod-shell` is basically the same thing

Sometimes - hopefully rarely - you need to elevate privilege to get to the bottom of what's going on. These two commands help with this process, assuming you have enough permissions to do so:

- `kubectl node-admin` is scary - it spins up a privileged `Pod` with the host node mounted, effectively acting as a "remote onto nodes" container - a bit like this -> `kubectl run node-admin -i -t --rm --restart=Never --image=debian:latest --overrides="${SPEC_JSON}"` (with SPEC_JSON built up with a whole load of scary privileges to do stuff). Powerful, if for some ungodly reason you need to get onto a node directly. `kubectl node-shell` does largely the same thing too (except with `alpine`, and no fancy node-picker)
- `kubectl spy` (`kubespy`) is a handy little utility for spinning up a (privileged) pod that can attach to another pod for debugging purposes. There is also `kubectl debug`, but this needs the 1.16 alpha feature `EphemeralContainers` switched on, which I don't have. `kubectl spy` looks a bit like this:

```sh
[~ (⎈ |sandpit-prod:mosstech)]$ kubectl spy mosstech-798ff8d85-4qfp5
loading spy pod...
If you don't see a command prompt, try pressing enter.
/ # ls -la
total 44
drwxr-xr-x    1 root     root          4096 Feb 16 21:13 .
drwxr-xr-x    1 root     root          4096 Feb 16 21:13 ..
-rwxr-xr-x    1 root     root             0 Feb 16 21:13 .dockerenv
drwxr-xr-x    2 root     root         12288 Dec 23 19:21 bin
drwxr-xr-x    5 root     root           360 Feb 16 21:13 dev
drwxr-xr-x    1 root     root          4096 Feb 16 21:13 etc
drwxr-xr-x    2 nobody   nogroup       4096 Dec 23 19:21 home
dr-xr-xr-x  230 root     root             0 Feb 16 00:09 proc
drwx------    1 root     root          4096 Feb 16 21:13 root
dr-xr-xr-x   12 root     root             0 Feb 16 00:08 sys
drwxrwxrwt    2 root     root          4096 Dec 23 19:21 tmp
drwxr-xr-x    3 root     root          4096 Dec 23 19:21 usr
drwxr-xr-x    4 root     root          4096 Dec 23 19:21 var
/ # ps -ef
PID   USER     TIME  COMMAND
    1 101       0:00 nginx: master process nginx -g daemon off;
    6 101       0:00 nginx: worker process
   13 root      0:00 sh
   19 root      0:00 ps -ef
```

Meanwhile, a new pod is running:

```sh
[~ (⎈ |sandpit-prod:mosstech)]$ kubectl neat pods spy-5092 -o yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    run: spy-5092
  name: spy-5092
  namespace: mosstech
  spec:  containers:
  - args:
    - docker
    - run
    - -it
    - --network=container:bc5a8e226cf664145bbf63eee55cf8413496e9df9eba4eff75756c35b2f543dd
    - --pid=container:bc5a8e226cf664145bbf63eee55cf8413496e9df9eba4eff75756c35b2f543dd
    - --ipc=container:bc5a8e226cf664145bbf63eee55cf8413496e9df9eba4eff75756c35b2f543dd
    - busybox:latest
    command:
    - /bin/chroot
    - /host
    image: busybox
    name: spy
    stdin: true
    stdinOnce: true
    tty: true
    volumeMounts:
    - mountPath: /host
      name: node
  hostIPC: true
  hostNetwork: true
  hostPID: true
  priority: 0
  restartPolicy: Never
  serviceAccountName: default
  tolerations:
  - effect: NoExecute
    key: node.kubernetes.io/not-ready
    operator: Exists
    tolerationSeconds: 300
  - effect: NoExecute
    key: node.kubernetes.io/unreachable
    operator: Exists
    tolerationSeconds: 300
  volumes:
  - hostPath:
      path: /
    name: node
```

---

## And finally ...

And, if you've made it this far, I've saved the best till last. Are you brave enough to try ...

`kubectl snap`

For Avengers fans hopefully what this is going to do is obvious :grin (It deletes half of ... everything)

/// balance joke/meme
