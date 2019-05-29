# Reflections on a trip to KubeCon/CloudNativeCon 2019

<!-- Unsplash - Barcelona: https://unsplash.com/photos/f-DvU93UhTs -->

---

I was lucky enough to make the trip over to Barcelona last week for [KubeCon + CloudNativeCon 2019](https://events.linuxfoundation.org/events/kubecon-cloudnativecon-europe-2019/). We are a big user of Kubernetes here in JL&P, and as we are increasingly pushing the boundaries of the easier out-the-box stuff you can do with it, it was very useful for myself and a colleague to learn more about what others are doing with these sort of technologies and hopefully learn some tricks - or at least share some challenges! - and see what might be coming down the pipeline that might helps us out.

---

## Any Stand-out Themes?

One personal theme I reflect on - in contrast to some of the "big vendor" conferences I've been to in the past - is how nice it is to actually hear from real customers and users of this stuff. Sure, I went to a few talks from the big tech companies, but even then they're presenting very experience-focused content; their speakers are focusing on how they're engaging with the community as a whole, and so on. Not at all sales-pitch-y. This I think is pretty brilliant - it meant the conference overall felt more valuable.

It's clear that the community around these technologies, with Kubernetes at its foundation, is vibrant. This is fantastic for the future of these sort of services, and it's also pretty clear to me that this is at least in part because the people behind it are working very hard on promoting this attitude.

From a technology perspective, people really are actually doing the hybrid & multi-cloud thing. I personally see this as something you should really only do if you absolutely have to (if lock-in scares you, perhaps think about your exit strategy rather than missing out on the benefits of backing just one host?). That said, I get that for some it isn't an option - compliancy for example. So if we *did* have to do this at some point in the future, at least other folks will have trodden down the path a bit.

<!-- /images/kubecon-expo-hall.jpg -->

---

As I work with - and am hopefully seen as someone who helps set the direction for! - a set of "platforms" teams, it was very reassuring to hear from several other end-users of this tech in the enterprise, all approaching things in more or less the same way that we have chosen to.

> This was definitely a community of platform builders using these tools as the raw blocks to build even better platforms tailored for their own users (software developers).

Many of these organisations also seem to have started from a place similar to ourselves (*"I need to go faster", "I need to put more powerful tools in the hands of our engineers"*) and they continue to face the same sort of challenges we do (*"how do I choose which option to go with?"*, *"how much do I do / get others to do / leave to product engineers?"*). We are clearly not alone in the questions we're facing.

And of course, the onus is then on us to continue to keep solving the next challenge and the challenge after that - and to make sure we stay on top of the best ways of building platforms so it doesn't become stale and ineffective. A trap we've fallen into in the past.

---

## This Kubernetes Thing Might Catch On

It was really nice to hear so much positivity about Kubernetes, the project that started this whole thing - which turned 5 years old at roughly the date of this conference. There were a few nods to this around the conference centre.

<!-- /images/kubecon-donuts.jpg -->

During one of the keynotes, [Janet Kuo](https://www.youtube.com/watch?v=w62T1SN4g6Y) did a very nice segment on the Kubernetes project in CNCF itself. I liked her timeline narrative:

- 2003 - Borg was created. No, not in Star Trek, this is Google's collective instead. That was, staggeringly enough, May 1989
- 2006 - Linux cgroups arrived
- 2008 - Linux adopted the "container" naming convention
- 2009 - Omega - the next-gen Borg at Google, which heavily influences Kubernetes design
- 2013 - Docker. Hooray for Docker, it's the best
- 2013 - Project 7 within Google (open sourcing container orchestrator - it comes from Seven of Nine from Star Trek of course - the Friendly Borg)
- 2014 - Kubernetes announced at DockerCon. There was much rejoicing
- 2015 - v1.0 of Kubernetes, CNCF formed, 1st KubeCon at the end of the year
- 2016 - Kubernetes SIGs came into being, KubeCon EU, industry adoption, in Production, at scale. This thing is a big deal
- 2017 - CRD introduced so can build your own Kubernetes resources on top. Kubernetes becomes a platform for building platforms *(if Kelsey Hightower had a pound for every time that was mentioned at this conference ...)*. Cloud adoption drives Kubernetes to become a *de facto* standard, many new native APIs introduced
- 2018 - Graduated from incubation in the CNCF. This is seen as a "stable in Production" thing
- Today - one of the highest velocity open source projects. It is number 2 on Pull Requests on Github (behind ... Linux), #4 on issues/authors (out of 10's of millions)

Now it is v1.14 and the most stable and mature ever. And hell you can even run it on Windows nodes. Now *that is* crazy ;)

<!-- /images/kubecon-books.jpg | It was also great to be reminded of the Kubernetes Comic Books - I didn't realise there were two of them! -->

---

> Right ... that's all well and good, but what about some actual technologies of interest?

---

## Service Meshes

Service Meshes have always been an interest for me since I first learnt about the concept in more detail at Google's Next conference last year. Up until now I'd always assumed that Istio would be our answer to that - we run on Google Kubernetes Engine after all, and Google are the main contributor to Istio - but I used the opportunity at KubeCon to get a bit more of a rounded view on this topic.

<!-- /images/kubecon-still-want-to-try.jpg | https://makeameme.org/meme/still-want-to-5b73f1 -->

First there was the introduction of [SMI](https://smi-spec.io/) (Service Mesh Interface). The [announcement was brief](https://www.youtube.com/watch?v=gDLD8gyd7J8), but welcome - having these things start to coalesce around a standard can only be a good thing.

I also listened to [1&1 talk](https://www.youtube.com/watch?v=vQ2IktsMlgQ) about some of the challenges they'd experienced with Istio which mimicked my own (I've had some real fun and games trying to get the GKE Addon to play nicely with [PodSecurityPolicy](https://kubernetes.io/docs/concepts/policy/pod-security-policy/) - I think I'll [wait for the Operator](https://discuss.istio.io/t/istio-operator-plans-for-1-2/2227)!).

Finally, I took in [a talk](https://www.youtube.com/watch?v=E-zuggDfv0A) on an alternative mesh technology, [Linkerd](https://linkerd.io), which I liked the sound of and am tempted to try as an alternative. That said, my suspicion is that, in the end, we'll still end up on Istio due to the richer features and fast-moving development resolving some of these challenges before we get round to needing it.

If you're interested, in my opinion there's a really nicely put together [comparison piece here](https://itnext.io/linkerd-or-istio-2e3ce781fa3a).

---

## Observability

I also deliberately took in a number of talks on observability tools - in large part because it's an area of focus for my team at the moment.

<!-- https://unsplash.com/photos/kSLNVacFehs -->

We actually have three different angles on this:

### 1. Prometheus at Scale

We have an emerging scaling problem with Prometheus. At the moment, our Promethei are relatively small but still by far the biggests pods we run and occasionally perturbed by increases in tenant usage. Their standalone nature feels wasteful and limits us to only holding data for (I think) 45 days. Teams want to keep it for longer, which is not unreasonable in my view - we need to do better.

I therefore was particularly interested to learn more about [M3](https://eng.uber.com/m3/), [Cortex](https://medium.com/weaveworks/what-is-cortex-2c30bcbd247d) and a little about [Thanos](https://improbable.io/blog/thanos-prometheus-at-scale) as open source solutions to this particular challenge. Whatever we end up picking will probably end up making an interesting blog post in its own right!

### 2. The Logging Problem

Logs are useful to debug problems. They're also in my view a trap for metrics to go into. We are doing a really good job lately at avoiding the latter, but that doesn't mean we can continue to get away with providing our engineers with sub-optimal tooling for log analysis.

We already know about the Elastic Stack, but I went to [a talk](https://www.youtube.com/watch?v=CQiawXlgabQ) on [Loki](https://grafana.com/loki) just to see what that was about and I liked what I heard - "built for engineers to solve problems" in particular - but I think it's just too early days for us to gamble on it. Probably. Maybe. We will see :)

It was also interesting to observe how crowded the booths for Elastic, Loki, Logz.io and DataDog were. Or maybe folks just wanted the stickers ... ;)

### 3. Distributed Tracing

We are pretty close to needing to provide some sort of distributed tracing tool in my view. We are doing microservice-y things and we are at the point now where those microservices are a little less **customer --> microservice --> legacy** and a bit more **customer --> microservice --> another microservice --> oh and maybe another microservice too --> eventually legacy**.

So the news that the two main open source ways of instrumenting that stuff - OpenTracing & OpenCensus - are joining together in a backwards-compatible way is very pleasing indeed.

I'm hoping we'll have an excuse to write more about this topic soon once we've gotten more hands-on with some of the tooling in this area. It looks very powerful.

---

## Wrap Up

And that wasn't all either - I went to some other really interesting specialist tracks, including a couple of sessions on [Open Policy Agent](https://www.openpolicyagent.org/) & [Gatekeeper](https://github.com/open-policy-agent/gatekeeper) (which I think we should use!), Multi-Cluster Ingress, dynamic pod autoscaling, a Custom Deployment Controller, plus some weird and wonderful Kubernetes failure scenarios.

All-in-all, it was a really fantastic conference with loads and loads of technical depth on an incredibly diverse platform ecosystem which is supported by an active and enthusiastic community.

/images/kubecon-party.jpg / Party time at KubeCon!

---

## More Detail on Some of These Topics

If you made it this far into this article, good job! If you're interested in a bit of a deeper look at some of these areas, I have written a few more words on [my own blog](https://mosstech.io/categories/kubecon/).
