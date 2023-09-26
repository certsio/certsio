# certsio
The command-line client and library for using the certs.io API.

Search the entire internet by data in TLS certificates.

# WARNING:
This is experimental and subject to change. Use at your own risk. 

# Resources
- [Usage](#usage)

## Usage:
**Do not implicitly trust results, certificates can be self-signed.**

|                                       |                                                                                                      | 
|---------------------------------------|------------------------------------------------------------------------------------------------------|
| **Search by root domain**             | `certsio search domain example.com`                                                                  |
| **Search by organization**            | `certsio search org "Uber Technologies, Inc."`                                                       |
| **Search in emails** | `certsio search emails @example.com`, `certsiosearch emails keyword`                                 |
| **Search by server + port**           | `certsio search server 3.20.2.223:443`                                                               |
| **Search by certificate serial**      | `certsio search serial 0C:1F:CB:18:45:18:C7:E3:86:67:41:23:6D:6B:73:F1`                              |
| **Search by certificate fingerprint** | `certsio search fingerprint_sha256 5ef2f214260ab8f58e55eea42e4ac04b0f171807d8d1185fddd67470e9ab6096` |
| **Search by ssl name**                | `certsio search ssl_names www.github.com`                                                            | 
| **Print version**                     | `certsio version`                                                                                    |

### Configuration File:

Get an API Key from [here](https://rapidapi.com/certsio-certsio-default/api/certs-io1/pricing)

certsio automatically looks for a configuration file at `$HOME/.certsio.toml` or `%USERPROFILE%\.certsio.toml`. You must specify your API key in this configuration file for the client to work. 

An example configuration file can be found in the repository.


## Installation:
### From source:
```
$ go install github.com/certsio/certsio/cmd/certsio@latest
```
### From github :
```
git clone https://github.com/certsio/certsio.git; \
cd certsio/certsio/cmd; \
go build; \
sudo mv certsio /usr/local/bin/; \
certsio --version;
```
