// SPDX-License-Identifier: Apache-2.0

// Copyright 2026 Cisco Systems, Inc. and their affiliates

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package html

import (
	"fmt"
	"html/template"
	"time"

	"github.com/ciscotools/docgen/internal/log"

	"github.com/segmentio/encoding/json"
	"github.com/tidwall/gjson"
)

func debugHelper(a any) string {
	res, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(gjson.GetBytes(res, "@pretty"))
	return ""
}

func dateHelper() string {
	return time.Now().Format("2 Jan 2006")
}

func (doc *Document) includeHelper(name string, ctxs ...any) template.HTML {
	tpl, ok := doc.Templates[name]
	if !ok {
		log.Warn().Str("name", name).Msg("include template  not found")
		return template.HTML("")
	}
	ctx := doc.ctx
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	}
	val, err := tpl.Render(ctx)
	if err != nil {
		log.Warn().Err(err).Str("name", name).Msg("error rendering template")
		return template.HTML("")
	}
	return template.HTML(val)
}

func (doc *Document) rootHelper() any {
	return doc.ctx
}
