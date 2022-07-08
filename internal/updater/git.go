package updater

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/labd/mach-composer/internal/model"
	"github.com/labd/mach-composer/internal/utils"
	"github.com/sirupsen/logrus"
)

// commit": "%H",
// author": "%aN <%aE>",
// date": "%ad",
// message": "%s",
const gitFormat = "%H|%aN <%aE>|%ad|%s"

type gitSource struct {
	URL        string
	Repository string
	Path       string
	Name       string
}

type gitCommit struct {
	Commit  string
	Author  string
	Date    string
	Message string
}

func GetLastVersionGit(ctx context.Context, c *model.Component, origin string) (*ChangeSet, error) {
	cacheDir := getGitCachePath(origin)
	source, err := parseGitSource(c.Source)

	if err != nil {
		return nil, fmt.Errorf("cannot check %s component since it doesn't have a Git source defined", c.Name)
	}

	branch := "HEAD"
	if c.Branch != "" {
		branch = c.Branch
	}
	fetchGitRepository(ctx, source, cacheDir)
	commits := loadGitHistory(ctx, source, c.Version, branch, cacheDir)

	cs := &ChangeSet{
		Changes:   commits,
		Component: c,
	}

	if len(commits) < 1 {
		cs.LastVersion = c.Version
	} else {
		cs.LastVersion = commits[0].Commit
	}

	return cs, nil
}

func getGitCachePath(origin string) string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	base := strings.TrimSuffix(origin, filepath.Ext(origin))
	path := filepath.Join(cwd, ".mach", base)
	if err := os.MkdirAll(path, 0700); err != nil {
		panic(err)
	}
	return path
}

// Parse a git url and return a gitSource reference
func parseGitSource(source string) (*gitSource, error) {
	re := regexp.MustCompile("^git::(?P<repo>https://.*?)(?://(?P<path>.*))?$")
	match := re.FindStringSubmatch(source)

	if match == nil {
		return nil, errors.New("invalid Git source defined")
	}

	result := &gitSource{
		URL: source,
	}
	for i, name := range re.SubexpNames() {
		if name == "repo" {
			result.Repository = match[i]
		}
		if name == "path" {
			result.Path = match[i]
		}
	}

	parts, err := url.Parse(result.Repository)
	if err != nil {
		panic(err)
	}
	result.Name = filepath.Base(parts.Path)
	return result, nil
}

// fetchGitRepository clones or updates the repository. We only need the history
// so clone using --bare
func fetchGitRepository(ctx context.Context, source *gitSource, cacheDir string) {
	dest := filepath.Join(cacheDir, source.Name)

	_, err := os.Stat(dest)
	if os.IsNotExist(err) {
		output := runGit(ctx, ".", "clone", "--bare", source.Repository, dest)
		logrus.Debug(string(output))
	} else {
		output := runGit(ctx, dest, "fetch", "-f", "origin", "*:*")
		logrus.Debug(string(output))
	}
}

func loadGitHistory(ctx context.Context, source *gitSource, baseRef string, branch string, cacheDir string) []gitCommit {
	dest := filepath.Join(cacheDir, source.Name)

	args := []string{
		"log", branch, fmt.Sprintf(`--pretty=%s`, gitFormat),
	}
	if baseRef != "" {
		args = append(args, fmt.Sprintf("%s...%s", baseRef, branch))
	}

	output := runGit(ctx, dest, args...)
	commits := []gitCommit{}

	for _, line := range SplitLines(string(output)) {
		parts := strings.SplitN(line, "|", 4)
		commits = append(commits, gitCommit{
			Commit:  parts[0][:7],
			Author:  parts[1],
			Date:    parts[2],
			Message: parts[3],
		})
	}
	return commits
}

func Commit(ctx context.Context, fileNames []string, message string) {
	args := []string{"commit"}
	args = append(args, fileNames...)
	args = append(args, "-m", message)

	runGit(ctx, ".", args...)
}

// runGit executes the git command
func runGit(ctx context.Context, cwd string, args ...string) []byte {
	logrus.Debugf("Running: git %s\n", strings.Join(args, " "))
	cmd := exec.CommandContext(
		ctx,
		"git",
		args...,
	)
	cmd.Dir = cwd
	utils.CmdSetForeground(cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, string(output))
		os.Exit(1)
	}

	return output
}
