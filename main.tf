resource "pagerduty_service" "example" {
  name              = "example"
  escalation_policy = pagerduty_escalation_policy.example.id
  type              = "generic_events"
  description       = "Example service"
}
