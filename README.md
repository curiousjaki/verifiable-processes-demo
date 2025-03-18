# Verifiable Computation in Business Processes

## Description

This project includes the implementation for the verifiable computing in business processes research paper.
I consists of 
- a Risc0 verifiable computing service written in Rust,
- a Zeebe workflow engine adapter service written in Go,
- kubernetes deployment files managed by a HELM chart that help to setup a local Camunda platform 8.14. as well as the above mentioned components, and
- a process instance generator, used for testing and performance evaluations.

## Installation

To install this project, you need to have a running kubernetes environment and deploy the camunda platfrom as well as our verifiable computing components.

Either reference to the according HELM Chart: (kubernetes-deplyoment/helm/)[kubernetes-deployment/helm/]

```bash
# Add installation commands here
```
## Caveats 
- The risczero host and guest versions must match