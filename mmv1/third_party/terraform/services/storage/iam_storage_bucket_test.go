// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestAccStorageBucket_basicIamBinding(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/storage.objectViewer",
		"admin_role":              "roles/storage.admin",
		"condition_title":         "expires_after_2019_12_31",
		"condition_expr":          `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
		"condition_desc":          "Expiring at midnight of 2019-12-31",
		"condition_title_no_desc": "expires_after_2019_12_31-no-description",
		"condition_expr_no_desc":  `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_basicIamBinding(context),
			},
			{
				ResourceName:      "google_storage_bucket_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("b/%s roles/storage.objectViewer", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccStorageBucket_updatedIamBinding(context),
			},
			{
				ResourceName:      "google_storage_bucket_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("b/%s roles/storage.objectViewer", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_basicIamMemeber(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/storage.objectViewer",
		"admin_role":              "roles/storage.admin",
		"condition_title":         "expires_after_2019_12_31",
		"condition_expr":          `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
		"condition_desc":          "Expiring at midnight of 2019-12-31",
		"condition_title_no_desc": "expires_after_2019_12_31-no-description",
		"condition_expr_no_desc":  `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccStorageBucket_basicIamMember(context),
			},
			{
				ResourceName:      "google_storage_bucket_iam_member.foo",
				ImportStateId:     fmt.Sprintf("b/%s roles/storage.objectViewer user:admin@hashicorptest.com", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_basicIamPolicy(t *testing.T) {
	t.Parallel()

	// This may skip test, so do it first
	sa := envvar.GetTestServiceAccountFromEnv(t)
	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/storage.objectViewer",
		"admin_role":              "roles/storage.admin",
		"condition_title":         "expires_after_2019_12_31",
		"condition_expr":          `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
		"condition_desc":          "Expiring at midnight of 2019-12-31",
		"condition_title_no_desc": "expires_after_2019_12_31-no-description",
		"condition_expr_no_desc":  `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
	}
	context["service_account"] = sa

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_basicIamPolicy(context),
				Check:  resource.TestCheckResourceAttrSet("data.google_storage_bucket_iam_policy.foo", "policy_data"),
			},
			{
				ResourceName:      "google_storage_bucket_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("b/%s", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccStorageBucket_emptyBindingIamPolicy(context),
			},
			{
				ResourceName:      "google_storage_bucket_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("b/%s", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_iamBindingWithCondition(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/storage.objectViewer",
		"admin_role":              "roles/storage.admin",
		"condition_title":         "expires_after_2019_12_31",
		"condition_expr":          `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
		"condition_desc":          "Expiring at midnight of 2019-12-31",
		"condition_title_no_desc": "expires_after_2019_12_31-no-description",
		"condition_expr_no_desc":  `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_withConditionIamBinding(context),
			},
			{
				ResourceName:      "google_storage_bucket_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("b/%s roles/storage.objectViewer %s", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"]), context["condition_title"]),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_iamBindingWithAndWithoutCondition(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/storage.objectViewer",
		"admin_role":              "roles/storage.admin",
		"condition_title":         "expires_after_2019_12_31",
		"condition_expr":          `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
		"condition_desc":          "Expiring at midnight of 2019-12-31",
		"condition_title_no_desc": "expires_after_2019_12_31-no-description",
		"condition_expr_no_desc":  `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_withAndWithoutConditionIamBinding(context),
			},
			{
				ResourceName:      "google_storage_bucket_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("b/%s roles/storage.objectViewer", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_storage_bucket_iam_binding.foo2",
				ImportStateId:     fmt.Sprintf("b/%s roles/storage.objectViewer %s", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"]), context["condition_title"]),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_storage_bucket_iam_binding.foo3",
				ImportStateId:     fmt.Sprintf("b/%s roles/storage.objectViewer %s", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"]), context["condition_title_no_desc"]),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_iamMemberWithCondition(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/storage.objectViewer",
		"admin_role":              "roles/storage.admin",
		"condition_title":         "expires_after_2019_12_31",
		"condition_expr":          `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
		"condition_desc":          "Expiring at midnight of 2019-12-31",
		"condition_title_no_desc": "expires_after_2019_12_31-no-description",
		"condition_expr_no_desc":  `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_withConditionIamMember(context),
			},
			{
				ResourceName:      "google_storage_bucket_iam_member.foo",
				ImportStateId:     fmt.Sprintf("b/%s roles/storage.objectViewer user:admin@hashicorptest.com %s", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"]), context["condition_title"]),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_iamMemberWithAndWithoutCondition(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/storage.objectViewer",
		"admin_role":              "roles/storage.admin",
		"condition_title":         "expires_after_2019_12_31",
		"condition_expr":          `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
		"condition_desc":          "Expiring at midnight of 2019-12-31",
		"condition_title_no_desc": "expires_after_2019_12_31-no-description",
		"condition_expr_no_desc":  `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_withAndWithoutConditionIamMember(context),
			},
			{
				ResourceName:      "google_storage_bucket_iam_member.foo",
				ImportStateId:     fmt.Sprintf("b/%s roles/storage.objectViewer user:admin@hashicorptest.com", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_storage_bucket_iam_member.foo2",
				ImportStateId:     fmt.Sprintf("b/%s roles/storage.objectViewer user:admin@hashicorptest.com %s", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"]), context["condition_title"]),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_storage_bucket_iam_member.foo3",
				ImportStateId:     fmt.Sprintf("b/%s roles/storage.objectViewer user:admin@hashicorptest.com %s", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"]), context["condition_title_no_desc"]),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucket_iamPolicyGeneratedWithCondition(t *testing.T) {
	t.Parallel()

	// This may skip test, so do it first
	sa := envvar.GetTestServiceAccountFromEnv(t)
	context := map[string]interface{}{
		"random_suffix":           acctest.RandString(t, 10),
		"role":                    "roles/storage.objectViewer",
		"admin_role":              "roles/storage.admin",
		"condition_title":         "expires_after_2019_12_31",
		"condition_expr":          `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
		"condition_desc":          "Expiring at midnight of 2019-12-31",
		"condition_title_no_desc": "expires_after_2019_12_31-no-description",
		"condition_expr_no_desc":  `request.time < timestamp(\"2020-01-01T00:00:00Z\")`,
	}
	context["service_account"] = sa

	// Test should have 3 bindings: one with a description and one without, and a third for an admin role. Any < chars are converted to a unicode character by the API.
	expectedPolicyData := acctest.Nprintf(`{"bindings":[{"members":["serviceAccount:%{service_account}"],"role":"%{admin_role}"},{"condition":{"description":"%{condition_desc}","expression":"%{condition_expr}","title":"%{condition_title}"},"members":["user:admin@hashicorptest.com"],"role":"%{role}"},{"condition":{"expression":"%{condition_expr}","title":"%{condition_title}-no-description"},"members":["user:admin@hashicorptest.com"],"role":"%{role}"}]}`, context)
	expectedPolicyData = strings.Replace(expectedPolicyData, "<", "\\u003c", -1)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucket_withConditionIamPolicy(context),
				Check: resource.ComposeAggregateTestCheckFunc(
					// TODO(SarahFrench) - uncomment once https://github.com/GoogleCloudPlatform/magic-modules/pull/6466 merged
					// resource.TestCheckResourceAttr("data.google_iam_policy.foo", "policy_data", expectedPolicyData),
					resource.TestCheckResourceAttr("google_storage_bucket_iam_policy.foo", "policy_data", expectedPolicyData),
					resource.TestCheckResourceAttrWith("data.google_iam_policy.foo", "policy_data", tpgresource.CheckGoogleIamPolicy),
				),
			},
			{
				ResourceName:      "google_storage_bucket_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("b/%s", fmt.Sprintf("tf-test-my-bucket%s", context["random_suffix"])),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStorageBucketIamPolicy_destroy(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucketIamPolicy_destroy(),
			},
		},
	})
}

func testAccStorageBucket_basicIamMember(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "default" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "foo" {
  bucket = google_storage_bucket.default.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}
`, context)
}

func testAccStorageBucket_basicIamPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "default" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "US"
  uniform_bucket_level_access = true
}

data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
  }
  binding {
    role = "%{admin_role}"
    members = ["serviceAccount:%{service_account}"]
  }
}

resource "google_storage_bucket_iam_policy" "foo" {
  bucket = google_storage_bucket.default.name
  policy_data = data.google_iam_policy.foo.policy_data
}

data "google_storage_bucket_iam_policy" "foo" {
  bucket = google_storage_bucket.default.name
  depends_on = [
    google_storage_bucket_iam_policy.foo
  ]
}
`, context)
}

func testAccStorageBucket_emptyBindingIamPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "default" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "US"
  uniform_bucket_level_access = true
}

data "google_iam_policy" "foo" {
}

resource "google_storage_bucket_iam_policy" "foo" {
  bucket = google_storage_bucket.default.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccStorageBucket_basicIamBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "default" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_binding" "foo" {
  bucket = google_storage_bucket.default.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}
`, context)
}

func testAccStorageBucket_updatedIamBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "default" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_binding" "foo" {
  bucket = google_storage_bucket.default.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
}
`, context)
}

