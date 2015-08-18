resource "execute_command" "commands" {
  command = "touch testfile"
  destroy_command = "rm testfile"
}