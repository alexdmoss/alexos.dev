---
title: "KubeCon 2019 - Other Cool Things"
date: 2019-05-23T23:00:00-01:00
author: "@alexdmoss"
description: "This post goes through some of the other tool things I learnt at KubeCon that don't fit in with the themes of the other articles!"
banner: "/images/kube-donuts.jpg"
tags: [ "KubeCon", "CNCF", "opinion", "reflections" ]
categories: [ "Conference", "Kubernetes", "CNCF", "Kubecon" ]
---

{{< figure src="/images/kube-donuts.jpg?width=600px&classes=shadow" title="Yes, this was a wall of donuts to celebrate Kubernetes 5th birthday!" >}}

I'm going to use this blog post to talk about a couple of other things I learnt about while at KubeCon 2019. These didn't fit so neatly into my other posts from my time at the conference, but I still felt they were interesting enough to call out separately!

---

## Compliance as Code

In some of the specialist tracks I went to, I learnt several other interesting things too.

For example, I started the day learning about some proposals around a Multi Cluster Ingress - Ingress v2 basically *(yes please, looked great!)*, an intro to SPIFFE and SPIRE *(which sounded useful but a bit early-days for folks like us, perhaps)*, some wildly sensible folks from Lyft & Triller talking about permissions options in Kubernetes *(reassuringly, we are doing a lot of things in this space already - I did learn some funky things about Admission Controllers and Open Policy Agent though)*.

