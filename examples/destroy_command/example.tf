resource "execute_command" "commands1" {
  command = "touch testfile"
  destroy_command = "rm testfile"
}

resource "execute_command" "commands2" {
  command = "echo 'hello world'"
  only_if = "true"
  destroy_command = "echo 'goodbye world'"
}