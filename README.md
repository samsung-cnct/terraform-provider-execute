# terraform-provider-exec

Provides an ability to execute arbitrary commands


## Usage

    resource "exec" "command" {
      command "/path/to/command"
    }

### Attribute reference
    
* `command` - (Required) Command to execute
* `only_if` - (Optional) Guard attribute, to create the resource (Execute) the command only if this guard is satisfied. If the command returns 0, the guard is applied. If the command returns any other value, then the guard attribute is not applied.
 

### Examples

The below example will run the command after creating VPCs, where `commands/peer-vpc` is shell script to add peering connection between VPCs and will be executed after the VPCs are created when `depends_on` attribute is specified.

    resource "aws_vpc" "primary" {
      cidr_block = "10.0.0.0/16"
    }

    resource "aws_vpc" "app" {
      cidr_block = "10.1.0.0/16"
    }

    resource "exec", "peer_vpcs" {
      command = "commands/peer-vpc ${aws_vpc.primary.id} ${aws_vpc.app.id}"
      depends_on = ["aws_vpc.primary.id","aws_vpc.app"]
    }

## Installation

    $ git clone https://github.com/gosuri/terraform-exec-provider.git
    $ cd terraform-exec-provider
    $ make install
