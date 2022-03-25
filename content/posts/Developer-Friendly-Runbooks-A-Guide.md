---
title: "Developer-Friendly Runbooks: A Guide"
panelTitle: "Developer-Friendly Runbooks: A Guide"
date: 2021-09-24T19:28:00-01:00
author: "@alexdmoss"
description: "A guide to creating an easy-to-use, developer-friendly runbook site to help improve operability"
banner: "/images/runbook-doggo.jpg"
tags: [ "Operability", "Runbooks", "DevOps", "Guide" ]
---

{{< figure src="/images/runbook-doggo.jpg?width=600px&classes=shadow" attr="Photo by Cookie the Pom on Unsplash" attrlink="https://unsplash.com/photos/gySMaocSdqs" >}}

**When** things go wrong - and yes, they will go wrong - it's extremely helpful to have easy access to a set of runbooks to guide the unfortunate engineer through the steps needed to mitigate the problem as swiftly as possible. In this post I'm going to describe the approach we use for this where I work, which we've found to work very well.

This post is therefore intended to be a guide to creating a site to host runbooks, as well as the sort of content it's useful to include within them. Read on if that sounds interesting to you!

> To help illustrate some of the points I'm making below, I've created my own simplified version of the site - i.e. a lot less content. You can visit it here: [runbooks.alexos.dev](https://runbooks.alexos.dev), and the source code for it lives [here](https://github.com/alexdmoss/runbooks).

---

## What do we need from a Runbooks repository?

{{< figure src="/images/runbook-library.jpg?width=600px&classes=shadow" attr="Photo by Susan Q Yin on Unsplash" attrlink="https://unsplash.com/photos/2JIvboGLeho" >}}

For a collection of runbooks to be successful, there are three essential things in my opinion for any runbook repository to have from the outset:

1. It needs to be **easy to access**. The user is often accessing this stuff in less than ideal circumstances. It needs to be straight-forward to get to, easy to navigate / search for the thing you need, and the pages need to load fast over a potentially limited connection.
2. It needs to **work when others things aren't working**. I can't imagine anything worse than being called out for a problem and finding the thing on fire is also the place where all the instructions on how to put it out are kept!
3. It needs to be **easy to update**, in particular by the people who know how to fix things. By that I mean - it mustn't require some software that not everyone will have, or hide the ability to update it behind complicated access controls or connectivity that aren't part of the engineer's normal day job.

There are some additional things that I think are also really helpful, but not strictly necessary:

- Hosting all the runbooks from **a shared location**. In reality this might be a collection for a single (complex) website, or grouped by business area perhaps. The benefits of this are numerous and not always obvious - "new" engineers always know where to look, and it provides inspiration through visibility of how other teams have approached it. It also helps ensure the contents are self-sustaining - standing on the shoulders of giants to cut down the effort for those that follow, and fostering continual improvement. It also helps in spotting obsolete content for removal.
- **Keep the tech simple**. You want something that's easy to host and unlikely to break, and not need too much specialist knowledge to maintain. In some ways this goes for the contents of the runbooks themselves too - although they need to be detailed enough to be useful, of course :smile:
- As a content author, always try to **think about whether you can remove a runbook** - or parts of it - through automation. Can the contents be simplified with a script to gather the relevant diagnostic data? Can something be done to (safely) automatically heal the service when an alert fires, instead of requiring a human to get out of bed and do something?

---

## What technology should I use, and where should I host it?

{{< figure src="/images/runbook-glue.jpg?width=600px&classes=shadow" attr="Photo by Erik MacLean on Unsplash" attrlink="https://unsplash.com/photos/RfkaDKptt-A" >}}

As mentioned above - the important thing here is that it works when other things aren't working. You want a **simple and reliable tech stack**, and hosting that does not have many moving parts of dependencies, and is **separate from where the stuff you need to fix lives**.

I've found a lot of success using static sites generated from markdown - and for this I recommend [Hugo](https://gohugo.io). As mine (and my work's) primary hosting technology is Kubernetes I **do not** deploy this to kube. The fact that it's just a bunch of HTML / other static assets means my choices of where to host are lengthy - I could in fact use an entirely different cloud provider, or Gitlab/Github Pages for example.

