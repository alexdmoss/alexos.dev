---
title: "Some Musings From KubeCon EU 2021 (Virtual Edition)"
panelTitle: "Musings From KubeCon EU 2021"
date: 2021-05-12T18:48:00+01:00
author: "@alexdmoss"
description: "Thoughts and scribbles after attending a whole bunch of talks at KubeCon EU 2021"
banner: "/images/conference-chairs.jpg"
tags: [ "KubeCon", "CNCF", "opinion", "reflections" ]
---

I've just finished binging on a whole host of talks from KubeCon EU 2021, and thought I'd scribble down some of my immediate thoughts - from the talks I "went" to, as well as more generally about virtual conferencing these days. Which is a bit of a change of pace, let's be honest.

{{< figure src="/images/conference-chairs.jpg?width=800px&classes=shadow" attr="Photo by Jonas Jacobsson on Unsplash" attrlink="https://unsplash.com/photos/2xaF4TbjXT0" >}}

I should caveat this post with the obvious - my experiences are heavily based around how _I_ chose to take in the conference - the sessions I picked, the things I got involved with and the things I didn't. Conferences for me tend to be focused on learning about or picking up tips on bits of technology or techniques that I already have some idea about in the first place, rather than totally new things that I don't think I'm going to have any practical use for in my current line of work. I'm also not that great at the "networking" thing :sweat_smile:

Anyway, let's get into it. Want to jump to a section? Here's how this is going to go down ...

