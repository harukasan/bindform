# bindform

Post data binding into struct, inspired from
[martini-contrib/binding](https://github.com/martini-contrib/binding).

Because bindform is aimed to only binding data, it does not cover data
validations (only type checking). You should validate data after (or before)
data binding.

## Usage

```go

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

