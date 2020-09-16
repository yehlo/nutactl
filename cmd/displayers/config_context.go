// Copyright © 2020 Joshua Leuenberger
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package displayers

import (
	"github.com/simonfuhrer/nutactl/pkg"
	"io"
	"strconv"
	"fmt"
)

// ConfigContexts ...
type ConfigContexts struct {
	config.ContextList
}

func (o ConfigContexts) JSON(w io.Writer) error {
	return displayJSON(w, o.Entities)
}

func (o ConfigContexts) JSONPath(w io.Writer, template string) error {
	return displayJSONPath(w, template, o.Entities)
}

func (o ConfigContexts) PP(w io.Writer) error {
	return displayPP(w, o.Entities)
}

func (o ConfigContexts) YAML(w io.Writer) error {
	return displayYAML(w, o.Entities)
}

func (o ConfigContexts) header() []string {
	return []string{
		"ID",
		"URL",
		"User",
		"Insecure",
	}
}

func (o ConfigContexts) TableData(w io.Writer) error {
	data := make([][]string, len(o.Entities))
	for i, configContext := range o.Entities {
		// converting to strings
		id := fmt.Sprint(configContext.ID)
		insecure := strconv.FormatBool(configContext.Insecure)
		data[i] = []string{
			id,
			configContext.URL,
			configContext.User,
			insecure,
		}
	}
	return displayTable(w, data, o.header())
}

func (o ConfigContexts) Text(w io.Writer) error {
	return nil
}
