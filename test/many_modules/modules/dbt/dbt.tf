resource "random_integer" "sleep" {
  count = 2
  min = 1
  max = 5
}

resource "time_sleep" "bar" {
  count = 2
  create_duration = "${random_integer.sleep[count.index].result}s"
}