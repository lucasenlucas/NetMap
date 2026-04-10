package discovery

import (
	"bufio"
	"os"
	"strings"
)

// Pack definitions based on SecLists-style discovery categories.
var (
	subdomainsBasic = []string{
		"www", "api", "mail", "remote", "blog", "test", "dev", "stage", "stg", "vpn",
		"m", "secure", "admin", "api-docs", "docs", "cdn", "shop", "store", "app", "portal",
		"support", "help", "webmail", "mx", "ns1", "ns2", "beta", "old", "new", "assets",
		"images", "js", "css", "files", "download", "uploads", "v1", "v2", "web",
		"autodiscover", "sip", "mobile", "alpha", "staging", "demo", "prod", "intranet",
		"cloud", "git", "devops", "jenkins",
	}

	subdomainsComprehensive = append(subdomainsBasic, []string{
		"internal", "lab", "dev-api", "beta-api", "qa", "testing", "secure-api",
		"gitlab", "bitbucket", "jira", "confluence", "wiki", "status", "monitor",
		"nagios", "zabbix", "grafana", "prometheus", "elastic", "kibana", "logstash",
		"splunk", "graylog", "sentry", "newrelic", "datadog", "aws", "azure", "gcp",
		"kubernetes", "k8s", "docker", "container", "registry", "harbor", "rancher",
		"db", "database", "sql", "rds", "redis", "cache", "storage", "media", "bucket",
		"s3", "pub", "private", "hidden", "secret", "vault", "auth", "sso", "identity",
		"login", "signin", "signup", "account", "profile", "billing", "payment",
	}...)

	subdomainsUltra = append(subdomainsComprehensive, []string{
		"dev-ws", "staging-ws", "prod-ws", "ws", "socket", "stream", "live", "chat",
		"beta-app", "alpha-app", "demo-app", "qa-app", "test-app", "mobile-api",
		"android", "ios", "blackberry", "desktop", "client", "customer", "partner",
		"reseller", "affiliate", "tracking", "analytics", "pixel", "telemetry",
		"metrics", "log", "logs", "error", "errors", "trace", "debug", "admin-api",
		"adm", "mgt", "management", "control", "cp", "panel", "direct", "origin",
		"edge", "lb", "proxy", "waf", "cdn-cgi", "cache-api", "static", "assets",
		"images-api", "media-server", "upload-api", "file-api", "download-api",
		"api01", "api02", "api03", "api04", "node01", "node02", "node03", "cluster1",
		"cluster2", "aws-eu", "aws-us", "gcp-eu", "gcp-us", "azure-eu", "azure-us",
		"backup-api", "db-api", "sql-api", "search-api", "index-api", "internal-api",
		"private-api", "secure-gateway", "identity-api", "auth-api", "sso-api",
		"mail-api", "smtp", "pop3", "imap", "exchange", "webmail-api", "sip-api",
		"voip", "rtc", "conference", "meeting", "events-api", "calendar", "tasks",
		"drive", "box", "sync", "share", "docs-api", "wiki-api", "notes", "tasks-api",
		"shop-api", "cart", "checkout", "billing-api", "payment-api", "wallet",
		"crypto", "coins", "exchange-api", "trade", "market", "data", "stats",
		"research", "training", "university", "academy", "learning", "portal-api",
		"help-api", "support-desk", "tickets", "crm-api", "erp-api", "hr-api",
		"fin", "legal", "compliance", "audit", "security-api", "soc-api",
	}...)

	pathsBasic = []string{
		"/", "/login", "/admin", "/dashboard", "/api", "/api/v1", "/graphql", "/auth",
		"/robots.txt", "/.env",
	}

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

	pathsComprehensive = append(pathsCommon, []string{
		"/.git/config", "/.git/head", "/.ssh/id_rsa", "/.aws/credentials",
		"/wp-config.php", "/config.php", "/settings.py", "/appsettings.json",
		"/web.config", "/.htaccess", "/phpinfo.php", "/info.php", "/status.php",
		"/.well-known/security.txt", "/swagger/index.html", "/prometheus", "/grafana",
		"/docker-compose.yaml", "/kubernetes.yaml", "/terraform.tfstate",
	}...)

	pathsUltra = append(pathsComprehensive, []string{
		"/api/v1/user", "/api/v1/admin", "/api/v1/config", "/api/v1/settings",
		"/api/v2/user", "/api/v2/admin", "/api/v2/config", "/api/v2/settings",
		"/api/v3/user", "/api/v3/admin", "/api/v3/config", "/api/v3/settings",
		"/v1.0", "/v2.0", "/v3.0", "/api/internal", "/api/private", "/api/secure",
		"/auth/v1/login", "/auth/v2/login", "/oauth2/token", "/oauth2/authorize",
		"/saml/metadata", "/.well-known/openid-configuration", "/.well-known/jwks.json",
		"/api/graphql", "/api/rest", "/api/json", "/api/xml", "/api/v1/status/health",
		"/management/health", "/management/info", "/management/metrics", "/actuator",
		"/actuator/health", "/actuator/info", "/actuator/prometheus", "/api/v1/debug",
		"/api/v1/trace", "/api/v1/log", "/api/v1/dump", "/admin/settings", "/admin/users",
		"/admin/api", "/admin/logs", "/admin/system", "/admin/config", "/admin/db",
		"/pma", "/phpmyadmin", "/myadmin", "/mysql", "/sql", "/dbadmin", "/web-console",
		"/console", "/shell", "/terminal", "/exec", "/run", "/cmd", "/ping",
		"/temp/config", "/tmp/config", "/backup/config", "/old/config", "/v1/config",
		"/v2/config", "/v3/config", "/api-docs/v1", "/api-docs/v2", "/api-docs/v3",
		"/swagger-resources", "/swagger-ui", "/v1/swagger.json", "/v2/swagger.json",
		"/.env.production", "/.env.staging", "/.env.local", "/.env.test", "/.env.dev",
		"/config/database.yaml", "/config/settings.yaml", "/config/auth.yaml",
		"/docker-compose.override.yml", "/docker-stack.yml", "/k8s/deployment.yaml",
		"/k8s/service.yaml", "/k8s/ingress.yaml", "/k8s/secret.yaml", "/k8s/configmap.yaml",
		"/.vscode/settings.json", "/.idea/workspace.xml", "/.circleci/config.yml",
		"/.github/workflows/main.yml", "/.gitlab-ci.yml", "/.travis.yml",
		"/Jenkinsfile", "/capfile", "/Procfile", "/Gemfile", "/Gemfile.lock",
		"/requirements.txt", "/pipfile", "/pipfile.lock", "/go.mod", "/go.sum",
		"/cargo.toml", "/cargo.lock", "/composer.json", "/composer.lock",
		"/package-lock.json", "/yarn.lock", "/pnpm-lock.yaml",
		"/api/v1/upload", "/api/v1/download", "/api/v1/file", "/api/v1/media",
		"/api/v1/image", "/api/v1/video", "/api/v1/stream", "/api/v1/ws",
		"/api/v1/socket", "/api/v1/rpc", "/api/v1/subscription", "/api/v1/billing",
		"/api/v1/payment", "/api/v1/order", "/api/v1/invoice", "/api/v1/report",
		"/api/v1/search", "/api/v1/query", "/api/v1/filter", "/api/v1/sort",
		"/api/v1/page", "/api/v1/limit", "/api/v1/offset", "/api/v1/detail",
		"/api/v1/list", "/api/v1/create", "/api/v1/update", "/api/v1/delete",
		"/api/v1/patch", "/api/v1/options", "/api/v1/head",
	}...)

	pathsAPI = []string{
		"/api", "/v1", "/v2", "/v3", "/graphql", "/swagger", "/swagger-ui.html",
		"/swagger/index.html", "/api/docs", "/docs", "/api-docs", "/v1/api",
		"/api/v1", "/api/v2", "/rest", "/rest-api", "/api/v1/users", "/api/v1/auth",
		"/api/v1/status", "/health", "/metrics", "/auth/token", "/oauth/token",
		"/api/login", "/api/signup", "/v1/oauth2/token", "/v1/authorize", "/v1/userinfo",
	}

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
	case "ultra":
		return subdomainsUltra, pathsUltra
	case "full":
		return subdomainsComprehensive, pathsComprehensive
	case "dns-extended":
		return subdomainsComprehensive, pathsBasic
	case "web-deep":
		return subdomainsBasic, pathsComprehensive
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
