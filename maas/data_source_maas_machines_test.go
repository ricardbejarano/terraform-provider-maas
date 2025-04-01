package maas_test

import (
	"fmt"
	"terraform-provider-maas/maas/testutils"
	"testing"

	"github.com/canonical/gomaasclient/entity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMaasMachines_basic(t *testing.T) {

	var machine entity.Machine
	domain := acctest.RandomWithPrefix("tf-domain-")
	hostname := acctest.RandomWithPrefix("tf-machine-")
	power_parameters := "{}"
	power_type := "manual"
	pxe_mac_address := testutils.RandomMAC()
	zone := "default"

	checks := []resource.TestCheckFunc{
		testAccMaasMachineCheckExists("maas_machine.test", &machine),
		resource.TestCheckResourceAttr("data.maas_machines.test", "machines.0.domain", domain),
		resource.TestCheckResourceAttr("data.maas_machines.test", "machines.0.hostname", hostname),
		resource.TestCheckResourceAttr("data.maas_machines.test", "machines.0.power_parameters", power_parameters),
		resource.TestCheckResourceAttr("data.maas_machines.test", "machines.0.power_type", power_type),
		resource.TestCheckResourceAttr("data.maas_machines.test", "machines.0.pxe_mac_address", pxe_mac_address),
		resource.TestCheckResourceAttr("data.maas_machines.test", "machines.0.zone", zone),
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.TestAccProviders,
		CheckDestroy: testAccCheckMaasMachineDestroy,
		ErrorCheck:   func(err error) error { return err },
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMaasMachines(domain, hostname, power_parameters, power_type, pxe_mac_address, zone),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

func testAccDataSourceMaasMachines(domain string, hostname string, power_parameters string, power_type string, pxe_mac_address string, zone string) string {
	return fmt.Sprintf(`
%s

data "maas_machines" "test" {
}
`, testAccMaasMachine(domain, hostname, power_parameters, power_type, pxe_mac_address, zone))
}
