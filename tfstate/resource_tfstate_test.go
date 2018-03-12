package tfstate

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/rocksolidlabs/jsonq"
)

func TestAcc_TFStateOutputsBasic(t *testing.T) {
	jsonMap := map[string]interface{}{}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTFStateOutputsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testValidJSONCheck("tfstate_outputs.outtahere", jsonMap),
					testJSONFieldExistsCheck(jsonMap, "version"),
					testJSONFieldExistsCheck(jsonMap, "serial"),
					testJSONFieldExistsCheck(jsonMap, "lineage"),
					testJSONFieldCheck(jsonMap, "modules.0.path.0", "root"),
					testJSONFieldCheck(jsonMap, "modules[0].outputs.wowow.value", "oohlala"),
					testJSONFieldCheck(jsonMap, "modules[0].outputs.bar.value", "baz"),
					testJSONFieldCheck(jsonMap, "modules[0].outputs.bar.sensitive", true),
				),
			},
		},
	})
}

func TestAcc_TFStateOutputsUpdate(t *testing.T) {
	jsonMap := map[string]interface{}{}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTFStateOutputsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testValidJSONCheck("tfstate_outputs.outtahere", jsonMap),
					testJSONFieldExistsCheck(jsonMap, "version"),
					testJSONFieldCheck(jsonMap, "serial", float64(1)),
					testJSONFieldExistsCheck(jsonMap, "serial"),
					testJSONFieldExistsCheck(jsonMap, "lineage"),
					testJSONFieldCheck(jsonMap, "modules.0.path.0", "root"),
					testJSONFieldCheck(jsonMap, "modules[0].outputs.wowow.value", "oohlala"),
					testJSONFieldCheck(jsonMap, "modules[0].outputs.bar.value", "baz"),
					testJSONFieldCheck(jsonMap, "modules[0].outputs.bar.sensitive", true),
				),
			},
			{
				Config: testAccTFStateOutputsModified,
				Check: resource.ComposeAggregateTestCheckFunc(
					testValidJSONCheck("tfstate_outputs.outtahere", jsonMap),
					testJSONFieldExistsCheck(jsonMap, "version"),
					testJSONFieldCheck(jsonMap, "serial", float64(2)),
					testJSONFieldExistsCheck(jsonMap, "lineage"),
					testJSONFieldCheck(jsonMap, "modules.0.path.0", "root"),
					testJSONFieldCheck(jsonMap, "modules[0].outputs.wowow.value", "doge"),
					testJSONFieldCheck(jsonMap, "modules[0].outputs.bar.value", "baz"),
					testJSONFieldCheck(jsonMap, "modules[0].outputs.bar.sensitive", true),
				),
			},
		},
	})
}

func testValidJSONCheck(id string, jsonMap map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		jsonStr := rs.Primary.Attributes["json"]

		err := json.Unmarshal([]byte(jsonStr), &jsonMap)
		if err != nil {
			return err
		}

		return nil
	}
}

func testJSONFieldExistsCheck(jsonMap map[string]interface{}, path ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		jq := jsonq.NewQuery(jsonMap)
		_, err := jq.Interface(path...)
		if err != nil {
			return err
		}

		return nil
	}
}

func testJSONFieldCheck(jsonMap map[string]interface{}, path string, value interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		jq := jsonq.NewQuery(jsonMap)
		val, err := jq.Interface(path)
		if err != nil {
			return err
		}
		if val != value {
			return fmt.Errorf("json value [%v] did not match expected value [%v]", val, value)
		}

		return nil
	}
}

const (
	testAccTFStateOutputsConfig = `
resource "tfstate_outputs" "outtahere" {
  output {
		name  = "wowow"
		value = "oohlala"
	}
  output {
		name  		= "bar"
		value 		= "baz"
		sensitive = true
	}
}
`

	testAccTFStateOutputsModified = `
resource "tfstate_outputs" "outtahere" {
output {
	name  = "wowow"
	value = "doge"
}
output {
	name  		= "bar"
	value 		= "baz"
	sensitive = true
}
}
`
)
