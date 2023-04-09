resource "random_integer" "sleep" {
  min = 0
  max = 5
}

resource "time_sleep" "bar" {
  create_duration = "${random_integer.sleep.result}s"

  provisioner "local-exec" {
    command = "if [ ${random_integer.sleep.result} -gt 4 ]; then return 1; fi"
  }
}