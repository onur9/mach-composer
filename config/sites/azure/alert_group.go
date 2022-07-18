package azure

type AlertGroup struct {
	Name        string   // TODO
	AlertEmails []string `yaml:"alert_emals"`
	WebhookURL  string   `yaml:"webhook_url"`
	LogicApp    string   `yaml:"logic_app"`
}
