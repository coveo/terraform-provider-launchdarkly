package launchdarkly

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)

	name := d.Get("name").(string)
	key := d.Get("key").(string)

	payload := map[string]string{
		"name": name,
		"key":  key,
	}

	_, err := client.Post(getProjectCreateUrl(), payload, []int{201})
	if err != nil {
		return err
	}

	// Default environments will be created, we want to get rid of those
	environmentKeys, err := getEnvironmentKeys(client, key)
	if err != nil {
		return err
	}

	err = ensureThereIsADummyEnvironment(client, key)
	if err != nil {
		return err
	}

	for _, environmentKey := range environmentKeys {
		err = client.Delete(getEnvironmentUrl(key, environmentKey), []int{204})
		if err != nil {
			return err
		}
	}

	d.SetId(key)
	d.Set("name", name)
	d.Set("key", key)

	return nil
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	key := d.Get("key").(string)

	client := m.(Client)

	raw, err := client.Get(getProjectUrl(key), []int{200})
	if err != nil {
		d.SetId("")
		return nil
	}

	payload := raw.(map[string]interface{})
	d.Set("name", payload["name"])
	d.Set("key", payload["key"])

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)

	name := d.Get("name").(string)

	payload := []map[string]string{{
		"op":    "replace",
		"path":  "/name",
		"value": name,
	}}

	_, err := client.Patch(getProjectUrl(d.Id()), payload, []int{200})
	if err != nil {
		return err
	}

	return nil
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(Client)

	err := client.Delete(getProjectUrl(d.Id()), []int{204, 404})
	if err != nil {
		return err
	}

	return nil
}
