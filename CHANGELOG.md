# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Unit tests and benchmarks
- Pre-submit checklist: `docs/checklist.md`
- CORS configurable whitelist (security hardening)
- Request body size limits: MaxHeaderBytes + MaxBytesReader (security hardening)
- XSS input sanitization in profile updates (security hardening)
- Timing attack protection in login (security hardening)
- Coverage reporting: `make coverage`
- `.dockerignore` for leaner builds

### Changed
- Provider architecture refactored to `Infra` + `Modules` builder pattern
- Health readiness checks perform real DB + Redis connectivity tests
- Package-level doc comments added across all packages
- `configs/config.yaml` removed from git tracking (contained secrets)

### Fixed
- Dockerfile Go version aligned with go.mod (1.25)
- `.gitignore` rules fixed and expanded (CORS secrets, test artifacts, IDE files)
- Cache benchmark TTL units corrected
