package config

type Config struct {
	JWTSecret          string
	JWTExpirationHours int
	JWTCookieName      string
	JWTCookieDomain    string
	JWTCookieSecure    bool
	JWTCookieHTTPOnly  bool
}

var DefaultConfig = Config{
	JWTSecret:          "MY_KEY", // Should be overridden via environment
	JWTExpirationHours: 24 * 7,            // 7 days
	JWTCookieName:      "auth_token",
	JWTCookieDomain:    "localhost",
	JWTCookieSecure:    true,
	JWTCookieHTTPOnly:  true,
}
