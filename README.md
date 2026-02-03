# Anvil
CLI wrapper around virsh & virt-manager used for cloning VMs based on a specified golden image and adjusting their hardware configurations based on the expected workload.

## Preconfigured Hardware Allocation Profiles
- Workload profiles are divided into `[Light, Medium, Heavy, Ultra, Max]` categories with an additional enhanced variation `+` for each category.
- Profile allocations are individually overridable at creation-time.
- Allocations are meant to span a system that has at least 24 CPU cores and 128GB of RAM.
- It is advisable that at least 2 cores & 16 GB should be reserved for the host machine for reliable stability & performance.

### Preconfigured Profile Allocations Table
| Profile | vCPU | RAM | Diskspace | RAM / vCPU |
| --- | --- | --- | --- | --- |
| Light | 1 | 2 GB | 12 GB | 2 GB / core |
| ***Light+*** | ***2*** | ***4 GB*** | ***24 GB*** |
| Medium | 4 | 8 GB | 36 GB |
| ***Medium+*** | ***6*** | ***12 GB*** | ***64 GB*** |
| Heavy | 8 | 24 GB | 128 GB | 3 GB / core |
| ***Heavy+*** | ***12*** | ***36 GB*** | ***256 GB*** |
| Ultra | 16 | 64 GB | 512 GB | ~4 GB / core |
| ***Ultra+*** | ***18*** | ***80 GB*** | ***512 GB*** |
| Max | 20 | 96 GB | 1 TB | ~5 GB / core |
| ***Max+*** | ***22*** | ***112 GB*** | ***1 TB*** |

## Profiling Policy
You might be asking yourself, *how do I know which profile would I need for `<insert specific task>`?* And *when would I prefer to use containers over virtual machines?*

Below is a breakdown that should hopefully shine a light on the stengths and weaknesses of each workload type and some example use-cases.

| Workload Type | Intent | Generalized Use Case Description |
| --- | --- | --- |
| Containers | Efficiency | Trusted code, speed, density, process-level isolation |
| Light | Isolation | Untrusted or infra-like tasks, real OS boundary with minimal footprint |
| Medium | Usability | Human-facing workspaces |
| Heavy | Throughput | Performance work |
| Ultra | Dedication | Dedicated domain machines |
| Max | Experimental Dedication | Host-takeover experimentation |


## Containers
Use a container when *you trust the code* and *don't require a real machine boundary*. Because containers are ephemeral they are great for isolated (atom-like) tasks that can be managed by higher level software machinery like k8s when it comes to their lifecycles & retry policies.

### Typical Use Cases
- Build jobs (CI-style)
- Linting / formatting
- Unit tests
- Language toolchains (Node, Python, Go, Rust)
- Terraform / Ansible execution
- Small dev databases (Postgres, Redis, MySQL)
- Local dev services (APIs, workers, queues)
- One-off scripts
- Data format conversion
- Web scrapers you trust
- Background utilities

### Strengths
- Fast startup
- High density
- Minimal overhead
- Easy reproducibility

### Why not a VM? (Weaknesses)
- No kernel separation required
- OS realisim not needed
- Blast radius is acceptable if something were to break out


## Light :star:
The Light profile is on the other side of the edge of the blade of containers. This is meant to run operating systems with the most minimal overhead and most likely without a GUI.

### Typical Use Cases
- Bastion host / jump host / access VM
- "Toolbox VM" for untrusted binaries, scripts, downloads
- Malware or exploit testing (basic level)
- Network tooling:
  - DNS / DHCP
  - routing experiments
  - firewall testing
- Automation controller nodes (Ansible/Terraform) when isolation matters
- Security tooling where you don’t trust dependencies
- Disposable infra nodes in a lab

### Why a VM instead of a container?
- Separate kernel
- Separate syscall surface
- Easy full teardown
- Accurate OS behavior

### Why not Medium
- No GUI
- No IDE
- No sustained interactive load


## Medium :star::star:
The Medium profile is geared towards general purpose computing. Allows for enough horsepower to run a GUI comfortably. More for general web browsing or light prototyping.

### Typical Use Cases
- General development (IDE + browser)
- Cloud learning labs (AWS/GCP/Azure)
- IaC work with local testing
- Light Docker usage (compose with a few services)
- “Safe browsing” VM
- Personal admin tasks (finance, docs, portals)
- GUI-based security tools
- Small data analysis (pandas, notebooks)

### Why Medium VM and not a Light
- GUI
- More broader toolset installed
- Long-lived state with user ergonomics

### Why not Heavy
- Not build or compute-heavy hardware allocation
- Meant more for responsiveness than throughput


## Heavy :star::star::star:
The Heavy profile is designed to be your go-to performance workstaion. Typically more data-intensive processes.

### Typical Use Cases
- Large codebases
- Parallel builds
- Docker-heavy dev stacks
- Multiple local services
- ML prototyping (CPU-heavy, GPU-light)
- Jupyter with moderate datasets
- Local observability stacks
- Databases with meaningful data volume
- Security labs with multiple tools running concurrently

### Why Heavy?
- CPU parallelism matters
- Memory caching matters
- Medium starts to feel sluggish


## Ultra :star::star::star::star:
The Ultra profile is for tasks that require the majority of the physical machine's capabilities. Ultra VMs are essentially dedicated domain machines.

### Typical Use Cases
- ML training or serious experimentation
- GPU passthrough development
- Large datasets in memory
- Long-running notebooks / experiments
- Big infra simulations (k8s, distributed systems)
- Integration test environments
- “Primary machine for this domain”

### Operational Reality
- Typically 1 Ultra VM active at a time
- Other VMs are Light/Medium or paused
- Meant for operational maximum while maintaining system stability

### Why not Heavy
- Want to avoid resource contention
- Require even more horsepower than what heavy provides


## Max :star::star::star::star::star:
Max is the experimental profile meant to yield the maximum amount of resources from the physical workstation while still providing some guardrails for basic reliability. Stability is traded for performance, and may require tweaking to get the right balance.

### Typical Use Cases
- Near–bare-metal experiments
- Stress testing host + storage + memory
- Huge in-memory workloads
- Large model training
- Performance characterization
- One-off, high-risk experiments

### Operational Expectations
- Host responsiveness may degrade
- Other VMs stopped or minimal
- Manual tuning likely required
