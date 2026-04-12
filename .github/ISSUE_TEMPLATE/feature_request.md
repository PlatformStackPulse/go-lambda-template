name: Feature Request
description: Suggest a feature for this project
title: "[Feature] "
labels: ["enhancement"]
body:
  - type: markdown
    attributes:
      value: |
        Thanks for suggesting a feature! Help us make this project better.

  - type: textarea
    id: description
    attributes:
      label: Description
      description: Clear description of the feature
      placeholder: "What feature would you like to see?"
    validations:
      required: true

  - type: textarea
    id: motivation
    attributes:
      label: Motivation
      description: Why is this feature needed?
    validations:
      required: true

  - type: textarea
    id: solution
    attributes:
      label: Proposed Solution
      description: How should this feature work?

  - type: textarea
    id: alternatives
    attributes:
      label: Alternative Solutions
      description: Other approaches you've considered
