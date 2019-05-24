---
title: "KubeCon 2019 - Other Cool Things"
date: 2019-05-23T23:00:00-01:00
author: "@alexdmoss"
description: "KubeCon 2019 - Other Cool Things"
banner: "/images/barcelona-main.jpg"
tags: [ "KubeCon", "CNCF", "opinion", "reflections" ]
categories: [ "Conference", "Kubernetes", "CNCF" ]
draft: true
---

First up though, there was the keynote. As keynotes go, it was pretty good. The main speakers that stuck in my mind were Spotify and Conde Nast. This fact reflects a theme from yesterday - that I am more interested in end users using this tech and what they're doing with it than listening to vendors of this tech talk to me about how they think I should be using it.

[David Xia](https://twitter.com/davidxia_?lang=en) @ Spotify talked about how he (and then his team) deleted some GKE clusters by accident and why that was a good thing. This I think was a great talk for a keynote - it's a conversation about learning:

- learning that having multiple tabs between Prod and Non-Prod open in your browser is dangerous - we've all been there, right?
- learning that your restore processes when you make a boo-boo might not actually be that great
- learning that migrating slowly so you have a backout plan or continuous availability is super-useful!
- learning that switching to new tooling (hi Terraform!) pays off in the end, but there's a learning curve that sometimes causes even more boo-boo's
- learning that Terraform's state file and declarative nature can just sometimes be a real a%$!
- learning that having a great team and culture that supports you is probably the most important thing of all

The keynote also featured [Katie Gamanji](https://twitter.com/k_gamanji?lang=en) from Conde Nast International - a company I confess to never having heard of, despite working in retail - they do digital publishing stuff and are pretty enormous.

I really liked this talk - they have some interesting regional challenges (wanting to run in Russia & China), a historic complexity of running decentralised IT in their local markets (that they wanted to unify in a good way), and they were very open about their tech stack and how they made their choices. Even though most of their choices are different to my own team's, it's really interesting to hear folks talking about *why* they made the choices they did.

For example - they're running 9 clusters globally now across 100 instances of AWS. This is their current choice, but they've self-hosted their k8s control plane using [Tectonic](https://coreos.com/tectonic/) to give them freedom to port their nodes to other providers in the future if there's good reasons why, pretty clever stuff.












Karol Golab & Beata Skiba - work on autoscaling at Google

## Intro

Initiative in the community to try to overcome limitation in the current spec.

You are expected to declare a static request (CPU & memory) in your deployment. This doesn't match real life!

- Daily/weekly patterns
- User base growth
- App lifecycle phases

Kubernetes Scheduler uses this as a contract - it won't schedule your workload without it. Workloads can get constrained basede on bin packing on the node

Vertical Pod Autoscaler - better!

Observe usage --> recommend resources --> updates resources (experimental)

The last bit is the hard bit!

VPA has three modes:

- off mode - dry run, does recommendations
- initial - updates requests of pods on creation
- auto (recreate) - will update by evicting and recreating

Why do we have to evict? It's because the PodSpec is immutable in etcd. Would even see it if you tried to Patch your resource limits directly - so not just a VPA thing.

So ... make it mutable?

- a breaking change
- there is a *lot* of code in Kubernetes that assumes a Pod is immutable
  - Scheduler (cache, extra work)
  - Kubelet
  - Quota

Downscaling is not easy - as what is reported by kubelet is not exactly the actual usage - cgroups are not perfect. Risk

There is a KEP out there pending approval

With this change - the PodSpec will describe a desired state, no longer the actual state - which I imagine matters for metrics? Instead you would look at PodStatus.ContainerStatus.ResourceAllocated

There will be a ResizePolicy on the podspec - NoRestart as default, but has RestartContainer (useful for Java where e.g. memory limits are set) and RestartPod (useful where some sort of InitContainer involved so RestartContainer can't be used) options.

Also RetryPolicy - NoRetry / RetryUpdate / Reschedule (RestartPod but for failure scenario)

## Summing Up

Eating the cake - minimised disruption + better stability with up-to-date resource requests

VPA can therefore offer: in-place only (e.g. long-running batch job that must not restart) or in-place preferred (hopefully new default - try to minimise disruption, but can fall back to a restart if needed)

VPA talk tomorrow - deep-dive on VPA if interested

## Q&A

Limits - will work the same as Request, move in-step with resize

JVM - confirmation that restart would still be required for Java containers. Comment also that VPA in general struggles a bit with Java stuff











> Went so so fast, hard to follow!

Laurent Bernaille & Robert Boll, Datadog infrastructure team

Datadog is a cloud-native service provider and cloud-native end user

Several thousand nodes of Kubernetes - biggest cluster 2k, typically 1-1.5k

## 1 - It's ~never~ always DNS

> Gets a round of applause!

Classic Kubernetes DNS - 3+ search domains, ndots:5

Means it searches - www.google.com.<namespace>.svc.cluster.local, www.google.com.svc.cluster.local, www.google.com.cluster.local, www.google.com.google.internal, www.google.com

- so 5 queries to make, load on CoreDNS

So, CoreDNS has an autopath option to strip the suffix off and realise it's a web address - one call. Good idea - turned it on

One day, DNS wasn't working - could see for CoreDNS metrics

Cause? Rate limited by the upstream

- autopath disables cache - hit their limit and all started failing

### 1.2 Another DNS example

CoreDNS rolling update and IPCS

Under high load, ports are reused faster than they expire against the TTL

Tunable, but kube-proxy doesn't set by default - net/ipv4/vs/expire_???

Graceful termination by setting weight in DNS to zero - good for TCP, not for UDP

77802 pull request to Kubernetes

## #2 Jobs are not starting, image pulls fail

Image Stampede - graph of number of image pulls - burst up at times

See 429's from image registry (too many requests)

What happened?

- Permission change on a bucket broke a DaemonSet
- DaemonSet in CrashLoopBackoff
- `imagePullPolicy: Always` ...
- ~ 1000 pods pulling through 3 NAT intances
- Daily quota reached for NAT IPs
- All NAT instances are impacted

So they replaced the impacted NAT instances

Forgot they use these for CloudSQL ... broke this connection. Tracking down fixed IPs

Afterwards:

- Admission Webhook denying "latest" tag
- Applied to Deployments, STS, DS, Jobs, Pods
  - Pods: Controllers can't create pods for existing workloads

## #3 I Can't Kubectl

Symptoms: Load on apiservers - >=50, can see apiserver getting OOM Killed, admin traffic much higher

kube2iam pod restarts - DaemonSet to handle AWS IAM permissions delegation. Restarts due to OOMKill

So they upped the limits by a lot

Patch to kube2iam - small typo - removed a node selector, so syncing **all** pods, not just all pods on the node (5-10)

Didn't see in test because much smaller cluster

## #4 New nodes aren't scheduling application pods

Symptom: see increase in events at the time, NotTriggerScaleUp - wouldn't fit if a new node is added

Single tenancy on a node

resources match node type

mins DaemonSet reserved

Added a DaemonSet with resource requests

Scheduler now couldn't fit applications on nodes

Even worse - DaemonSet had a critical PodPriority. Lucky - cluster was running k8s 1.10 new one, lucky didn't take the whole of DataDog down (would've evicted every pod)

## #5 Log intake volume

1mill/min to 15mill/min for one account

Enabled audit with a DaemonSet

Nodes running Kubernetes generate a hugo amout of audit logs (eexec, clones, iptables)

Kernel Audit DDOS

## #6 Where did all my pods go?

HorizontalPodAutoscaler fun

Controller manages replica count of dfeplomeny based on the value of a metric

Must remove explicit replica count!

Removing spec.replicas of the Deploemtn reset replicas count to single replica - annotation added. Need to kubectl edit to remove annotation at the time of removing replicas

## #7 There's a ghost in Cassandra

120 node Casandra, deployed fine

Next morning - broken

25% pods pending - "volume affinity issue"

25% nodes had been delted + loca lvolumes bound to deleted nodes

Nodes per AZ > ASG rebalance

Capacity issue in an AZ so AWS created he nodes in another AZ at the time of deployment

But ... ASG is clever, and when capacity in AZ is back in, it reprovisioned, broke things

ASG in separate zones now

## #8 Slow deploy heartbeat

Seeing deployments getting steadily slower

Cause - 4000 pending pods

apiserver etc pods seen rising in step

Cronjob evy 10 mins, wrong toelration - no nodes could satisfy, scheduling loop having a hard time

Every loop of scheduler trying to place these 4000 pods - slows it down

## #9 "Contained"

We expect containers to be contained ...

Broken runtime - containerd

### Broken Runtime 1 - ZOmbies

16018 zombie threads is `psauxf | grep -c defunct`

How did it happen? Readiness Probe doing a redis ping - this is slow on slave nodes (few seconds) compared to master. Timeout of 1s in the readinessProbe, so get killed, but nothing to reap the processes

--> Be careful with exec probes!

### Broken Runtime 2

Can't kill a process:

- Blocked IO (nvme issue)
- Deletion starts with cgroup freeze
- Freezing was hanging (process stuck on IO)

Very low level disk issue

### Performance Issue

This one is not so rare!

Symptoms: deployment slow, pod took +1 min to start

Cause: load on nodes running deployment high, because instance was writing a lot of IO

Why so many IOs? Core DNS queries

- nodes were running coredns pods
- An app started 5k+ DNS queries per scond
- coredns logs filled the disk

Changed App to use local DNS cache (as they'd had so many DNS issues in the past)

But still not quite right ...

DaemonSet had been removed, but audit was still enabled and writing to disk

But one node still had more iops ...

- can see these iops in 1 minute bursts
- CronJob running every minute
- but scheduler is very predictable - it picks the same node for the Job every time as image already there, space already available

This job was just synchronising between AWS instances & Consul - why would this generate a lot of IO?

Turns out kubectl does quite a lot behind the scenes!

- lots of queries, saves to $HOME/.kube

## #10 "Graceful" Termination

Queue consumer autoscaled on queue depth

On scale down, job must finish (hours - can be up to 48) - pod enters Terminating state

This is so long, things like kubelet restarts (cert rotation every 24hrs)

## Key Takeaways

1. Careful with DaemonSets in large clusters
2. DNS is hard
3. Cloud infra not transparent
4. Containers not really contained

For 4, they use pools per app to really contain (with taints and tolerations)




> This is a talk about kube-hunter, a tool that's been on my Friyay list for a while

---

Liz Rice, AquaSec - used to also sit on CNCF committee

## Intro

Pen testing specialist skills that you pay for.

Attack your own deployment

- nmap
- metasploit
- nessus

If you want to do this yourself, do you need to skill up on these tools?

https://github.com/aquasecurity/kube-hunter

- see what an attacker would see
- not a replacement for a full pen test
- but better than nothing, and free! Pen tests cost a lot ...

Generally what it does - find some open ports, try some API requests

Deliberately created a badly configured cluster to showcase.

Going to step through the things that kube-hunter does, but outside of it

## Demo

### Insecure Port

Curl an IP, try some typical ports - get a clue what the server is running (can see it is k8s)

From the API list, get some pod names and their images

With my image names, I know some vulnerabilities about them

> Do not leave the insecure port open

### Secure Port

Good news - I get a 403 here as I'm not authorised to make the request

Without auth, you're acting as the k8s RBAC service account `system:anonymous`, which hasn't been given any roles. Things like the /healthz endpoint are allowed in anonyously

Pods have access to a default (or custom, if you've done it) service account though

```sh
kubectl exec -it <pod> bash
export TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
curl -k -H "Authorization: Bearer ${TOKEN}" https://kubernetes/api/v1/namespaces/default/pods
```

So being able to get a reverse shell on a pod can be super dangerous! Your app being compromised doesn't just mean your app is the only thing impacted.

Also through secret leakage, same deal - e.g. ability to read secrets. Careful with your RBAC for secrets, may not just be app secrets.

Not just the api-server API ...

### etcd

`curl -k https://<ip>:2379/version` <-- etcd

### kubelet

`curl -k https://<ip>:10250/metrics` (or /pods) - insecure port example. Can pass bearer token in here also

This gives you all sorts of events and running pods

---

## kube-hunter

kube-hunter automates a lot of this manual crafting of curl commands.

Showcase it running against the api-server. It's python, CLI, spits out a load of interesting things.

Can run it inside a pod (assumes pod compromised). Results depend on your RBAC settings and the service account you use for the pod - so may need to run under multiple contexts to see all.

---

All this stuff is passive hunting - not tried to change the state.

What about active hunters? Trying to change the state of the cluster - create/destroy things to demonstrate the potential impact.

`./kube-hunter.py --active` - created a clusterrole and a privileged pod

### how does it work?

Python, observer pattern

- iterate on ports, where open, publish an event that it's open
- classes subscribe to these events to then perform tests, e.g. get metrics, published as an event
  - this can chain to then try other things based on what it found
- event per IP to test

Why show the code? They want more hunters! Aqua are encouraging contributions to build up this list - both ideas for hunters and implementations of them. See github link earlier

---

## Q&A

Can you run it in CI? *Sure, why not - not seen people do it but you could*

Planning to add detection for known CVEs, threat disclosures? *Yeah you could - contribute! Nothing to attack app code for example*

Do you see these insecure ports open with customers? *Don't have data, but it is definitely the case that it happens - especially older clusters that are forgotten about. More common is RBAC settings being too loose*






Danny & Alex - software engineers at Intuit

They use Argo rollouts

Intuit is about finance/accounting software

- 4k developers (1.2k now on Kubernetes)
- 21 locations
- $6b revenue
- 50m customers
- 120 clusters, 200-300 nodes
- 1300 deploys a day
- 2400 namespaces
- 27k pods

---

Kubernetes adoption started slightly slow, but is starting to look a bit more exponential now

Most engineers use RollingUpdate - version B rolling out slowly to replace version A

- so they started looking at B/G, Canary, A/B testing

Intuit follow a GitOps methodology

They developed Argo CD (open sourced it) to manage this - show and concile from what's in Git

How do you implement B/G or Canary while following GitOps - no built-in support in k8s for it

- declarative vs imperative rollout strategies

Also wanted a clean integration in CI, keep it simple

---

Attempt 1 - they started with Jenkins. It worked, but ...

- did not fit GitOps model (state also needed in Jenkins)
- Not idempotent
- very brittle (lots of assumptions and edge cases)
- Jenkisn requires k8s credentials to deploy (risk)
- Painful to setup
- and more ...

Attempt 2 - Deployment Hooks

- still not idempotent and not transparent
- requires a lot of work to start using it
- still not following GitOps - not declarative

Attempt 3 - Custom Controller

Argo Rollouts Design Considerations:

- Codifies the deployment orchestration in the controller
- GitOps-friendly (idempotent)
- Runs inside the k8s cluster (creds)
- Easy adoption and migration from deployments
- Feature parity with deployments

6 months work later ...

- advanced open source k8s deployment controller
- kubernetes native
- support blue-green & canary deployments
- https://github.com/argoproj/argo-rollouts

How it works

- handles ReplicaSet creation, scaling, and deletion
- single desired state as a Pod Spec
- Support manual and automated promotions
- Integrates with HPA

---

From Deployment to Argo Rollout:

- for a standard deployment - change the apiVersion & Kind to a Rollout, then add a `strategy:`

```yaml
strategy:
    blueGreen:
        activeService: active-svc
        previewService: preview-svc
        previewREeplicaCount: 1 # optional
        autoPromotionSeconds: 30 # optional
        scaleDownDelaySeconds: 30 # optional
```

- Manages an active and preview services slecttor to provide a service level cutover
- sizing control over preview env
- manual or automatic promotions

Canary:

```yaml
strategy:
    canary:
        maxSurge: 10%
        maxUnavailable: 1
    steps:
    - setWeight: 10 # percentage
    - pause:
        duration: 60 # seonds
```

- Declative promotion
- No servie modification
- Traffic split based on replcia ratio between versions of an application
- Superset of RollingUpdate Strategy

---

## Demo

They use ArgoCD - unsurprisingly - and also Kustomize.

Demo of Canary first - updated image tag in kustomize patch, CD triggers - deploys 20% (one pod out of five)

ArgoCD has a really nice visualisation for the running pods to show this working - can see relationship from pods to services

Then ues CLI or UI in ArgoCD to signal to complete the rollout - goes to 100% then winds down the old version

Nice and visual - test app red --> green, UI for ArgoCD looks nice

Then moves on to Blue-Green demo.

Very cool visual network view in ArgoCD also - showing blue and green traffic flowing

## What's Next?

- Service Mesh integration (to avoid need for large number of replicas to canary)
- Decision-based Promotion CRD
  - Deliberately didn't add any analysis into the Controller - keep it clean and tidy, composable
- A/B testing, experimentation strategies






Theme: Cloud Native / Kubernetes is a Journey and not a Destination

> T-shirts are in!

## Kubernetes - Don't Stop Believin'

Brian Liles is back.

Sadly he's not going to sing some Journey :(

Stat attack - 5 years of Kubernetes

- 31k contributors
- 164k commits
- 1.2m comments
- Graduated March 6th, 2018. CNCF saying "we've shipped this so much, it's stable"

Kubernetes is not a product - keep hammering the cloud-native platform for building platforms

Brian's job - what does optimal developer experience look like with Kubernetes?

People don't talk about v2 of deployment, service, etc very much

But can't stop - new paths to explore?

Ideas for where else to run Kubernetes other than servers and cloud

- in cars, in stores, IoT

Developer experience: client-go isn't for mortals. Barrier to entry

Go is a fine language, but ... there are a lot more JS/Py/Java/C developers out there

Kubernetes as a platform - Custom Resource Definitions

- lots are writing custom APIs for Kubernetes
- Operators as a pattern, but do they feel like a crutch? Do we need to think about it?
- Upbound: create objects in the cluster to expose things outside the cluster, like databases

Cluster API - tell Kubernetes what to do with itself!

- tell Kubernetes how to upgrade, change machine types etc - mutate itself

Build on vs Build with

- Operators/CRDs are "on"

KIND: Kubernetes in Docker - just need Docker to create a multi-node Kubernetes cluster

- super powerful for e.g. running on desktop, testing, etc

Talk later

Declarative Config Management

- haven't really solved this yet. Things like Helm act like rpms, debs
- kustomize is start of this - could be more out there

Bad advice - don't run stateful workloads on Kubernetes

- change the No to Maybe. It's still complex

Community - turn end-users into contributors, contributors into maintainers

## From COBOL to Kubernetes

A 250 year old bank's cloud-native journey - ABN-AMRO

Highly regulated but highly competitive industry. 20k employees, 400+ development teams, 3000+ apps

Usual reasons for containers - Dev speed, flexibility, consistency, cost-efficient

They are also seeing software suppliers deliver software in containers

l.cncf.io - so many choices. 300 teams would pick everything different - need guidelines, help them make these things more consumable

Team Stratus - "low-level clouds characterised by horizontal layering with a uniform base"

Sounds familiar ...

Stratus mission:

Enable development teams to quickly deliver secure and high quality software by providing them with:

- easy-to-use-platform
- security
- portability across clouds on enterprise levels
- reusable software components. Lego for DevOps teams ... hmmmmmmm ... Some stuff around "compliant by default"

They have platform teams specialising in AWS and Azure

Tech stack - Terraform, Nexus, Docker, CNI, Azure AKS, Amazon EKS, Helm, CloudBees, Azure DevOps, Twistlock, Vault, OPA (in pipelines & in cluster), Splunk (already had, reused compliance dashboards), Prometheus

They have SonarQube and Fortify in the software build pipeline also

Slide is quite good at illustrating:

- software pipeline --> docker pipeline --> deploy pipeline

OPA for compliance-as-code emphasised - benefits of software engineering practices like pairing, code review etc - same argument as Infra-as-code

- example of early feedback e.g. not allowed public load balancer. Feedback in response, look up code in pipeline for advice on how to get compliant

Stratus timeline:

- Q4 2018 - team created
- Q1 2019 - first MVP, EKS, Twistlock, docker pipeline, compliant
- Q2 2019 - working on platform govenrance, training & education, infra & compliance as code, metrics / telemetry

Very pleased how fast they've gone in a compliance world

Lessons learnt:

- focus on MVPs - new tools, new regulations - gotta stay on track
- holistic approach - don't only focus on the technical aspects, but also create clear governance
- platform capabilities - over just choosing tooling. Talk to teams about needs not tools
- iteration - start small, fail fast
- automation - help prevent and correct mistakes, and of course allow scale

## Metrics, Logs & Traces - What Does The Future Hold for Observability

Tom Wilkie from Grafana, Frederic Brancyzk from Red Hat

Three predictions ...

1 - Correlation

- more correlation between the pillars
- by three pillars - they mean metrics/logs/tracing

Being able to look at a latency graph and immediately switch to the logs for it, or looking at some logs and dive into the trace for a particular request

Loki - open source log aggregation, launched at Seattle KubeCon in December - use Prometheus service discovery to match up to log stream data

ElasticSearch -> Zipkin - put a link using a filter formatter to Zipkin tracing

OpenCensus - exemplars - example traces associated to a histogram

2 - New Signals & New Analysis

- Example - OOMKill of container - container comes up as new time-series as new process
- Google-wide profiling - a continuous profiling infrastructure for datacentres
- I think he's talking about time-series of Profiling your code - memory leak detection etc within your app
- ConProf - his new project

Tom to Frederic - continue the theme of implementing Google ideas as open source. Ouch!

3 - Rise of Index-free Log Aggregation

"Just give me log files and grep" ... old is new

OK Log - discontinued but cool ideas. Distributing grep

`kubectl logs` - but it'd be nice to get logs for pods that are missing, right?

Loki is an index-free log aggregator. It doesn't have the power of Elastic - it's there for developer troubleshooting, not business analytics. Hmm

> The idea of something focused on developer experience, and deliberately not something that's failing into a business analytics trap, is very appealing. Will go to talk later ...
