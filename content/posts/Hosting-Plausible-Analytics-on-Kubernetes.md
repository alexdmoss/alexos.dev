---
title: "Hosting Plausible Analytics on Kubernetes"
panelTitle: "Hosting Plausible Analytics on Kubernetes"
date: 2022-03-26T19:21:00-01:00
author: "@alexdmoss"
description: "How to run Plausible Analytics - a less angry version of Google Analytics - self-hosted on your own Kubernetes cluster"
banner: "/images/analytics.jpg"
tags: [ "Plausible", "Google", "Analytics", "Kubernetes", "Postgres", "Databases" ]
---

{{< figure src="/images/analytics.jpg?width=1000px&classes=shadow" attr="Photo by Markus Winkler on Unsplash" attrlink="https://unsplash.com/photos/IrRbSND5EUc" >}}

In this blog post I'll take you through how I set up the web analytics software [Plausible](https://plausible.io) on my own Kubernetes cluster. Whilst Plausible do publish guidance on how to do this, I found that it needed a few tweaks and I wanted to make a few enhancements to get it working reliably enough for my needs.

---

## A Brief Sidebar: What is Plausible?

If you haven't heard of Plausible, I'll keep it brief. It's a lightweight [more privacy-conscious](https://plausible.io/vs-google-analytics) version of Google Analytics. I personally use it on the handful of non-profit websites I host outside of my work _(including this one!)_ and think it's great. It does not have all the bells and whistles of GA but what it does have is:

