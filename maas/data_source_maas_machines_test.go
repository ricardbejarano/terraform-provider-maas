package maas_test

import (
	"fmt"
	"terraform-provider-maas/maas/testutils"
	"testing"

	"github.com/canonical/gomaasclient/entity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMAASMachines_basic(t *testing.T) {

	var machine entity.Machine
	domain := acctest.RandomWithPrefix("tf-domain-")
	hostname := acctest.RandomWithPrefix("tf-machine-")
	power_parameters := "{}"
	power_type := "manual"
	pxe_mac_address := testutils.RandomMAC()
	zone := "default"

	checks := []resource.TestCheckFunc{
		testAccMAASMachineCheckExists("maas_machine.test", &machine),
		resource.TestCheckResourceAttr("data.maas_machines.test", "machines.0.hostname", hostname),
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.TestAccProviders,
		CheckDestroy: testAccCheckMAASMachineDestroy,
		ErrorCheck:   func(err error) error { return err },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMAASMachines(domain, hostname, power_parameters, power_type, pxe_mac_address, zone),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

func testAccDataSourceMAASMachines(domain string, hostname string, power_parameters string, power_type string, pxe_mac_address string, zone string) string {
	return fmt.Sprintf(`
%s

data "maas_machines" "test" {
}
`, testAccMAASMachine(domain, hostname, power_parameters, power_type, pxe_mac_address, zone))
}