The static content satisfies my desire for the pages load really fast (unless you overdose on the frontend wizardry) and keeps the moving parts (and therefore risk of breakage) to a minimum.

The use of markdown satisfies my desire to ensure it's easy to update. The vast majority of engineers are going to be familiar with markdown - and even if they're not its an extremely simple syntax to pick up / copy from someone else. Text files are easily stored in git, which fits nicely with an engineer's normal workflow, too.

> A tangential perk of git is that it's also common and/or easy for those people on call to have a local copy of the runbooks cloned on their machine. This can (usually) be pulled even if something bad does happen to the runbook hosting when it's needed. Although as noted above this shouldn't happen. *Cough* shouldn't *cough*.

What Hugo brings to the mix is the templating of these markdown files into something that's easy to use through the theme you create for it. For [runbooks.alexos.dev](https://runbooks.alexos.dev) I've customised the [Hugo Learn Theme](https://themes.gohugo.io/themes/hugo-theme-learn/), which is itself based on [Netlify's Grav Learn Theme](http://learn.getgrav.org/). It's worked well at work and I don't recall hearing any complaints in the three or so years we've been running with it. It's also been fairly easy to extend with new functionality when the odd request has come in.

---

## A little more on the hosting ...

{{< figure src="/images/runbook-sign.jpg?width=600px&classes=shadow" attr="Photo by Jamie Templeton on Unsplash" attrlink="https://unsplash.com/photos/6gQjPGx1uQw" >}}

At work, due to the need to still have some authentication & authorisation to get to the content, we stick with the same cloud provider, but use one of their serverless compute options - Google AppEngine - instead of Kubernetes.

Whist AppEngine is arguably a little overkill (compared to say, using a storage bucket), it does give some freebies that we've become used to on Kubernetes - automatic SSL certificates, and most crucially automatic integration (through Google's Identity Aware Proxy product) with our corporate identity service. This makes it trivially simple to make our runbooks accessible to large numbers of employees with almost no effort or ongoing maintenance.

