package maas_test

import (
	"fmt"
	"strconv"
	"terraform-provider-maas/maas"
	"terraform-provider-maas/maas/testutils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceMAASMachines_basic(t *testing.T) {
	checks := []resource.TestCheckFunc{
		func(s *terraform.State) error {
			conn := testutils.TestAccProvider.Meta().(*maas.ClientConfig).Client

			machines, err := conn.Machines.Get(nil)
			if err != nil {
				return err
			}

			if err := resource.TestCheckResourceAttr("data.maas_machines.test", "machines.#", strconv.Itoa(len(machines)))(s); err != nil {
				return err
			}
			// According to MAAS' API documentation, machines are returned "sorted by id (i.e. most recent last)"
			for i, machine := range machines {
				if err := resource.TestCheckResourceAttr("data.maas_machines.test", fmt.Sprintf("machines.%v.system_id", i), machine.SystemID)(s); err != nil {
					return err
				}
				if err := resource.TestCheckResourceAttr("data.maas_machines.test", fmt.Sprintf("machines.%v.hostname", i), machine.Hostname)(s); err != nil {
					return err
				}
			}

			return nil
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { testutils.PreCheck(t, nil) },
		Providers:  testutils.TestAccProviders,
		ErrorCheck: func(err error) error { return err },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMAASMachines(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

func testAccDataSourceMAASMachines() string {
	return `data "maas_machines" "test" {}`
}
