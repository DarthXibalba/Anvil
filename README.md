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


## Light VM - 1 vCPU / 2 GB
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