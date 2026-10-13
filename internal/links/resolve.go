package links

import (
	"bytes"
	"errors"
	"net/url"
	"regexp"
	"strings"
	"text/template"
)

const MaxDestination = 64 * 1024

var helpers = template.FuncMap{
	"PathEscape": url.PathEscape, "QueryEscape": url.QueryEscape,
	"TrimPrefix": strings.TrimPrefix, "TrimSuffix": strings.TrimSuffix,
	"ToLower": strings.ToLower, "ToUpper": strings.ToUpper,
	"Match": func(pattern, s string) bool { ok, _ := regexp.MatchString(pattern, s); return ok },
}

func parse(s string) (*template.Template, error) {
	return template.New("destination").Funcs(helpers).Parse(s)
}

type boundedBuffer struct{ bytes.Buffer }

func (b *boundedBuffer) Write(p []byte) (int, error) {
	if b.Len()+len(p) > MaxDestination {
		return 0, errors.New("Destination exceeds 64 KiB")
	}
	return b.Buffer.Write(p)
}

func Expand(destination string, env Environment) (*url.URL, error) {
	if err := ValidateDestination(destination); err != nil {
		return nil, err
	}
	if !strings.Contains(destination, "{{") {
		if strings.HasSuffix(destination, "/") {
			destination += "{{.Path}}"
		} else {
			destination += "{{with .Path}}/{{.}}{{end}}"
		}
	}
	tmpl, err := parse(destination)
	if err != nil {
		return nil, err
	}
	var out boundedBuffer
	env.Now = env.Now.UTC()
	if err := tmpl.Execute(&out, env); err != nil {
		return nil, errors.New("Could not expand destination: " + err.Error())
	}
	u, err := destinationURL(out.String())
	if err != nil {
		return nil, err
	}
	if len(env.Query) > 0 {
		query, err := url.ParseQuery(u.RawQuery)
		if err != nil {
			return nil, errors.New("Destination has an invalid query string")
		}
		for key, values := range env.Query {
			for _, value := range values {
				query.Add(key, value)
			}
		}
		u.RawQuery = query.Encode()
	}
	if len(u.String()) > MaxDestination {
		return nil, errors.New("Destination exceeds 64 KiB")
	}
	return u, nil
}
