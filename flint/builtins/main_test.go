package builtins

var yml = `---
title: "Hello World"
date: 2012-12-12
date_bad: 12121212
nested:
  slug: "foo"
tags:
  - foo
  - bar
  - baz
categories:
  - foo
  - bar
  - baz
key:
  nested: value_enum

date_RFC3339: 2012-12-12T12:12:12+00:00
date_RFC3339_bad: 2012-99-12T12:12:12+00:00
date_RFC850: Tuesday, 12-Dec-12 12:12:12 UTC
date_RFC850_bad: Tuesday, 99-Dec-12 12:12:12 UTC
---`
