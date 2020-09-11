// Copyright © 2020 Serdar KÖKTAŞ <contact@serdarkoktas.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package esc

import (
	"context"
	"github.com/elastic/go-elasticsearch/v7"
)

//Search sends the search request to API.
func Search(es *elasticsearch.Client, size int, index string, query string, fields []string, sort []string) (string, error) {
	//Send request
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithSize(size),
		es.Search.WithIndex(index),
		es.Search.WithQuery(query),
		es.Search.WithFilterPath("hits.hits._source"),
		es.Search.WithSource(fields...),
		es.Search.WithSort(sort...),
	)
	//Read request
	if err == nil {
		json := jread(res.Body)
		defer res.Body.Close()
		return json, err
	}
	return "", err
}
