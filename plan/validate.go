package plan

import (
	"bytes"
	"html/template"
	"log"
	"strconv"
	"strings"
)

func Validate(tmpl string, input interface{}) bool {
	t := template.Must(template.New("validate").Parse(tmpl))
	b := new(bytes.Buffer)
	if err := t.Execute(b, input); err != nil {
		log.Panic(err)
	}

	out := strings.TrimSpace(b.String())
	if out == "" {
		return false
	}

	valid, err := strconv.ParseBool(out)
	if err != nil {
		log.Panic(err)
	}

	return valid
}
