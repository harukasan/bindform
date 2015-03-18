# bindform

[![Godoc](https://img.shields.io/badge/Godoc-references-blue.svg?style=flat-square)](https://godoc.org/github.com/harukasan/bindform/bindform)
[![Build Status](https://img.shields.io/travis/harukasan/bindform.svg?style=flat-square)](https://travis-ci.org/harukasan/bindform)

Because bindform is aimed to only binding data, it does not cover data
validations (only type checking). You should validate data after (or before)
data binding.

## Usage

```
go get -u github.com/harukasan/bindform/bindform
```

```go
package main

import (
    "net/http"
    "github.com/harukasan/bindform/bindform"

type ContactForm struct {
  Name    string `form:"name"`
  Email   string `form:"email"`
  Message string `form:"message"`
}

func HandlePostContactForm(w http.ResponseWriter, r *http.Request) {
	form := &ContactForm{}

	err := bindform.BindPostForm(r, form)
	if err != nil {
		http.Error(w, "Bad Request", 400)
	}

  // Debug print
  // log.Println(form)

  // needs more validations ...
}
```

## Supported types

- bool
- int
- uint
- float
- string

## TODO

- [ ] Supports array type
- [ ] Supports require tag

## LICENSE

Copyright (c) 2015, MICHII Shunsuke.

See [License](./LICENSE)
