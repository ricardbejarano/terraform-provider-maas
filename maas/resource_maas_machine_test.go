package maas_test

import (
	"fmt"
	"strings"
	"terraform-provider-maas/maas"
	"terraform-provider-maas/maas/testutils"

	"github.com/canonical/gomaasclient/entity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccMAASMachine(domain string, hostname string, power_parameters string, power_type string, pxe_mac_address string, zone string) string {
	return fmt.Sprintf(`
resource "maas_dns_domain" "test" {
	name          = "%s"
	ttl           = 3600
	authoritative = true
}

resource "maas_machine" "test" {
	domain           = maas_dns_domain.test.name
	hostname         = "%s"
	power_parameters = "%s"
	power_type       = "%s"
	pxe_mac_address  = "%s"
	zone             = "%s"
}
`, domain, hostname, power_parameters, power_type, pxe_mac_address, zone)
}

func testAccMAASMachineCheckExists(rn string, machine *entity.Machine) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s\n %#v", rn, s.RootModule().Resources)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		conn := testutils.TestAccProvider.Meta().(*maas.ClientConfig).Client
		gotMachine, err := conn.Machine.Get(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error getting machine: %s", err)
		}

		*machine = *gotMachine

		return nil
	}
}

func testAccCheckMAASMachineDestroy(s *terraform.State) error {
	// retrieve the connection established in Provider configuration
	conn := testutils.TestAccProvider.Meta().(*maas.ClientConfig).Client

	// loop through the resources in state, verifying each maas_machine
	// is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "maas_machine" {
			continue
		}

		// Retrieve our maas_machine by referencing it's state ID for API lookup
		response, err := conn.Machine.Get(rs.Primary.ID)
		if err == nil {
			if response != nil && response.SystemID == rs.Primary.ID {
				return fmt.Errorf("MAAS Machine (%s) still exists.", rs.Primary.ID)
			}

			return nil
		}

		// If the error is equivalent to 404 not found, the maas_machine is destroyed.
		// Otherwise return the error
		if !strings.Contains(err.Error(), "404 Not Found") {
			return err
		}
	}

	return nil
}
