---
title: "KubeCon - Service Meshes"
date: 2019-05-23T20:00:00-01:00
author: "@alexdmoss"
description: "A potential solution for dealing with your lovingly crafted and delightfully distributed service mess"
banner: "/images/sailboat.jpg"
tags: [ "KubeCon", "CNCF", "service mesh", "Istio", "Linkerd", "Google" ]
categories: [ "Conference", "Kubernetes", "CNCF", "Service Mesh", "Kubecon" ]
aliases:
- /2019/05/23/kubecon-2019-service-meshes/
---

{{< figure src="/images/sailboat.jpg?width=600px&classes=shadow" >}}

So, why service meshes? ~~Because I am architect, and architects are born to love service meshes.~~ Because actually I can totally see the value of a service mesh, I just can't quite seem to convince anyone that they're worth biting the bullet on yet.

It's a design pattern I'm totally sold on and just need to find the killer need for so we can start reaping some of those sweet, sweet benefits. Or crying at the complexity and brittle overlay we've subjected ourselves too. Maybe.

{{< figure src="/images/still-want-to-try.jpg?width=400px" title="This image is also eeriely similar to my old London commute" >}}

So ... will these talks help me out I wonder?

---

Firstly, I'd like to talk about [the keynote announcement](https://www.youtube.com/watch?v=gDLD8gyd7J8) of [SMI](https://smi-spec.io/) (Service Mesh Interface). This is hopefully going to prove to be highly relevant, as we aren't running a service mesh in my team currently but I'm convinced we will be in not too long, so taking some of the pressure off making the right decision is helpful.

