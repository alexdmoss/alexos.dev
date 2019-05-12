# Themes

Colour Themes I experimented with are below, in case they are useful later.

## Restructure of Theme CSS

This CSS allowed me to test colour combos more easily, using the styling below:

```css
:root {
    --main-bg-color: var(--bg-color);

    --main-font-color: var(--font-color);
    --post-title-color: #var(--font-color);
    --header-link-color: var(--font-color);

    --panel-bg-color: var(--shade-color);
    --button-color: var(--shade-color);
    --box-color: var(--shade-color);

    --header-bg-color: var(--post-color);
    --article-bg-color: var(--post-color);

    --border-color: var(--line-color);
    --sidebar-border-color: var(--line-color);
    --button-border-color: var(--line-color);

    --header-shadow-color: var(--highlight-color);
    --button-hover-color: var(--highlight-color);
    --header-link-hover-color: var(--highlight-color);
    --link-hover-color: var(--highlight-color);
    --article-shadow-color:  var(--highlight-color);

    --shadow-color: rgba(0, 0, 0, 0.15);
    --deep-shadow-color: rgba(0, 0, 0, 0.6);

    --code-bg-color: #222;
    --code-color: #e6e6e6;
    --back-to-top-color: #333;

    --table-shadow-1:rgba(0, 0, 0, 0.12);
    --table-shadow-2: rgba(0, 0, 0, 0.24);
    --table-head-bg-color: #888;
    --table-body-bg-color: #e0e0e0;
    --table-alt-row-color: #f4f4f4;
}
```

## Useful Colours

These colours were pulled from themes below - may still be useful:

```css
:root {
    --post-color: #1A1A1C;
    --bg-color: #233237;
    --bg-color: #20232A;
    --bg-color: #1F2833;
    --bg-color: #151516;
    --font-color: #F5F5F5;
    --shade-color: #4B4A4B;
    --shade-color: #393F4D;
    --highlight-color: #57CBCC;
    --highlight-color: #3D8AA6;
    --highlight-color: #57CC8A;
    --highlight-color: #F7CE3E;
    --highlight-color: #A01D26;
}
```

### Hugo Themes

- Future Imperfect - clean, nice icons - https://themes.gohugo.io/future-imperfect/
- Tranquilpeak - very nice soft theme, left menu - https://themes.gohugo.io/hugo-tranquilpeak-theme/
- Arabica - nice font - https://themes.gohugo.io//theme/arabica/
- Meghna - very professional, dark grey with animations - https://themes.gohugo.io/meghna-hugo/
- m10c - quite nice colours - https://themes.gohugo.io/hugo-theme-m10c/
- Massively - the way Contacts is done is clever - https://themes.gohugo.io/hugo-theme-massively/

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

---

## Themes Experimented With

I eventually settled on a variant of "Tango", although very nearly went with "Bumblebee, and liked most of them to be honest!

```css

/* Armour */
:root {
    --bg-color: #353B43;
    --highlight-color: #57CBCC;
    --shade-color: #50514F;
    --post-color: #000000;
    --font-color: #e6e6e6;
    --line-color: #57CBCC;
}

/* Matrix */
:root {
    --bg-color: #1F1F21;
    --highlight-color: #57CC8A;
    --shade-color: #4ECDC4;
    --post-color: #0F0A0A;
    --font-color: #F4EFED;
    --line-color: #57CC8A;
}

/* Futuristic */
:root {
    --bg-color: #1E1F26;
    --highlight-color: #4D648D;
    --line-color: #4D648D;
    --shade-color: #283655;
    --post-color: #000000;
    --font-color: #D0E1F9;
}

/* Armory */
:root {
    --bg-color: #20232A;
    --highlight-color: #A01D26;
    --line-color: #A01D26;
    --shade-color: #ACBEBE;
    --post-color: #000000;
    --font-color: #F4F4EF;
}

/* Lime */
:root {
    --bg-color: #232122;
    --highlight-color: #A5C05B;
    --line-color: #A5C05B;
    --shade-color: #7BA4A8;
    --post-color: #000000;
    --font-color: #DDDEDE;
}

/* Tron */
:root {
    --bg-color: #1F2833;
    --highlight-color: #66FCF1;
    --line-color: #66FCF1;
    --shade-color: #45A29E;
    --post-color: #0B0C10;
    --font-color: #C5C6C7;
}

/* Bumblebee */
:root {
    --bg-color: #1D1E22;
    --highlight-color: #FEDA6A;
    --line-color: #FEDA6A;
    --shade-color: #393F4D;
    --post-color: #000000;
    --font-color: #D4D4DC;
}

/* Tango */
:root {
    --bg-color: #393939;
    --highlight-color: #FF5A09;
    --line-color: #FF5A09;
    --shade-color: #393F4D;
    --post-color: #000000;
    --font-color: #D4D4DC;
}

/* Industrial */
:root {
    --bg-color: #233237;
    --highlight-color: #EAC67A;
    --shade-color: #0F0A0A;
    --post-color: #000000;
    --font-color: #F4EFED;
    --line-color: #EAC67A;
}

/* Steampunk */
:root {
    --bg-color: #116466;
    --highlight-color: #D9B08C;
    --line-color: #D9B08C;
    --shade-color: #2C3531;
    --post-color: #0B0C10;
    --font-color: #D1E8E2;
}
```

## Saved Colours

I saved these colours as potentially useful later on:

```css
/* May have useful extra colours */
:root {
    --bg-color: #262626;
    --bg-color: #363739;
    --bg-color: #1F1F21;
    --bg-color: #6A687A;
    --bg-color: #151516;

    --highlight-color: #00635D;
    --highlight-color: #DCDCDC;
    --highlight-color: #3D8AA6;
    --highlight-color: #7B7C82;
    --highlight-color: #F7CE3E;

    --post-color: #1D2024;
    --post-color: #1A1A1C;
    --post-color: #000000;
    --post-color: #080708;


    --line-color: #EAC67A;
    --line-color: #7B7C82;

    --shade-color: #594E36;
    --shade-color: #3F3F3F;
    --shade-color: #4B4A4B;

    --font-color: #F5F5F5;
    --font-color: #C5C1C0;
    --font-color: #C4C3C4;

    --link-color: #F25F5C;
}
```

## Original

The following was the original colour theme, based on the Future Imperfect Hugo Theme:

```css
:root {
    --main-bg-color: #e6e6e6;
    --main-font-color: #333;
    --panel-bg-color: #f4f4f4;
    --header-bg-color: #fff;
    --article-bg-color: #fff;
    --border-color: rgba(160, 160, 160, 0.7);
    --sidebar-border-color: rgba(160, 160, 160, 0.3);
    --button-border-color: rgba(160, 160, 160, 0.3);
    --button-hover-color: #2ebaae;
    --button-color: #333;
    --code-bg-color: #222;
    --code-color: #e6e6e6;
    --box-color: #555;
    --header-link-color: #999;
    --header-link-hover-color: #000;
    --post-title-color: #000;
    --back-to-top-color: #333;
    --table-shadow-1:rgba(0, 0, 0, 0.12);
    --table-shadow-2: rgba(0, 0, 0, 0.24);
    --table-head-bg-color: #888;
    --table-body-bg-color: #e0e0e0;
    --table-alt-row-color: #f4f4f4;
    --shadow-color: rgba(0, 0, 0, 0.15);
    --deep-shadow-color: rgba(0, 0, 0, 0.6);
}
```
