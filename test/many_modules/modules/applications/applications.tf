module "airflow" {
    source = "../airflow"
    count = 10
}

module "dbt" {
    source = "../dbt"
    count = 5
}