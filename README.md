# mosstech.io

Building a new blog theme using Hugo.

---

## To Do

### Version 1

- [x] Convert px to rem/em
- [x] Share menu
- [x] Social share links not working - because localhost?
- [x] Make search work
- [x] Search wants 'tags' pages
- [x] Mobile appearance
- [x] Share menu links are active even when menu not showing
- [x] Blog Posts index is showing up in recents
- [x] Bring in images and markup from Google blog post
- [x] Sweep out any comments
- [ ] CI/CD
- [ ] 404 page - nginx config
- [ ] Contacts page processing
- [ ] Google Analytics
- [ ] alexmoss.co.uk redirect
- [ ] Other domains?

### Version 2

- [ ] JS - auto-complete.js
- [ ] JS - html5shiv-printshiv.min.js
- [ ] JS - jquery.sticky.js
- [ ] JS - modernizr.custom.71422.js
- [ ] JS - perfect-scrollbar.jquery.min.js + perfect-scrollbar.min.js
- [ ] Harden nginx build
- [ ] Move to Github so there's some proper stuff there
- [ ] Link to splash photos directly
- [ ] Not sure about placement of social links on blog post - make a bit more subtle?
- [ ] Print the tags/categories somewhere on posts
- [ ] How about some tests?
- [ ] Reclick on share button hides it
- [ ] Experiment with fonts
- [ ] Hover effects equivalent for thumbs?
- [ ] Auto-minify css/js
- [ ] Accessibility testing
- [ ] Contacts form maybe does not need the labels at all with the placeholders?
- [ ] Recent posts to use thumbnail rather than description (see Future Imperfect)
- [ ] Tranquil Peak - hide sidebar when reading a post is nice
- [ ] Captcha on the Contacts form
- [ ] Printable version - colours especially need inverting
- [ ] Comments
- [ ] Sticky footer - https://codepen.io/BretCameron/pen/oVNYKR
- [ ] Tags cloud
- [ ] Portfolio page (inc where host?)
- [ ] Archives page for whole taxonomy (Tranquil Peak)

### Fancier Formatting

- [ ] Keys to navigate left and right
- [ ] Gradient backgrounds: https://gradienthunt.com/gradient/1955
- [ ] Tranquil Peak - animation load of the about page
- [ ] Slight zoom-in effect on panels / image banners
- [ ] meghna has some cool animations - especially social buttons at bottom
- [ ] meghna contact form at bottom?
- [ ] Pinch animations from https://www.demisto.com/community/

### Testing

- [ ] Social share links
- [ ] Safari
- [ ] Firefox
- [ ] Edge
- [ ] Mobile
- [ ] Analytics
- [ ] RSS / Atom

---

## Nice-to-Have Features

- Nicer emoji
- Add open graph protocol - meta property="og:title" - see Future Imperfect
- Add twitter cards integration - meta name="twitter:card" - see Future Imperfect
- Captcha
- Highlighting for comments, like Medium
- Related Posts popping up from the bottom

## Tech

- Move to github
- Build some CI in github
- GCP Load Balancing
- CDN
- NGINX in Kubernetes
- Stackdriver Dashboards
- Stackdriver Alerts

---

### Icons

https://www.webfx.com/tools/emoji-cheat-sheet/

```html
<i class="fas fa-code-branch"></i>
<i class="far fa-comments"></i>
<i class="far fa-copy"></i>

<i class="fab fa-gitlab"></i>

<i class="fas fa-home"></i>
<i class="fas fa-layer-group"></i>
<i class="fab fa-pinterest"></i>
<i class="fab fa-reddit"></i>

<i class="fas fa-tags"></i>
```

### Colours

https://coolors.co

- https://coolors.co/1d2024-353b43-57cbcc-50514f-f25f5c
- https://coolors.co/1A1A1C-363739-3D8AA6-706C61-F8F4E3
- https://coolors.co/1F1F21-131315-B3A683-C8C3BD-CEE0DC
- https://coolors.co/0F0A0A-F4EFED-57CC8A-1A535C-4ECDC4
- https://coolors.co/18121E-233237-984B43-EAC67A-000000
- https://coolors.co/C5C1C0-0A1612-1A2930-F7CE3E-000000
- https://coolors.co/1E1F26-283655-4D648D-D0E1F9-000000
- https://coolors.co/756867-D5D6D2-353C3F-FF8D3F-000000
- https://coolors.co/20232A-ACBEBE-F4F4EF-A01D26-000000
- https://coolors.co/080706-EFEFEF-D1B280-594D46-000000
- https://coolors.co/DDDEDE-232122-A5C05B-7BA4A8-000000
- https://coolors.co/0B0C10-1F2833-C5C6C7-66FCF1-45A29E
- https://coolors.co/2C3531-116466-D9B08C-FFCB9A-D1E8E2
- https://coolors.co/161617-090A0A-151516-000000-000000
- https://coolors.co/080708-C4C3C4-7B7C82-4B4A4B-7D8284
- https://coolors.co/feda6a-d4d4dc-393f4d-1d1e22-000000
- https://coolors.co/393939-FF5A09-ec7f37-be4f0c-000000
- https://coolors.co/262626-3f3f3f-f5f5f5-dcdcdc-000000

### Hugo Themes

Future Imperfect - clean, nice icons - https://themes.gohugo.io/future-imperfect/
Tranquilpeak - very nice soft theme, left menu - https://themes.gohugo.io/hugo-tranquilpeak-theme/
Arabica - nice font - https://themes.gohugo.io//theme/arabica/
Meghna - very professional, dark grey with animations - https://themes.gohugo.io/meghna-hugo/
m10c - quite nice colours - https://themes.gohugo.io/hugo-theme-m10c/
Massively - the way Contacts is done is clever - https://themes.gohugo.io/hugo-theme-massively/
