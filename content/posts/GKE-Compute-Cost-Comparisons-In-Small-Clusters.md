---
title: "GKE Compute Cost Comparisons In Small Clusters"
panelTitle: "GKE Compute Cost Comparisons"
date: 2023-05-14T18:00:00-01:00
author: "@alexdmoss"
description: "Comparing GKE's Manual pools, Node-AutoProvisioning & Autopilot to see which achieves the lowest cost for the least inconvenience"
banner: "/images/savings.jpg"
tags: [ "GKE", "Kubernetes", "Autopilot", "Node Auto-Provisioning", "GCP", "Vertical Pod Autoscaling", "VPA", "Autoscaling", "Cost" ]
draft: true
---

{{< figure src="/images/savings.jpg?width=800px&classes=shadow" attr="Photo by Micheile Henderson on Unsplash" attrlink="https://unsplash.com/photos/SoT4-mZhyhE" >}}

I've been running a short-term experiment across the GKE Compute options within my "Home Lab" _(read: GCP-hosted Google Kubernetes Engine)_ with a view to optimise for a balance between ease of operation and cost efficiency. If that sounds useful to you, then by all means read on!

---

First, lets start with the Billing view. The chart below shows the view of costs over the course of my two-month experiment for this small GKE Cluster. I'll be stepping across this graph to detail each experiment I ran and what led to these outcomes, as well as some of the trade-offs along the way.

{{< figure src="/images/compute-costs/gcp-spend.png?width=1000px&classes=shadow" attr="GCP Spend - Two Months" >}}

The collection of workloads deployed did not fluctuate much over the course of the experiment - we're talking about 13 or so `Deployments` with maybe the odd replica dropped here and there, but nothing significant.

---

## The As-Is

Let's start by describing the steady state. Prior to the experiments this cluster ticked along happily on 3 x `e2-medium` Spot Instance machines in a Zonal cluster. The initial part of the graph shows this from a cost perspective:

{{< figure src="/images/compute-costs/3-node-normal.png?width=1000px&classes=shadow" attr="GCP Spend - 3 x Spot Nodes, Normal Running" >}}

The remaining costs are standing charges for things like Secret Manager, Network Load Balancing and PD Storage - these do not fluctuate very much throughout the course of the experiment - pennies here and there, as the load on these apps is consistent.

---

## Squeeze Me Seymour

My first act was to try to squeeze this down to two nodes - which I managed, ish. I was aggressive with setting of resource requests & limits for my Pods whilst still ensuring everything was running, yielding the following saving:

{{< figure src="/images/compute-costs/2-node-squeeze.png?width=1000px&classes=shadow" attr="GCP Spend - 2 x Spot Nodes" >}}

