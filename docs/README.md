# GoWrench
GoWrench is a platform that makes life easier for developers and companies by offering a simple and efficient solution to solve common problems without the need to write repetitive code. The key feature is the use of YAML files, allowing for quick, clear, and flexible configuration and customization of services and APIs. This saves time and resources, enabling you to focus on what truly matters: delivering value and innovation.

The application will be packaged as a Docker image, ensuring seamless deployment and scalability, making it even easier to integrate into any development environment.


## Install

```
docker run -d -p 8080:8080 --name gowrench wrenchapps/gowrench:latest
```


## Features

- **Declarative Development**: Create applications without writing code, only through YAML configurations
- **Service Integration**:
  - HTTP requests
  - AWS SNS integration
  - Mock responses for testing
- **Data Transformation**:
  - Header mapping
  - Request/response body transformation
  - Contract mapping
- **Configuration Management**:
  - YAML-based configuration
  - Environment variable support
  - Multiple environment profiles

## Project Structure

```
.
├── app/                    # Main application code
│   ├── auth/              # Authentication related code
│   ├── clients/           # External service clients
│   ├── cmd/               # Application entry points
│   ├── contexts/          # Context management
│   ├── handlers/          # Request handlers
│   ├── json_map/          # JSON mapping utilities
│   ├── manifest/          # Application manifest
│   └── startup/           # Application startup code
├── docs/                  # Documentation
├── configApp.yaml         # Main configuration file
├── configAppProxy.yaml    # Proxy configuration
├── docker-compose.yaml    # Docker compose configuration
├── Dockerfile            # Docker build configuration
└── go.mod                # Go module definition
```

## Configuration

The application is configured through `configApp.yaml`, which supports:

- Service configuration
- API endpoints and routes
- Authentication settings
- Action definitions
- Contract mappings

## Environment Variables

The application supports multiple environment profiles:
- `.ENV` - Default environment variables
- `.ENV.development` - Development-specific variables

## Docker Support

The application can be containerized using Docker:

```bash
docker-compose up
```

## Getting Started

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Configure environment variables
4. Run the application:
   ```bash
   ./startup.sh
   ```

## Available Endpoints

The application provides several configurable endpoints:

- `GET /api/mock` - Mock endpoint for testing
- `GET /api/customer` - Retrieve customer information
- `POST /api/customer` - Create new customer
- `GET /api/hellos` - Example endpoint
- `POST /api/sns` - AWS SNS integration endpoint
- `GET /api/json` - JSON file reading endpoint

## Security

- JWT-based authentication
- Role-based access control
- Scope-based permissions
- Header-based security controls

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the GPL-3.0 - see the LICENSE file for details.