1. No really invasive privacy stuff, like depending on cookies and allowing Google to use your data for other things that you aren't in control over.
2. Far simpler to use. I find GA to be much too complex for what I need. You can [have a look at their demo](https://plausible.io/plausible.io) to see what I mean.
3. The lack of dependence on your personal data or cookies means you can forego the consent forms, privacy notices and other stuff like that which is pervasive on the internet these days. For simple websites needing to do all that is overkill if you're not using the advanced features of GA anyway.
4. It can be self-hosted if you really want to make the commitment to your users that it's not gone off to who knows where.
5. It's open source, adding more confidence that it's not up to anything dodgy.
6. It's lean enough that adblockers (and Firefox) don't auto-block it. So the data is probably more reliable.

Okay that sounded like a bit of a sales pitch. I didn't mean it to come across that way - I just like [their philosophy](https://plausible.io/privacy-focused-web-analytics) and think they're trying to do a good thing.

{{< figure src="/images/tangent-meme.jpg?width=400px&classes=shadow" >}}

---

## Running It On Kubernetes

They provide [guidance in their docs](https://plausible.io/docs/self-hosting) and [their hosting repo](https://github.com/plausible/hosting) has some examples, but as I mentioned above I found I needed to make some tweaks.

{{< figure src="/images/k8s-yaml.jpg?width=400px&classes=shadow" >}}

[My working code can be found on my Github](https://github.com/alexdmoss/plausible-on-kubernetes/). The CI is put together using Gitlab (where its mirrored from), but can be adapted quite easily to your chosen CI tool, I'm sure.

Broadly, we need to solve for the following:

1. Get the basic components running - the main app, its Postgres database and its Clickhouse database.
2. Come up with a way of generating and storing the relevant secrets for it to work - there's quite a few of them: the credentials for Plausible itself, Postgres credentials, Clickhouse credentials, email settings, and your Twitter token if you use that.
3. Getting the emailed reports working.
4. Dealing with the fact that Plausible is behind some proxying load balancers.
5. Data backups and recovering into another cluster.
6. (Optional) Proxying the tracking Javascript (optional, depending on your morals).

I'll step through each of these sections in turn. Spoiler alert: the backups is the most complicated bit.

---

## The Basic Components & Secret Handling

{{< figure src="/images/secrets.png?width=400px&classes=shadow" >}}

I started with [their Kubernetes manifests](https://github.com/plausible/hosting/tree/master/kubernetes) and adapted these to suit my needs. I'm a fan of [`kustomize`](https://kustomize.io/) to patch existing manifests and there's nothing too fancy needed here as it turns out. You can see the structure I adopted in [this part of my repo](https://github.com/alexdmoss/plausible-on-kubernetes/tree/main/k8s).

> Note that I dropped the mail server - I had low confidence this would work in GCP where I run things anyway, as Google block that sort of thing as standard for probably obvious reasons. Without a service like AWS' SES, I configured Plausible to use Sendgrid as a mail relay instead.

The one bit of kustomize magic that's worth talking through is my use of a [`SecretGenerator`](https://github.com/alexdmoss/plausible-on-kubernetes/blob/main/k8s/base/kustomization.yaml#L10) in the base. I opted to use a single secret to hold all the required config data for all three of Plausible's components - slightly less good from a least-privilege perspective, but given they're all residing with their own namespace anyway, and the edge-facing app needs most of the juicy stuff anyway, this was a trade-off I could live with.

The generator ensures that a new secret is created on each pipeline run, guaranteeing any changes are picked up due to the pod restart. The secret values themselves are held in Google Secret Manager and patched in at deploy time through lines like this:

```bash
export ADMIN_USER_EMAIL=$(gcloud secrets versions access latest --secret="PLAUSIBLE_ADMIN_USER_EMAIL" --project="${GCP_PROJECT_ID}")`
```

Followed by an `envsubst`:

```bash
cat plausible-conf.env | envsubst "\$ADMIN_USER_EMAIL" > k8s/base/plausible-conf.env.secret
# don't persist this file, as it contains the unencrypted secret values!
```

In practice this part is a lot longer due to the many more secrets than this that need to be handled - see [deploy.sh](https://github.com/alexdmoss/plausible-on-kubernetes/blob/main/deploy.sh).

I used a similar technique to substitute in the target Plausible / Postgres / Clickhouse version - using `latest` tags in Kubernetes is a risky business, don't do that!

> Another tweak you'll likely want to make is thinking about the amount of storage to dish out to the Postgres & Clickhouse databases. This is going to depend on how heavily your sites are used but a word of warning - I found it failed silently when it filled up with the default sizing - you'll want a separate alert to track your disk usage to avoid being caught out by this!

---

## Email Reports

As mentioned above, given the constraints on GCP I opted not to muck about trying to persuade a mail server to work within GKE, and instead signed up for Sendgrid. At the scale I'm working with, I'm well within the free tier usage. I'm not going to go through setting up your Sendgrid account itself in detail as writing that down is unlikely to age well and I found it to be pretty straight-forward - the important part is to generate yourself an **API Key** in there, which we then inject into the Plausible config.

The following represents the particular combination of variables that I found did the trick:

```conf
MAILER_EMAIL=THE_EMAIL_ADDRESS_YOU_CONFIGURE_IN_SENDGRID
SMTP_HOST_ADDR=smtp.sendgrid.net
SMTP_HOST_PORT=465
SMTP_HOST_SSL_ENABLED=true
SMTP_USER_NAME=apikey
SMTP_USER_PWD=$SENDGRID_KEY
SMTP_RETRIES=2
```

The secret to substitute in is the `SMTP_USER_PWD` of course - the rest is standard stuff. With that in place, it "just works". You can test it by creating an extra user in Plausible and using the forgotten password link, as opposed to waiting for the weekly reporting to kick in.

---

## X-Forwarded-For Behind a Load Balancer

For Plausible to work behind a reverse proxy load balancer like the Kubernetes `nginx-ingress-controller`, some further tweaks are needed. You'll know this is needed if you find your Plausible setup seems to be ok, but the visitor countries / unique visitor tracking is not working - this means you are not forwarding the `X-Forwarded-For` header onto Plausible correctly. I found that the following configuration for the `nginx-ingress-controller` did the trick (these go in its `ConfigMap`):

```yaml
hsts: "true"
ssl-redirect: "true"
use-forwarded-headers: "false"      # not needed as not behind L7 GCLB, but YMMV
enable-real-ip: "true"
compute-full-forwarded-for: "true"
```

I also needed to edit my `Service` of `Type: LoadBalancer` to have `spec.externalTrafficPolicy: Local`. This affects evenness of load balancing a little, but was required for this to work and I didn't particularly mind this downside at my scale.

---

## Data Backups

{{< figure src="/images/kubernetes-storage.jpg?width=400px&classes=shadow" >}}

I wanted something in place here to allow me to recover historic data in the two databases in the event of needing to recreate my cluster (or something going wrong with it). I use this cluster to experiment with Kubernetes features fairly regularly, so there's a certain inevitability in this for me at least :grin:

My original plan was to use GCP's Cloud SQL (where you can run Postgres, and setting up of backups is trivial there), until I realised there was the Clickhouse database also. I thought I'd be able to use some of Clickhouse's own tooling for this, but as per [this github discussion](https://github.com/plausible/analytics/discussions/1226), the way its setup for Plausible means these don't work. Bummer.

So, what I ended up with was running both Postgres + Clickhouse inside Kubernetes (as per the recommendations by Plausible anyway) and using a trusted tool I've used elsewhere - [Velero](https://velero.io/) - to handly my backups. Velero allows you to snapshots Kubernetes resources and save them in object storage - in this case we can backup the Persistent Volumes and store in a GCS bucket.

> Note that Velero is a block storage snapshot of the disks in use. This is not a foolproof operation! At the small scale I am running this at I haven't had a problem on the two occasions I've used it to restore, but if your Plausible setup is under high load the risk of this not taking well increases. Always test your backups!

The setup of Velero itself is not part of my public Github repo, in part due to the high privileges it runs under as well as the lack of time to split it out. I plan to follow this up soon and will likely write a separate blog post about it when I do, so watch this space!

The setup is reasonably simple, with the key steps to know about once you've got it [installed and running](https://velero.io/docs/v1.8/basic-install/#install-and-configure-the-server-components) being:

1. Terraforming a GCS bucket to use for the backups, and a Service Account with credentials to use it (I suspect I could improve this further with a bit of Workload Identity usage).
2. Configure the `VoumeSnapshotLocation` and `BackupStorageLocation` for the GCS bucket above and ensuring the credentials are available to it.
3. Configuring a backup schedule that includes the PVs. Mine looks like this:

```yaml
---
apiVersion: velero.io/v1
kind: Schedule
metadata:
  name: pv-backup
spec:
  schedule: 0 */1 * * *
  template:
    includedResources:
    - persistentvolumeclaims
    - persistentvolumes
    includeClusterResources: true
    includedNamespaces:
    - '*'
    ttl: 168h0m0s
```

As you can see, I'm not very sensitive about the freshness of the data here - snapshotting hourly is fine by me, but you should tweak according to your needs of course.

Restores rely on you setting up the Velero client (although you could wrap this up in a CI script also if you needed to, of course). To test it, you then end up with something like this:

```bash
velero client config set namespace=velero

# get newest backup (it's a good idea to `describe` this to make sure it is valid)
BACKUP_NAME=$(velero backup get --output=json | jq -r '[ .items[] | select(.status.phase=="Completed") | {"name": .metadata.name, "startTimestamp": (.status.startTimestamp | fromdateiso8601)} ]| sort_by(.startTimestamp)[-1].name')

kubectl create ns plausible-test
# skip the `namespace-mapping` if doing this for real in a fresh cluster
velero restore create --from-backup "${BACKUP_NAME}" --include-resources persistentvolumeclaims,persistentvolumes --include-namespaces=plausible --namespace-mappings plausible:plausible-test --restore-volumes=true

# ... and then I run my Plausible's deploy.sh as normal against the right namespace
```

I've tested this a couple of times in practice without issue, although could do with wrapping some code around it to make it easier to execute next time. As always with backup solutions - make sure you keep a few versions, test it frequently, and set up alerts on failed backups!

---

## Proxying the Request

{{< figure src="/images/privacy.jpg?width=1000px&classes=shadow" attr="Photo by Lianhao Qu on Unsplash" attrlink="https://unsplash.com/photos/LfaN1gswV5c" >}}

Now then, decision time. Whether you do this or not is a bit of a moral choice. By disguising the Plausible tracking javascript like this, you are being a bit disingenuous with your users - although keep in mind that [their code is respectful](https://plausible.io/privacy-focused-web-analytics). Despite their approach, some browsers / browser extensions are sensitive enough to block the Plausible tracker, assuming it is just as naughty as the Google one. This technique helps you avoid that for more accurate analytics capture, if you're ok with that.

In my case, the majority of my sites are served via NGINX, so the [guidance here](https://plausible.io/docs/proxy/guides/nginx) covers what I need. You can see one of my [examples of this customisation here](https://gitlab.com/alexos-public/alexos.dev/-/blob/main/config/default.conf?ref_type=heads#L68).

> I've also recently started using Caddy for some website hosting. [Here is an example](https://gitlab.com/alexos-public/alexmoss-co-uk/-/blob/main/Caddyfile?ref_type=heads#L41) of how I've done it using that web server too.

Stripping this down:

```bash
# my cache path was different
proxy_cache_path /var/cache/nginx/data/jscache levels=1:2 keys_zone=jscache:100m inactive=30d  use_temp_path=off max_size=100m;

server {

  # proxy to plausible script - my self-hosted copy
  location = /js/visits.js {
      # you may want to change the filename here if you're not using the outbound link tracking feature
      proxy_pass https://YOUR.PLAUSIBLE.HOSTNAME/js/plausible.outbound-links.js;
      proxy_buffering on;

      # Cache the script for 6 hours, as long as plausible returns a valid response
      proxy_cache jscache;
      proxy_cache_valid 200 6h;
      proxy_cache_use_stale updating error timeout invalid_header http_500;
      add_header X-Cache $upstream_cache_status;

      proxy_set_header Host YOUR.PLAUSIBLE.HOSTNAME;
      proxy_ssl_name YOUR.PLAUSIBLE.HOSTNAME;
      proxy_ssl_server_name on;
      proxy_ssl_session_reuse off;
  }

  # proxy to plausible API - my self-hosted copy
  location = /api/event {
      proxy_pass https://YOUR.PLAUSIBLE.HOSTNAME/api/event;
      proxy_buffering on;
      proxy_http_version 1.1;
      
      proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto $scheme;
      proxy_set_header X-Forwarded-Host  $host;

      proxy_set_header Host YOUR.PLAUSIBLE.HOSTNAME;
      proxy_ssl_name YOUR.PLAUSIBLE.HOSTNAME;
      proxy_ssl_server_name on;
      proxy_ssl_session_reuse off;
  }

}
```

And then you update your tracking code to point to `/js/visits.js` instead, something like this:

```html
<script defer data-api="/api/event" data-domain="yourwebsite.com" src="/js/visits.js"></script>
```

You could of course host the Javascript itself within your side code too (skipping the first location block) - although you still need to make sure the `/api/event` calls are proxied on to your self-hosted Plausible instance to capture the visits. You'd also not pick up upgrades to the tracking code automatically if you did this.

---

## Conclusions

With all that done you should be in a position to have a usable Plausible setup self-hosted in your Kubernetes cluster. From here you can do a variety of other useful things to get more value out of it, like configuring [outbound link tracking](https://plausible.io/docs/outbound-link-click-tracking), [404 tracking](https://plausible.io/docs/404-error-pages-tracking) and [Google Search integration](https://plausible.io/docs/google-search-console-integration). These all worked fine for me just by following their instructions - there's nothing special that needs doing for a self-hosted version here.

I hope you found that run-through useful - it's been ticking along quite nicely for around 3 months now without issue. I've yet to attempt to upgrade it, although I'm not expecting any issues. I'll be keeping [my repo](https://github.com/alexdmoss/plausible-on-kubernetes) up to date with any tweaks I make along the way, so hopefully that will continue to prove a useful resource for anyone attempting to do what I have.

Oh, and if you work for a company looking to use Plausible, do consider their [hosted option](https://plausible.io/#pricing) if you don't need to keep the data on your own servers. It looks pretty sensibly priced to me, and helps them to keep improving their product!
