package config

type GoogleOauthConfig struct {
	ClientID         string `env:"GOOGLE_CLIENT_ID,notEmpty"`
	ClientSecret     string `env:"GOOGLE_CLIENT_SECRET,notEmpty"`
	RedirectEndpoint string `env:"GOOGLE_OAUTH_REDIRECT_ENDPOINT" envDefault:"/auth/google/callback"`
	ServerAddress    string `env:"SERVER_ADDRESS"`
}
