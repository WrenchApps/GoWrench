package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"wrench/app"
	"wrench/app/manifest/application_settings"
	"wrench/app/startup"
	"wrench/app/startup/connections"
	"wrench/app/startup/token_credentials"
)

func main() {
	ctx := context.Background()
	app.SetContext(ctx)

	startup.LoadEnvsFiles()

	byteArray, err := startup.LoadYamlFile(getFileConfigPath())
	startup.LoadAwsSecrets(byteArray)
	if err != nil {
		app.LogError2(fmt.Sprintf("Error loading YAML: %v", err), err)
	}

	byteArray = startup.EnvInterpolation(byteArray)
	applicationSetting, err := application_settings.ParseToApplicationSetting(byteArray)

	if err != nil {
		app.LogError2(fmt.Sprintf("Error parse yaml: %v", err), err)
	}

	application_settings.ApplicationSettingsStatic = applicationSetting

	lp := startup.InitLogProvider()
	app.InitLogger(lp)

	traceShutdown := startup.InitTracer()
	if traceShutdown != nil {
		defer traceShutdown(ctx)
	}

	metricShutdown := startup.InitMeter()
	if metricShutdown != nil {
		defer metricShutdown(ctx)
	}
	app.InitMetrics()

	loadBashFiles()

	connErr := connections.LoadConnections()
	if connErr != nil {
		app.LogError2("Error connections: %v", connErr)
	}

	go token_credentials.LoadTokenCredentialAuthentication()
	hanlder := startup.LoadApplicationSettings(ctx, applicationSetting)
	port := getPort()
	app.LogInfo(fmt.Sprintf("Server listen in port %s", port))
	http.ListenAndServe(port, hanlder)
}

func loadBashFiles() {
	envbashFiles := os.Getenv(app.ENV_RUN_BASH_FILES_BEFORE_STARTUP)

	if len(envbashFiles) == 0 {
		envbashFiles = "wrench/bash/startup.sh"
	}

	bashFiles := strings.Split(envbashFiles, ",")
	bashRun(bashFiles)
}

func bashRun(paths []string) {
	for _, path := range paths {
		path = strings.TrimSpace(path)
		if _, err := os.Stat(path); err != nil {
			app.LogInfo(fmt.Sprintf("file bash %s not found", path))
			continue
		}

		app.LogInfo(fmt.Sprintf("Will process file bash %s", path))
		cmd := exec.Command("/bin/sh", "./"+path)

		output, err := cmd.Output()
		if err != nil {
			app.LogError2(err.Error(), err)
			return
		} else {
			app.LogInfo(string(output))
		}
	}
}

func getFileConfigPath() string {
	configPath := os.Getenv(app.ENV_PATH_FILE_CONFIG)
	if len(configPath) == 0 {
		configPath = "../../configApp.yaml"
	}
	return configPath
}

func getPort() string {
	port := os.Getenv(app.ENV_PORT)
	if len(port) == 0 {
		port = ":9090"
	}

	if port[0] != ':' {
		port = ":" + port
	}

	return port
}
