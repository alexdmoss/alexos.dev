Down with the Krew — my favourite kubectl plugins

Photo by Quino Al on Unsplash
In this post, I’ll briefly introduce krew, a plugin manager for kubectl, and then run through some of the plugins I’ve discovered with it that I’ve found particularly useful.
Yes! Made it through the intro without a clickbait title … “top 10 krew plugins to help you devop harder”
If you’d like to see some of the commands and their output in more detail — as well as a little bit more of me attempting to be funny … 😏 — then check out my more detailed post on this here.
Introducing krew

The official krew logo. It always important to have a good logo
Krew has actually been around since the end of 2018 — mostly likely around the time that kubectl plugins became a thing (v1.12). Krew’s job is to make installing these plugins easier, and was created by the same awesome person responsible for the indispensable kubectx/kubens.
Installation is a breeze, and the basic usage is intuitive:
kubectl krew list — list your plugins installed through krew. If you only use krew to do this, this should be the same as kubectl plugin list.
kubectl krew search — find plugins. Without <string>, all are listed.
kubectl krew info <plugin>— more info about what a plugin does. I always check this to make sure a plugin isn’t going to do something dubious — it invariably contains a link to the source code. And let’s be honest, you’re basically running someone else’s bash on your PC with these plugins.
kubectl krew install <plugin> — installs it. It can then be used through kubectl <plugin>.
kubectl krew upgrade — upgrade your plugins.
At time of writing, there are now ~70 plugins available through it. I feel like when I first started looking at krew in mid-2019, there were only a couple of dozen, so usage is clearly accelerating! 👍
I’ve grouped a few useful plugins I’ve found into five areas:
Viewing resources
Resource usage
Housekeeping
Security
Debugging
As mentioned earlier, if you want to see these in action with a bit more detail check out my more detailed blog post on the subject.
Viewing resources
First up in this section we have kubectl tail (usually kail). This is a handy plugin for tailing logs — it’s just nice and simple. There are a number of alternatives out there (I know some folks like Stern, which I’ve never got round to trying). The main feature I like is being able to target a deployment/service/ingress (e.g. kubectl tail --ing=mosstech) and be able to watch all the logs from the pods behind it, without necessarily needing to know how the resource has been labelled.
kubectl get-all does exactly what it says on the tin. They are not kidding — this really gets everything. I’ve found this to be increasingly useful with our uptake of CustomResourceDefinition — i.e. where it’s trickier to remember every type of object in a namespace.
Almost the opposite is kubectl neat which neatens up verbose output. If you find your mind glazing over the system-injected annotations and such when get you get a pod’s details, this one might be for you — it strips the gumpf from kubectl get — for example:

Lastly for this section are a couple of plugins that help visualise the relationship between a resource and other things in the cluster. kubectl pod-dive is good if you want to know what surrounds a particular pod, e.g.
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
             ├── kube-proxy-gke-sandpit-prod-np-0-947e5a45-84z
             ├── prometheus-to-sd-cjkwx
             ├── node-exporter-7ml5r
WAITING:
prometheus crashloopbackoff (Back-off 5m0s restarting failed container=prometheus pod=prometheus-sandpit-0_prometheus)
TERMINATION:
prometheus error (code 1)

Somewhat similarly, kubectl tree shows the hierarchy for a particular resource — especially useful when working with CustomResourceDefinition.
Resource usage

There are a whole bunch of plugins in this area — I’ve picked out a couple that I think do a bit more than what is easily achieved with basic kubectl usage. Some of these were particularly helpful when trying to squeeze more out of my tiny ‘home’ GKE cluster.
First, kubectl resource-capacity offers a nice and simple view for an overall cluster usage.
Occasionally, info about disk utilisation is useful, and not super-obvious from other sources … enter kubectl df-pv.
Finally, kubectl view-allocations is handy for spotting pods without resource requests/limits:
kubectl view-allocations
Recorded by alex
asciinema.org
Housekeeping
If you’re running a cluster that’s been around for a while, there are a few plugins that help you track down old or unused object:
kubectl deprecations — as the name suggests, warns you of objects in your cluster using older versions of objects. On my home setup, it was pretty good at telling me about all the APIs I was using that I’d specified using v1beta1 for example 😃
kubectl prune-unused configmaps is a handy one — it wipes out (… in other words you should definitely run with --dry-run first …) unused Secrets or ConfigMaps. Great for clusters I use for experimentation — which often leads to trashing things that don’t work out
kubectl outdated is a very neat idea — it finds images that running older versions than what’s out there in public — you can see it in action below:
kubectl outdated
Recorded by alex
asciinema.org
Note that as it needs to be able to connect to the registry anonymously, it can’t check things that are in my private GCR — but then again, they should have healthy CI/CD pipelines pushing out the latest image automatically. If you pull and re-tag images privately though (e.g. to run them through a vulnerability scanner) then this is a bit less useful to you (or at least, it is without some fiddling).
Security
I’ll be the first to admit I didn’t “get” the RBAC setup in Kubernetes straight away. Even though I feel like I now do, having a couple of plugins to visualise permissions I find really useful, especially at work where our setup is needfully more complex.
The first of these is kubectl access-matrix, which in an RBAC-enabled cluster has two modes:
kubectl access-matrix [-n=namespace] displays a table showing what the current user can do against each resource type. With no namespace, it looks at the cluster scope
kubectl access-matrix for <resource> displays, for a given resource (e.g. pod), which users/groups/serviceaccounts can perform which roles in which namespaces if applicable.
It is the latter I find particularly helpful in our PodSecurityPolicy-enabled, multi-tenant (i.e. permissions per namespace) cluster at work.
Related is kubectl who-can, which as the name suggests shows which subjects can perform what actions on which objects.

If you’re worried about the hardening of your pods and willing to trust the folks at kubesec.io, then kubectl kubesec-scan is quite interesting. Here’s some output from an old deployment I know still has some dodgy stuff in it 😜
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
And finally in this section, a quick one for the lazy / keyboard efficient — kubectl view-secret <secret> <key|-all> — no need to base64 decode things — tedious!
Debugging
In the last section, I cover a few plugins that can help with debugging issues.
First up, another one for those who don’t like typing — kubectl iexec — simplifies the kubectl exec -it <pod> /bin/sh wrapper by offering an interactive menu to pick the pod (+ container) you want to exec onto. kubectl pod-shell does basically the same thing too.
Sometimes — hopefully rarely — you need to elevate privilege to get to the bottom of what’s going on. There are a few krew plugins that help with this process, assuming you have enough permissions to do so:
kubectl node-admin is scary — it spins up a privileged Pod with the host node mounted, effectively acting as a “remote onto nodes” container. kubectl node-shell does largely the same thing too (except with Alpine, and no fancy node-picker).
kubectl spy (kubespy) is a handy little utility for spinning up a (privileged) busybox pod that can attach to another pod for debugging purposes. There is also kubectl debug, but this needs the 1.16 alpha feature EphemeralContainers switched on, which I don’t have.
And finally, if you’ve made it this far, something to bring perfect balance — kubectl snap. For Avengers fans hopefully what this is going to do is obvious 😁 (It deletes half of … everything).

In reality, it is somewhat gentle — see recording below:
kubectl snap
Recorded by alex
asciinema.org

Kubernetes
Krew
Plugins
DevOps



Alex Moss
WRITTEN BY
