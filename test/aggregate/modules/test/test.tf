variable "input" {

}

resource "time_sleep" "count" {
  count = 3
  create_duration = "${count.index}s"
}
