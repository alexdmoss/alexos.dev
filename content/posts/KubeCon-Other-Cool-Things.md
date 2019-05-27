---
title: "KubeCon - Other Cool Things"
date: 2019-05-23T23:00:00-01:00
author: "@alexdmoss"
description: "Some of the other tool things I learnt at KubeCon that don't fit neatly into other posts!"
banner: "/images/kube-donuts.jpg"
tags: [ "KubeCon", "CNCF", "opinion", "reflections" ]
categories: [ "Conference", "Kubernetes", "CNCF", "Kubecon" ]
aliases:
- /2019/05/23/kubecon-2019-other-cool-things/
---

{{< figure src="/images/kube-donuts.jpg?width=600px&classes=shadow" title="Yes, this was a wall of donuts to celebrate Kubernetes 5th birthday!" >}}

I'm going to use this blog post to talk about a couple of other things I learnt about while at KubeCon 2019. These don't slow neatly into [my](/2019/05/23/kubecon-the-keynotes/) [other](/2019/05/23/kubecon-service-meshes/) [posts](/2019/05/23/kubecon-observability/) from my time at the conference, but I still felt they were interesting enough to call out separately.

---

## How Changing Kubernetes Itself Works

I went to a couple of talks focusing on some proposals for improving Kubernetes itself. As well as liking the content being proposed, it was interesting to learn a little bit about how the process of making changes to Kubernetes - what is a giant open source behemoth now - actually works.

### M-M-M-M-Multi-Cluster Ingress

