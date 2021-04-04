---
title: "Tools To Visualise Your Kubernetes Cluster"
panelTitle: "Tools To Visualise Your Kubernetes Cluster"
date: 2021-03-28T18:00:00-01:00
author: "@alexdmoss"
description: "My experience running a couple of open source tools to help visualise workloads deployed on Kubernetes"
banner: "/images/???"
tags: [ "Kubernetes", "Weave", "Scope", "Ops View", "kube-ops-view", "Graph", "Observability" ]
categories: [ "Kubernetes", "Observability" ]
---

{{< figure src="/images/???.jpg?width=800px&classes=shadow" attr="Photo by ??? on Unsplash" attrlink="???" >}}

In this post I'm going to discuss a couple of tools I've used to try to visualise workloads deployed on my Kubernetes Clusters. The tools I'll be looking at are:

1. Kube Ops View
2. WeaveWorks Scope

I'll also finish by discussing some of the other approaches I am yet to try fully, but which might be better alternatives for larger clusters.

---

## Background

I have been on-off interested in tools to help visualise all the things running on a Kubernetes cluster since I first got involved in Kubernetes. A bash script that hacks together a bit of `kubectl` to dump out the pods on each node (later replaced by the built-in tools in the Google Cloud Console, since I use GKE) was an early creation that did a good enough job of satisfying my curiosity, but this is naturally hard to read and maintain.

My initial interest was motivated by trying to better understand how Kubernetes' scheduler makes decisions on where to put workloads, and potentially to try to get better bin packing (I've [written about this before](https://alexos.dev/2019/09/28/squeezing-gke-system-resources-in-small-clusters/)). That interest later evolved into wanting to understand the interconnectedness of workloads running on the cluster - and I'm going to describe two tools I've used to attempt to do this below.

/// medium - invite comments on alternatives here

## So Just How Loaded Are My Nodes, Anyway?

The first tool I'm going to mention is [kube-ops-view](https://codeberg.org/hjacobs/kube-ops-view). You can see my implementation of it [here](https://github.com/alexdmoss/kube-ops-view) which is a simple implementation of some CI to deploy the upstream into my cluster, alongside its Redis pod. We've found when running this on a large cluster at work (my definition of large being 150+ nodes and 3500+ pods).

{{< figure src="/images/kube-ops-view.png?width=800px&classes=shadow" attr="kube-ops-view running in my cluster" >}}

In the screenshot above we see the kube-ops-view UI running in my tiny cluster for running personal websites. I've highlighted a few of the interesting features.

The link to my github repo above contains the Kubernetes objects I use to deploy it from the upstream image. There's some basic wiring, plus the recommended redis backend for caching the results, plus a simple basic-auth implementation that relies on my nginx ingress-controller. The app supports OAuth but I just haven't got round to trying it out yet.

And just to prove it can, here's a couple of screenshots of it running against one of our much larger clusters at work - 154 nodes and 3681 pods at time of writing.

{{< figure src="/images/kube-ops-view-large-cluster.png?width=800px&classes=shadow" attr="kube-ops-view running in a larger cluster" >}}

It still feels pretty usable at this scale, especially combined with the filtering. You can see we've got our occupancy reasonably well sorted, particularly from a CPU perspective - as well as a few low occupancy nodes (they're running the telemetry stack on dedicated nodes).

> To work reliably the resource requests/limits needed to be increased - in this case a 600Mi limit on memory with unbounded CPU, and requesting half a core + 300Mi of memory (and similar for the redis pods). This is 6x the size of the pods in the small cluster (and the example manifests provided by its creator)

## I Need More. How Does It All Fit Together?

For this, I tried out WeaveWorks Scope. My original thoughts were to look to [Netflix's Vizceral](https://github.com/Netflix/vizceral) as I had seen that previously look pretty cool, but this is sadly no longer maintained. A little bit of searching led me to Scope though, so I decided to give it a try, with mixed results.

> You can see the adjustments I made to get it running in [this github repo](https://github.com/alexdmoss/weave-scope).

On the plus side, it did what I was after. I could use its UI to interrogate what was running on a node, and show the interconnectedness of workloads on the cluster. But it needs an awful lot of permission to work successfully - certainly in a work context rather than for a hobbyist cluster.

/// screenshots and showcase its powers more

My interest was peeked enough to spend a little time seeing whether I could reduce this down. The [repo](https://github.com/alexdmoss/weave-scope) elaborates on this more detail - I had some success, but unsurprisingly, it needed a fair bit of privilege to be able to get at the data it needed:

1. Disabled the ability to change configuration with `'--probe.no-controls=true'`. Change should be made through code, always, in my strongly-held opinion!
2. Drop privileges from its ClusterRole to only have get/list/watch
3. Tighten its PodSecurityPolicy. This was more successful for two of its three microservices
   1. For `weave-scope-app` and `weave-scope-cluster-agent` I could rebuild it from its upstream but modify its Dockerfile to run as non-root and my default restrictive PSP was fine
   2. This was no good for the `weave-scope-agent` DaemonSet, which is mounting host paths and sniffing the host network. I went through a few rounds of trial and error attempting to drop its system capabilities§. It ultimately needs access to `hostPID`, `hostNetwork` and `hostPath` as root - at which point you may as well just let it run in privileged mode anyway! On the plus side, at least this the DaemonSet component rather than the user-facing bit. Still, no way of avoiding it, it's a potential vector for a) gaining a huge amount of information on what's running on the cluster and where, and b) extremely powerful if successfully broken into

So, do we run it in our 150 node clusters at work? Sadly no - on top of the security concerns above, when spinning up against clusters of the sizes described earlier in this post, it didn't seem to work too well - services blipping in and out of the UI, very rarely able to render something as complex as that many microservices talking to each other.

To be fair, I did not try to resolve this - logs were silent, but most likely it just needed more juice (beyond the few extra pods I gave it). If it turned out the DaemonSet needed more than that would obviously have a more significant impact.

> Was the tool useful enough to be worth that extra cost? Would we be happy enough to open it up through an ingress for our users, or leave it running in the background as an exploitable tool for an attacker? In the end the conclusion was no.

For smaller clusters with less sensitive workloads though, Weave Scope is undeniably really cool though.

## The Nirvana?

/// check right usage

/// node graphing in grafana?

/// mesh
