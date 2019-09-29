---
title: "Squeezing GKE Resources In Small Clusters"
date: 2019-08-28T18:00:00-01:00
author: "@alexdmoss"
description: "Trying to create as much usable Compute in a tiny GKE cluster as possible, by squeezing its system resources"
banner: "/images/squeeze-1.jpg"
tags: [ "GKE", "Kubernetes", "GCP", "VPA", "Autoscaling", "Cost" ]
categories: [ "GKE", "Kubernetes", "Cost" ]
draft: true
---

{{< figure src="/images/squeeze-1.jpg?width=600px&classes=shadow" attr="Photo by Davide Ragusa on Unsplash" attrlink="https://unsplash.com/photos/cDwZ40Lj9eo" >}}

**Spoiler Alert!** This blog is *really* about Vertical Pod Autoscaling and patching of `kube-system` workloads in GKE. It just might not sound like it at the start :smile: If you're not interested in how I got there and just want to jump to the good stuff - I've summarised it at the bottom of the post.

---

So, this all started because I wanted to try out Elasticsearch. I'm not the biggest fan of GCP's Stackdriver Logging, and I wanted to poke around with some alternatives. I have a "personal" GKE cluster that I use to run a few websites (including this blog!), and where I also experiment with these sort of things outside of work. Should be straight-forward I thought to myself. Sure, Elastic is Java - but I don't need it to do that much, I'll just shrink it down enough that it'll start up ok, play around a bit, and then maybe see if I want to do more with it or sign up to a managed logging service instead.

Well, I was dead wrong about that! I quickly realised that there's no way it was going to fit - even with only a collection of tiny workloads, too much of the machine resources were gobbled up by things running in `kube-system` and I really didn't want to extend beyond my 3 x g1-small instances if I could possibly help it.

So, while I could make them temporarily larger while I experiment (yay Cloud), where's the fun in that?! What follows is a look at what's using up those resources and how I can shrink them down to smaller sizes - because let's be honest, the handful of services I run on my GKE cluster aren't going to need them to do a huge amount.

---

## So, What's Going On?

*Visualise the Problem*. No, not some sort of zen meditative technique. What are we actually dealing with here? Surely a dozen or so deployments of tiny web-apps isn't gobbling up 3x1.7Gb?

As with all things in life, we start with `bash`:

```bash
kubectl get po $1 \
-o custom-columns=NAME:.metadata.name,CPU_REQ:.spec.containers[].resources.requests.cpu,CPU_LIMIT:.spec.containers[].resources.limits.cpu,MEMORY_REQ:.spec.containers[].resources.requests.memory,MEMORY_LIMIT:.spec.containers[].resources.limits.memory
```

