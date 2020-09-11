package esc

/*
Copyright © 2020 Serdar KÖKTAŞ <contact@serdarkoktas.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

import (
	"bytes"
	"github.com/elastic/go-elasticsearch/v7"
	"io"
)

// Stats sends request to cluster stat API.
func Stats(es *elasticsearch.Client, nodeFilter string) (string, error) {
	// Send request
	res, err := es.Cluster.Stats(es.Cluster.Stats.WithNodeID(nodeFilter), es.Cluster.Stats.WithHuman())
	// Read request
	if err == nil {
		json := jread(res.Body)
		defer res.Body.Close()
		return json, err
	}
	return "", err
}

// Read the values.
func jread(r io.Reader) string {
	var b bytes.Buffer
	b.ReadFrom(r)
	return b.String()
}
