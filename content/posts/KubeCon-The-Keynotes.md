---
title: "KubeCon - The Keynotes"
date: 2019-05-23T09:00:00-01:00
author: "@alexdmoss"
description: "My summary & thoughts from some of the keynote presentations whilst over at KubeCon 2019 in Barcelona"
banner: "/images/opening-keynote.jpg"
tags: [ "KubeCon", "CNCF" ]
categories: [ "Conference", "Kubernetes", "CNCF", "Kubecon" ]
aliases: 
- /2019/05/21/kubecon-2019-opening-keynote/
- /2019/05/21/kubecon-2019-keynote-1/
---

{{< figure src="/images/opening-keynote.jpg?width=600px&classes=shadow" >}}

---

## Intro by the CNCF Director

[Dan Kohn's](https://www.dankohn.com/) talk starts with a screenshot of the game Civ VI - *who doesn't want to see a turn-based strategy game reference in a talk to nearly 8,000 ~~IT professionals~~ geeks* - he is emphasising the idea that in the game you could not train Knights without first discovering Stirrups.

{{< figure src="/images/civ.jpg?width=600px&classes=shadow" title="Civ VI Tech Tree" >}}

He builds on this with further "real" examples of two people independently discovering things - he uses the examples of Facebook *(watch the film the Social Network if you're not sure about this reference!)*, [Natural Evolution](https://www.npr.org/2013/04/30/177781424/he-helped-discover-evolution-and-then-became-extinct?t=1558564125403), [Calculus](https://en.wikipedia.org/wiki/History_of_calculus) - the idea being that simultaneous invention is not an uncommon thing, it is triggered by a spark of creativity and - crucially for his point to land - the pre-requisite things need to be in place (such as computers and the internet, in the case of Facebook).

And now we get to the crux of it - web-scale platform implementations. It wasn't always just Kubernetes you know! Mesos, Swarm, etc all were coming about at the same time. The list he brings up beyond these is actually quite staggering - a lot more than you'd think. In the end, "everything is a Remix".

So why did Kubernetes win?

1. Well, umm, it works really well!
2. Vendor-neutral open source. There's a suggestion that enterprise want to see multiple companies backing it as they are scared of lock-in or paying a platform tax
   - I would argue, it depends on the enterprise. Sure, I talk about those things a lot now, but it didn't seem to come up so often in the past when we were getting into bed with big vendors, just sayin'
3. Oh, and also, it's the people. The ease of the steps to go from user --> contributor --> reviewer --> etc with a large and complicated thing like Kubernetes is incredibly impressive
   - I can't speak to this personally, but I'll take his word for it. We'll see later in the day just how huge the contributor base for Kube actually is.

---

## Cheryl Hung - Director of Ecosystem, CNCF

Next up we have Cheryl, whose name I recognise as she runs the [CloudNative London meetup](https://www.meetup.com/Cloud-Native-London/), which I'm a silent member of and am reminded that I really should start going to this one!

Apparently the CNCF now has 400+ members, with 88 end user community members. The latter are apparently more or less companies that pay to be members (sponsors, you might say?). I suppose it needs to get it's money from somewhere!

According to Github, for the CNCF there are now 2.66 million contributions and 56,214 contributors. That is, you have to say, quite a lot, you can't argue with that! Sadly though, apparently only 3% of contributors are women, so there is still a lot of work to do. If this is interesting to folks reading - look out for the talk [Navigating the Cloud Native Ecosystem for End Users](https://kccnceu19.sched.com/event/MPZh/navigating-the-cloud-native-community-for-end-users-cheryl-hung-cncf) amongst the KubeCon YouTube entries post-conference, which Cheryl references.

---

## Brian Liles - Senior Staff Engineer, VMWare

Brian works for VMWare. I'm going to be totally honest, I admit to raising my eyebrows when big vendor folks get up on stage at this sort of conference ... but he totally won me over by singing! He's bringing the energy for sure - his job seems to be to compere some CNCF project updates.

> Maybe I'm being harsh on the big vendors these days - they're trying to change and I've got to be honest - I'm biased.

With that in mind, we have ...

### Linkerd

A service mesh that isn't Istio. I hear someone joking later in the conference that Istio is the most common open source software to be talked about at a CNCF conference that isn't actually a CNCF project. It does make you wonder why - one to look into later ...

<blockquote class="twitter-tweet" align="center"><p lang="en" dir="ltr">Data point: Istio still hasn&#39;t been submitted to the CNCF and knative relies on Istio. Pretty clear that Google is playing OSS closer to the vest these days.</p>&mdash; Joe Beda (@jbeda) <a href="https://twitter.com/jbeda/status/1021793060420632576?ref_src=twsrc%5Etfw">July 24, 2018</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>

In any event, Linkerd is getting more popular, and now has some much-asked for features (that Istio already has, /cough), like zero-config mTLS and traffic shifting. I feel like I may go to a Linkerd talk later in the week just to see how it looks, or maybe swing by the booth. It has a reputation for being *much* easier to get up and running, and I've certainly found Istio to be painful in a Kubernetes cluster with good security :wink:

### Helm

Helm 3.0 is in alpha (*FINALLY*) and has some big things in it. The most relevant for my team is surely no Tiller (having a super-high privilege thing running is just not a good story), as well as some other features around namespace scoping, validation and chart libraries. To be honest we seem to be managing okay without Helm, but it's clearly very popular so this is still interesting to hear about - just in case we want to head down this route in future.

---

Project updates then move on to Harbor (a Registry), Rook (storage abstraction), and CRI-O (a Kubernetes Runtime Interface implementation). These are pretty neat bits of technology in general, but don't feel at all relevant to us right now - largely because we are all-in on Google Cloud and they have stuff that deals with this sort of thing for us.

Moving on to ...

### OpenCensus + OpenTracing

> This is more relevant for me - we have folks in our product teams asking us about these sort of tools already, and as we start to get more and more microservices talking to more and more other teams' microservices it'll become more and more relevant.

These guys have a great analogy that I like. Up pops a picture of the tip of an iceberg - "this is what your user sees". *Chuckles from the audience, I think we know what's coming*.

The pictures then move onto - "This is what your CFO sees" - cue underneath the iceberg, several shattered boats resting at the bottom. *Chuckles turn to laughs*

What do devs see? The iceberg starting to break up into smaller chunks under the surface. Everyone would rather have a Monoservice Iceberg than a Monolith Iceberg, right?

They then have a rather snappy definition of **Observability** - how well you can understand a system **given only the telemetry data** to work it out. Telemetry must become a built-in feature of cloud-native software - you can't have to hand-craft it all the time. That makes sense.

The good news is the confusion between these two tools is being lifted - OpenCensus and OpenTracing are merging into [OpenTelemetry](https://www.cncf.io/blog/2019/05/21/a-brief-history-of-opentelemetry-so-far/). I knew about this already but it's still great to hear and also to have a name for it. Now please bring on the snazzy logo.

This will deliver a single set of APIs for frameworks, libraries, etc to bind to. Also, crucially for us perhaps, it will be sdrawkcab compatible with both existing projects via software bridges if you have existing code, without change.

> This means we can just pick one now, safe in the knowledge things should still work when they've merged. Splendid!

---

The opening keynote rounds out with some conversation about graduated projects (which are up to 6 now, and an extended session on Fluentd). I chose to step out at this point to nose around the exhibition area for a bit before the Sessions started - you know, meet some vendors, ~~find solutions to our business problems~~ blag some T-shirts and stickers, etc.

{{< figure src="/images/kubecon-vendors.jpg?width=600px&classes=shadow" >}}

See below for some further highlights from the keynotes from the remaining days too!

---

## Spotify's Culture of Failure

> I'm kidding with the title - I :heart: Spotify :wink:

You really just to have to stand up and clap for folks who are willing to jump up on a stage in front of thousands of people and talk about the mistakes they've made. In this case I'm referring to Spotify ([David Xia](https://twitter.com/davidxia_?lang=en)), who had an extended segment on how they deleted their Kubernetes clusters by accident and why that's a good thing.

Basically, it's a conversation about learning:

- Learning that having multiple tabs in your browesr - some Prod, some Non-Prod -  is dangerous. We've all been there man
- Learning that your restore processes when you make a boo-boo might not actually be that great, and you probably should've practiced it before
- Learning that when you're changing things that are big, migrating users gradually so you have a backout plan or continuous availability is super-useful
- Learning that switching to new tooling (hi Terraform!) pays off in the end, but there's a learning curve that sometimes causes even more boo-boo's
- Learning that Terraform's state file and declarative nature can just sometimes be a real a%$!
- Learning that having a great team and culture that supports you when you're bricking it is probably the most important thing of all

---

## Other Customers Are Doing What We Are

There were two other customer talks at the keynotes that resonated with me - one because of the similarities and one due to fresh perspectives.

The first of these was ABN-AMRO - a large bank. They have a similar-ish narrative to ourselves - trying to transform to get ahead in a competitive industry and using containers to help them do that. I pulled out a couple of interesting observations from this talk:

- Firstly, they are starting to see software vendors supply their software in containers. This is something we've talked about internally but not seen come to pass (that I know of, at least). It was one of the things that was cited as a trigger for "we'll have to containers on-premise eventually".
- They find [the amount of choices out there](https://l.cncf.io) a challenge. My personal opinion is that I like choice, and find myself wrestling on what feels like a daily basis with "we'll pick this technology and make it work for you" versus giving development teams the freedom to choose. I think constantly challenging ourselves on that question is probably a healthy thing - which is why I like [the paved road analogy](https://medium.com/netflix-techblog/how-we-build-code-at-netflix-c5d9bd727f15) so much (because you can go off it!).
- They've built a dedicated platforms team, like us. They called theirs Stratus. We called ours Nimbus. They have the advantage of a poetic definition - *"low-level clouds characterised by horizontal layering with a uniform base"* - whereas I just think Nimbus sounds cooler :grin:
- They use - or possibly plan to use - [OPA](https://www.openpolicyagent.org/) for "compliance-as-code", and bringing the same arguments for infra-as-code (pairing, code reviews, etc) to that work. I like this idea. They have a good example of how this can provide early feedback to engineers, with policy that prevents them from creating a public load balancer. In the response the developer gets there's a clear message indicating why, and a place to go for more information and suggested alternatives. Clever stuff.
- Despite the analogies to my own world, they do however trigger me a little bit with their comments about **software components** being building blocks. This takes this too far in my opinion - this is beyond the platform boundary for me -  engineers need some freedom to innovate! I have some sympathy though that a bank might have some security challenges to face into (they mention "compliant by default"), but I do hope it's not forced ... just paved. And I wish them all the luck in the world in their potential for dependency hell :wink:

The second company I wanted to mention featured [Katie Gamanji](https://twitter.com/k_gamanji?lang=en) from Conde Nast International - a company I confess to never having heard of, despite working in retail - they do digital publishing stuff and are pretty enormous.

I really liked this talk - they have some interesting regional challenges (wanting to run in Russia & China), a historic complexity of running decentralised IT in their local markets (that they wanted to unify in a good way), and they were very open about their tech stack and how they made their choices. Even though most of their choices are different to my own team's, it's really interesting to hear folks talking about *why* they made the choices they did.

For example - they're running 9 clusters globally now across 100 instances of AWS. This is their current choice, but they've self-hosted their k8s control plane using [Tectonic](https://coreos.com/tectonic/) to give them freedom to port their nodes to other providers in the future if there's good reasons why. Pretty clever stuff.

---
