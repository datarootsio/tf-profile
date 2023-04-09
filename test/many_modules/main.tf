// Create a bunch of toplevel resources with a timeout to make them fail occasionally.
resource "time_sleep" "foo" {
  count = 150
  create_duration = "${random_integer.sleep[count.index].result}s"
}

resource "random_integer" "sleep" {
  count = 300
  min = 1
  max = 5
}

resource "time_sleep" "bar" {
  count = 150
  create_duration = "${random_integer.sleep[150+count.index].result}s"
}


module "core" {
  count = 3
  source = "./modules/core_infrastructure"
}

module "applications" {
  count = 10
  source = "./modules/applications"
}