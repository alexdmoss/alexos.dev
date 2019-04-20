---
title: "Markdown Testing"
description: "This post is all about markdown testing - hopefully it works!"
date: 2019-02-21T15:28:51Z
author: "@alexdmoss"
tags: hugo markdown
banner: ""
---

## Heading 2

### Heading 3

#### Heading 4

##### Heading 5

---

Finally a normal paragraph. I wonder if I can manage to make up enough text to turn this into a "genuine" paragraph. I know! Let's try some **BOLD TEXT**. Oh that was exciting, do you *feel the italics* add to the narrative? I think so too.

Here's a second paragraph. This one is a lot less ~~boring~~ interesting but at least it makes sure things are pretty sensibly spaced. I think we can all sleep easier now, eh?

> Oh, and here is a blockquote. It should be emphatically emphatic

---

## Tables

| Column 1 | Column 2   |          Column 3 |
| -------- | ---------- | ----------------: |
| Row 1    | A          |                 B |
| Row 2    | Words here | With **markdown** |
| Row 3    | C          |                 D |
| Row 4    | E          |                 F |

---

## Code

This is some `inline code`, and then some code paragraph:

```
$ ~:> this is a shell
$ ~:> echo "hello"
```

It may even support GHM:

```js
grunt.initConfig({
  assemble: {
    options: {
      assets: 'docs/assets',
      data: 'src/data/*.{json,yml}',
      helpers: 'src/custom-helpers.js',
      partials: ['src/partials/**/*.{hbs,md}']
    },
    pages: {
      options: {
        layout: 'default.hbs'
      },
      files: {
        './': ['src/templates/pages/index.hbs']
      }
    }
  }
};
```

---

## Bullets 

Unordered bullets look like this:

- One
- Two
- Three
  - Three A
  - Three B
- I thought you said unordered?

Whereas ordered ones look like this:

1. A
2. B
3. C
4. 4

---

## Images

Images can be resized, bordered and shadowed:

![Minion](https://octodex.github.com/images/minion.png?width=10pc&classes=shadow)

![stormtroopocat](https://octodex.github.com/images/stormtroopocat.jpg?height=100px&classes=border)

Shortcode can even be used to caption them:

{{< figure src="https://octodex.github.com/images/dojocat.jpg?width=200px&classes=shadow" title="A cool picture" >}}