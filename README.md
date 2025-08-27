# Verifiable Computation in Business Processes

## Abstract

Ensuring the integrity of business processes without disclosing confidential business information is a major challenge in inter-organizational processes. 
This paper introduces a zero-knowledge proof (ZKP)-based approach for the verifiable execution of business processes while preserving confidentiality. 
We integrate ZK virtual machines (zkVMs) into business process management engines through a comprehensive system architecture and a prototypical implementation. 
Our approach supports chained verifiable computations through proof compositions. 
On the example of product carbon footprinting, we model sequential footprinting activities and demonstrate how organizations can prove and verify the integrity of verifiable processes without exposing sensitive information. 
We assess different ZKP proving variants within process models for their efficiency in proving and verifying, and discuss the practical integration of ZKPs throughout the Business Process Management (BPM) lifecycle. 
Our experiment-driven evaluation demonstrates the automation of process verification under given confidentiality constraints.

This is the accompanying repository for the 

## Description

This project includes the implementation for the verifiable computing in business processes research paper.
It consists of:
- a Risc0 verifiable computing service written in Rust,
- a Zeebe workflow management ambassador written in Java,
- Camunda deployment files for running Camunda locally.



## Installation

### On bare metal and Docker Compose:

*Prerequisites:* Before you begin, ensure you have the following installed:
- Docker and Docker Compose
- cargo-risczero @ 2.3.1, cpp @ 2024.1.5, r0vm @ 2.3.1, rust @ 1.88.0

To install this project, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/verifiable-processes-demo.git
   cd verifiable-processes-demo
   ```

2. Run the Camunda Platform within Docker:
   ```bash
   docker compose -p zkvm4bpm down
   docker compose -f docker-compose/docker-compose.yaml -p zkvm4bpm --env-file ./docker-compose/.env up -d --pull always 
   ```

3. Deploy the Camunda platform:
   - Navigate to the `camunda` directory and follow the deployment instructions provided in the `README.md` file there.

4. Build the Risc0 verifiable computing service:
   ```bash
   cd risc0-service
   cargo build --release
   ```

5. Start the Zeebe workflow management ambassador:
   ```bash
   cd zeebe-ambassador
   ./gradlew bootRun
   ```

## Usage

1. Access the Camunda platform at `http://localhost:8080`.
2. Deploy your BPMN workflows using the Camunda Modeler.
3. Interact with the Risc0 service and Zeebe ambassador as part of your workflow execution.

## Caveats

## Contributing

Contributions are welcome! To contribute:
1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Commit your changes and open a pull request.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

## Citation

```bibtex
@inproceedings{kiesel_confidentiality-preserving_2025,
	title = {{Confidentiality}-{Preserving} {Verifiable} {Business} {Processes} through {Zero}-{Knowledge} {Proofs}},
	booktitle = {Enterprise {Design}, {Operations}, and {Computing}},
	publisher = {Springer Nature Switzerland},
	author = {Kiesel, Jannis and Heiss, Jonathan},
    doi = {},
    note = {Preprint Version},
	year = {2025},
}
```