1. [Overall Thoughts](#overall-thoughts) on the conference content
2. [Technical Sessions](#technical-sessions) - some more specific thoughts around the main themes of the talks I went for
3. [Virtual Conferences](#virtual-conferences) and whether they're worth the money

---

## Overall Thoughts

I managed to take in quite a large number of talks (26!) - more than I would if I were there in person (I tend to top out at about 6-7 per day in the flesh, and I think that's generally considered quite a lot!). Some of the themes that came out for me were:

- It was good to have a sense that at work we are on the right path - many of the things I heard about were things we are aware of and either already doing or looking to adopt. It feels like we're on the curve now, rather than trailblazing or (perish the thought) lagging behind. It's nice to have the reassurance :relieved:
- It's not really **Kube**Con any more, it feels a lot more like the rest of its name, i.e. CNCFCon (which happens to include Kubernetes). I guess this is a reflection that Kubernetes is kinda the slam-dunk on which the whole CNCF thing is based around. I went to a couple of talks that were about the core features of kube, but frankly not that many. I'm not sure if that's a reflection of the fact that it's seen more as a given these days, or perhaps that it's changing more rarely _(probably not the latter, if Kubernetes release notes lately are anything to go by!)_.
- With the above in mind, it's nice to see some of the newer sandbox/incubator projects getting more air time (even in Keynotes). I haven't had this vibe before, where the feel was more of the already big-hitters. The example I have in my head would be something like [KEDA](https://keda.sh).
- By the end, I swear if someone mentioned **GitOps** one more time I'm going to ... :angry:. Seriously, this was definitely the buzzword bingo winner of the talks I ended up going to at least. Perhaps I'm going to notice this more because I really haven't drunk that particular Kool-Aid - to me GitOps is just a pattern of Continuous Deployment and nothing fancier than that (especially when it's push-based), and how you choose to do it is not as important as the fact you are doing it at all.

{{< figure src="/images/yo-dawg-gitops.jpg?width=600px" >}}

---

## Technical Sessions

**Multi-Tenancy and Multi-Cluster**. Little bit of a hot topic at work - we're exploring options to run multiple copies of our large multi-tenant cluster. I took in a few sessions around this topic, enough to make me feel like quite a few others are looking into it too. There was a decent talk explaining some of the thinking going on in SIG multi-cluster and why solving some of these problems is hard (TL:DR - broad range of expectations / use-cases is hard to design for). It was interesting hearing them talk about some of the mistakes of the past, avoiding the rush/temptation to build something too early but also follow the sensible approach of slicing off useful functionality that can be effectively described and implemented, rather than solving the whole thing in one go (hello Agile!).

> There was also a great related talk about the new [Gateway API](https://kubernetes.io/blog/2021/04/22/evolving-kubernetes-networking-with-the-gateway-api/) which ticked a lot of boxes for me, solving some meaningful problems like control of traffic distribution and making the separation between Product Developers and Platform Developers more explicit and managable. Me likey :thumbs_up:

**Sidecars**. Came up a few times. One talk in particular talked about a company with 7+ "standard" sidecars across their fleet, which feels like overkill. Sure, they were solving a lot of operational concerns so the separation from the "business logic" makes a lot of sense, but really, 7? :open_mouth:. That seems like a bit of an organisational problem (i.e. unnecessary separation of ownership, but I am speculating). Surely you'd start looking at service meshes (which solved a lot of the ones they listed) or different patterns to get the result you need? We seem to be doing okay with one sidecar and occasionally a second, plus a few DaemonSets.

{{< figure src="/images/sidecar.jpg?width=400px&classes=border" attrlink="https://devrant.com/rants/3775630/saw-this-on-linkedin" >}}

**Service Mesh**. Most of these talks were naturally Linkerd-focused as a CNCF project, so less directly useful "this is how to do X" given I am about to embark on an Istio journey. Still, pretty interesting to see some of the use-cases, and frankly a bit of reinforcement of my confirmation bias before we start on a technically complex stream of work to get it rolled out. In these talks I also picked up a couple of interesting definitions of what cloud-native could actually mean:

> "Cloud Native means all your resources and dependencies are network attached"

Which is a riff on this quote from Duncan Winn (Director of SRE at Google Cloud), which I rather like to be honest!

> "Cloud Native is a term describing software designed to run at scale reliably and predictably, on top of potentially unreliable cloud-based infrastructure"

**Security & Networking**. I found myself somewhat unexpectedly going to a few security & networking related sessions. Always a topic worth staying on top of, and newer/better tools for helping monitoring for security issues is of particular interest lately. There were a few sessions advocating the benefits of eBPF and Cilium and some of the awesome tooling you can build on top of that. It did look rather handy. Interestingly this is also the technology behind GKE's [dataplane v2](https://cloud.google.com/kubernetes-engine/docs/concepts/dataplane-v2) which is new but I have my eye on as a shop running on GKE already. Coincidentally, that technology went GA [a few days ago](https://cloud.google.com/blog/products/containers-kubernetes/bringing-ebpf-and-cilium-to-google-kubernetes-engine).

**Autoscaling**. After some recent interesting challenges at work that can be summed up as "how can I make Kubernetes autoscale fast enough to deal with very sudden bursts of traffic" :chart_with_upwards_trend:, I went to a few talks on this to see if there were any new tricks. The answer to that seems to be no, but this in of itself was useful to know, and I learnt about some of the ways other companies have dealt with the problem. One presentation (from StockX) talked through an example we've talked about at work before - the idea of a proactive scaler that can be told about upcoming activities (like a marketing push or new product launch) and upscale based on preconfigured known-good configuration (and bring it back down again). That this works for other companies is interesting - but for us, getting to the point of knowing its coming and telling someone about it is the problem statement :pensive:

> Also on the autoscaling point - great to see decent coverage for [KEDA](https://keda.sh), as we've recently started using this ourselves after all sorts of problems with our previous custom-metrics autoscaling. So far our experience has been positive, and I really like the potential for scale-to-zero in Kubernetes too

This wasn't all I covered - I also took in sessions on chaos engineering tools like Litmus, as well as Jaeger, Falco, Gatekeeper, and some positive reinforcement about recent decisions we've taken on how to scale up our telemetry stack. Some of those could turn into blog posts in their own right later!

---

## Virtual Conferences

As mentioned towards the top, one of the things I find with Virtual Conferences is I'm able to take in more talks overall. Whilst they're all recorded anyway when you go in person, I find that the ability to dip in and out of sessions to see if they're something I'm going to get value out of is really valuable. Sometimes for me the intro is too basic when I'm already familiar with the product so being able to skip along to the new stuff or demo suits me. Other times it goes into more detail than I can handle and I can bow out without the need to awkwardly shuffle past people or distract the presenter :persevere:

> In person I'm way too polite to just up and leave, even if the talk isn't at all what I was hoping for - although I probably should!

Where virtual conferences loose out for me - and this is despite what is clearly a great amount of effort put in by the organisers - is the social interaction. For me personally, I'm not a "chat to random people about shared interests in the corridor" kind of person, although after a few sweeps around the Expo Hall I will start stopping in on booths for vendors/products I'm interested in. Perhaps for future conferences I'll use the pre-recorded vendor presentations as a hook to jump onto Slack to ask a question, but certainly for this KubeCon and the other virtual conferences I attended last year, I ended up just skipping out on this part of the conference experience, more or less. This is similar with some of the other trimmings - there were loads of "virtual events" and wellbeing stuff on the event platform this time round (games, trivia, etc.) - if it helps others, I'm all for it of course but this just isn't my bag.

So was it worth the money? Well objectively you'd probably expect to say no - the talks are online for free a few days later (I'm writing this on the Wednesday after the conference - the talks are going to be on YouTube on Friday). Except for me, it actually is. Why's that?

1. It was only a tenner for the early bird. That's three posh coffees :coffee:, if they were actually still open. I don't know how many coffees it would take to make it not feel worthwhile, but this wasn't a hard sell for me.
2. The platform used for these has gotten a lot better (certainly in comparison to KubeCon last year at least). It's much easier to navigate, find similar courses, catch up where you left off, see the Q&A and so on. With some effort I'm sure this could be replicated from YouTube playlists and the conference website, but this is a handy convenience in my opinion. There were even prompts to help you reach out to the presenters when they were online, which I didn't use but is a nice touch.
3. The most important reason though - it forces me to actually get some value :money_with_wings:. This is the really powerful one for me. Think like gym membership - coz you know you're paying for it, and it's over a certain time window, you're motivated to get your value out of it. In practice this means I politely declined or rescheduled work meetings to take in the talks, I reviewed what I'd covered at the end of each day to help me shape the talks I'd hit the following day, and so on. The same sort of things I'd do if I were there in person.

{{< figure src="/images/gym-membership.png?width=600px&classes=border" >}}

In other words, it creates a great excuse to carve out the time for learning :mortar_board:. As long as the price of admittance remains low, I think it's well worth the few days a couple of times a year to get the value out of it.
