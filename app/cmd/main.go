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
	"wrench/app/middleware"
	"wrench/app/startup"
	"wrench/app/startup/connections"
	"wrench/app/startup/token_credentials"
)

func main() {
	ctx := context.Background()
	app.SetContext(ctx)

	startup.LoadEnvsFiles()

	pathFiles, err := startup.GetFileConfigPath()
	if err != nil {
		app.LogError2(fmt.Sprintf("Error to load files env config: %v", err), err)
	}

	byteArray, err := startup.LoadYamlFile(pathFiles)
	if err != nil {
		app.LogError2(fmt.Sprintf("Error loading YAML: %v", err), err)
	}

	err = startup.LoadAwsSecrets(byteArray)
	if err != nil {
		app.LogError2(fmt.Sprintf("Error loading YAML: %v", err), err)
	}

	byteArray = startup.EnvInterpolation(byteArray)
	applicationSetting, err := application_settings.ParseMapToApplicationSetting(byteArray)

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
	http.ListenAndServe(port, middleware.CaseInsensitiveMux(hanlder))
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
