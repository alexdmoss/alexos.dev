# mosstech.io

Building a new blog theme using Hugo.

---

## To Do

### Version 1

- [x] Restructure CSS
- [x] Panel links - make obvious its a link
- [x] Extra links appearing in sidebar
- [x] Sidebar recent posts formatting
- [x] Contacts page formatting
- [x] About page content
- [x] Links to subheading in posts possible
- [ ] Colours from m10c / make it look a bit different
- [ ] Contacts page processing
- [x] Post links to include date in URL / permalinking
- [ ] Share menu
- [x] Code highlighting for GHM not working - needs styles
- [x] Copy-to-clipboard
- [ ] Make search work
- [x] 404 page - content
- [ ] Mobile appearance
- [x] Credits for theme
- [x] Deep link to header should scroll up a little to allow for top-menu
- [ ] 404 page - nginx config
- [ ] CI/CD
- [ ] Google Analytics
- [ ] alexmoss.co.uk redirect

### Version 2

- [ ] Contacts form maybe does not need the labels at all with the placeholders?
- [ ] Recent posts to use thumbnail rather than description (see Future Imperfect)
- [ ] Tranquil Peak - hide sidebar when reading a post is nice
- [ ] Captcha on the Contacts form
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

- 1D2024 353B43 57CBCC 50514F F25F5C
- 1D2024 353B43 57CBCC B4B8AB 284B63
- 1D2024 353B43 57CBCC 00635D 08A4BD
- 1D2024 353B43 57CBCC 84828F 6A687A
- 1D2024 353B43 57CBCC 594E36 7E846B
- 1A1A1C 363739 3D8AA6 706C61 F8F4E3
- 1F1F21 131315 B3A683 C8C3BD CEE0DC
- #0F0A0A / #F4EFED / #57CC8A / #1A535C / #4ECDC4
- #18121E / #233237 / #984B43 / #EAC67A. Navy + Gunmetal + Rusty + Warm Yello
- #C5C1C0 / #0A1612 / #1A2930 / #F7CE3E. Screen + Steel + Denim + Marigold
- #494E6B / #98878F / #985E6D / #192231. Stormy + Cloud + Sunset + Evening
- #E14658 / #22252C / #3F3250 / #C0B3A0. Coral + Navy + Mountain + Scrub
- #1E1F26 / #283655 / #4D648D / #D0E1F9. Midnight + Indigo + Blueberry + Periwinkle
- #756867 / #D5D6D2 / #353C3F / #FF8D3F. Wood + Sand + Charcoal + Orange
- #20232A / #ACBEBE / #F4F4EF / #A01D26. Ink + Aluminium + Paper + Ruby Red
- #080706 / #EFEFEF / #D1B280 / #594D46. Black Steel + Paper + Gold Leaf + Silver
- #DDDEDE / #232122 / #A5C05B / #7BA4A8. Gray + Blackish + Houseplant + Blue-Gray
- #8EE4AF / #EDF5E1 / #5CDB95 / #907163 / #379683 (with black back)
- #0C0032 / #190061 / #240090 / #3500D3 / #282828 (with black back)
- #0B0C10 / #1F2833 / #C5C6C7 / #66FCF1 / #45A29E
- #2C3531 / #116466 / #D9B08C / #FFCB9A / #D1E8E2
- #161617 / #090A0A / #151516
- #080708 / #C4C3C4 / #7B7C82 / #4B4A4B / #7D8284
- #feda6a #d4d4dc #393f4d #1d1e22
- #393939 #FF5A09 #ec7f37 #be4f0c
- #262626 #3f3f3f #f5f5f5 #dcdcdc

### Hugo Themes

Future Imperfect - clean, nice icons - https://themes.gohugo.io/future-imperfect/
Tranquilpeak - very nice soft theme, left menu - https://themes.gohugo.io/hugo-tranquilpeak-theme/
Arabica - nice font - https://themes.gohugo.io//theme/arabica/
Meghna - very professional, dark grey with animations - https://themes.gohugo.io/meghna-hugo/
m10c - quite nice colours - https://themes.gohugo.io/hugo-theme-m10c/
Massively - the way Contacts is done is clever - https://themes.gohugo.io/hugo-theme-massively/