The goal with SMI is ([I'm paraphrasing massively](https://cloudblogs.microsoft.com/opensource/2019/05/21/service-mesh-interface-smi-release/)):

- common portable APIs across different service mesh technologies - works with the three bigger players (Istio, Linkerd, Consul)
- apps, tools, the wider ecosystem will all be able to integrate through these standardised APIs

It is being done to help folks get started quickly - i.e. choice is believed to be a barrier, simpler is better, tapping into the ecosystem is good

> My personal opinion on this is focused more around the "simpler is better" comment. I think the barrier for us when it comes to meshes is the complexity it introduces - we would be able to make a choice, and we would trust the wider ecosystem to adapt, but we (and in particular our wider developer community) are ourselves still only just getting our heads round Kubernetes itself, if we're being honest about it

It was a brief segment in the end-of-day keynote, so I would imagine a "watch this space" situation. The video demo was snazzy though if you want to check it out.

---

The first [deeper-dive talk](https://www.youtube.com/watch?v=-t2BfT59zJA) I went to covered Istio multi-cluster gateway networking - it had a demo involving running Istio's BookInfo sample app across three "cloud" providers - some of its services were spread onto GKE, some in AKS, and some on a bit of bare metal in the guy's house!

The demo was great, and my takeaway from this talk was more-or-less that you can quite easily join up Istio with its Gateways, but you need to create lots of `ServiceEntry` resources for your services to do this which sucks ... so they're going to do some work to generate this automatically for you. Thanks!

However, I followed this one up with a less techie but more customer-centric talk on Istio, which was particularly interesting ...

---

First though, some context, I've been playing a little with Istio in GKE. Nothing major, just trying out the new GKE Addon for it to see if it can solve some problems we have around where to run "interactive computing" workloads - in other words, workbenches for data scientists.

That topic is a blog post in its own right, but the important point for context is that I found the addon in GKE a pretty troublesome experience. While I got it working eventually, I needed to sacrifice some of our good security things along the way - it hated `PodSecPolicy` in a big way - and I wasn't left feeling too impressed ...

---

... so [a talk from some folks](https://www.youtube.com/watch?v=vQ2IktsMlgQ) who were also battling with Istio was pretty interesting. These guys (David Meder-Marouelli, Stephan Fudeus) work for 1&1 Mail & Media, a big German ISP / Collab+Productivity tools provider. They have a load of legacy tin scattered around the place as well as a shiny new GKE cluster they like, and they were doing mesh expansion things to join their new Kubernetes workloads with their bare metal legacy. Daunting.

> This was a super-popular session - the most oversubscribed one I attended - standing room only!

These guys really want those sweet service mesh features - as they put it, to "deal with their service mess of microservices". Might steal that :wink:

They also referenced a post on Twitter (don't know the source) - suggesting that folks love to start with [Linkerd](https://linkerd.io) and see if it works, then switch to Istio when you need all the features and are willing to spend 9 months on getting it working. Ha, love it!

So, these guys wanted to skip that cycle and went straight for Istio (v1.0.x, like me with my experiments).

Skipping to the crux of it, they were basically having the same security context misery that I was - it needs to run as root, needs a writable root f/s, needs NET_ADMIN, high privileges for serviceaccounts, and so on. They also experienced pain in the networking requirements which they should totally see a doctor about - but I think this was a lot to do with their setup (the mesh expansion and restrictions on not making their pods routable).

Intriguingly, they also observed problematic ordering of automatic sidecar injection vs PodSecPolicy evaluation - this was fascinating for me as I suspect this was exactly the same weird and wonderful behaviour I was seeing, I just never got to the bottom of it! Sometimes the PSP was being applied before the admission controller does its thing, sometimes after :open_mouth: Major confusion, breaking  deployments that were working, fun for the whole family.

{{< figure src="/images/security-violations.jpg?width=600px" >}}

They did note that many of these things are improved in Istio 1.1 apparently, and I actually had a chat with Google later on that suggested things are going to be even better in Istio 1.2 for our particular painpoints, so that's reassuring - see the bottom of this blog post ...

---

Somewhat phased by my Istio conversations, I got a bit rebellious. We may be running on Google Cloud but I thought to myself *"I don't need no Google solution, bring on Linkerd!"*

So I went to a [deep-dive talk from their CTO](https://www.youtube.com/watch?v=E-zuggDfv0A). Who, by the way, seemed like a very nice chap and it was really refreshing to hear from someone clearly so passionate about the product he'd come up with.

A little history - Linkerd has been around for a while, relatively speaking. It started in Feb 2016 and had some sort of Mesosphere relationship going on. It joined CNCF in 2017 and went through a product known as Conduit (from the company Buoyant, whom the speaker - [Oliver Gould](https://twitter.com/olix0r?lang=en) - is CTO) before by the sounds of it being reborn as [Linkerd2](https://blog.linkerd.io/2018/09/18/announcing-linkerd-2-0/) - now with even more Kubernetes.

> Or, to be more accurate, "very coupled" to Kubernetes. Good move!

Why the Kubernetes love-in (apart from the obvious?). Kubernetes has Pods. Pods = sidecar pattern. Sidecars + a service mesh = :heart:

It's goal is to be lightweight, simple, and focused on the essentials. That sounds good. It also has some big users - from a retail perspective the most eye-catching is probably eBay. Who, coincidentally, are I think a GCP shop.

His view? Service Meshes are too complicated. Kubernetes is hard enough, we want simple. I hear ya bud!

The speaker immediately endears himself to me because he's got photos of scraps of paper with architectural diagrams scribbled on them. I like the style. What follows is a deeper dive into how it works - without going into loads of details, it's broadly similar in construct to the Mesh-That-Shall-Not-Be-Named. It's got some mTLS in there, its proxy is written in Rust and its control plane in Go.

It doesn't have a built-in tracing capability because it doesn't believe that code should have to be changed (instrumented) to use a mesh. Instead it does a load of metadata injection by the sounds of it. Its UI is very visual and very appealing. There's a lot of focus on telemetry in here. Apparently a common pattern is to let linkerd use the prometheus it deploys, and have your own Prometheus scrape from it.

There's clearly some super-clever stuff going on within it. It upgrades connections to HTTP/2 to get a single connection between proxies to drive efficiency (especially for the TLS handshaking now). It's got a "Tap Server" which allows for on-demand tracing, e.g. "In this namespace, show me the requests that are GETs to this service", which returns sampled metadata (full headers and payloads coming soon - just need to sort out some RBAC to prevent you being able to see payloads you shouldn't be able to!). It was also early to get in on the SMI spec (mesh standardisation APIs) that was announced in the keynotes yesterday.

Some of the questions that follow show that the project does have gaps in comparison to "that other" service mesh. Some of those gaps are, by the sounds of it, quite deliberate - "ingresses aren't part of a mesh, they get traffic into the mesh", "auth is an ingress problem", "meshes can't add a lot of value for non-HTTP traffic so why do you want it exactly?"

All in all I'm left with a "this would be cool to try" impression from the whole thing. Maybe one of our Friyay's (days to experiment, off-ticket).

I have a gut feeling that Istio is going to be the eventual answer though.

---

And on that note ... I did get a chance to chat with some folks from Google on the Istio roadmap and some of my ~~misery~~ experiences with Istio lately.

In summary, I was left with the distinct impression that the GKE Addon was not the way to go. Why? Because it sounds like Istio 1.2 is going to bring with it a [lovely](https://github.com/istio/istio/issues/9333) [Operator](https://discuss.istio.io/t/istio-operator-plans-for-1-2/2227). Operators are the best. Banzai Cloud even have an [early working version of it](https://github.com/banzaicloud/istio-operator).

The operator should make deployments altogether more flexible, composable, and all-round easier. You'll be able to run two of them to help you manage upgrades - or indeed use Google Traffic Director as a managed implementation of this control plane. It also sounds like things are getting simplified with the removal of Mixer and such. The GKE Addon is very much *not* configurable in any way, which makes it feel really hard to interact with or manage when it goes wrong.

This is pretty fab if it does come to pass. I feel reassured that we can wait a little longer - Istio is currently on a quarterly cadence and 1.1 is out already. So, fingers crossed!

---
