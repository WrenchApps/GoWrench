package app

import "go.opentelemetry.io/otel"

const ENV_PORT = "PORT"
const ENV_PATH_FILE_CONFIG string = "PATH_FILE_CONFIG"
const ENV_PATH_FOLDER_ENV_FILES string = "PATH_FOLDER_ENV_FILES"
const ENV_APP_ENV string = "APP_ENV"
const ENV_RUN_BASH_FILES_BEFORE_STARTUP string = "RUN_BASH_FILES_BEFORE_STARTUP"

var Tracer = otel.Tracer("trace")