The [first of these](https://www.youtube.com/watch?v=Ne9UJL6irXY) was about ingresses from a couple of Googlers - Rohit Ramkumar & Bowei Du. I do enjoy a good load balancer talk, I do. Topics included:

- Should we even have ingresses at all?
- Assuming yes, how portable, extensible, featured should it be? The current one has been dumbed down to ensure it's portable - at what cost? They actually got divisive responses from users when they surveyed on this topic too
- What are others doing - in particular Istio & Contour by the sounds of it

In the end, they are essentially proposing two new kinds - `MultiClusterIngress` and `MultiClusterService`. Both sound very sensible, and are modelled on existing primitives, but with the typical stuff you include nested within a `template`, allowing for the extension to also include things like `cluster.select.matchLabels(region=eu)` or similar. There's also mention of being able to route to storage bucket (**do want!!**).

They are deliberately trying to minimise how opinionated they are - allowing ingresses to use Extended or Custom APIs on top, such that users can trade off portability for feature-set if they wish to.

It's early days so may take some time to hit prime-time, but it was very interesting to learn about their thought processes and the direction things are (probably) heading - it's still a draft for a proposal at present.

### Dynamic Autoscaling (and Cake)

At the moment, Kubernetes is limited. No way! Well, about this specific thing it is. The `PodSpec` is immutable, which means if you want to dynamically set your resource requests or limits, you're SOL unless you're ok with restarting your pod. And who wants to do that? This is not so great when you have an app that has daily/weekly/seasonal patterns of use, steady growth of user base, varying app lifecycles, and so on.

> So, this talk was about, well, umm, changing it to be mutable instead. Sounds simple enough, right?

Well [apparently doing that sort of thing ain't so simple](https://www.youtube.com/watch?v=58uRFofXUyw) - there's loads of code in the Kubernetes Core (in particular, the Scheduler and Kubelet) that assume it is.

At the moment, the Vertical Pod Autoscaler will both recommend and (if you enable it) recreate your pods based on its sizing recommendations. This change would allow it to do this without restarting them. For certain kinds of pods (not Java, /sigh). Still, we have other stuff too that could benefit from this. And I'm sure at some point we'll rewrite everything in Golang when Google's completely taken over the world. Or if they don't, then maybe Rust?

> If you're wondering why we aren't using it already - I think it has a lot to do with our workloads not being bound by the obvious stuff like CPU and memory. We'd need to spend some time working out the best (custom) metrics to use as the scaling measure, and we're just not at the scale where it's worth doing that yet.

---

## Intuit, Argo & The Custom Deployment Controller

I really want to give a special mention to Intuit (of finance/accounting software fame). They had a couple of [their engineers (Danny & Alex) on stage talking about](https://www.youtube.com/watch?v=yeVkTTO9nOA) a customer controller they've created and open sourced *(to a big round of applause - lovely!)*.

So, why did they write their own controller? Well, they have a pure GitOps philosophy going on, and really wanted their deployments to follow that methodology - even when their development teams were making use of things like [Blue/Green](https://martinfowler.com/bliki/BlueGreenDeployment.html) and [Canary](https://martinfowler.com/bliki/CanaryRelease.html)  deployments. There's no built-in approach in Kubernetes to follow these sort of techniques.

Their story started with Jenkins (long list of problems), via some Deployment hooks (not transparent for engineering teams) - neither of which were truly declarative for these style of deployments. So, they did it themselves. They liked the Controller pattern because it codifies the orchestration, can be built from the ground up to match their approach, runs inside their cluster for easier identity management, and allowed them to make the migration path smooth for teams by iterating on the standard Kubernetes constructs (i.e. a `Deployment`), with just a couple of tweaks to their yaml.

6 months later, they launched this: https://github.com/argoproj/argo-rollouts. It takes care of things like ReplicaSet creation, scaling and deletion - just like a Deployment, with a single desired state described in the PodSpec, supporting manual and automated promotions and the use of HPA.

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

{{< figure src="/images/argocd-ui.gif?width=600px&classes=shadow" title="Simpler version of the Argo CD UI" >}}

<!-- https://github.com/argoproj/argo-cd -->

> As an aside, I popped by the [CodeFresh](https://codefresh.io/) booth while at KubeCon, which also had something similar

They've still got some work to do - for example they're interested in Service Meshes to avoid the need for a large number of replicas for Canary (the classic "I need 10 pods to do 10%") and other deployment strategies to support experimentation, A/B testing, etc.

All in all, really impressed with their talk and the things they've achieved here. Very cool!

---

## Compliance as Code

In some of the specialist tracks I went to, I learnt several other interesting things too.

For example, I started the day learning about some proposals around a Multi Cluster Ingress - Ingress v2 basically *(yes please, looked great!)*, an [intro to SPIFFE and SPIRE](https://www.youtube.com/watch?v=Rx6PMptyEtg) *(which sounded useful but a bit early-days for folks like us, perhaps)*, some wildly [sensible folks from Lyft & Triller](https://www.youtube.com/watch?v=TZ73EBP2a9Q) talking about permissions options in Kubernetes *(reassuringly, we are doing a lot of things in this space already - I did learn some funky things about Admission Controllers and Open Policy Agent though)*.

The latter inspired me sufficiently to go to [a more thorough talk](https://www.youtube.com/watch?v=Yup1FUc2Qn0) on [Open Policy Agent Gatekeeper](https://github.com/open-policy-agent/gatekeeper) at the expense of some other interesting stuff - presented by folks from Google & Microsoft. It works through a set of CRDs so feels more config-y than code-y - although there is code buried in the templates you produce. This sounds like something we should get into - it's only in Alpha at the mo. They worked through a demo of four policy enforcement use-cases:

- All namespaces must have a label that lists a point of contact
- All pods must have an upper bound for resource usage
- All images must be from an approved repository
- Services must all have globally unique selectors

These sort of things sounds highly relevant to me. I know our team have tried it out a bit and been put off by the [Rego language](https://www.openpolicyagent.org/docs/latest/how-do-i-write-policies), but I understand that the tooling support for this is improving considerably, and as the popularity of OPA increases, examples of how to do this well will make our lives easier in this area anyway. They're definitely shooting for making it easier to write, test, automate into pipelines, and re-use - all good things I like to hear!

> Most importantly, I learnt that we should be pronouncing it "Oh-pa" not "Oh-Pea-Ay" ... +1 for more pronunciation guides. Now that's a valuable reflection right there.

{{< figure src="/images/opa.png?width=600px&classes=shadow,border" title="OPA Architecture, from their website" >}}

---

## Kubernetes Failure Scenarios

One of the more technically deep sessions I went to [was by Datadog](https://www.youtube.com/watch?v=QKI-JRs2RIE) (Laurent Bernaille & Robert Boll) of monitoring tools fame. Their talk was on a top 10 of "surprising" Kubernetes failure scenarios.

{{< figure src="/images/failure.jpg?width=600px&classes=shadow" title="Uh-oh, need coffee. Photo by Nathan Dumlao on Unsplash" >}}

<!-- https://unsplash.com/photos/aZ9X3L1Va2Y -->

> They run a very large fleet (biggest cluster has 2k nodes, they average 1-1.5k) - for comparison ours are less than 100 ...

 I've summarised this down **massively**, if you're interested in the detail, be sure to check out their talk online:

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
3. Cloud infrastructure is not transparent
4. Containers are not really contained

---

## All The Sessions I Went To

If you're curious what sessions I went to, here's a full list, with video links:

[Playlist with every single session from KubeCon](https://www.youtube.com/playlist?list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3
) - see you next year!

### Tuesday

- [Keynote: Stitching Things Together](https://www.youtube.com/watch?v=lmGFgZ889kY&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=19&t=0s) – Dan Kohn, Executive Director, Cloud Native Computing Foundation
- [Keynote: 2.66 Million](https://www.youtube.com/watch?v=w62T1SN4g6Y&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=20&t=0s) - Cheryl Hung, Director of Ecosystem, Cloud Native Computing Foundation
- [Keynote: CNCF Project Update](https://www.youtube.com/watch?v=vdxcaR3I2ic&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=21&t=0s) - Bryan Liles, Senior Staff Engineer, VMware
- [Sponsored Keynote: Network, Please Evolve](https://www.youtube.com/watch?v=KmCfIQFllOM&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=22&t=0s) – Vijoy Pandey, VP/CTO Cloud, Cisco
- [Ingress V2 and Multicluster Services](https://www.youtube.com/watch?v=Ne9UJL6irXY&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=331&t=0s) - Rohit Ramkumar & Bowei Du
- [Intro: SPIFFE](https://www.youtube.com/watch?v=Rx6PMptyEtg&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=36&t=0s) - Emiliano Bernbaum & Scott Emmons, Scytale
- [Fine-Grained Permissions in Kubernetes: What’s Missing, and How to Fix That](https://www.youtube.com/watch?v=TZ73EBP2a9Q&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=329&t=0s) - Vallery Lancey, Lyft & Seth McCombs, Triller
- [Intro: Open Policy Agent](https://www.youtube.com/watch?v=Yup1FUc2Qn0&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=61&t=0s) - Rita Zhang, Microsoft & Max Smythe, Google
- [The Multicluster Toolbox](https://www.youtube.com/watch?v=Fv2PKKDgjIQ&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=33&t=0s) - Adrien Trouillaud, Admiralty
- [Keynote: Welcome Remarks](https://www.youtube.com/watch?v=Npgx6g3Fbds&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=106&t=0s) - Janet Kuo, Software Engineer, Google
- [Sponsored Keynote: Democratizing Service Mesh on Kubernetes](https://www.youtube.com/watch?v=gDLD8gyd7J8&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=105&t=0s) - Gabe Monroy, Microsoft
- [Keynote: Kubernetes Project Update](https://www.youtube.com/watch?v=jISu86XmkHE&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=104&t=0s) - Janet Kuo, Software Engineer, Google
- [Sponsored Keynote: Recursive Kubernetes: Cluster API and Clusters as Cattle](https://www.youtube.com/watch?v=OXSRfl8mYyo&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=103&t=0s) - Joe Beda, Principal Engineer, VMware
- [Keynote: Reperforming a Nobel Prize Discovery on Kubernetes](https://www.youtube.com/watch?v=CTfp2woVEkA&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=102&t=0s) - Ricardo Rocha, Computing Engineer & Lukas Heinrich, Physicist, CERN
- [Sponsored Keynote: Expanding the Kubernetes Operator Community](https://www.youtube.com/watch?v=KPOEnFwspiY&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=101&t=0s) - Rob Szumski, Principal Product Manager for OpenShift, Red Hat

### Wednesday

- [Keynote: Opening Remarks](https://www.youtube.com/watch?v=5IvT80d8YVU&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=115&t=0s) - Bryan Liles, Senior Staff Engineer, VMware
- [Keynote: How Spotify Accidentally Deleted All its Kube Clusters with No User Impact](https://www.youtube.com/watch?v=ix0Tw8uinWs&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=114&t=0s) - David Xia, Infrastructure Engineer, Spotify
- [Sponsored Keynote: Building a Bigger Tent: Cloud Native, Cultural Change and Complexity](https://www.youtube.com/watch?v=hi5jXcauQE4&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=113&t=0s) - Bob Quillin, VP Developer Relations, Oracle Cloud
- [Keynote: A Journey to a Centralized, Globally Distributed Platform](https://www.youtube.com/watch?v=D7pbISekc8g&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=112&t=0s) – Katie Gamanji, Cloud Platform Engineer, Condé Nast International
- [Sponsored Keynote: What I Learned Running 10,000+ Kubernetes Clusters](https://www.youtube.com/watch?v=HXF0QzxUBTw&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=111&t=0s) - Jason McGee, IBM Fellow, IBM
- [Keynote: Debunking the Myth: Kubernetes Storage is Hard](https://www.youtube.com/watch?v=169w6QlWhmo&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=110&t=0s) - Saad Ali, Senior Software Engineer, Google
- [Keynote: Closing Remarks](https://www.youtube.com/watch?v=w3wN0PHwgUo&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=109&t=0s) - Bryan Liles, Senior Staff Engineer, VMware
- [Monitoring at Planet Scale for Everyone](https://www.youtube.com/watch?v=EFutyuIpFXQ&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=200&t=0s) - Rob Skillington, Uber
- [Resize Your Pods w/o Disruptions aka How to Have a Cake and Eat a Cake](https://www.youtube.com/watch?v=58uRFofXUyw&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=147&t=0s) - Karol Gołąb & Beata Skiba, Google
- [Benefits of a Service Mesh When Integrating Kubernetes with Legacy Services](https://www.youtube.com/watch?v=vQ2IktsMlgQ&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=176&t=0s) - Stephan Fudeus & David Meder-Marouelli, 1&1 Mail & Media
- [Reinventing Networking: A Deep Dive into Istio's Multicluster Gateways](https://www.youtube.com/watch?v=-t2BfT59zJA&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=191&t=0s) - Steve Dake, Independent
- [Deep Dive: Cortex](https://www.youtube.com/watch?v=mYyFT4ChHio&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=162&t=0s) - Tom Wilkie, Grafana Labs & Bryan Boreham, Weaveworks
- [10 Ways to Shoot Yourself in the Foot with Kubernetes, #9 Will Surprise You](https://www.youtube.com/watch?v=QKI-JRs2RIE&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=195&t=0s)- Laurent Bernaille & Robert Boll, Datadog

### Thursday

- [Keynote: Opening Remarks](https://www.youtube.com/watch?v=VljLVMMtSLk&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=327&t=0s) - Janet Kuo, Software Engineer, Google
- [Keynote: Kubernetes - Don't Stop Believin'](https://www.youtube.com/watch?v=Rbe0eNXqCoA&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=326&t=0s) – Bryan Liles, Senior Staff Engineer, VMware
- [Keynote: From COBOL to Kubernetes: A 250 Year Old Bank's Cloud-Native Journey - Laura Rehorst, Product Owner](https://www.youtube.com/watch?v=uRvKGZ_fDPU&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=325&t=0s) - Stratus Platform, ABN AMRO Bank NV & Mike Ryan, DevOps Consultant, backtothelab.io
- [Keynote: Metrics, Logs & Traces; What Does the Future Hold for Observability?](https://www.youtube.com/watch?v=MkSdvPdS1oA&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=324&t=0s) - Tom Wilkie, VP Product, Grafana Labs & Frederic Branczyk, Software Engineer, Red Hat
- [Keynote: Closing Remarks](https://www.youtube.com/watch?v=5SToNEA9vgk&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=323&t=0s) - Bryan Liles, Senior Staff Engineer, VMware & Janet Kuo, Software Engineer, Google
- [DIY Pen-Testing for Your Kubernetes Cluster](https://www.youtube.com/watch?v=fVqCAUJiIn0&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=240&t=0s) - Liz Rice, Aqua Security
- [Grafana Loki: Like Prometheus, but for Logs](https://www.youtube.com/watch?v=CQiawXlgabQ&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=245&t=0s) - Tom Wilkie, Grafana Labs
- [Deep Dive: Linkerd](https://www.youtube.com/watch?v=E-zuggDfv0A&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=256&t=0s) - Oliver Gould, Buoyant
- [How Intuit Does Canary and Blue Green Deployments with a K8s Controller](https://www.youtube.com/watch?v=yeVkTTO9nOA&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=315&t=0s) - Daniel Thomson & Alex Matyushentsev, Intuit
- [Delivering Serverless Experience on Kubernetes: Beyond Web Applications](https://www.youtube.com/watch?v=VCGlGlBdr-o&list=PLj6h78yzYM2PpmMAnvpvsnR4c27wJePh3&index=259&t=0s) - Alex Glikson, Carnegie Mellon University
