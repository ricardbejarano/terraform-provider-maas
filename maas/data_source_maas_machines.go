package maas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMAASMachines() *schema.Resource {
	return &schema.Resource{
		Description: "Lists MAAS machines visible to the user.",
		ReadContext: dataSourceMachinesRead,

		Schema: map[string]*schema.Schema{
			"machines": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "A set of machines visible to the user.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"architecture": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The architecture type of the machine.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The domain of the machine.",
						},
						"hostname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The machine hostname.",
						},
						"min_hwe_kernel": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The minimum kernel version allowed to run on this machine.",
						},
						"pool": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource pool of the machine.",
						},
						"power_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The power management type (e.g. `ipmi`) of the machine.",
						},
						"pxe_mac_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The MAC address of the machine's PXE boot NIC.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The zone of the machine.",
						},
					},
				},
			},
		},
	}
}

func dataSourceMachinesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ClientConfig).Client

	machines, err := client.Machines.Get(nil)
	if err != nil {
		return diag.FromErr(err)
	}

	items := []map[string]interface{}{}
	for _, machine := range machines {
		items = append(items, map[string]interface{}{
			"architecture":    machine.Architecture,
			"domain":          machine.Domain.Name,
			"hostname":        machine.Hostname,
			"min_hwe_kernel":  machine.MinHWEKernel,
			"pool":            machine.Pool.Name,
			"power_type":      machine.PowerType,
			"pxe_mac_address": machine.BootInterface.MACAddress,
			"zone":            machine.Zone.Name,
		})
	}

	if err := d.Set("machines", items); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(machines[0].SystemID)

	return nil
}
