// Create some resources with foreach and count
// tf-profile will agreggate these

module "test" {
    source = "./modules/test"
    count = 2
    input = "foo"
}

resource "time_sleep" "count" {
  count = 5
  create_duration = "${count.index}s"
}

resource "time_sleep" "foreach" {
    for_each = {"a": 1, "b": 2, "c": 3}
    create_duration = "${each.value}s"
}
