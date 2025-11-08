package app

type Application struct {
	Config  Config
	Secrets Secrets
}

func LoadApplication(configPath string) (*Application, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	secrets, err := LoadSecrets()
	if err != nil {
		return nil, err
	}

	app := Application{
		Config:  *config,
		Secrets: *secrets,
	}
	return &app, nil
}
