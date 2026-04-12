name: Bug Report
description: Report a bug to help us improve
title: "[Bug] "
labels: ["bug"]
body:
  - type: markdown
    attributes:
      value: |
        Thanks for reporting a bug! Please provide as much detail as possible.

  - type: textarea
    id: description
    attributes:
      label: Description
      description: Clear description of the bug
      placeholder: "What happened?"
    validations:
      required: true

  - type: textarea
    id: expected
    attributes:
      label: Expected Behavior
      description: What should have happened?
    validations:
      required: true

  - type: textarea
    id: steps
    attributes:
      label: Steps to Reproduce
      description: Steps to reproduce the behavior
      placeholder: |
        1. Run `...`
        2. Execute `...`
        3. See error
    validations:
      required: true

  - type: textarea
    id: environment
    attributes:
      label: Environment
      description: Your environment details
      placeholder: |
        - OS: 
        - Go Version: 
        - Application Version:
    validations:
      required: true

  - type: textarea
    id: logs
    attributes:
      label: Logs/Error Output
      description: Relevant log output
      render: shell

  - type: textarea
    id: additional
    attributes:
      label: Additional Context
      description: Any other context
