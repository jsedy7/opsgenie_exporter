# Opsgenie Exporter Metrics Documentation

Opsgenie Exporter provides a set of Prometheus metrics that offer insights into various aspects of Opsgenie configuration and usage. Below is a detailed description of each metric exported by the Opsgenie Exporter.

## General Metrics

### `opsgenie_last_update_timestamp_seconds`
- **Description**: Timestamp of the last successful update of the metrics from Opsgenie.
- **Type**: Gauge
- **Labels**: None

## User Metrics

### `opsgenie_users{key="total"}`
- **Description**: The total number of users in Opsgenie.
- **Type**: Gauge
- **Labels**:
  - `key`: Fixed label with value "total".

### `opsgenie_users{key="blocked"}`
- **Description**: The number of users currently blocked in Opsgenie.
- **Type**: Gauge
- **Labels**:
  - `key`: Fixed label with value "blocked".

### `opsgenie_users{key="unverified"}`
- **Description**: The number of users with unverified accounts in Opsgenie.
- **Type**: Gauge
- **Labels**:
  - `key`: Fixed label with value "unverified".

### `opsgenie_user_verified_status{username="<username>"}`
- **Description**: Indicates the verification status of a specific Opsgenie user. `1` for verified, `0` for unverified.
- **Type**: Gauge
- **Labels**:
  - `username`: The username of the Opsgenie user.

## Team Metrics

### `opsgenie_teams_total`
- **Description**: The total number of teams configured in Opsgenie.
- **Type**: Gauge
- **Labels**: None

## Account Metrics

### `opsgenie_account{key="<attribute>"}` 
- **Description**: Provides various attributes of the Opsgenie account.
- **Type**: Gauge
- **Labels**:
  - `key`: The specific attribute of the account being reported. Possible values include "userCount", "maxUserCount", and "isYearly".

## Integration Metrics

### `opsgenie_integrations_total{type="<integration_type>"}` 
- **Description**: The total number of Opsgenie integrations by type.
- **Type**: Gauge
- **Labels**:
  - `type`: The type of the integration (e.g., "email", "API", "slack").

## Heartbeat Metrics

### `opsgenie_heartbeats_total`
- **Description**: The total number of heartbeats configured in Opsgenie.
- **Type**: Gauge
- **Labels**: None

### `opsgenie_heartbeats_enabled_total`
- **Description**: The total number of enabled heartbeats in Opsgenie.
- **Type**: Gauge
- **Labels**: None

### `opsgenie_heartbeats_expired{team="<team_name>"}`
- **Description**: Indicates whether a specific Opsgenie heartbeat is expired. `1` for expired, `0` for not expired.
- **Type**: Gauge
- **Labels**:
  - `team`: The name of the team owning the heartbeat.

This documentation should be included in your project repository to assist users in understanding the metrics provided by the Opsgenie Exporter.