func testAccStorageBucket_withConditionIamBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "default" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_binding" "foo" {
  bucket = google_storage_bucket.default.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
  condition {
    title       = "%{condition_title}"
    description = "%{condition_desc}"
    expression  = "%{condition_expr}"
  }
}
`, context)
}

func testAccStorageBucket_withAndWithoutConditionIamBinding(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "default" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_binding" "foo" {
  bucket = google_storage_bucket.default.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
}

resource "google_storage_bucket_iam_binding" "foo2" {
  bucket = google_storage_bucket.default.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
  condition {
    title       = "%{condition_title}"
    description = "%{condition_desc}"
    expression  = "%{condition_expr}"
  }
}

resource "google_storage_bucket_iam_binding" "foo3" {
  bucket = google_storage_bucket.default.name
  role = "%{role}"
  members = ["user:admin@hashicorptest.com"]
  condition {
    # Check that lack of description doesn't cause any issues
    # Relates to issue : https://github.com/hashicorp/terraform-provider-google/issues/8701
    title       = "%{condition_title_no_desc}"
    expression  = "%{condition_expr_no_desc}"
  }
}
`, context)
}

func testAccStorageBucket_withConditionIamMember(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "default" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "foo" {
  bucket = google_storage_bucket.default.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
  condition {
    title       = "%{condition_title}"
    description = "%{condition_desc}"
    expression  = "%{condition_expr}"
  }
}
`, context)
}

func testAccStorageBucket_withAndWithoutConditionIamMember(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "default" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_iam_member" "foo" {
  bucket = google_storage_bucket.default.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
}

resource "google_storage_bucket_iam_member" "foo2" {
  bucket = google_storage_bucket.default.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
  condition {
    title       = "%{condition_title}"
    description = "%{condition_desc}"
    expression  = "%{condition_expr}"
  }
}

resource "google_storage_bucket_iam_member" "foo3" {
  bucket = google_storage_bucket.default.name
  role = "%{role}"
  member = "user:admin@hashicorptest.com"
  condition {
    # Check that lack of description doesn't cause any issues
    # Relates to issue : https://github.com/hashicorp/terraform-provider-google/issues/8701
    title       = "%{condition_title_no_desc}"
    expression  = "%{condition_expr_no_desc}"
  }
}
`, context)
}

