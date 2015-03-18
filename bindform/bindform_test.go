package bindform_test

import (
	"bytes"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/harukasan/bindform/bindform"
)

type params map[string]string

func makeQueryString(p params) string {
	form := url.Values{}
	for key, val := range p {
		form.Set(key, val)
	}
	return form.Encode()
}

func makeRequestStub(url, body string) *http.Request {
	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(body))
	req.Header.Add("User-Agent", "go-test")
	req.Header.Add("Host", "localhost:8080")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-LEngth", string(len(body)))

	return req
}

func ExampleBindPostForm() {
	type ContactForm struct {
		Name    string `form:"name"`
		Email   string `form:"email"`
		Message string `form:"message"`
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		form := &ContactForm{}

		err := bindform.BindPostForm(r, form)
		if err != nil {
			http.Error(w, "Bad Request", 400)
		}

		// Write more code...
	})
}

func TestBindForm(t *testing.T) {
	type target struct {
		V1 string `form:"v1"`
		V2 string `form:"v2"`
	}

	tests := map[string][]string{
		makeQueryString(params{"v2": "world"}):                 []string{"hello", "world"},
		makeQueryString(params{"v1": "goodby", "v2": "world"}): []string{"goodby", "world"},
	}

	for query, expect := range tests {
		req := makeRequestStub("http://localhost/?v1=hello", query)
		target := &target{}
		err := bindform.BindForm(req, target)

		t.Logf("%-10v: %+v", query, target)
		if err != nil {
			t.Errorf("Got Error: %v", err)
		}
		if got := target.V1; got != expect[0] {
			t.Errorf("Got: %v, But Expect: %v", got, expect[0])
		}
		if got := target.V2; got != expect[1] {
			t.Errorf("Got: %v, But Expect: %v", got, expect[1])
		}
	}
}

func TestBindPostForm(t *testing.T) {
	type target struct {
		V1 string `form:"v1"`
		V2 string `form:"v2"`
	}

	tests := map[string][]string{
		makeQueryString(params{"v2": "world"}): []string{"", "world"},
	}

	for query, expect := range tests {
		req := makeRequestStub("http://localhost/?v1=hello", query)
		target := &target{}
		err := bindform.BindPostForm(req, target)

		t.Logf("%-10v: %+v", query, target)
		if err != nil {
			t.Errorf("Got Error: %v", err)
		}
		if got := target.V1; got != expect[0] {
			t.Errorf("Got: %v, But Expect: %v", got, expect[0])
		}
		if got := target.V2; got != expect[1] {
			t.Errorf("Got: %v, But Expect: %v", got, expect[1])
		}
	}
}

func TestBindBool(t *testing.T) {
	type target struct {
		V bool `form:"v"`
	}

	tests := map[string]bool{
		makeQueryString(params{"v": "TRUE"}):  true,
		makeQueryString(params{"v": "True"}):  true,
		makeQueryString(params{"v": "true"}):  true,
		makeQueryString(params{"v": "T"}):     true,
		makeQueryString(params{"v": "t"}):     true,
		makeQueryString(params{"v": "1"}):     true,
		makeQueryString(params{"v": "FALSE"}): false,
		makeQueryString(params{"v": "False"}): false,
		makeQueryString(params{"v": "false"}): false,
		makeQueryString(params{"v": "f"}):     false,
		makeQueryString(params{"v": "f"}):     false,
		makeQueryString(params{"v": "0"}):     false,
		makeQueryString(params{"v": ""}):      false,
		makeQueryString(params{}):             false,
	}

	for query, expect := range tests {
		req := makeRequestStub("http://localhost/", query)
		target := &target{}
		err := bindform.BindPostForm(req, target)

		t.Logf("%-10v: %+v", query, target)
		if err != nil {
			t.Errorf("Got Error: %v", err)
		}
		if got := target.V; got != expect {
			t.Errorf("Got: %v, But Expect: %v", got, expect)
		}
	}
}

func TestBindBoolFail(t *testing.T) {
	type target struct {
		V bool `form:"v"`
	}

	tests := map[string]bool{
		makeQueryString(params{"v": "2"}):       false,
		makeQueryString(params{"v": "invalid"}): false,
	}

	for query, expect := range tests {
		req := makeRequestStub("http://localhost/", query)
		target := &target{}
		err := bindform.BindPostForm(req, target)

		t.Logf("%-10v: %+v with error `%+v`", query, target, err)
		if err == nil {
			t.Errorf("Got no error")
		}
		if got := target.V; got != expect {
			t.Errorf("Got: %v, But Expect: %v", got, expect)
		}
	}
}

func TestBindInt(t *testing.T) {
	type target struct {
		V int `form:"v"`
	}

	tests := map[string]int{
		makeQueryString(params{"v": "12345"}):  12345,
		makeQueryString(params{"v": "-12345"}): -12345,
		makeQueryString(params{"v": ""}):       0,
		makeQueryString(params{}):              0,
	}

	for query, expect := range tests {
		req := makeRequestStub("http://localhost/", query)
		target := &target{}
		err := bindform.BindPostForm(req, target)

		t.Logf("%-10v: %+v", query, target)
		if err != nil {
			t.Errorf("Got Error: %v", err)
		}
		if got := target.V; got != expect {
			t.Errorf("Got: %v, But Expect: %v", got, expect)
		}
	}
}

