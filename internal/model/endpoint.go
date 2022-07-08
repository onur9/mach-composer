package model

import (
	"log"
	"net/url"
	"strings"
)

type AzureEndpoint struct {
	SessionAffinityEnabled bool   `yaml:"session_affinity_enabled"`
	SessionAffinityTTL     int    `yaml:"session_affinity_ttl_seconds"`
	WAFPolicyID            string `yaml:"waf_policy_id"`
	InternalName           string `yaml:"internal_name"`
}

type AWSEndpoint struct {
	ThrottlingBurstLimit int  `yaml:"throttling_burst_limit"`
	ThrottlingRateLimit  int  `yaml:"throttling_rate_limit"`
	EnableCDN            bool `yaml:"enable_cdn"`
}

type Endpoint struct {
	URL   string         `yaml:"url"`
	Key   string         `yaml:"key"`
	Zone  string         `yaml:"zone"`
	AWS   *AWSEndpoint   `yaml:"aws"`
	Azure *AzureEndpoint `yaml:"azure"`

	Components []SiteComponent
}

func (e *Endpoint) SetDefaults() {
	e.URL = stripProtocol(e.URL)
	if e.Zone == "" && e.URL != "" {
		e.Zone = zoneFromURL(e.URL)
	}
}

func (e *Endpoint) IsRootDomain() bool {
	return e.URL == e.Zone
}

func (e Endpoint) Subdomain() string {
	if e.URL == "" {
		return ""
	}
	return subdomainFromURL(e.URL)
}

func zoneFromURL(value string) string {
	u, err := url.Parse(value)
	if err != nil {
		log.Fatal(err)
	}
	var domains []string
	if !strings.Contains(value, "://") {
		parts := strings.SplitN(value, "/", 2)
		domains = strings.Split(parts[0], ".")
	} else {
		domains = strings.Split(u.Hostname(), ".")
	}
	if len(domains) < 3 {
		return strings.Join(domains, ".")
	} else {
		return strings.Join(domains[1:], ".")
	}
}

func subdomainFromURL(value string) string {
	u, err := url.Parse(value)
	if err != nil {
		log.Fatal(err)
	}
	var domains []string
	if !strings.Contains(value, "://") {
		parts := strings.SplitN(value, "/", 2)
		domains = strings.Split(parts[0], ".")
	} else {
		domains = strings.Split(u.Hostname(), ".")
	}
	return domains[0]
}

func stripProtocol(value string) string {
	if strings.HasPrefix(value, "http://") {
		return strings.TrimPrefix(value, "http://")
	}
	if strings.HasPrefix(value, "https://") {
		return strings.TrimPrefix(value, "https://")
	}
	return value
}