func testAccStorageBucket_withConditionIamPolicy(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "default" {
  name                        = "tf-test-my-bucket%{random_suffix}"
  location                    = "US"
  uniform_bucket_level_access = true
}

data "google_iam_policy" "foo" {
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
    condition {
      # Check that lack of description doesn't cause any issues
      # Relates to issue : https://github.com/hashicorp/terraform-provider-google/issues/8701
      title       = "%{condition_title_no_desc}"
      expression  = "%{condition_expr_no_desc}"
    }
  }
  binding {
    role = "%{role}"
    members = ["user:admin@hashicorptest.com"]
    condition {
      title       = "%{condition_title}"
      description = "%{condition_desc}"
      expression  = "%{condition_expr}"
    }
  }
  binding {
    role = "%{admin_role}"
    members = ["serviceAccount:%{service_account}"]
  }
}

resource "google_storage_bucket_iam_policy" "foo" {
  bucket = google_storage_bucket.default.name
  policy_data = data.google_iam_policy.foo.policy_data
}
`, context)
}

func testAccStorageBucketIamPolicy_destroy() string {
	return fmt.Sprintf(`
resource "google_service_account" "accessor" {
  account_id = "pub-sub-test-service-account"
}

resource "google_storage_bucket" "test_bucket" {
  name          = "sd-pubsub-test-bucket"
  location      = "US"
  storage_class = "STANDARD"

  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"

  force_destroy = true
}

data "google_iam_policy" "bucket_policy_data" {
  binding {
    role = "roles/storage.admin"

    members = ["serviceAccount:${google_service_account.accessor.email}"]
  }
}

resource "google_storage_bucket_iam_policy" "bucket_policy" {
  bucket      = google_storage_bucket.test_bucket.name
  policy_data = data.google_iam_policy.bucket_policy_data.policy_data
}

resource "google_pubsub_topic" "topic" {
  name = "sd-pubsub-test-bucket-topic"
}

resource "google_storage_notification" "storage_notification" {
  bucket         = google_storage_bucket.test_bucket.name
  payload_format = "JSON_API_V1"
  topic          = google_pubsub_topic.topic.id

  depends_on = [google_pubsub_topic_iam_policy.topic_policy]
}

data "google_storage_project_service_account" "gcs_account" {}

data "google_iam_policy" "topic_policy_data" {
  binding {
    role = "roles/pubsub.publisher"
    members = [
      "serviceAccount:${data.google_storage_project_service_account.gcs_account.email_address}"
    ]
  }
}

resource "google_pubsub_topic_iam_policy" "topic_policy" {
  topic       = google_pubsub_topic.topic.name
  policy_data = data.google_iam_policy.topic_policy_data.policy_data
}
`)
}
