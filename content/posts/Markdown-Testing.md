---
title: "Markdown Testing"
description: "This post is all about markdown testing - hopefully it works!"
date: 2019-02-21T15:28:51Z
author: "@alexdmoss"
draft: true
banner: ""
tags:
- Markdown
categories:
- Testing
---

## Heading 2

### Heading 3

#### Heading 4

##### Heading 5

---

Finally a normal paragraph. I wonder if I can manage to make up enough text to turn this into a "genuine" paragraph. I know! Let's try some **BOLD TEXT**. Oh that was exciting, do you *feel the italics* add to the narrative? I think so too.

Here's a second paragraph. This one is a lot less ~~boring~~ interesting but at least it makes sure things are pretty sensibly spaced. I think we can all sleep easier now, eh?

> Oh, and here is a blockquote. It should be emphatically emphatic

And finally, [a link](/posts/).

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

This is some `inline code that is long enough to have a copy button`, and then a code paragraph:

```sh
~:> this is a shell
~:> echo "hello"
```

It may even support GHFM - lets try some Javascript:

```js
var getUrlParameter = function getUrlParameter(sPageURL) {
	var url = sPageURL.split('?');
	var obj = {};
	if (url.length == 2) {
		var sURLVariables = url[1].split('&'),
			sParameterName,
			i;
		for (i = 0; i < sURLVariables.length; i++) {
			sParameterName = sURLVariables[i].split('=');
			obj[sParameterName[0]] = sParameterName[1];
		}
		return obj;
	} else {
		return undefined;
	}
};
```

Python?

```python
from modules.joke import joke

joke = joke()

print(joke)
```

Golang?

```golang
func main() {

	var userMsg string

	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "-help" {
		displayHelp()
		os.Exit(0)
	} else {
		userMsg = os.Args[1]
	}

	msg := formatMessage(userMsg)

	fmt.Println(msg)

}
```

YAML?

```yaml
apiVersion: v1
kind: Service
metadata:
  name: moss-work-svc
  namespace: moss-work
  labels:
    name: moss-work-svc
spec:
  type: NodePort
  ports:
  - port: 80
    name: http
  selector:
    name: moss-work
---
```

Terraform?

```terraform

resource "google_container_cluster" "cluster_1" {

  name                            = "mw-prod"
  region                          = "europe-west1"

  master_auth {
    client_certificate_config {
      issue_client_certificate    = false
    }
    username                      = ""
    password                      = ""
  }

  initial_node_count              = 1
  min_master_version              = "1.10.7-gke.6"
  remove_default_node_pool        = true

  monitoring_service              = "monitoring.googleapis.com/kubernetes"
  logging_service                 = "logging.googleapis.com/kubernetes"

  lifecycle {
    ignore_changes = ["node_pool"]
  }

}
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

## Icons

https://www.webfx.com/tools/emoji-cheat-sheet/

---

## Images

Images can be resized, bordered and shadowed:

![Minion](https://octodex.github.com/images/minion.png?width=10pc&classes=shadow)

![stormtroopocat](https://octodex.github.com/images/stormtroopocat.jpg?height=100px&classes=border)

Shortcode can even be used to caption them:

{{< figure src="https://octodex.github.com/images/dojocat.jpg?width=200px&classes=shadow" attr="A cool picture" >}}