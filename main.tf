variable "ld_access_token" {}

provider "launchdarkly" {
  access_token = "${var.ld_access_token}"
}

resource "launchdarkly_project" "my-project" {
  key = "my-project-key"
  name = "test"
}

resource "launchdarkly_environment" "dev" {
  project_key = "${launchdarkly_project.my-project.key}"
  name = "Development"
  key = "dev"
  color = "FF00FF"
}

resource "launchdarkly_environment" "hipaa" {
  project_key = "${launchdarkly_project.my-project.key}"
  name = "HIPAA"
  key = "hipaa"
  color = "FF00FF"
}

resource "launchdarkly_feature_flag" "my-flag" {
  project_key = "${launchdarkly_project.my-project.key}"
  key = "my-flag"
  name = "My Super Flag"
  description = "description!!"
  tags = ["foo", "bar", "spam"]
  custom_properties = [{
    key = "some.property"
    name = "Some Property"
    value = ["value1", "value2", "value3"]
  }]
}
