package ud

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func Test_findId(t *testing.T) {
	cases := []struct {
		in  string
		exp string
	}{
		{`<em id="who">github</em>`, "who"},
		{`<em>github</em>`, ""},
	}
	for _, c := range cases {
		got := findId([]byte(c.in))
		if got != c.exp {
			t.Error(got, c.exp)
		}
	}
}

func TestReplace(t *testing.T) {
	cases := []struct {
		desc         string
		doc          string
		id           string
		frag         string
		exp          string
		replaceChild bool
	}{
		{
			desc:         "replace element by given id when element is empty",
			doc:          `<b><i id="x"></i></b>`, // empty start
			id:           "x",
			frag:         `content`, // ok id
			exp:          `<b><i id="x">content</i></b>`,
			replaceChild: true,
		},
		{
			desc:         "element when id is found in fragment",
			doc:          `<b><i id="x">a</i></b>`,
			id:           "",
			frag:         `<em id="x">A</em>`, // ok id
			exp:          `<b><em id="x">A</em></b>`,
			replaceChild: false,
		},
		{
			desc:         "do nothing if no matching id is found",
			doc:          `<b><i id="x">a</i></b>`,
			id:           "",
			frag:         `<i id="X"></i>`, // wrong id
			exp:          `<b><i id="x">a</i></b>`,
			replaceChild: false,
		},
		{
			desc:         "remove element when fragment is empty",
			doc:          `<b><i id="x">a</i></b>`,
			id:           "x",
			frag:         "",
			exp:          `<b></b>`,
			replaceChild: false,
		},
		{
			desc:         "remove all sub elements when fragment is empty",
			doc:          `<b><i id="x"><span>here</span></i></b>`,
			id:           "x",
			frag:         "",
			exp:          `<b></b>`,
			replaceChild: false,
		},
		{
			desc:         "child when id given",
			doc:          `<b><i id="x">a</i></b>`,
			id:           "x",
			frag:         `A`,
			exp:          `<b><i id="x">A</i></b>`,
			replaceChild: true,
		},
		{
			desc:         "Empty id cannot replace child",
			doc:          `<b><i id="x">a</i></b>`,
			id:           "", // empty
			frag:         `<i id="x">A</i>`,
			exp:          `<b><i id="x">A</i></b>`,
			replaceChild: true, // no effect when no id is given
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			w := bytes.NewBufferString("")
			hr := NewHtmlRewriter(c.id, c.replaceChild, []byte(c.frag))
			hr.Rewrite(w, strings.NewReader(c.doc))
			got := w.String()
			if got != c.exp {
				t.Logf(`
Document to modify         %s
given id                   %q
The fragment               %s
Expected to see            %s
but got                    %s`, c.doc, c.id, c.frag, c.exp, got)
				t.Fail()
			}
		})
	}
}

func Test_skip(t *testing.T) {
	cases := []struct {
		part string
		exp  string
	}{
		{
			part: `<i></i></b>`,
			exp:  `b`,
		},
		{
			part: `</b>`,
			exp:  `b`,
		},
	}
	for _, c := range cases {
		z := html.NewTokenizer(strings.NewReader(c.part))
		got := skip(z)
		if got.Data != c.exp {
			t.Error(got)
		}
	}
}

func ExampleHtmlRewriter_Rewrite_element() {
	r := strings.NewReader(`
<html>
 <body>
  <h1>Title</h1>
  <div id="ABC">Old content</div>
 </body>
</html>
`)
	hr := NewHtmlRewriter("ABC", false, []byte(`new stuff`))
	hr.Rewrite(os.Stdout, r)
	// output:
	// <html>
	//  <body>
	//   <h1>Title</h1>
	//   new stuff
	//  </body>
	// </html>
}

func ExampleHtmlRewriter_Rewrite_child() {
	r := strings.NewReader(`
<html>
 <body>
  <h1>Title</h1>
  <div id="ABC">Old content</div>
 </body>
</html>
`)
	hr := NewHtmlRewriter("ABC", true, []byte(`new stuff`))
	hr.Rewrite(os.Stdout, r)
	// output:
	// <html>
	//  <body>
	//   <h1>Title</h1>
	//   <div id="ABC">new stuff</div>
	//  </body>
	// </html>
}