> I later discovered [`kube-capacity`](https://github.com/robscott/kube-capacity) through `krew`, which is a neater way of doing this

This ugly-as-hell one-liner produces something that looks a bit nicer than you'd think:

```bash
[~ (⎈ |my-cluster:default)]$ get-pod-resource-requests.sh kube-system
NAME                                                        CPU_REQ   CPU_LIMIT   MEMORY_REQ   MEMORY_LIMIT
event-exporter-v0.2.5-7d99d74cf8-9jdn4                      <none>    <none>      <none>       <none>
fluentd-gcp-scaler-55bdf597c-xsxwr                          <none>    <none>      <none>       <none>
fluentd-gcp-v3.1.1-6r6x5                                    10m       100m        100Mi        250Mi
fluentd-gcp-v3.1.1-mzgdn                                    10m       100m        100Mi        250Mi
heapster-v1.6.1-5b5df74474-vjvhd                            13m       13m         120Mi        120Mi
kube-dns-6987857fdb-42kzx                                   100m      <none>      70Mi         170Mi
kube-dns-6987857fdb-w8bd4                                   100m      <none>      70Mi         170Mi
kube-dns-autoscaler-bb58c6784-rsklz                         20m       <none>      10Mi         <none>
kube-proxy-gke-mw-prod-np-1-2cfd6748-k8kg                   100m      <none>      <none>       <none>
kube-proxy-gke-mw-prod-np-1-68476602-7zlg                   100m      <none>      <none>       <none>
l7-default-backend-fd59995cd-vcmd4                          10m       10m         20Mi         20Mi
metrics-server-v0.3.1-57c75779f-zmbt6                       43m       43m         55Mi         55Mi
prometheus-to-sd-rgdvx                                      1m        3m          20Mi         20Mi
prometheus-to-sd-smrbm                                      1m        3m          20Mi         20Mi
stackdriver-metadata-agent-cluster-level-55dfd764dd-8q674   40m       <none>      50Mi         <none>
```

Useful, but not sexy. I wanted a view of the packing onto my nodes. For this I dusted down a copy of [kube-ops-view](https://github.com/hjacobs/kube-ops-view/), which is something I'd played around with before. [My deployment](https://github.com/alexdmoss/kube-ops-view) of it is largely based on the sample yaml's with a view tweaks to secure it and make the Redis pod a little happier. This gives something visually more appealing:

![EFK won't fit!](/images/squeeze-gke-elastic-wont-fit.png?width=600px)

The red circle shows the unscheduled Elasticsearch pod. The vertical bars - in particular the amber ones, which show memory - tell us that we're not going to be able to fit something this size on here (and Kubernetes won't re-pack the existing workloads as it is honouring their anti-affinity rules). I actually found that even with a third node, it'd struggle because the kube-system workloads expand over time - many of them don't have caps on resource utilisation as we see above from the sexy-bash earlier.

```sh
[~/personal/github/elastic-on-gke (⎈ |mw-prod:elastic)]$ kubectl get pods
NAME                          READY   STATUS     RESTARTS   AGE
efk-elasticsearch-0           0/2     Pending    0          7m34s
efk-fluentd-es-bjp4m          1/1     Running    0          7m34s
efk-fluentd-es-zn5r7          1/1     Running    0          7m34s
efk-kibana-5bf6df5fc7-6p9h5   1/1     Running    0          7m34s
efk-kibana-7877d95c-nc4cp     0/1     Running    1          7m34s
efk-kibana-init-job-f2jph     0/1     Init:0/1   0          7m34s
```

```sh
kubectl describe pod efk-elasticsearch-0

Events:
  Type     Reason            Age                  From               Message
  ----     ------            ----                 ----               -------
  Warning  FailedScheduling  10m (x5 over 10m)    default-scheduler  pod has unbound immediate PersistentVolumeClaims
  Warning  FailedScheduling  108s (x21 over 10m)  default-scheduler  0/2 nodes are available: 1 Insufficient memory, 1 node(s) had volume node affinity conflict
```

---

## Someone Else Has Solved This, Right?

As is standard in our lovely industry, the first thing I did was of course to Google it. Someone else must've had this problem before, right? Turns out, Google themselves - or rather, they've been asked the question before and produced [a handy document as a starting point](https://cloud.google.com/kubernetes-engine/docs/how-to/small-cluster-tuning).

This sounded perfect! I wouldn't have been writing a blog post if that were true however :disappointed:

There *is* some good advice in here - but a lot of it didn't really apply for me. I wanted to keep Stackdriver Logging/Monitoring switched on in the meantime while I was experimenting - and might need the metric exports anyway later. The Kubernetes Dashboard was already switched off because it [needs high security privileges](https://cloud.google.com/kubernetes-engine/docs/how-to/hardening-your-cluster#disable_kubernetes_dashboard).

One thing from this documentation that I did leap on was the `ScalingPolicy` for fluentd - I followed this recommendation with a ScalingPolicy as follows and it worked a treat:

```yaml
---
apiVersion: scalingpolicy.kope.io/v1alpha1
kind: ScalingPolicy
metadata:
  name: fluentd-gcp-scaling-policy
  namespace: kube-system
spec:
  containers:
  - name: fluentd-gcp
    resources:
      requests:
      - resource: cpu
        base: 10m
      - resource: memory
        base: 100Mi
      limits:
      - resource: cpu
        base: 100m
      - resource: memory
        base: 250Mi
```

I was really impressed with this - and thought I'll just do the same thing for other stuff, like `kube-dns`, `kube-proxy`, etc. Job done right? NOPE! :thumbsdown: Turns out looking a little closer, this relies on an additional component running inside GKE - `fluentd-gcp-scaler` - which is [based on this](https://github.com/justinsb/scaler).

Cool concept - so I did spend a little bit of time trying to run my own equivalent of this in my cluster ... and in so doing, realised there was a better way ...

---

## Enter Vertical Pod Autoscaling

VPA is a feature that Google [recently announced](https://cloud.google.com/blog/products/containers-kubernetes/using-advanced-kubernetes-autoscaling-with-vertical-pod-autoscaler-and-node-auto-provisioning) in Beta. It is talked about in the context of GKE Advanced / Anthos, so I may need to keep an eye on whether it becomes a chargeable product :moneybag: - but in the meantime, it seemed worth experimenting with.

The [google_container_cluster Terraform resource](https://www.terraform.io/docs/providers/google/r/container_cluster.html#vertical_pod_autoscaling) already contains the option to enable GKE's VPA addon, so turning this on was a breeze:

```terraform
resource "google_container_cluster" "cluster" {
  # enabling VPA needs Beta
  provider = "google-beta"

  # [... other important stuff ...]

  vertical_pod_autoscaling {
    enabled = true
  }
}
```

As part of this I also flipped on Cluster Autoscaling just to see if the behaviours involved here caused my cluster to flex in size (I'd removed my EFK deployment at this point, so we were back at square one).

Getting VPA to do its thing involves applying some straight-forward policy, which looks a bit like this:

```yaml
apiVersion: autoscaling.k8s.io/v1beta2
kind: VerticalPodAutoscaler
metadata:
  name: kube-dns-vpa
  namespace: kube-system
spec:
  targetRef:
    apiVersion: extensions/v1beta1
        kind: Deployment
        name: kube-dns
  updatePolicy:
    updateMode: "Off"
```

With this yaml, we create a VPA policy in Recommendation mode - describing the VPA then tells us what it thinks the resource bounds should be. Google are reasonably cagey on how exactly it works things out, but it doesn't seem to take too long to start making recommendations for you.

```sh
[~ (⎈ |mw-prod:kube-system)]$ kubectl describe vpa metrics-server-v0.3.1
Name:         metrics-server-v0.3.1
Namespace:    kube-system
# [ ... snip ... ]
Status:
  Recommendation:
    Container Recommendations:
      Container Name:  metrics-server
      Lower Bound:
        Cpu:     12m
        Memory:  131072k
      Target:
        Cpu:     12m
        Memory:  131072k
      Uncapped Target:
        Cpu:     12m
        Memory:  131072k
      Upper Bound:
        Cpu:     12m
        Memory:  131072k
```

That's all well and good, but more fun is changing this to `updateMode: "auto"` and letting it actually perform the resizing these pods for you. A handy extension to your VPA definitions that can be made here is to set your own upper/lower bounds - particularly useful for situations where workloads can be spiky or extra resource is needed for pod initialisation. For example:

```yaml
spec:
  resourcePolicy:
    containerPolicies:
    - containerName: '*'
      maxAllowed:
        memory: 2Gi
      minAllowed:
        memory: 100Mi
```

I set some VPA definitions up for all the things in `kube-system` and left it for a short while to do its thing. I ended up with the following:

Original:

```sh
NAME                                                        CPU_REQ   CPU_LIMIT   MEMORY_REQ   MEMORY_LIMIT
event-exporter-v0.2.5-7d99d74cf8-6blm5                      <none>    <none>      <none>       <none>
fluentd-gcp-v3.1.1-mzgdn                                    10m       100m        100Mi        250Mi
fluentd-gcp-v3.1.1-sjztp                                    10m       100m        100Mi        250Mi
fluentd-gcp-scaler-55bdf597c-xsxwr                          <none>    <none>      <none>       <none>
heapster-v1.6.1-74885ff99d-wkjl9                            13m       13m         120Mi        120Mi
kube-dns-6987857fdb-7jm2z                                   100m      <none>      70Mi         170Mi
kube-dns-6987857fdb-w8bd4                                   100m      <none>      70Mi         170Mi
metrics-server-v0.3.1-57c75779f-zmbt6                       43m       43m         55Mi         55Mi
```

New:

```sh
NAME                                                        CPU_REQ   CPU_LIMIT   MEMORY_REQ   MEMORY_LIMIT
event-exporter-v0.2.5-7d99d74cf8-drbv2                      12m       <none>      131072k      <none>
fluentd-gcp-v3.1.1-27gs7                                    23m       230m        225384266    563460665
fluentd-gcp-v3.1.1-8bhv4                                    23m       230m        203699302    509248255
fluentd-gcp-v3.1.1-rmh47                                    23m       230m        203699302    509248255
fluentd-gcp-scaler-55bdf597c-dzxf5                          63m       <none>      262144k      <none>
heapster-v1.6.1-74885ff99d-8gksk                            11m       11m         87381333     87381333
kube-dns-6987857fdb-7jh59                                   11m       <none>      100Mi        254654171428m
kube-dns-6987857fdb-dqtft                                   11m       <none>      100Mi        254654171428m
metrics-server-v0.3.1-57c75779f-mpdvq                       12m       12m         131072k      131072k
```

All kinds of wacky units involved! As well as occasionally recommending things at a larger size than I hoped, I also couldn't get it to target certain resources - namely the `kube-proxy` pods which aren't DaemonSets as expected (or Deployments/StatefulSets), but individual pods in the world of GKE (weird, right?) - this meant the VPA resource couldn't target them, as it relies on the `targetRef` field (rather than something like a label selector, which it looks like it used to support but now no longer does).

---

## AutoScale ALL THE THINGS ... Oops!

This seemed effective enough that I fancied rolling it out to my actual workloads too (rather than just the stuff in kube-system) - with that in mind I created a lightweight controller (inspired by some of the work my colleagues have done) - code for it is here: https://github.com/alexdmoss/right-sizer. This will skim through Deployments every 10 mins and create VPA Policies for any new workloads it spots. This had rather comedic effects with `updateMode: Auto` :smile:

![Uh-oh - BIG cluster!](/images/squeeze-gke-all-auto-vpa.png?width=600px)

As can be seen by the red bars - this happened a few minutes after the VPA policies were created and isn't super-suprising - all those tiny pods of nginx etc were rapidly set with requests of 200-500Mi. For nodes with only 1Gb of spare RAM available, there was no choice but for the Cluster Autoscaler to kick in!

I did this as a bit of an experiment - but it's obvious that we need to be careful with this stuff. The VPA only has limited info to go on, and unless you set `resourcePolicies` to cap it to sensible values (not so practical for a Controller that applies to all Deployments!) it can do some wacky things. For this reason, I switched things back to recommend-only mode for all my workloads, and then used this data to set sensible defaults that I was happy with for my pods.

---

## Right-sizing kube-system

I still had a problem. I had some recommendations from VPA - some I didn't agree with - but the beefier workloads on my tiny cluster were mostly the ones residing in `kube-system`, for which I don't control the Deployments (etc).

It's here that things get hacky. I extended my Controller (used earlier to create the VPA policies) to also set some resource requests/limits on the kube-system resources that were on the larger side. [The code for it](https://github.com/alexdmoss/right-sizer/blob/master/main.py#L102) is really quite awful (it was a quick proof-of-concept!) but has been ticking away for a few weeks now and seems to be working out ok, so I really should clean it up, parameterise it, etc.

As a general approach, this is also pretty uncomfortable for something important in Production. Google will no doubt be reconciling these things themselves - we're effectively overwriting their chosen settings on a more frequent basis - hence my "hacky" comment. But, for a small cluster used for less important things, this approach does seem to work out okay.

---

## Conclusions

{{< figure src="/images/squeeze-2.jpg?width=600px&classes=shadow" attr="Photo by Josh Appel on Unsplash" attrlink="https://unsplash.com/photos/NeTPASr-bmQ" >}}

So, can I run Elastic now?

![Final Squeeze](/images/squeeze-gke-final.png?width=600px)

As can be seen above, we've just about got it squeezed onto a pair of n1-standard-1's (so slightly bigger machines than when we started, but only two of them - so cost neutral, more or less).

More relevant is that I learnt some interesting things about VPA and right-sizing of workloads, and some idiosyncrasies
 in how GKE manages pods in the `kube-system` namespace!

---

## Evolution

I was sufficiently impressed with VPA that it seemed worth a closer look in a work context too. At a larger scale the savings in Compute could become quite significant - depending on the maturity of testing practices in teams (some of our teams at work are *very* good at right-sizing their workloads already, whereas some could probably use the help).

We've therefore recently enabled it in Recommendation mode and started bringing the results into Prometheus and visualising them against current utilisation in Grafana - early days, but it looks really cool and I'd like to replicate the same in my home setup. Some of the recommendations are pretty quirky though, so it may need a bit more time to bed in ... and letting it auto-resize may not be viable for us given the amount of JVM-based workloads we run (it can't also set -Xmx, for example).

---

## In Summary

1. Get a tool that helps you visualise your resource requests/limits/utilisation. I like [kube-ops-view](https://github.com/hjacobs/kube-ops-view/) as it's simple but effective
2. Google have [an article](https://cloud.google.com/kubernetes-engine/docs/how-to/small-cluster-tuning) with specific advice for GKE. Many of the recommendations won't be suitable, but things like disabling the Kubernetes Dashboard and taking advantage of their fluentd autoscaler are good quick wins
3. In GKE, enable the VerticalPodAutoscaler addon and apply some VPA policies targeting the deployments you are interested in. I started in "Recommend" mode to see what it was going to do first
4. If you'd like a custom controller to setup VPA for all your deployments, [have a nose at this for inspiration](https://github.com/alexdmoss/right-sizer)
5. If you're comfortable with the recommendations and that your workloads can tolerate the restarts - switch on update mode
6. To really squeeze things down, you can update kube-system resources with a [custom controller](https://github.com/alexdmoss/right-sizer/blob/master/main.py#L102) which does the equivalent of `kubectl patch` on the resource requests/limits on a regular basis
