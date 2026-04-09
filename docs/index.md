---
layout: home

hero:
  name: opentelemetry-go
  text: Extensions for OpenTelemetry Go
  tagline: Beautiful terminal exporters, semantic conventions, and test helpers for the OpenTelemetry Go SDK
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/foomo/opentelemetry-go
  image:
    src: /logo.png
    alt: opentelemetry-go

features:
  - title: Glossy Exporters
    details: Lipgloss-styled terminal output for traces and metrics. Tree views, flamegraphs, and color-coded durations make debugging a pleasure.
    link: /guide/exporters
    linkText: Learn more
  - title: Semantic Conventions
    details: Reusable attribute constructors for GoTSRPC, Keel, HTTP, and more. Follow the upstream OTel pattern with typed key-value helpers.
    link: /guide/semconv
    linkText: Learn more
  - title: Testing Helpers
    details: One-line ReportTraces and ReportMetrics wiring for testing.T and TestMain. Automatic cleanup, zero boilerplate.
    link: /guide/testing
    linkText: Learn more
---
