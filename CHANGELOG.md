# RegistryMate Changelog

All notable changes to this project will be documented in this file. The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Released]

## [1.1.0] - 2026-02-02

### Added
- There is now an About-Page that you can access by clicking on the info button in the upper right corner

---

## [1.0.0] - 2026-01-09
**Starting the journey of RegistryMate [First Stable Release]**

### Added
- Desktop application for creating Kubernetes ImagePullSecrets
- YAML generation for `kubernetes.io/dockerconfigjson` secrets
- Copy generated secrets to clipboard
- Save secrets as YAML files
- Registry and secret metadata history (up to 100 entries)
  - Automatic cleanup of oldest entries when the limit is reached
- Base64 encode and decode utility for Docker auth strings
- Cross-platform support (Windows, Linux)
- Desktop UI built with Go and Fyne