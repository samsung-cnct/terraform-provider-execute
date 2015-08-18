# terraform-provider-execute

Terraform plugin mostly based on https://github.com/gosuri/terraform-exec-provider (thanks!!!). Provides an ability to execute arbitrary commands on Terraform create and destroy.

## Usage

    resource "exec" "command" {
      command "/path/to/command"
      destroy_command "/path/to/command"
    }

### Attribute reference

* `command` - (Required) Command to execute on terraform Create
* `destroy_command` - (Optional) Command to execute on terraform destroy
* `only_if` - (Optional) Guard attribute, to create the resource (Execute) the command only if this guard is satisfied. If the command returns 0, the guard is applied. If the command returns any other value, then the guard attribute is not applied.
* `timeout` - (Optional) Create/Destroy max timeout


### Examples

The below example will create a 'testfile' file when you run 'terraform apply' and delete the 'testfile' file when you run 'terraform destroy'

    resource "execute_command" "commands" {
      command = "touch testfile"
      destroy_command = "rm testfile"
    }

## Installation

    $ git clone https://github.com/gosuri/terraform-exec-provider.git
    $ cd terraform-exec-provider
    $ go get
