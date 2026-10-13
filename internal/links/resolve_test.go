package links

import (
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestShortNames(t *testing.T) {
	if got := Normalize("Meeting-Notes"); got != "meetingnotes" {
		t.Fatal(got)
	}
	for _, name := range []string{"a", "Meeting-Notes", "a.b", "_test", "_api", "api", "API", "a-p-i"} {
		if err := ValidateShort(name); err != nil {
			t.Errorf("%q: %v", name, err)
		}
	}
	for _, name := range []string{"", ".help", "a/b", "a b", "a+", "-a", "_", "_-", "_--"} {
		if ValidateShort(name) == nil {
			t.Errorf("accepted %q", name)
		}
	}
}

func TestExpand(t *testing.T) {
	cases := []struct{ name, long, path, query, want string }{
		{"escaped", "https://host/a%2f/b", "extra", "", "https://host/a%2f/b/extra"},
		{"slash", "https://host/a/", "b//c", "", "https://host/a/b//c"},
		{"query", "/rel?a=1", "", "a=2&b=2", "/rel?a=1&a=2&b=2"},
		{"legacy query append", "/rel?a=1", "path", "", "/rel?a=1/path"},
		{"legacy fragment append", "https://host/#frag", "path", "", "https://host/#frag/path"},
		{"path escape", "https://host/{{PathEscape .Path}}", "a/b+c", "", "https://host/a%2Fb+c"},
		{"query escape", "https://host/?q={{QueryEscape .Path}}", "a/b+c", "x=1", "https://host/?q=a%2Fb%2Bc&x=1"},
		{"prefix", `https://host/{{TrimPrefix .Path "BUG-"}}`, "BUG-12", "", "https://host/12"},
		{"suffix", `https://host/{{TrimSuffix .Path "/"}}`, "abc/", "", "https://host/abc"},
		{"lower", "https://host/{{ToLower .Path}}", "ABC", "", "https://host/abc"},
		{"upper", "https://host/{{ToUpper .Path}}", "abc", "", "https://host/ABC"},
		{"match", `https://host/{{if Match "^[0-9]+$" .Path}}id/{{else}}search/{{end}}{{.Path}}`, "12", "", "https://host/id/12"},
		{"bad regex", `https://host/{{if Match "[" .Path}}yes{{else}}no{{end}}`, "", "", "https://host/no"},
		{"date", `https://host/{{.Now.Format "2006-01-02"}}`, "", "", "https://host/2026-10-13"},
		{"relative", "/other", "sub", "", "/other/sub"},
		{"protocol relative", "//host/", "", "", "//host/"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			query, _ := url.ParseQuery(tc.query)
			got, err := Expand(tc.long, Environment{Path: tc.path, Query: query, Now: time.Date(2026, 10, 13, 0, 0, 0, 0, time.UTC)})
			if err != nil {
				t.Fatal(err)
			}
			if got.String() != tc.want {
				t.Fatalf("got %q, want %q", got.String(), tc.want)
			}
		})
	}
}

func TestInvalidDestinations(t *testing.T) {
	for _, destination := range []string{"", "javascript:alert(1)", "data:text/plain,a", "https://host/\n", "https://host/{{"} {
		if ValidateDestination(destination) == nil {
			t.Errorf("accepted %q", destination)
		}
	}
	for _, destination := range []string{"https://host/{{.Unknown}}", "{{.Path}}:alert(1)", "https://host/{{.Path}}", "{{range .Path}}x{{end}}"} {
		path := "javascript"
		if strings.Contains(destination, "host/{{.Path}}") {
			path = "\n"
		}
		if _, err := Expand(destination, Environment{Path: path}); err == nil {
			t.Errorf("expanded unsafe %q", destination)
		}
	}
	if _, err := Expand("https://host/{{.Path}}", Environment{Path: strings.Repeat("a", 65536)}); err == nil {
		t.Fatal("unbounded expansion")
	}
}

func TestTemplateDoesNotProvideIdentity(t *testing.T) {
	if _, err := Expand("https://host/{{.User}}", Environment{Path: "docs"}); err == nil {
		t.Fatal("templates must not provide a user field")
	}
}
