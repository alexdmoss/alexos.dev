---
title: "GKE Compute Cost Comparisons In Small Clusters"
panelTitle: "GKE Compute Cost Comparisons In Small Clusters"
date: 2023-05-15T18:00:00-01:00
author: "@alexdmoss"
description: "Comparing Autopilot, Node-AutoProvisioning to controlling the Nodes yourself, to see which can achieve the lowest cost for the least inconvenience"
banner: "/images/squeeze-1.jpg"
tags: [ "GKE", "Kubernetes", "Autopilot", "Node Auto-Provisioning", "GCP", "Vertical Pod Autoscaling", "VPA", "Autoscaling", "Cost" ]
draft: true
---

{{< figure src="/images/???.jpg?width=800px&classes=shadow" attr="Photo by ??? on Unsplash" attrlink="https://unsplash.com/photos/???" >}}

I've been running a short-term experiment across the GKE Compute options within my "Home Lab" _(read: GCP-hosted Google Kubernetes Engine)_ with a view to optimise for a balance between ease of operation and cost efficiency. If that sounds useful to you, then by all means read on!

---

First, lets start with the Billing view. The chart below shows the view of costs over the course of my two-month experiment for this small GKE Cluster. I'll be stepping across this graph to detail each experiment I ran and what led to these outcomes, as well as some of the trade-offs along the way.

![GCP Spend - Two Months](/images/gcp-spend.png)

The collection of workloads deployed did not fluctuate much over the course of the experiment - we're talking about 13 or so `Deployments` with maybe the odd replica dropped here and there, but nothing significant.

---

## The As-Is

Let's start by describing the steady state. Prior to the experiments this cluster ticked along happily on 3 x `e2-medium` Spot Instance machines in a Zonal cluster. The initial part of the graph shows this from a cost perspective:

![GCP Spend - 3 x Spot Nodes, Normal Running](/images/3-node-normal.png)

The remaining costs are standing charges for things like Secret Manager, Network Load Balancing and PD Storage - these do not fluctuate very much throughout the course of the experiment - pennies here and there, as the load on these apps is consistent.

---

## Squeeze Me Seymour

My first act was to try to squeeze this down to two nodes - which I managed, ish. I was aggressive with setting of resource requests & limits for my Pods whilst still ensuring everything was running, yielding the following saving:

![GCP Spend - 2 x Spot Nodes](/images/2-node-squeeze.png)

> The nodes ended up as a custom machine type still on Spot E2, as it was easier to squeeze the CPU than the memory (these apps can tolerate being throttled as they're very low usage, but of course have minimum levels of memory that they need to function).

This worked for a little while. The memory pressure created did ultimately lead to scheduling issues though - especially for the chunkier workloads such as the databases behind Plausible, which I [self-host](https://alexos.dev/2022/03/26/hosting-plausible-analytics-on-kubernetes/).

The instability caused here led me to try one of GKE's funkier features ...

---

## Nap Time Boys n Girls

Node Auto-provisioning can be enabled per node-pool, and various options exist to set the behaviour you'd like. With appropriate tuning of behaviour this ended up around 25% more expensive than my extremely-squeezed setup, but around 33% cheaper than the three-node rig I started with:

![GCP Spend - Node Auto-provisioning + VPA](/images/nap.png)

There's several bits n bobs going on here that are worth elaborating on:

1. Enabling this via Terraform is straight forward. I can afford to take risks with this cluster as the workloads are not important, but if you wanted to you could create a new node pool with NAP enabled, and use taints and tolerations to progressively move your workloads over to the new setup based on risk.
1. **Vertical Pod Autoscaling**. Rather than continuing
2. 
// NB VPA

---

## Whose Node Is It Anyway?

![GCP Spend - Autopilot](/images/autopilot.png)

---

## Conclusions

My conclusions from this - caveat: brief - bit of experimentation?

If you really, really need to minimise the damage to your bank balance, and are willing to risk the impact to service availability by getting things wrong or pushing things too hard, then right-sizing your workloads as tightly as possible gets you the biggest saving.

If you want to lower our cognitive load, then just don't worry about the nodes at all and enable GKE in Autopilot Mode. Check the feature compatibility needs against your workloads before doing so however, especially if you need to run privileged pods (ideally, you should not!). Even though the minimum pod spec sizes are generous for small apps, this seems to **nearly** balance out by avoiding the cost of "other stuff" on the node, but is slightly more expensive than you could achieve if you used ...

... Node Auto-provisioning fulfils a nice sweetspot if you have workloads that need the full-fat GKE Standard features, but you don't want to think about the Nodes _that_ much. I combined this with Vertical Pod Autoscaling to right-size my Pods for me, combined with a fairly tight overall ceiling on the size of compute available so that I didn't get a nasty surprise.

And of course one final point - if you choose either of the automatic compute options NAP or Autopilot, be sure to follow the guidance on scheduling onto Spot Instances to maximise your savings - this is significant!
