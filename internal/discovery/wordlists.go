package discovery

import (
	"bufio"
	"os"
	"strings"
)

// Pack definitions based on SecLists-style discovery categories.
var (
	// Subdomains basic (Top 50)
	subdomainsBasic = []string{
		"www", "api", "mail", "remote", "blog", "test", "dev", "stage", "stg", "vpn",
		"m", "secure", "admin", "api-docs", "docs", "cdn", "shop", "store", "app", "portal",
		"support", "help", "webmail", "mx", "ns1", "ns2", "beta", "old", "new", "assets",
		"images", "js", "css", "files", "download", "uploads", "v1", "v2", "web",
		"autodiscover", "sip", "mobile", "alpha", "staging", "demo", "prod", "intranet",
		"cloud", "git", "devops", "jenkins",
	}

	// Subdomains deep (Top 100+)
	subdomainsDeep = append(subdomainsBasic, []string{
		"internal", "lab", "dev-api", "beta-api", "qa", "testing", "secure-api",
		"gitlab", "bitbucket", "jira", "confluence", "wiki", "status", "monitor",
		"nagios", "zabbix", "grafana", "prometheus", "elastic", "kibana", "logstash",
		"splunk", "graylog", "sentry", "newrelic", "datadog", "aws", "azure", "gcp",
		"kubernetes", "k8s", "docker", "container", "registry", "harbor", "rancher",
	}...)

	// Paths basic (Core 10)
	pathsBasic = []string{
		"/", "/login", "/admin", "/dashboard", "/api", "/api/v1", "/graphql", "/auth",
		"/robots.txt", "/.env",
	}

	// Paths common (Top 100)
	pathsCommon = []string{
		"/", "/index.html", "/index.php", "/home", "/about", "/contact",
		"/login", "/signin", "/logout", "/signup", "/register", "/auth",
		"/api", "/v1", "/v2", "/admin", "/administrator", "/dashboard", "/panel",
		"/controlpanel", "/cp", "/wp-admin", "/wp-login.php", "/config", "/config.json",
		"/config.yml", "/settings", "/setup", "/install", "/user", "/users", "/profile",
		"/account", "/search", "/blog", "/news", "/posts", "/archive", "/media",
		"/assets", "/static", "/images", "/img", "/css", "/js", "/lib", "/vendor",
		"/node_modules", "/package.json", "/package-lock.json", "/composer.json",
		"/composer.lock", "/docker-compose.yml", "/Dockerfile", "/Makefile", "/README.md",
		"/LICENSE", "/CONTRIBUTING.md", "/CHANGELOG.md", "/docs", "/documentation",
		"/api-docs", "/swagger", "/v1/api", "/v1/users", "/v1/auth", "/graphql",
		"/metrics", "/health", "/status", "/info", "/ping", "/test", "/temp", "/tmp",
		"/backup", "/backups", "/old", "/dev", "/staging", "/test", "/public",
		"/private", "/secure", "/secret", "/auth/login", "/auth/register",
		"/auth/callback", "/oauth", "/callback", "/api/v1/users", "/api/v1/auth",
		"/.git", "/.gitignore", "/.env", "/.env.example", "/.env.local",
		"/.gitlab-ci.yml", "/.travis.yml", "/favicon.ico", "/robots.txt", "/sitemap.xml",
		"/ads.txt",
	}

	// Paths API
	pathsAPI = []string{
		"/api", "/v1", "/v2", "/v3", "/graphql", "/swagger", "/swagger-ui.html",
		"/swagger/index.html", "/api/docs", "/docs", "/api-docs", "/v1/api",
		"/api/v1", "/api/v2", "/rest", "/rest-api", "/api/v1/users", "/api/v1/auth",
		"/api/v1/status", "/health", "/metrics", "/auth/token", "/oauth/token",
		"/api/login", "/api/signup", "/v1/oauth2/token", "/v1/authorize", "/v1/userinfo",
	}

	// Paths Admin
	pathsAdmin = []string{
		"/admin", "/administrator", "/dashboard", "/panel", "/controlpanel", "/cp",
		"/wp-admin", "/wp-login.php", "/manage", "/management", "/backend",
		"/server-status", "/server-info", "/phpmyadmin", "/pma", "/admin/login",
		"/admin/dashboard", "/setup", "/install", "/config", "/settings", "/cpanel", "/whm",
	}
)

// GetPack returns the subdomain and path sets for a given pack name.
func GetPack(name string) ([]string, []string) {
	switch strings.ToLower(name) {
	case "dns-extended":
		return subdomainsDeep, pathsBasic
	case "web-deep":
		return subdomainsBasic, pathsCommon
	case "api-focused":
		return subdomainsBasic, pathsAPI
	case "admin-stealth":
		return subdomainsBasic, pathsAdmin
	case "standard":
		fallthrough
	default:
		return subdomainsBasic, pathsBasic
	}
}

// LoadWordlistFromFile reads a local text file and returns a slice of lines.
func LoadWordlistFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}
