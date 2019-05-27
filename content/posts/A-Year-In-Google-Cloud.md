---
title: "A Year in Google Cloud"
date: 2019-02-23T15:28:51Z
author: "@alexdmoss"
description: "An overview of our journey from migrating our first enterprise-scale application to Google Cloud, through to multi-tenant Kubernetes"
banner: "/images/platforms.jpg"
tags: [ "GCP", "Google", "Kubernetes", "GKE", "DevOps", "Migration", "Platforms" ]
categories: [ "Cloud", "Google" ]
---

> This blog entry was originally [posted on Medium](https://medium.com/john-lewis-software-engineering/a-year-in-google-cloud-4586a117f352) for my employer.

This time last year, our newly-formed Platforms Team in John Lewis Online were putting the finishing touches to a brand new **Kubernetes** platform designed to run the frontend of [johnlewis.com](https://www.johnlewis.com/) in Google Cloud. Twelve months later, we’ve passed through Black Friday without a hitch and built a raft of new capabilities along the way. What follows is a post reflecting on the journey so far — if that sounds interesting, then read on!

---

## Doing That Strategic Thing

{{< figure src="/images/strategic-thing.jpg?width=600px&classes=shadow" title="Photo by rawpixel on Unsplash" >}}

The frontend of johnlewis.com — what we call ‘Browse’ — wasn’t the first thing we built in Google Cloud. There were a couple of teams deliberately given the freedom to experiment in GCP, and they built a number of smaller apps that could quickly get into Production. This helped cultivate the idea that this was good technology to be working with, and we should start using it for bigger things. That, plus the fact that our engineers were chomping at the bit to get their hands on it, really helped generate the initial push it needed. Adopting the cloud for [johnlewis.com](https://www.johnlewis.com/) really felt like an engineer-led venture — more so than any other piece of work I’ve been involved with in my time at JLP.

Choosing to host something business-critical in the cloud still felt like a big step though. As is often the way in large enterprises, we had to do a fair bit of convincing across IT that this was the right thing to do — or more accurately, that we were going about it the right way (the truth was that there was little doubt that cloud was the future, but were we doing the right things in cloud, and in the right way?). I used to refer to this work as the “lightning rod” — if we were comfortable putting the very first thing our customers see onto cloud infrastructure, then surely we’d be ok putting other things there too? This really helped us move from something perceived as tactical to something we could brand internally as “strategic” and bring the various teams on the journey with us.

## Managed Open Source

{{< figure src="/images/open-source.jpg?width=600px&classes=shadow" title="Photo by Émile Perron on Unsplash" >}}

<!-- https://unsplash.com/photos/xrVDYZRGdw4?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText -->

One of the conscious choices we made was not to over-stretch ourselves. We were migrating an application that we knew would be great in cloud — it was built with SpringBoot & Nginx and had no requirement for state (in other words, it was designed for cloud). We also knew from early experimentation that it would fit well into Docker containers, and so Kubernetes to manage it (we’d need to run a lot of them!) was an obvious choice.

We therefore had to decide whether to use [managed Kubernetes from Google](https://cloud.google.com/kubernetes-engine/) (GKE), or roll our own. We knew from others who’d worked with the technology that while Kubernetes was great when it was up and running, it was also pretty complicated when you’re new to it (… it’s gotten a bit easier these days). Given that Google invented the thing, it seemed like a no-brainer to start there!

This turned out to be a great choice. We get all the lovely features Google give you out the box with very low effort to maintain the core platform — which of course we have automated through CI/CD pipelines executing [Terraform](https://www.terraform.io/docs/providers/google/getting_started.html) code (automate all the things — not just the stuff you do a lot — was a principle we were clear on from the start). This lets us focus on what our end-users (and we as a team) want, rather than spending lots of time keeping things running.

In fact, this approach worked so well that we carried this pattern forward as we expanded into other technology stacks. Does it make sense to run your own Kafka, Mongo or PostgreSQL when there’s a capable, low-cost option within your cloud provider’s ecosystem? We continuously challenge ourselves as to whether using the cloud-native offering (in our case, provided by Google) holds us back in terms of the features we need, or whether it just helps us get there faster. We need to always be thinking about the work we’d need to do if we had to switch to another provider of an equivalent service, but try to make sure the trade-off feels worth it. Most of the time we choose the Google offering as the on-ramp is so smooth, but this is not always the case — we aren’t big users of Stackdriver or their development tools for example.

## A Platform for Growth

{{< figure src="/images/platform-growth.jpg?width=600px&classes=shadow" title="Photo by Sven Scheuermeier on Unsplash" >}}

<!-- https://unsplash.com/photos/YhdEgF-qWlI?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText -->

So the new johnlewis.com platform went live in February with all our mobile traffic (~40% of the site at the time), and a few months later we flipped all the desktop traffic over too. Thanks to a lot of hard work this went without a hitch, and we were done and dusted by the summer. Great! So, now what?

Enter the JL Digital Platform, version 2.

{{< figure src="/images/bigger-platform.jpg?width=600px" >}}

Our first iteration of the platform was very much built for the Browse application. We were new to it, we focused solely on what that app needed, and made some decisions along the way that were ultimately to make that initial migration easier.

But now, with the success of the switch and the promise of what cloud could do for JL (and also partly coz we shouted about how great it was!) we found that our platform had caught the eye of other teams — who were keen to adopt our tech stack and move away from some of the constraints of their on-premise infrastructure.

We therefore had a think about this, ran some experiments and settled on a design for a new platform — one that used the same base building blocks (Docker, Kube), but this time was designed from the ground up to support different workloads across different teams — in other words a multi-tenant Kubernetes cluster. Fantastic though GKE is, we didn’t want to be spinning up clusters left right and centre for every service that came along, especially as there was now significant momentum behind an “aren’t microservices great” strategy.

{{< figure src="/images/say-microservices.jpg?width=600px" >}}

So, we spent a month playing with some of the newer features of Kubernetes and convinced ourselves that we could still have one cluster for all the things — even if all the things were being developed and operated by completely different Product Teams.

If you’re familiar with Kubernetes, you’ll recognise that we leaned heavily on its [Namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/) concept to achieve this separation, building on top of that with extra operational layers and security capabilities to give us the control we needed but still preserve the freedom for development teams to concentrate on building and running their own services. This for me is what makes Kubernetes so great — it finds an awesome sweet spot between empowering teams with self-service goodness while giving Platform Engineers appropriate tools to make workloads, well, work.

## A Paved Road, or a Dusty Path?

{{< figure src="/images/paved-road.jpg?width=600px&classes=shadow" title="Photo by Jesse Bowser on Unsplash" >}}

<!-- https://unsplash.com/photos/c0I4ahyGIkA?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText -->

Most of the effort with version 2 has revolved around streamlining the on-boarding process for teams, and it is still an area we are continuously improving upon.

We love the concept of [the paved road](https://medium.com/netflix-techblog/how-we-build-code-at-netflix-c5d9bd727f15) and our focus is on building something that teams want to use. We check our thinking by talking to tenants — during on-boarding and through joint team retros, trying to make sure they have the tools they need already there and ready for them.

One of the powerful things about cloud that really changes the game compared to my experiences of working within your own datacenter is how easy it is to try things in a no regrets way. Whether it is when Google release a cool-new-thing or someone has a brilliant idea that they then open source, you can usually spin it up in minutes (*... sometimes hours*) and see whether it’s something that will work for you. Our team are great at running spikes and we allocate specific days to just “try new stuff” (Friyays … yes, the name could do with work …) specifically for this purpose — and it often puts us in a good position when a Product Team arrives with a new requirement we haven’t figured out a paved road for yet.

That said, we also value contributions from Product Engineers when they’re really keen to try something — especially when our small team doesn’t have the bandwidth to pick it up themselves. Things like our Grafana/Prometheus stack and API Gateway are now part of the core platform, having been born from other teams experimentation inside Kubernetes. We recognise when teams have had a good idea and try to industrialise it so that others can benefit. In this way, our platform becomes a product in its own right — continually gaining new features to stay modern and relevant, rather than withering on the vine.

## What’s Next?

{{< figure src="/images/whats-next.jpg?width=600px&classes=shadow" title="Photo by Djim Loic on Unsplash" >}}

<!-- https://unsplash.com/photos/ft0-Xu4nTvA?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText -->

As I’ve alluded to, the use of our new platform is expanding rapidly and the challenge for us is staying on top of all the new features we want to offer to teams as they onboard — as well as making that on-boarding process as slick as possible.

What does that mean in practice? Well, we are working on things like:

- Delivering visualisation of cloud costs per team and service, so that teams can make smart choices about the size of their containers and how they consume other services.
- Providing blueprints for telemetry, alerting and dashboards, so that every service has these by default without exception and teams are thinking about this sort of thing early on.
- Enabling log analytics at scale, rather than being dependent on tools that scale poorly with large volumes of operational log data.
- Continually improving our resilience and scalability, because why wouldn’t you want to do that!
- Introducing new cross-cutting capabilities such as API management, service meshes and distributed caching, because we’re going to need them and — let’s be honest — they’re cool!
- Beyond that, who knows? As we design more event-based architectures, I can see things like serverless technologies (perhaps through the emerging [Knative](https://cloud.google.com/knative/)) playing a big role.

Let’s just say, it’ll be really interesting to see what a similar post at the end of 2019 has to say about things ...

---
