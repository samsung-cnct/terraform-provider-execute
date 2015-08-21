package main

import (
	"fmt"
  "strconv"
	"io"
	"os/exec"
  "runtime"
  "log"

  "github.com/armon/circbuf"
	"github.com/hashicorp/terraform/helper/schema"
  "github.com/hashicorp/terraform/helper/hashcode"
  "github.com/mitchellh/go-linereader"
)

const (
  // maxBufSize limits how much output we collect from a local
  // invocation. This is to prevent TF memory usage from growing
  // to an enormous amount due to a faulty process.
  maxBufSize = 8 * 1024
)

// ExecCmd holds data necessary for a command to run
type ExecCmd struct {
	Cmd string
}

// Terraform schema for the 'exec' resource that is
// used in the provider configuration
func resource() *schema.Resource {
	return &schema.Resource{
		Create: Create,
		Read:   Read,
		Update: Update,
		Delete: Delete,

		Schema: map[string]*schema.Schema{
			"command": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
      "destroy_command": &schema.Schema{
        Type:     schema.TypeString,
        Optional: true,
      },
			"only_if": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
      "output": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
      },
		},
	}
}

func Create(d *schema.ResourceData, m interface{}) error {
	return PrepareCommand(d, m, true)
}

func Update(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("command") {
		// Set the id of the resource to destroy the resource
		d.SetId("")
	}
	return PrepareCommand(d, m, true)
}

func Read(d *schema.ResourceData, m interface{}) error {
	return nil
}

func Delete(d *schema.ResourceData, m interface{}) error {
  return PrepareCommand(d, m, false)
}

func PrepareCommand(d *schema.ResourceData, m interface{}, create bool) error {

  cmd := &ExecCmd{
		Cmd:     d.Get("command").(string),
	}

  if !create {
    cmd = &ExecCmd{
      Cmd:     d.Get("destroy_command").(string),
    }
  } else {
    onlyIf := &ExecCmd{
      Cmd:     d.Get("only_if").(string),
    }

    if onlyIf.Cmd != "" {
      onlyIfOut, err := ExecuteCommand(onlyIf)
      if err != nil {
        d.Set("output", onlyIfOut)
        return fmt.Errorf("Error running command '%s': %v. Output: %s", onlyIf.Cmd, err, onlyIfOut)
      }
    }
  }

	// run the actual command
	out, err := ExecuteCommand(cmd)
  d.Set("output", out)

	if err != nil {
		return fmt.Errorf("Error running command '%s': %v. Output: %s", cmd.Cmd, err, out)
	}

	// Set the id of the resource
	d.SetId(hash(cmd.Cmd))
	return nil
}

func ExecuteCommand(command *ExecCmd) (output string, err error) {

  // Execute the command using a shell
  var shell, flag string
  if runtime.GOOS == "windows" {
    shell = "cmd"
    flag = "/C"
  } else {
    shell = "/bin/sh"
    flag = "-c"
  }

  // Setup the reader that will read the lines from the command
  pr, pw := io.Pipe()
  copyDoneCh := make(chan struct{})
  go copyOutput(pr, copyDoneCh)

  // Setup the command
  cmd := exec.Command(shell, flag, command.Cmd)
  out, _ := circbuf.NewBuffer(maxBufSize)
  cmd.Stderr = io.MultiWriter(out, pw)
  cmd.Stdout = io.MultiWriter(out, pw)

  // Run the command to completion
  runErr := cmd.Run()
  pw.Close()
  <-copyDoneCh

  if runErr != nil {
    return string(out.Bytes()), fmt.Errorf("Error running command '%s': %v. Output: %s", command, runErr, out.Bytes())
  }

  return string(out.Bytes()), nil
}

func copyOutput(r io.Reader, doneCh chan<- struct{}) {
  defer close(doneCh)
  lr := linereader.New(r)
  for line := range lr.Ch {
    log.Printf("%s\n", line)
  }
}

// GenerateSHA1 generates a SHA1 hex digest for the given string
func hash(str string) string {
	return strconv.Itoa(hashcode.String(str))
}
