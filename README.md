# Opsgenie Exporter

Opsgenie Exporter is a Prometheus exporter for sourcing metrics from Opsgenie. It aims to provide real-time insights into Opsgenie's data such as users, teams, integrations, and heartbeats. This tool is designed to help teams monitor their Opsgenie configuration and operation statuses directly from a Prometheus setup.

## Installation

To install Opsgenie Exporter, you will need to clone this repository and build the binary from source. Ensure you have Go installed on your system.

```bash
git clone https://github.com/jsedy7/opsgenie-exporter.git
cd opsgenie-exporter
go build -o opsgenie-exporter .
```

## Usage

Run the exporter using the following command:

```bash
./opsgenie-exporter
```

By default, the exporter runs on port 8080 and updates metrics every 600 seconds. These settings can be customized using command-line flags.

### Command-Line Flags

- `--http.port`: Specifies the port on which the exporter server runs. Default is `8080`.
- `--refresh`: Sets the interval for metrics update in seconds. Default is `600`.

Example usage with flags:

```bash
./opsgenie-exporter --http.port=9090 --refresh=300
```

This command runs the exporter on port 9090 and updates metrics every 5 minutes.

## Configuration

The Opsgenie Exporter requires an Opsgenie API key for fetching data. Set your Opsgenie API key as an environment variable:

```bash
export OPSGENIE_API_KEY='your_opsgenie_api_key_here'
```

Ensure this environment variable is set before running the exporter.

## Metrics

Go to [metrics documentation](docs/metrics.md).

## Contributing

Contributions to the Opsgenie Exporter are welcome and appreciated. Here are ways you can contribute:

- Submit bugs and feature requests.
- Review the source code and improve code quality.
- Add new features or enhance existing ones.

### Submitting Pull Requests

1. Fork the repository.
2. Create a new branch for your feature or fix.
3. Commit your changes with meaningful commit messages.
4. Push your branch and submit a pull request against the main branch.

## License

This project is licensed under the MIT License - see the LICENSE file for details.


