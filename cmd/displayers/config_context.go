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
	"io"
)

// ConfigContext ...
type ConfigContext struct {
	ID, URL, User, Insecure	string
}

// ConfigContextSlice slice of ConfigContext structs
type ConfigContextSlice []ConfigContext

func (o ConfigContextSlice) JSON(w io.Writer) error {
	return displayJSON(w, o)
}

func (o ConfigContextSlice) JSONPath(w io.Writer, template string) error {
	return displayJSONPath(w, template, o)
}

func (o ConfigContextSlice) PP(w io.Writer) error {
	return displayPP(w, o)
}

func (o ConfigContextSlice) YAML(w io.Writer) error {
	return displayYAML(w, o)
}

func (o ConfigContextSlice) header() []string {
	return []string{
		"ID",
		"URL",
		"User",
		"Insecure",
	}
}

func (o ConfigContextSlice) TableData(w io.Writer) error {
	data := make([][]string, len(o))
	for i, configContext := range o {
		data[i] = []string{
			configContext.ID,
			configContext.URL,
			configContext.User,
			configContext.Insecure,
		}
	}
	return displayTable(w, data, o.header())
}

func (o ConfigContextSlice) Text(w io.Writer) error {
	return nil
}
