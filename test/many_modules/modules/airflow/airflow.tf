resource "random_integer" "sleep" {
  min = 1
  max = 5
}

resource "time_sleep" "bar" {
  create_duration = "${random_integer.sleep.result}s"
}