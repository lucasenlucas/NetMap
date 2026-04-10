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
		"invoice", "order", "tracking", "status", "news", "press", "events", "jobs",
		"career", "hiring", "team", "staff", "client", "customer", "partner", "dist",
		"reseller", "api1", "api2", "api3", "v3", "mobile-api", "iot", "sensor",
		"camera", "print", "printer", "scan", "scanner", "wifi", "guest", "office",
		"branch", "hq", "corp", "local", "home", "backoffice", "erp", "crm", "hr",
		"payroll", "finance", "legal", "compliance", "audit", "security", "soc",
		"noc", "proxy", "lb", "gateway", "edge", "cdn1", "cdn2", "static1", "static2",
		"upload", "upload-api", "media-api", "stream", "video", "audio", "voice",
		"chat", "support-api", "helpdesk", "ticketing", "service", "services",
		"catalog", "inventory", "search", "search-api", "elastic-api", "solr",
		"api-v1", "api-v2", "api-v3", "endpoint", "endpoints", "router", "switch",
		"firewall", "vpn-gate", "vps", "server", "host", "node", "cluster", "master",
		"worker", "slave", "back", "front", "app1", "app2", "app3", "web1", "web2",
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
		"/.git/config", "/.git/head", "/.ssh/id_rsa", "/.ssh/id_dsa", "/.ssh/authorized_keys",
		"/.aws/credentials", "/.aws/config", "/.docker/config.json", "/.npmrc", "/.yarnrc",
		"/.bash_history", "/.zsh_history", "/.mysql_history", "/.pgsql_history",
		"/phpinfo.php", "/info.php", "/status.php", "/health_check", "/debug",
		"/.well-known/security.txt", "/.well-known/assetlinks.json", "/.well-known/apple-app-site-association",
		"/wp-config.php", "/wp-config.php.bak", "/wp-config.php.swp", "/configuration.php",
		"/config.php", "/settings.py", "/local_settings.py", "/appsettings.json",
		"/web.config", "/htaccess", "/.htaccess", "/nginx.conf", "/proxy.conf",
		"/php.ini", "/etc/passwd", "/etc/shadow", "/etc/group", "/etc/hosts",
		"/boot.ini", "/servlets/config.xml", "/WEB-INF/web.xml", "/WEB-INF/config.xml",
		"/META-INF/maven/pom.xml", "/pom.xml", "/build.gradle", "/package.json",
		"/npm-debug.log", "/yarn-error.log", "/error_log", "/access_log", "/logs/access.log",
		"/logs/error.log", "/logs/api.log", "/storage/logs/laravel.log", "/tmp/session",
		"/sessions", "/cache", "/temp", "/tmp", "/uploads", "/attachments", "/downloads",
		"/export", "/import", "/scripts", "/bin", "/src", "/lib", "/include", "/tests",
		"/spec", "/features", "/vendor", "/node_modules", "/bower_components",
		"/dist", "/build", "/out", "/public/dist", "/public/build", "/client/dist",
		"/server/dist", "/admin/assets", "/admin/css", "/admin/js", "/admin/img",
		"/api/v1/health", "/api/v1/metrics", "/api/v1/status", "/api/v1/debug",
		"/api/v2/health", "/api/v3/health", "/api/health", "/actuator/health",
		"/actuator/metrics", "/actuator/env", "/actuator/beans", "/actuator/info",
		"/prometheus", "/graph", "/grafana", "/zabbix", "/nagios", "/kibana",
		"/elastic", "/solr", "/jenkins", "/gitlab", "/travis", "/circleci",
		"/docker-compose.yaml", "/kubernetes.yaml", "/k8s.yaml", "/manifest.yaml",
		"/chart.yaml", "/values.yaml", "/terraform.tfstate", "/terraform.tfvars",
		"/secrets.yaml", "/vault.yaml", "/credentials.json", "/key.pem", "/cert.pem",
		"/id_rsa", "/id_rsa.pub", "/authorized_keys", "/known_hosts", "/config.bak",
		"/database.sql", "/db.sql", "/dump.sql", "/backup.sql", "/data.sql",
		"/backup.zip", "/backup.tar.gz", "/backup.tgz", "/site.zip", "/www.zip",
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
