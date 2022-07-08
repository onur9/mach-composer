package updater

import (
	"fmt"
	"strings"

	"github.com/labd/mach-composer/internal/model"
)

type UpdateSet struct {
	filename string
	updates  []ChangeSet
}

type ChangeSet struct {
	LastVersion string
	Changes     []gitCommit
	Component   *model.Component
	Forced      bool
}

func (cs *ChangeSet) HasChanges() bool {
	return cs.Component.Version != cs.LastVersion
}

func OutputChanges(cs *ChangeSet) string {
	var b strings.Builder

	if cs.Forced && len(cs.Changes) == 0 {
		fmt.Fprintf(&b, "Update %s to %s\n", cs.Component.Name, cs.LastVersion)
		return b.String()
	}

	fmt.Fprintf(&b, "Updates for %s (%s...%s)\n", cs.Component.Name, cs.Component.Version, cs.LastVersion)

	if !cs.HasChanges() {
		fmt.Fprintln(&b, "  No updates...")
		fmt.Fprintln(&b, "")
		return b.String()
	}

	for _, commit := range cs.Changes {
		fmt.Fprintf(&b, "  %s: %s <%s>\n", commit.Commit, commit.Message, commit.Author)
	}
	fmt.Fprintln(&b, "")

	return b.String()
}

func (u *UpdateSet) ChangeLog() string {
	var b strings.Builder

	fmt.Fprintf(&b, "Updated %d components\n\n", len(u.updates))
	for _, cs := range u.updates {
		content := OutputChanges(&cs)
		b.WriteString(content)
	}

	return b.String()
}

func (u *UpdateSet) ComponentChangeLog(component string) string {
	var b strings.Builder

	for _, cs := range u.updates {
		if strings.EqualFold(cs.Component.Name, component) {
			content := OutputChanges(&cs)
			b.WriteString(content)
		}
	}
	return b.String()
}

func (u *UpdateSet) HasChanges() bool {
	for _, cs := range u.updates {
		if cs.HasChanges() || cs.Forced {
			return true
		}
	}
	return false
}