func TestBindIntFail(t *testing.T) {
	type target struct {
		V int `form:"v"`
	}
	tests := map[string]int{
		makeQueryString(params{"v": "1.5"}): 0,
		makeQueryString(params{"v": "a"}):   0,
	}

	for query, expect := range tests {
		req := makeRequestStub("http://localhost/", query)
		target := &target{}
		err := bindform.BindPostForm(req, target)

		t.Logf("%-10v: %+v with error `%+v`", query, target, err)
		if err == nil {
			t.Errorf("Got no error")
		}
		if got := target.V; got != expect {
			t.Errorf("Got: %v, But Expect: %v", got, expect)
		}
	}
}

func TestBindUInt(t *testing.T) {
	type target struct {
		V uint `form:"v"`
	}
	tests := map[string]uint{
		makeQueryString(params{"v": "12345"}): 12345,
		makeQueryString(params{"v": ""}):      0,
		makeQueryString(params{}):             0,
	}

	for query, expect := range tests {
		req := makeRequestStub("http://localhost/", query)
		target := &target{}
		err := bindform.BindPostForm(req, target)

		t.Logf("%-10v: %+v", query, target)
		if err != nil {
			t.Errorf("Got Error: %v", err)
		}
		if got := target.V; got != expect {
			t.Errorf("Got: %v, But Expect: %v", got, expect)
		}
	}
}

func TestBindUIntFail(t *testing.T) {
	type target struct {
		V uint `form:"v"`
	}
	tests := map[string]uint{
		makeQueryString(params{"v": "-12345"}): 0,
		makeQueryString(params{"v": "1.5"}):    0,
		makeQueryString(params{"v": "a"}):      0,
	}

	for query, expect := range tests {
		req := makeRequestStub("http://localhost/", query)
		target := &target{}
		err := bindform.BindPostForm(req, target)

		t.Logf("%-10v: %+v with error `%+v`", query, target, err)
		if err == nil {
			t.Errorf("Got no error")
		}
		if got := target.V; got != expect {
			t.Errorf("Got: %v, But Expect: %v", got, expect)
		}
	}
}

func TestBindFloat32(t *testing.T) {
	type target struct {
		V float32 `form:"v"`
	}

	maxE := strconv.FormatFloat(1e+38, 'e', 46, 64)
	maxF := strconv.FormatFloat(1e+38, 'f', 46, 64)
	minE := strconv.FormatFloat(1e-45, 'e', 45, 64)
	minF := strconv.FormatFloat(1e-45, 'f', 45, 64)
	verysmall := strconv.FormatFloat(1e-46, 'f', 46, 64)

	tests := map[string]float32{
		makeQueryString(params{"v": "12345"}):   12345.0,
		makeQueryString(params{"v": "-12345"}):  -12345.0,
		makeQueryString(params{"v": "123.45"}):  123.45,
		makeQueryString(params{"v": "-123.45"}): -123.45,
		makeQueryString(params{"v": maxE}):      1e+38,
		makeQueryString(params{"v": maxF}):      1e+38,
		makeQueryString(params{"v": minE}):      1e-45,
		makeQueryString(params{"v": minF}):      1e-45,
		makeQueryString(params{"v": verysmall}): 0.0,
		makeQueryString(params{"v": ""}):        0.0,
		makeQueryString(params{}):               0.0,
	}

	for query, expect := range tests {
		req := makeRequestStub("http://localhost/", query)
		target := &target{}
		err := bindform.BindPostForm(req, target)

		t.Logf("%-10v: %+v", query, target)
		if err != nil {
			t.Errorf("Got Error: %v", err)
		}
		if got := target.V; got != expect {
			t.Errorf("Got: %v, But Expect: %v", got, expect)
		}
	}
}

func TestBindFloat32Fail(t *testing.T) {
	type target struct {
		V float32 `form:"v"`
	}

	tests := map[string]float32{
		makeQueryString(params{"v": "a"}): 0.0,
	}

	for query, expect := range tests {
		req := makeRequestStub("http://localhost/", query)
		target := &target{}
		err := bindform.BindPostForm(req, target)

		t.Logf("%-10v: %+v with error `%+v`", query, target, err)
		if err == nil {
			t.Errorf("Got no error")
		}
		if got := target.V; got != expect {
			t.Errorf("Got: %v, But Expect: %v", got, expect)
		}
	}
}

func TestBindString(t *testing.T) {
	type target struct {
		V string `form:"v"`
	}

	tests := map[string]string{
		makeQueryString(params{"v": "hello"}): "hello",
		makeQueryString(params{"v": ""}):      "",
		makeQueryString(params{}):             "",
	}

	for query, expect := range tests {
		req := makeRequestStub("http://localhost/", query)
		target := &target{}
		err := bindform.BindPostForm(req, target)

		t.Logf("%-10v: %+v", query, target)
		if err != nil {
			t.Errorf("Got Error: %v", err)
		}
		if got := target.V; got != expect {
			t.Errorf("Got: %v, But Expect: %v", got, expect)
		}
	}
}