The latter inspired me sufficiently to go to a more thorough talk on [Open Policy Agent Gatekeeper](https://github.com/open-policy-agent/gatekeeper) at the expense of some other interesting stuff - presented by folks from Google & Microsoft. It works through a set of CRDs so feels more config-y than code-y - although there is code buried in the templates you produce. This sounds like something we should get into - it's only in Alpha at the mo. They worked through a demo of four policy enforcement use-cases:

- all namespaces must have a label that lists a point of contact
- all pods must have an upper bound for resource usage
- all images must be from an approved repository
- Services must all have globally unique selectors

These sort of things sounds highly relevant to me. I know the team have tried it out a bit and been put off by the Rego language, but I understand that the tooling support for this is improving considerably, and as the popularity of OPA increases, examples of how to do this well will make our lives easier in this area anyway. They're definitely shooting for making it easier to write, test, automate into pipelines, and re-use - all good things I like to hear!

> Most importantly, I learnt that we should be pronouncing it "Oh-pa" not "Oh-Pea-Ay" ... +1 for more pronunciation guides. Now that's a valuable reflection right there.

---

## How Changing Kubernetes Itself Works

I went to a couple of talks focusing on some proposals for improving Kubernetes itself. As well as liking the content being proposed, it was interesting to learn a little bit about how the process of making changes to Kubernetes - what is a giant open source behemoth now - actually works.

### M-M-M-M-Multi-Cluster Ingress

The first of these was about ingresses from a couple of Googlers - Rohit Ramkumar & Bowei Du. I do enjoy a good load balancer talk, I do. Topics included:

- Should we even have ingresses at all?
- Assuming yes, how portable, extensible, featured should it be? The current one has been dumbed down a lot at the expense of features to ensure it's portable - at what cost? They actually got divisive responses from users when they surveyed on this topic too
- What are others doing - in particular Istio & Contour by the sounds of it

In the end, they are essentially proposing two new kinds - `MultiClusterIngress` and `MultiClusterService`. Both sound very sensible, and are modelled on existing primitives, but with the typical stuff you include nested within a `template:`, allowing for the extension to also include things like `cluster.select.matchLabels(region=eu)` or similar. There's also mention of being able to route to storage bucket (**do want!!**).

They are deliberately trying to minimise how opinionated they are - allowing ingresses to use Extended or Custom APIs on top, such that users can trade off portability for feature-set if they wish to.

It's about to become a KEP (Kubernetes Enhacement Proposal), so very much a "watch this space" for the inevitable debate - which they encourage. By the sounds of it they have working code to demonstrate it, it's just a case of letting the Kubernetes governance wheels turn ... which personally I'm quite interested in seeing how long it takes for these sort of big ticket things to hit prime time.

Which leads somewhat on to ...

### Have the cake, eat the cake - or, dynamic autoscaling

At the moment, Kubernetes is limited. No way! Well, about this specific thing it is. The `PodSpec` is immutable, which means if you want to dynamically set your resource requests or limits, you're SOL unless you're ok with restarting your pod. And who wants to do that? This is not so great when you have an app that has daily/weekly/seasonal patterns of use, steady growth of user base, varying app lifecycles, and so on.

> So, this talk was about, well, changing it to be mutable instead. Sounds simple enough, right?

Well apparently doing that sort of thing ain't so simple - there's loads of code in the Kubernetes Core (in particular, the Scheduler and Kubelet) that assume it is.

At the moment, the Vertical Pod Autoscaler will both recommend and (if you enable it) recreate your pods based on its sizing recommendations. And to be clear, for certain kind of workloads, being able to do this on-the-fly isn't going to fly - *I'm writing that while on a plane to, so funny!* - Java being the obvious example, as you're expected to set min/max memory args so Java knows what it's got to work with (and there's no way to dynamically change those at the same time).

> That's a shame for us, as we have a lot of JVM-based services (Kotlin, if you're curious).

But other stuff - this would be great for. And I'm sure at some point we'll rewrite everything in Golang at some point when Google's completely taken over the world. Or if they don't then maybe Rust.

Aaaaaanyway ... this stuff is all pretty cool, and there's a KEP for this as well out there already awaiting proposal. And then probably a shedload of testing.

There's also mention in the Q&A that VPA in general struggles with Java stuff. That'll be one to watch out for when we play around with it.

> If you're wondering why we aren't using it already - I think it has a lot to do with our workloads not being bound by the obvious stuff like CPU and memory. We'd need to spend some time working out the best (custom) metrics to use as the scaling measure, and we're just not at the scale where it's worth doing that yet.

---

## Intuit, Argo & The Custom Deployment Controller

I really want to give a special mention to Intuit (of finance/accounting software fame). They had a couple of their engineers (Danny & Alex) on stage talking about a customer controller they've created and open sourced *(to a big round of applause - lovely!)*.

So, why did they write their own controller? Well, they have a pure GitOps philosophy going on, and really wanted their deployments to follow that methodology - even when their development teams were making use of things like [Blue/Green](https://martinfowler.com/bliki/BlueGreenDeployment.html) and [Canary](https://martinfowler.com/bliki/CanaryRelease.html)  deployments. There's no built-in approach in Kubernetes to follow these sort of techniques.

Their story started with Jenkins, but they found it too brittle, painful to setup, had some security concerns, and ultimately just not feeling idempotent to match their GitOps philosophy. They iterated on this with some Deployment hooks, but this was still not what they were after and it became less transparent for engineering teams - and was still not truly declarative.

So, they busted out their own code for it. They liked the Controller pattern because it codifies the orchestration, can be built from the ground up to be GitOps-friendly, runs inside their cluster for easier identity management, and allowed them to make the migration path smooth for teams by iterating on the standard Kubernetes constructs (i.e. a Deployment), with just a couple of tweaks to their yaml.

6 months later, they launched this: https://github.com/argoproj/argo-rollouts. It takes care of things like ReplicaSet creation, scaling and deletion - just like a Deployment, with a single desired state described in the PodSpec and supports manual and automated promotions and the use of HPA.

They showed a couple of example yaml files to illustrate the difference from a `Deployment`, which are minor. Besides the obvious adjustments to the apiVersion & kind at the top, you add a `strategy` stanza, such as:

```yaml
strategy:
    blueGreen:
        activeService: active-svc
        previewService: preview-svc
        previewReplicaCount: 1 # optional
        autoPromotionSeconds: 30 # optional
        scaleDownDelaySeconds: 30 # optional
```

... for Blue/Green (note the pre-canned active & preview environments, and some toggles), or Canary:

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

... which feels a lot more like a Deployment with its more transient environment and no service modification.

They do a demo, which I have to say looks super-whizzy - they use ArgoCD for their tooling and it has some lovely visualisations. I really liked how it detected and displayed the changes to the Service/Pods in real-time - great for radiator dashboards (rather than `kubectl get pods -w`). There's also a traffic flow visualisation showing the traffic split between blue and green (... in blue or green arrows). Slick.

> As an aside, I popped by the [CodeFresh](https://codefresh.io/) booth while at KubeCon, which also had something similar

They've still got some work to do - for example they're interested in Service Meshes to avoid the need for a large number of replicas for Canary (the classic "I need 10 pods to do 10%") and other deployment strategies to support experimentation, A/B testing, etc.

All in all, really impressed with their talk and the things they've achieved here. Very cool!

---

## Kubernetes Failure Scenarios

One of the more technically deep sessions I went to was by Datadog (Laurent Bernaille & Robert Boll) of monitoring tools fame.

> They run a very large fleet (biggest cluster has 2k nodes, they average 1-1.5k) - for comparison ours are less than 100 ...

Their talk was on "surprising" failure scenarios - a top 10. I've summarised this down **massively**, if you're interested in the detail, be sure to check out their talk online:

1. **It's ~~never~~ always DNS**. Gets a round of applause! :grin:
2. **Image Stampede**. Pods in CrashLoopBackoff with `imagePullPolicy: Always` when your image registry sits behind some NAT instances can cause your whole registry to get blacklisted
3. **"I Can't kubectl"**. Crashing the apiserver by accidentally deleting a node selector so that a search was running over the entire cluster rather than a subset
4. **Pods not scheduling**. Tight limits on resource quotas being accidentally breached by a DaemonSet that was supposed to not have them set. Pod Eviction due to this being a critical PodPriority also!
5. **Log intake volume**. They basically DDoS'd themselves by enabling audit logging on a DaemonSet
6. **Where did my pods go?**. A subtle bug after enabling `HorizontalPodAutoscaler` when removing the hard-coded number of replicas - you also need to `kubectl edit` to remove an annotation at the time you remove this to avoid the Deployment reseting replica count to 1.
7. **There's a ghost in Cassandra**. Their cloud provider being a bit too smart - lack of capacity in a preferred AWS Availability Zone resulting in new nodes being scheduled in another, and then moved back when the capacity was freed up - detaching the pods from their AZ-bound disk.
8. **Slow deploy heartbeat**. A misconfigured CronJob which couldn't be placed on a node causing every loop of the scheduler to try to place it on one of 4,000 pods.
9. **"We expect containers to be contained"**. Various examples of dodgy container behaviour, including zombie threads to a readinessProbe with an exec being too slow for the timeout, one due to a low level disk blocking issue, and another due to slow slow pod startup from excessive I/O from DNS, audit logs and kubectl API calls. They use pools per app (with taints and tolerations) to truly contain their workloads. That's pretty scary to me!
10. **"Graceful" Termination**. Constraints caued by long-running jobs - beyond the length of kubelet restart for example, due to batch-y workloads.

Their key takeaways were:

1. Careful with DaemonSets in large clusters
2. DNS is hard
3. Cloud infra is not transparent
4. Containers not really contained
