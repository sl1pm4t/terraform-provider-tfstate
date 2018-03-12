package tfstate

import (
	"bytes"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/hashicorp/terraform/terraform"
)

const FAKE_LINEAGE = "00000000-1111-2222-3333-444444444444"

func resourceTFStateOutputs() *schema.Resource {
	return &schema.Resource{
		Create:        resourceStateOutputsCreate,
		Read:          ReadJSON,
		Update:        resourceStateOutputsCreate,
		Delete:        schema.RemoveFromState,
		CustomizeDiff: resourceTFStateOutputsCustomDiff,

		Schema: map[string]*schema.Schema{
			"output": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"sensitive": {
							// not sure if this is useful, but we have the info so provide it in case it's useful downstream
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"type": {
							// This is included for future use, currently only possible to use string.
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "string",
							ValidateFunc: validation.StringInSlice([]string{"string"}, false),
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"serial": {
				Type:         schema.TypeInt,
				Computed:     true,
				ComputedWhen: []string{"output"},
			},
		},
	}
}

func generateState(outputs []interface{}, serial int) (string, error) {
	buf := new(bytes.Buffer)
	if len(outputs) > 0 {
		state := &terraform.State{}
		state.Init()
		state.Lineage = FAKE_LINEAGE
		state.Serial = int64(serial)
		for _, rawOut := range outputs {
			out := rawOut.(map[string]interface{})
			state.RootModule().Outputs[out["name"].(string)] = &terraform.OutputState{
				Type:      out["type"].(string),
				Sensitive: out["sensitive"].(bool),
				Value:     out["value"],
			}
		}

		if err := terraform.WriteState(state, buf); err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}

func resourceStateOutputsCreate(d *schema.ResourceData, meta interface{}) error {
	state, err := generateState(d.Get("output").([]interface{}), d.Get("serial").(int))
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] State: %s", state)

	d.Set("json", state)

	d.SetId(strconv.Itoa(hashcode.String(state)))
	return nil
}

func ReadJSON(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceTFStateOutputsCustomDiff(d *schema.ResourceDiff, meta interface{}) error {
	if d.HasChange("output") {
		d.SetNew("serial", d.Get("serial").(int)+1)

		newJSON, err := generateState(d.Get("output").([]interface{}), d.Get("serial").(int))
		if err != nil {
			return err
		}
		d.SetNew("json", newJSON)
	}

	return nil
}
