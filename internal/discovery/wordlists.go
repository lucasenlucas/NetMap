package discovery

// TopSubdomains contains high-impact subdomain prefixes for discovery.
var TopSubdomains = []string{
	"www", "api", "mail", "remote", "blog", "test", "dev", "stage", "stg", "vpn",
	"m", "secure", "admin", "api-docs", "docs", "cdn", "shop", "store", "app", "portal",
	"support", "help", "webmail", "mx", "ns1", "ns2", "beta", "old", "new", "assets",
	"images", "js", "css", "files", "download", "uploads", "v1", "v2", "web",
	"autodiscover", "sip", "mobile", "alpha", "staging", "demo", "prod", "intranet",
	"cloud", "git", "devops", "jenkins",
}

// CommonWebPaths contains high-impact directory and file paths for discovery.
var CommonWebPaths = []string{
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

// GetSubdomains returns the subset of wordlist based on mode (placeholder for now).
func GetSubdomains() []string {
	return TopSubdomains
}

// GetPaths returns the subset of wordlist based on mode (placeholder for now).
func GetPaths() []string {
	return CommonWebPaths
}
