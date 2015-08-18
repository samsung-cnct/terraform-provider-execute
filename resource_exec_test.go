package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"execute": testAccProvider,
	}
}

func TestResourceExecCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccResourceExecDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourceExecConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("execute_command.foo", "output", "success\n"),
				),
			},
		},
	})
}

func TestResourceExecUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccResourceExecDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourceExecConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("execute_command.foo", "output", "success\n"),
				),
			},
			resource.TestStep{
				Config: testAccResourceExecConfig_basic_2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("execute_command.foo", "output", "success2\n"),
				),
			},
		},
	})
}

func TestResourceExecCreateTestFail(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccResourceExecDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourceExecConfig_test_fail,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExecResourceIsNil("execute_command.failing"),
				),
			},
			resource.TestStep{
				Config: testAccResourceExecConfig_test_pass,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("execute_command.foo", "output", "success\n"),
				),
			},
		},
	})
}

func testAccCheckExecResourceIsNil(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[r]
		if ok {
			return fmt.Errorf("Resource exists: %s", r)
		}
		return nil
	}
}

func testAccResourceExecDestroy(s *terraform.State) error {
	return nil
}

const testAccResourceExecConfig_basic = `
resource "execute_command" "foo" {
	command = "echo 'success'"
}
`
const testAccResourceExecConfig_basic_2 = `
resource "execute_command" "foo" {
	command = "echo 'success2'"
}
`
const testAccResourceExecConfig_test_pass = `
resource "execute_command" "foo" {
	command = "echo 'success'"
	only_if = "true"
}
`
const testAccResourceExecConfig_test_fail = `
resource "execute_command" "failing" {
	command = "echo 'success'"
	only_if = "false"
}
`
const testAccResourceExecConfig_timeout = `
resource "execute_command" "foo" {
	command = "sleep 2 && echo 'success'"
	timeout = 1
}
`
const testAccResourceExecConfig_fail = `
resource "execute_command" "foo" {
	command = "echo 'failure' >&2 && exit 1"
}
`