> The nodes ended up as a custom machine type still on Spot E2, as it was easier to squeeze the CPU than the memory (these apps can tolerate being throttled as they're very low usage, but of course have minimum levels of memory that they need to function).

This worked for a little while. The memory pressure created did ultimately lead to scheduling issues though - especially for the chunkier workloads such as the databases behind Plausible, which I [self-host](https://alexos.dev/2022/03/26/hosting-plausible-analytics-on-kubernetes/).

The instability caused here led me to try one of GKE's funkier features ...

---

## Nap Time Boys n Girls

Node Auto-provisioning can be enabled per node-pool, and various options exist to set the behaviour you'd like. With appropriate tuning of behaviour this ended up around 25% more expensive than my extremely-squeezed setup, but around 33% cheaper than the three-node rig I started with:

{{< figure src="/images/compute-costs/nap.png?width=1000px&classes=shadow" attr="GCP Spend - Node Auto-provisioning + VPA" >}}

There's several bits n bobs going on here that are worth elaborating on:

1. Enabling this via Terraform is straight-forward. The block below has most of what you'd need. You have the option of doing this at the cluster level or creating a NAP-enabled node pool alongside your existing, and using taints & tolerations to progressively move your workloads over to the new setup, based on risk. This is handy!

```terraform
cluster_autoscaling {

  enabled = true

  resource_limits {
    resource_type = "cpu"
    minimum = "2"
    maximum = "10"
  }

  resource_limits {
    resource_type = "memory"
    minimum = "10"
    maximum = "20"
  }

  autoscaling_profile = "OPTIMIZE_UTILIZATION"      # this makes turn-down more aggressive

  auto_provisioning_defaults {
    service_account = google_service_account.cluster.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring.write",
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/trace.append",
      "https://www.googleapis.com/auth/servicecontrol",
      "https://www.googleapis.com/auth/service.management.readonly",
    ]
    disk_size = 30
  }
  
}
```

2. **[Vertical Pod Autoscaling](https://cloud.google.com/kubernetes-engine/docs/concepts/verticalpodautoscaler)**. Rather than continuing with my constrained Pod sizing, I enabled VPA for all workloads in "Auto" mode, which allows the pods to be restarted to set recommended sizes based on observed behaviour (none of these are on the JVM). I actually didn't monitor this as closely as I thought I would beforehand (!), just trusting the system to do its thing. The graph above is the result of that.

> Rather than create these for each workload independently, I used [Kyverno](https://kyverno.io) to generate these. That's probably a blog post in its own right, but here's the snippet of policy to generate these, in case useful:

{{< gist alexdmoss dc19d041a9b4a483b1de4abacdb20b8a >}}

3. Some [minor trickery](https://cloud.google.com/kubernetes-engine/docs/how-to/node-auto-provisioning#support_for_spot_vms) is required to ensure that NAP nodes are scheduled as **Spot Instances**, for the saving. This is because you're no longer in control of that from the cluster perspective, instead having to specify on the workload itself that it can tolerate being on Spot. I again used Kyverno for this, see this snippet _(note that I also had to target resources in `kube-system` for this to remove all non-Spot nodes):

{{< gist alexdmoss c48298c76096677eba687973adf30a50 >}}

> Fair warning - this wasn't as foolproof as I'd like. If you look closely at the cost chart above, they rise slightly towards the end alongside the appearance of some dark red blocks - this is non-Spot Compute nodes appearing. I did not investigate this closely, but something to be aware of and/or alert for, if you use this extensively.

4. Finally, ensure that the [safe-to-evict annotation](https://github.com/kubernetes/autoscaler/issues/3183) - which I already had set - is present, to help the Cluster Autoscaler turn-down nodes containing pods with ephemeral storage. Again, Kyverno can handle this:

{{< gist alexdmoss 2b6998256fa99d6ec5417afc064dd24e >}}

This worked well enough and I would've been content to stop here, but I was in a good position now to try out Autopilot - GKE's not-so-new option for going "nodeless" entirely.

---

## Whose Node Is It Anyway?

[Autopilot mode](https://cloud.google.com/kubernetes-engine/docs/concepts/autopilot-overview) for GKE takes away configuration of the node pools entirely. This is an either/or situation - you need to build a new cluster entirely to switch over to this, so needs careful adoption to bring about for workloads that are important to you.

You can see also from the graph below that this was not as cheap as the NAP option - settling into 15-20% more after I had gotten the setup stabilised on Spot _(you can see from the graph that it took me a little time to get things stable)_:

{{< figure src="/images/compute-costs/autopilot.png?width=1000px&classes=shadow" attr="GCP Spend - Autopilot" >}}

I set this up as a parallel cluster and then moved the workloads over to it bit by bit, due to some of the fundamental differences between Standard & Autopilot and not being willing to take the plunge on all at once. The terraform for the cluster is much simpler - `enable_autopilot = true` and running a plan will guide you into what's incompatible _(I found I could delete around 50 lines or so of config from the NAP-based one above)_.

As well as the need to create a new cluster, some other things to be aware of with Autopilot are:

1. Autopilot has [minimum required](https://cloud.google.com/kubernetes-engine/docs/concepts/autopilot-resource-requests#min-max-requests) CPU and Memory for a Pod as well as ratios between them - and it will raise the values to these as needed, depending on how you've configured your Pod spec. This meant that I actually chose to disable the Pod zonal anti-affinity rules I had set for some workloads (again - these aren't that important to me, I was just being a good boy) as this is one of the things that forces the minimums to be higher.

> A chose to leave VPA enabled even though I'm not sure it's actually effective at my scale, given the limits above :shrug:

2. Autopilot itself doesn't allow certain things, which may affect what you can run on it. [See this document](https://cloud.google.com/kubernetes-engine/docs/concepts/autopilot-security) for more detail - privileged pods for example may be an issue if you're using some weird-n-funky workloads (I used to run my Gitlab Runner out of GKE, which would've been an issue here if I wanted to use the docker-in-docker service, rather than kaniko)

I also chose to explicitly specify Spot in my Pod spec this time - based on the observation from the NAP setup not working perfectly via the Kyverno mutation. The docs recommend the following in your `.spec.template.spec`:

```yaml
spec:
  template:
    spec:
      nodeSelector:
        cloud.google.com/gke-spot: "true"
      terminationGracePeriodSeconds: 25
      # ... rest of pod spec
```

All that said though, I was pretty happy with the outcome. It feels leaner and simpler to just simply "not care" about the nodes any more (see [this list](https://cloud.google.com/kubernetes-engine/docs/resources/autopilot-standard-feature-comparison) for example of the differences). I wouldn't be worried about reverting back to NAP, but for the small increase in cost, I felt the trade-off was acceptable. It also lines up with Google's default behaviour going forward for GKE - I believe it became the default mode for new clusters from GKE v1.25.

---

## Conclusions

My conclusions from this - caveat: brief - bit of experimentation?

If you really, really need to minimise the damage to your bank balance, and are willing to risk the impact to service availability by getting things wrong or pushing things too hard, then right-sizing your workloads as tightly as possible gets you the biggest saving.

If you want to lower our cognitive load, then just don't worry about the nodes at all and enable GKE in Autopilot Mode. Check the feature compatibility needs against your workloads before doing so however, especially if you need to run privileged pods (ideally, you should not!). Even though the minimum pod spec sizes are generous for small apps, this seems to **nearly** balance out by avoiding the cost of "other stuff" on the node, but is slightly more expensive than you could achieve if you used ...

... Node Auto-provisioning fulfils a nice sweetspot if you have workloads that need the full-fat GKE Standard features, but you don't want to think about the Nodes _that_ much. I combined this with Vertical Pod Autoscaling to right-size my Pods for me, combined with a fairly tight overall ceiling on the size of compute available so that I didn't get a nasty surprise.

And of course one final point - if you choose either of the automatic compute options NAP or Autopilot, be sure to follow the guidance on scheduling onto Spot Instances to maximise your savings - this is significant!