As with most serverless platforms, getting a build & deployment pipeline up and running is also super-simple - see my [two-stage Gitlab CI](https://github.com/alexdmoss/runbooks/blob/main/.gitlab-ci.yml), for example - which more or less boils down to running hugo to generate the HTML and `gcloud app deploy` to publish it.

It's incredibly low cost too :money_with_wings:, which is a bonus.

> I'm planning to write more about how to host static sites on a variety of GCP hosting options in some future posts - so stay tuned if that's of interest!

---

## So, what sort of content goes into the runbook site?

{{< figure src="/images/runbook-tools.jpg?width=600px&classes=shadow" attr="Photo by Eugen Str on Unsplash" attrlink="https://unsplash.com/photos/CrhsIRY3JWY" >}}

I don't think it's crucial or really beneficial to be super-opinionated about what goes into a runbook personally - I think it's better to trust the owning team's engineers to put the things they need to know about in there _(although I would **strongly** encourage [testing your runbooks](https://alexos.dev/2020/03/29/team-nimbus-and-the-agents-of-chaos/) to make sure you've got that right!)_.

That said, what is useful is having a bit of a template (for mine, see [here](https://runbooks.alexos.dev/examples/02-runbook-master-template/) and [here](https://runbooks.alexos.dev/examples/01-service-index-template/)) to get things started and inspire those who are staring at the proverbial blank piece of paper.

What follows are some examples of what I think is good to include - linking to some examples where applicable to illustrate further:

| Entry | Example | Description |
| ------------------------------------------- | ------- | ----------- |
| Key information about the owning team       | [Example](https://runbooks.alexos.dev/application-2/#team-details) | Team Name, channels to contact you by, key team members like your PO/PM if applicable. Whilst you might think you know these things, keep in mind that your runbook might be in use by other teams. |
| A summary of what the service actually does | [Example-1](https://runbooks.alexos.dev/application-1/#service-information) <br />[Example-2](https://runbooks.alexos.dev/application-3/components/) | This doesn't have to be down to the minute detail of course - things change! - but a bit of an overview (especially again for those outside your team). |
| Important links in an emergency             | [Example](https://runbooks.alexos.dev/application-2/#useful-support-links) | Your key dashboards, application logs, support rota if applicable, perhaps where your code lives. A low-key useful thing I'd add is any A/B tests or experiments you might be running. New stuff breaks. |
| Your upstreams and downstreams              | [Example](https://runbooks.alexos.dev/application-2/general-information/dependencies/) | Try to [align on this terminology](https://reflectoring.io/upstream-downstream/) if you can :smile:. Regardless, listing our who your consumers/clients are, and what other services you depend on to work (internal and external) is useful info in complex architectures. |
| Instructions for responding to Alerts       | [Example](https://runbooks.alexos.dev/application-1/#content-renderer-microservice-has-all-pods-down) | Deep-link to these from wherever your alerts fire from if at all possible to cut down your time-to-respond. Having to navigate to them is just wasteful. Also as you tune and/or remove your alerts, it's a good prompt to tidy up the runbook too. |
| General How-To's                            | [Example-1](https://runbooks.alexos.dev/application-3/processes/resolving-dead-letters/) <br />[Example-2](https://runbooks.alexos.dev/application-3/processes/scaling-resources/) | In other words, these are instructions on how to repair or remediate The Thing, as opposed to instructions on how to respond to an alert. |
| How to release and rollback                 | [Example-1](https://runbooks.alexos.dev/application-2/investigating-issues/release-and-rollback/) <br />[Example-2](https://runbooks.alexos.dev/application-3/processes/build-and-deployment/) | Whilst ideally standard, rollbacks are not usually used frequently. Also don't underestimate the value in having this written down for others - I've certainly been in the position of describing this over the phone, and having a guide for them to follow would've been handy, believe me! |
| Debugging Guide                             | [Example](https://runbooks.alexos.dev/application-2/investigating-issues/) | Particularly helpful for newer team members - a few tips on where to look for clues is often appreciated. This can be things like where the logs are and how to increase the verbosity, or how to connect to and look at the running services and view their events, and so on. |
| Anticipated Failure Scenarios               | [Example](https://runbooks.alexos.dev/application-2/investigating-issues/anticipated-failure-scenarios/) | If you've spent the time thinking about what could go wrong - a useful exercise in my opinion - then including those in your runbooks is a good idea. |
| Third Party Support instructions            | [Example](https://runbooks.alexos.dev/application-3/third-party-support/) | If your service depends on third party systems, then how to get in touch with them is extremely useful. The majority of incidents I see at work tend to be this - and often there's little you can do except chase (although you absolutely should be doing something like [this](https://martinfowler.com/bliki/CircuitBreaker.html) if you can too). |
| Callout and internal support procedures     | [Example-1](https://runbooks.alexos.dev/application-3/processes/callouts/) <br />[Example-2](https://runbooks.alexos.dev/application-3/processes/hitperson/) | How you as a team have agreed to work and prompts for the things you need to do - especially during a major incident where you'll naturally be distracted trying to fix The Thing. |

---

## Summary

{{< figure src="/images/runbook-email.png?width=400px&classes=shadow" attr="The I.T Crowd is a masterpiece" >}}

I hope you found something of use in this post. Whilst hosting a static site with some instructions in it might not be the most technologically-complex thing in the world - I hope the post above shows that putting a little thought into it in advance is worthwhile when that site is the thing you really, really need when all kinds of [chaos is unleashed](https://alexos.dev/2020/03/29/team-nimbus-and-the-agents-of-chaos/).
