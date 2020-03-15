// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package template // import "miniflux.app/template"

import (
	"bytes"
	"html/template"
	"time"

	"miniflux.app/errors"
	"miniflux.app/locale"
	"miniflux.app/logger"

	"github.com/gorilla/mux"
)

// Engine handles the templating system.
type Engine struct {
	templates map[string]*template.Template
	funcMap   *funcMap
}

func (e *Engine) parseAll() {
	commonTemplates := ""
	for _, content := range templateCommonMap {
		commonTemplates += content
	}

	for name, content := range templateViewsMap {
		logger.Debug("[Template] Parsing: %s", name)
		e.templates[name] = template.Must(template.New("main").Funcs(e.funcMap.Map()).Parse(commonTemplates + content))
	}
}

// Render process a template.
func (e *Engine) Render(name, language string, data interface{}) []byte {
	tpl, ok := e.templates[name]
	if !ok {
		logger.Fatal("[Template] The template %s does not exists", name)
	}

	printer := locale.NewPrinter(language)

	// Functions that need to be declared at runtime.
	tpl.Funcs(template.FuncMap{
		"elapsed": func(timezone string, t time.Time) string {
			return elapsedTime(printer, timezone, t)
		},
		"t": func(key interface{}, args ...interface{}) string {
			switch k := key.(type) {
			case string:
				return printer.Printf(k, args...)
			case errors.LocalizedError:
				return k.Localize(printer)
			case *errors.LocalizedError:
				return k.Localize(printer)
			case error:
				return k.Error()
			default:
				return ""
			}
		},
		"plural": func(key string, n int, args ...interface{}) string {
			return printer.Plural(key, n, args...)
		},
	})

	var b bytes.Buffer
	err := tpl.ExecuteTemplate(&b, "base", data)
	if err != nil {
		logger.Fatal("[Template] Unable to render template: %v", err)
	}

	return b.Bytes()
}

// NewEngine returns a new template engine.
func NewEngine(router *mux.Router) *Engine {
	tpl := &Engine{
		templates: make(map[string]*template.Template),
		funcMap:   &funcMap{router},
	}

	tpl.parseAll()
	return tpl
}
