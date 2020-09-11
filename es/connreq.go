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
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
	"log"
)

// Connect creates es client.
func Connect() *elasticsearch.Client {
	// Read config file
	elkusr := viper.GetString("username")
	elkpass := viper.GetString("password")

	//Create client with basic auth
	if len(elkusr) > 1 {
		// Password prompt
		if len(elkpass) < 1 {
			var err error
			templates := &promptui.PromptTemplates{
				Prompt: "{{ . }} ",
				Valid:  "{{ . | red }} ",
			}

			prompt := promptui.Prompt{
				Label:     "Password for " + elkusr + ": ",
				Mask:      '*',
				Templates: templates,
			}

			elkpass, err = prompt.Run()

			if err != nil {
				log.Fatalf("ERROR: %v", err)
			}
		}
		// Connection with auth configuration.
		cfg := elasticsearch.Config{
			Addresses: viper.GetStringSlice("urls"),
			Username:  elkusr,
			Password:  elkpass,
		}
		es, _ := elasticsearch.NewClient(cfg)
		res, err := es.Info()

		if err != nil {
			log.Fatalf("ERROR: %v", err)
		}
		if res.IsError() == true {
			log.Fatalf("ERROR: %v", res.Status())
		}
		defer res.Body.Close()
		return es
	}
	// Create config.
	cfg := elasticsearch.Config{
		Addresses: viper.GetStringSlice("urls"),
	}

	//Create client.
	es, _ := elasticsearch.NewClient(cfg)
	res, err := es.Info()

	//Handle errors
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	if res.IsError() == true {
		log.Fatalf("ERROR: %v", res.Status())
	}

	//Close response body
	defer res.Body.Close()
	return es

}
