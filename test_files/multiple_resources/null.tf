locals {
  foreach = {
    a = "1s"
    b = "2s"
    c = "1s"
    b = "2s"
  }
}

resource "time_sleep" "count" {
  count = 10
  create_duration = "${count.index}s"
}

resource "time_sleep" "for_each" {
  for_each = local.foreach
  create_duration = each.value

}