package sites

import (
	"github.com/labd/mach-composer/internal/config/global"
	"github.com/labd/mach-composer/internal/config/sites/azure"
)

type Azure struct {
	global.Azure
	AlertGroup    azure.AlertGroup `yaml:"alert_group"`
	ResourceGroup string
}
