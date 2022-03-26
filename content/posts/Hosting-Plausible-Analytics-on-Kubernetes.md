---
title: "Hosting Plausible Analytics on Kubernetes"
panelTitle: "Hosting Plausible Analytics on Kubernetes"
date: 2022-03-26T19:21:00-01:00
author: "@alexdmoss"
description: "How to run Plausible Analytics - a less angry version of Google Analytics - self-hosted on your own Kubernetes cluster"
banner: "/images/containers.jpg"
tags: [ "Plausible", "Google", "Analytics", "Kubernetes", "Postgres", "Databases" ]
---

//TODO: banner image
//TODO: some funny if you can

{{< figure src="/images/containers.jpg?width=1000px&classes=shadow" attr="Photo by Ian Taylor on Unsplash" attrlink="https://unsplash.com/photos/jOqJbvo1P9g" >}}

In this blog post I'll take you through how I set up Plausible Analytics on my own Kubernetes cluster. Whilst Plausible do publish some guidance and manifests on how to do this, I found that it needed a few tweaks and I wanted to make a few enhancements to get it working reliably enough for my needs.

## Sidebar: What is Plausible?

A brief tangent into what Plausible is if you haven't heard of it. It's a lightweight [more privacy-conscious](https://plausible.io/vs-google-analytics) version of Google Analytics. I personally use it on a handful of website I host outside of work _(including this one!)_ and think it's great. It does not have all the bells and whistles of GA but what it does have is:

1. No really invasive privacy stuff, like depending on cookies and Google being able to use your data for other things. It's not shared with third parties or advertising companies etc etc.
2. Way simpler to use. I find GA to be far too complex for what I need. You can [have a look at their demo](https://plausible.io/plausible.io).
3. It doesn't rely on cookies or other personal data to work. Because of this, you can forego some of the consent form / privacy notice stuff on sites that you see everywhere these days and are getting really rather tiresome.
4. The option to self-host it if you really want to make the commitment to your users that it's not gone off to who knows where
5. It's open source.
6. It's lean enough that adblockers (and Firefox) don't auto-block it. So the data is probably more reliable to boot.
7. It's a lot more lightweight to run on your site - lower page weight and lower loading times.

Sounds like a bit of a sales pitch. Didn't mean it to come across that way - I just really like [their philosophy](https://plausible.io/privacy-focused-web-analytics) on this. You can [r]

---

## Running It

Okay enough with the elevator pitch - lets get on with the Kubernetes bits!

If you just want to skip to the TL:DR, then [my working code lives on my Github](https://github.com/alexdmoss/plausible-on-kubernetes/). The CI is put together using Gitlab (where the original lives), but can likely be adapted quite easily to your chosen CI tool.

Lets go through it though in a bit more detail. Broadly, we need to solve for the following:

1. Get the basic components running - the main app, its Postgres database and its Clickhouse database
2. Come up with a way of generating and storing the relevant secrets for it to work - there's quite a few of them: the credentials for Plausible itself, Postgres credentials, Clickhouse credentials, email settings, and your Twitter token if you use that
3. Getting the emailed reports working
4. Dealing with the fact that Plausible is behind some proxying load balancers
5. Data backups and recovering into another cluster
6. Proxying the tracking Javascript (optional, depending on your morals)

### The Basic Components & Secret Handling

I started with [their Kubernetes manifests](https://github.com/plausible/hosting/tree/master/kubernetes) and adapted these to suit my needs. I'm a fan of [`kustomize`](https://kustomize.io/) to patch onto existing manifests - and there's nothing too fancy here - you can see the structure I adopted in [this part of my repo](https://github.com/alexdmoss/plausible-on-kubernetes/tree/main/k8s).

> I dropped the mail server - I had low confidence this would work in GCP where I run things anyway, they block that sort of thing as standard for probably obvious reasons. Without a service like AWS' SES, I configured Plausible to use Sendgrid as a mail relay instead - see below

The only slightly fancy bit is the [`SecretGenerator`](https://github.com/alexdmoss/plausible-on-kubernetes/blob/main/k8s/base/kustomization.yaml#L10) in the base. I opted to use a single secret to hold all the required config data for all three components - slightly less good from a least-privilege perspective, but given they're all residing with their own namespace anyway, and the edge-facing app needs most of the juicy stuff anyway, this was a trade-off I could live with for much simpler configuration.

The generator ensures that a new secret is created on each pipeline run, guaranteeing any changes are picked up due to the pod restart. The secret values are held in Google Secret Manager and patched in at deploy time, for each through lines like this:

```bash
export ADMIN_USER_EMAIL=$(gcloud secrets versions access latest --secret="PLAUSIBLE_ADMIN_USER_EMAIL" --project="${GCP_PROJECT_ID}")`
```

Followed by an envsubst:

```bash
cat plausible-conf.env | envsubst "\$ADMIN_USER_EMAIL" > k8s/base/plausible-conf.env.secret
```

In practice this part is a lot longer due to the many more secrets than this that need to be handled - see [deploy.sh](https://github.com/alexdmoss/plausible-on-kubernetes/blob/main/deploy.sh).

I used a similar technique to substitute in the target Plausible / Postgres / Clickhouse version - using `latest` tags in Kubernetes is a risky business, don't do that!

> Another tweak you'll likely want to make is thinking about the amount of storage to dish out to the Postgres & Clickhouse databases. This is going to depend on how heavily your sites are used but a word of warning - I found it failed silently when it filled up with the default sizing - you'll want a separate alert to track your disk usage to avoid being caught out by this!

### X-Forwarded-For Behind a Load Balancer

For Plausible to work behind a reverse proxy load balancer like the Kubernetes nginx-ingress-controller, some further tweaks are needed.
Whilst things mostly appear to be working, if you find that visitor countries / unique visitor tracking is not working, then you are not forwarding the `X-Forwarded-For` header onto Plausible correctly. I found that the following configuration for the nginx-ingress-controller did the trick:

```json
hsts: "true"
ssl-redirect: "true"
use-forwarded-headers: "false"      # not needed as not behind L7 GCLB, but YMMV
enable-real-ip: "true"
compute-full-forwarded-for: "true"
# use-proxy-protocol: "true"        # breaks things
```

Do **not** set `use-proxy-protocol` to true. As I found, this breaks things :grin:

I also needed to edit my `Service` of `Type: LoadBalancer` to have `spec.externalTrafficPolicy: Local`. This affects evenness of load balancing a little, but was required for this to work.

### Email Reports

As mentioned above, given the constraints on GCP I opted not to muck about trying to persuade a mail server to work within GKE, and instead signed up for Sendgrid. At the scale I'm working with, I'm well within the free tier usage. I'm not going to go through setting up your Sendgrid account itself in detail as writing that down is unlikely to age well - the important part is to generate yourself an **API Key** in there, which we then inject into the Plausible config.

This stuff goes into secret I'm generating and mounting as environment variables, and the particular combination of variables you need are as follows:

```conf
MAILER_EMAIL=THE_EMAIL_ADDRESS_YOU_CONFIGURE_IN_SENDGRID
SMTP_HOST_ADDR=smtp.sendgrid.net
SMTP_HOST_PORT=465
SMTP_HOST_SSL_ENABLED=true
SMTP_USER_NAME=apikey
SMTP_USER_PWD=$SENDGRID_KEY
SMTP_RETRIES=2
```

The secret to substitute in is the `SMTP_USER_PWD` of course - the rest is standard stuff. With that in place, it "just works". You can test it by creating an extra user in Plausible and using the forgotten password link, as opposed to waiting for the weekly reporting.

### Data Backups

I wanted something in place here to allow me to recover historic data in the two databases in the event of needing to recreate my cluster (or something going wrong with it). I use this cluster to fiddle around with Kubernetes features all the time so that's somewhat likely to be honest!

Before I realised the Clickhouse element, my original plan was to use GCP's Cloud SQL for the Postgres part. Backups are trivial there. However given that the Clickhouse bit is A Thing, I ended up running both databases in-cluster. High Availability wasn't a particular concern for me, but data integrity was. Without a tickboxy solution here, I turned to a trusted tool I've used elsewhere - [Velero](https://velero.io/)  ...

Note that Velero is a block storage snapshot of the disks in use by the services running in Kubernetes. This is not a foolproof operation! At the small scale I am running this at I haven't had a problem on the two occasions I've used it to restore, but if your Plausible setup is under high load the risk of this not taking well increases. Always test your backups!

> If you're not okay with this, then plenty of other options exist to backup/restore Postgres cleanly. However the Clickhouse DB is a little more unusual and the type of tables used prevent some of their tools being used. [This github discussion](https://github.com/plausible/analytics/discussions/1226) provides some good background.

The setup of Velero itself is not part of my public Github repo _(it runs out of a private one in part due to the high privilege it needs to be installed with - although I'll think about splitting it out to make this part of the post better!)_, but is quite straight-forward to be honest. On GCP it boils down to:

1. Terraforming a GCS bucket to use for the backups, and a Service Account with credentials to use it (I suspect I could improve this further with a bit of Workload Identity usage)
2. Deploying Velero [according to its docs](https://velero.io/docs/v1.8/basic-install/#install-and-configure-the-server-components), configuring a  `BackupStorageLocation` for the bucket above and ensuring the credentials are available to it.

//TODO: can we move Velero setup to github also for better guidance?

My backup schedule for Persistent Volumes across the cluster looks like this:

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

Restores rely on you setting up the Velero client (althuogh you could wrap this up in a CI script also if you needed to, of course). You then end up with something like this:

```bash
kubectl create ns plausible-test
velero restore create --from-backup $BACKUP_NAME --include-resources persistentvolumeclaims,persistentvolumes --include-namespaces=plausible --namespace-mappings plausible:plausible-test --restore-volumes=true
# and then I run my Plausible's deploy.sh against the plausible-test namespace instead
```

In a real disaster/failover, you wouldn't be changing the namespace name like this, as you'd be in a clean new cluster instead, but you probably get the idea (just skip the `namespace-mapping`). Things to keep in mind if restoring into a new GCP project entirely:

- the new Service Account you create for Velero needs access to the old storage bucket with backups in it, and access to the compute snapshots its taken in the old project. Choose your project carefully for this! You can make your life simpler than I did by having a dedicated project for backup storage
- make sure your `VolumeSnapshotLocation` is aware of the old project by setting `.spec.config.project=old-project`

I ended up with something like this to deal with two locations:

```bash
velero client config set namespace=velero
# note storageLocation==old must match backuplocation above if not using shared bucket
BACKUP_NAME=$(velero backup get --output=json | jq -r '[ .items[] | select(.spec.storageLocation=="old") | select(.status.phase=="Completed") | {"name": .metadata.name, "startTimestamp": (.status.startTimestamp | fromdateiso8601)} ]| sort_by(.startTimestamp)[-1].name')

# ... depending on state of previous deployment, may need to delete old PVs
velero restore create --from-backup "${BACKUP_NAME}" --include-resources persistentvolumeclaims,persistentvolumes --include-namespaces=plausible --restore-volumes=true
```

I've tested this a couple of times in practice without issue, although could do with wrapping some code around it to make it easier to execute next time. As always with backup solutions - make sure you keep a few versions, test it frequently, and set up alerts on failed backups!

//TODO: can we add alert config for failed backups here

### Proxying the Request

Now then. Whether you do this or not is a bit of a moral choice. By disguising the Plausible tracking javascript like this, you are being a bit disingenuous with your users - although keep in mind that [their code is respectful of this](https://plausible.io/privacy-focused-web-analytics). Despite this, some browsers / browser extensions are sensitive enough to block the Plausible tracker, assuming it is just as naughty as the Google one. This technique helps you avoid that for more accurate analytics capture, if you're ok with that.

In my case, the majority of my sites are served via NGINX, so the [guidance here](https://plausible.io/docs/proxy/guides/nginx) covers more or less all you need. You can see one of my [examples of this customisation here](https://gitlab.com/alexos-dev/moss-work/-/blob/master/config/default.conf).

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

And then you update the tracking code to point to `/js/visits.js` instead:

```html
<script defer data-api="/api/event" data-domain="website.com" src="/js/visits.js"></script>
```

You could of course host the Javascript itself within your side code too (skipping the first location block) - although you still need to make sure the `/api/event` calls are proxied on to your self-hosted Plausible instance to capture the visits.

---

## Conclusions

With all that done you should be in a position to have a usable Plausible setup self-hosted in your Kubernetes cluster. From here you can do a variety of other useful things to get more value out of it, like configuring [outbound link tracking](https://plausible.io/docs/outbound-link-click-tracking), [404 tracking](https://plausible.io/docs/404-error-pages-tracking) and [Google Search integration](https://plausible.io/docs/google-search-console-integration). These all worked fine for me just by following their instructions - there's nothing special that needs doing for a self-hosted version here.

I hope you found that run-through useful - it's been ticking along quite nicely for a around 3 months now without issue.
