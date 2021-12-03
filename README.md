![Build Status](https://github.com/pokornyIt/cucm_performance_exporter/workflows/Full%20release%20workflow/badge.svg)
[![GitHub](https://img.shields.io/github/license/pokornyIt/cucm_performance_exporter)](/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/pokornyIt/cucm_performance_exporter)](https://goreportcard.com/report/github.com/pokornyIt/cucm_performance_exporter)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/pokornyit/cucm_performance_exporter?label=latest)](https://github.com/pokornyIt/cucm_performance_exporter/releases/latest)

# CUCM Performace Exporter

Prometheus exporter for CISCO Unified Communication Manager performance metrics. Project
utilize [CISCO CUCM Performance API](https://developer.cisco.com/site/sxml/discover/overview/perfmon/). System tested on
version 10+. Detail API description is on [CISCO DEVNET](https://developer.cisco.com/docs/sxml/#!perfmon-api-reference)
servers.

# Configuration

Program need configuration file with next structure and information.

```yaml
monitor_names: [ 'publisher.name','subscriber01.name' ]
metrics:
  callsActive: true
  callsInProgress: true
  callsCompleted: true
  partiallyRegisteredPhone: true
  registeredHardwarePhones: true
  gatewaysSessionsActive: true
  gatewaysSessionsFailed: true
  phoneSessionsActive: true
  phoneSessionsFailed: true
  goCollector: true
  processStatus: true
port: 9719
apiAddress: publisher.name
apiUser: api_allowed_user
apiPwd: password
ignoreCertificate: true
log:
  level: info
  fileName: ''
  jsonFormat: false
  logProgramInfo: false
  maxSize: 50
  maxBackups: 5
  maxAge: 30
  quiet: false
```

- **monitor_names** - name of CUCM servers, use same names as in system CUCM configuration
- **metrics** - allowed or disabled metrics collected from CUCM cluster
- **port** - port where program start HTTP server with metrics
- **apiAddress** - FQDN or IP address of publisher server
- **apiUser** - user with rights to read performance metrics
- **apiPwd** - password for apiUser
- **ignoreCertificate** - system ignore certificate validity
- **log** - setup logging from system

## Actual supported metrics

- **callsActive** - This represents the number of voice or video streaming connections that are currently in use (
  active).
- **callsInProgress** - This represents the number of voice or video calls that are currently in progress on this
  CallManager, including all active calls.
- **callsCompleted** - This represents the number of calls that were actually connected (a voice path or video stream
  was established) through this CallManager.
- **partiallyRegisteredPhone** - This represents the number of partially registered SIP Phones.
- **registeredHardwarePhones** - This represents the number of Cisco hardware IP phones (for example, models 7960, 7940,
  7910, etc.) that are currently registered in the system.
- **gatewaysSessionsActive** - This is a real-time counter which specifies the total number of active recording sessions
  between a recording-enabled gateway and a recording server.
- **gatewaysSessionsFailed** - This is a cumulative counter which specifies the total number of gateway-preferred
  recording sessions which failed since the last restart of the Cisco Unified Communications Manager service.
- **phoneSessionsActive** - This is a real-time counter which specifies the total number of active recording sessions
  between a Cisco IP Phone and a recording server.
- **phoneSessionsFailed** - This is a cumulative counter which specifies the total number of phone-preferred recording
  sessions which failed since the last restart of the Cisco Unified Communications Manager service.

Program allow enabling/disabling standard GO client metrics. Detail about this metrics are described in [Exploring Prometheus GO client Metrics](https://povilasv.me/prometheus-go-metrics/#).

- **goCollector** - enable/disable internal program GO metrics
- **processStatus** - enable/disable internal program status metrics

## Log setup

- **level** - Logging level, default Info, valid: Fatal, Error, Warning, Info, Debug, Trace
- **fileName** - File name for actual log file, empty doesn't log to file
- **jsonFormat** - Use JSON formatting, default false, valid: false, true
- **logProgramInfo** - Include in log line and source file from source program
- **maxSize** - maximal size of one log file in MB, default 50 MB, minimal1 MB, max 5 000 MB
- **maxBackups** - maximal backup log files, default 5, minimal 0 file, max 100 files
- **maxAge** - maximal log file age in day, default 30, minimal 1 day, maximal 365 days
- **quiet** - don't log any message to std output, default false, valid: false, true

# Start parameters

Program support CLI parameters. All parameters are optional and overwrite same configuration values.

- **--version** - show actual program version
- **--config.show** - show actual configuration and ends
- **--config.file file_name** - start program with configuration file file_name, when omitted system use name "
  server.yml" in current directory
- **--api.address fqdn_or_ip** - overwrite value `apiAddress` from configuration file with fqdn_or_ip
- **--api.user user_name** - overwrite value `apiUser` from configuration file with user_name
- **--api.pwd user_pwd** - overwrite value `apiPwd` from configuration file with user_pwd

# Contribute

We welcome any contributions. Please fork the project on GitHub and open Pull Requests for any proposed changes.

Please note that we will not merge any changes that encourage insecure behaviour. If in doubt please open an Issue first
to discuss your proposal.

# Sponsoring

Many thanks **Elevēo** for support this project.

[![Elevēo - powered by ZOOM International](.github/eleveo-logo.png)](https://eleveo.com)