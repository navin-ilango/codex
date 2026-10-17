# Kubernetes and Amazon EKS architecture

1. Design a Kubernetes platform on AWS that supports more than 1,000 microservices across multiple EKS clusters. How would you approach architecture, networking, security, scaling, and observability?
    - `Cluster Architecture & Multi-Tenancy`
        - `Account Structure` - Separate AWS accounts for Management/Shared Services, Development, Staging, and Production.
        - `Deploy EKS clusters per environment and region`. Microservices are grouped into clusters by domain/business unit (e.g., prod-payment-us-east-1, prod-search-us-east-1) rather than a single monolithic cluster.
        - `Intra-Cluster Isolation`: Within each cluster, use Kubernetes Namespaces as the boundary for individual microservices or small sub-domains.
        - `Compute Plane`: Use AWS Karpenter instead of traditional Cluster Autoscaler. Karpenter provisions nodes dynamically based on pending pod requirements (gpus, architectures, spot vs. on-demand), optimizing costs and drastically reducing scheduling latency.
    - `Networking & Service Mesh`
        - Scaling to 1,000+ microservices introduces massive East-West traffic volumes and complex service discovery.
        - `Pod Networking`: Use the AWS VPC CNI in prefix delegation mode to prevent IP address exhaustion. Allocate secondary IPv4 CIDR blocks specifically for EKS Pods.
        - `East-West Traffic & Service Mesh`: Istio or Cilium.
        - `North-South Traffic`: Deploy the AWS Load Balancer Controller. It automatically provisions AWS Application Load Balancers (ALBs) or Network Load Balancers (NLBs) based on Kubernetes Ingress or Service resources. Integrate this with an API Gateway (like Envoy) to handle global rate-limiting, authentication, and request routing to the underlying microservices.
    - `Enterprise Security & Governance`
        - `Identity & Access Management`: Eliminate long-lived IAM keys by using EKS Pod Identities (the modern successor to IRSA). This maps AWS IAM roles directly to Kubernetes ServiceAccounts, enforcing the principle of least privilege for microservices accessing S3, DynamoDB, or RDS.
        - Amazon EKS user authentication with Okta using OpenID Connect (OIDC).
        - `Network Segmentation`: Istio - Block All External Traffic by Default by setting the outboundTrafficPolicy.mode: REGISTRY_ONLY and Once REGISTRY_ONLY is active, applications will receive 502 Bad Gateway errors when trying to reach external websites. You must use a ServiceEntry to whitelist specific external domains. f your platform team explicitly defines a custom Sidecar resource in a namespace to restrict configurations and optimize memory usage across large proxy fleets, it usually limits visibility using the egress.hosts block. Because the Sidecar resource explicitly whitelists which namespaces it pulls configuration from, your app sidecar will ignore the external ServiceEntry unless - The ServiceEntry lives in the same namespace (./*) or The ServiceEntry is in a dedicated namespace (e.g., egress-controls/*).
        - `Policy Enforcement` - Deploy Kyverno or Open Policy Agent (OPA) Gatekeeper as a validating/mutating admission webhook. Block pods that attempt to run as root, lack standard resource limits, or pull images from untrusted registries.
        - `Secrets Management`: Do not store plain text secrets in Git. Use External Secrets Operator (ESO) integrated with AWS Secrets Manager. ESO automatically syncs secrets into native Kubernetes secret objects at runtime.
    - `Advanced Scaling & Deployment`
        - `Deployment Automation`: Standardize deployments using ArgoCD in a hub-and-spoke model. A centralized control-plane ArgoCD cluster monitors a Git repository and orchestrates application deployments out to the spoke EKS clusters.
        - `Application Packaging`: Standardize microservice definitions using Helm charts or Kustomize templates managed by a central platform team. Developers supply minimal configuration parameters (like image tags and environment variables), while the platform handles the underlying Kubernetes boilerplate.
        - `Microservice Auto-Scaling`: **Horizontal Pod Autoscaler (HPA)**: Scale pods using custom application metrics (e.g., Prometheus request rates or SQS queue depth) via the KEDA (Kubernetes Event-driven Autoscaling) operator, instead of relying purely on generic CPU/Memory metrics. **Vertical Pod Autoscaler (VPA)**: Run VPA in Recommendation mode to continuously analyze production metrics and guide developers on sizing their resource requests/limits accurately.`
    - `Observability & Telemetry`
        - 1,000+ microservices will generate millions of metrics and spans. The observability stack must scale horizontally.
        - Metrics: Prometheus, Victoria Metrics
        - Logs: Sumologic
        - Traces: Otel Collectors and traces are sent to Honeycomb.
        - Dashboards - Grafana

2. How do you decide whether to run multiple workloads in one cluster or distribute them across multiple clusters?
    - Criterias to think about,
        - Isolation Boundary - Soft isolation (Namespaces, NetworkPolicies, RBAC) (or) Hard isolation (AWS account, VPC, compute level) 
        - Operational Overhead - Low (Manage one control plane, fewer upgrades) (or) High (Manage multiple control planes, repeating configs)
        - Resource Efficiency - High (Better packing, shared idle capacity) (or) Low (Fragmented resources, duplicated tooling)
        - Blast Radius - High (Cluster-wide issues affect all workloads) (or) Low (Issues are contained to a single cluster)
        - Compliance/Security - Complex (Requires stringent controls and auditing) (or) Simple (Naturally meets strict data boundaries)
        - Compute Cost
    
    - `When to Choose a Single Multi-Tenant Cluster`
        - Opt for one large cluster divided by Kubernetes Namespaces when workloads share a similar trust level and operational lifecycle.
        - `Cost Optimization`: Shared infrastructure reduces the overhead of running multiple EKS control planes ($0.10/hour per cluster) and avoids underutilized worker nodes.
        - `Streamlined Operations`: You only need to maintain, patch, and upgrade a single control plane. Global add-ons (like Prometheus, FluentBit, or Karpenter) are installed just once.
        - `Microservices Architectures`: If your services frequently communicate with each other, keeping them in the same cluster minimizes cross-VPC or cross-cluster networking complexities and latency.
        - `Essential Best Practices for Multi-Tenancy`
            - Use `NetworkPolicies` to restrict traffic between namespaces.
            - `Enforce ResourceQuotas and LimitRanges` to prevent a single workload from consuming all cluster memory or CPU ("noisy neighbor" effect).
            - `Leverage IAM Roles for Service Accounts (IRSA) or EKS Pod Identities` to give namespaces scoped AWS permissions.
    - `When to Choose Multiple Dedicated Clusters`
        - `Strict Security & Compliance`: If you handle highly sensitive data (e.g., PCI-DSS, HIPAA) alongside non-regulated applications, separating them into different clusters (ideally in separate AWS accounts) simplifies auditing.
        - `Blast Radius Mitigation`: A misconfigured cluster-wide resource (like a bad ClusterRole or an aggressive CoreDNS failure) in a shared cluster can take down everything. Separate clusters isolate failures.
        - `Geographic or Latency Constraints`: If workloads must run close to users in different AWS Regions or Local Zones, you must deploy distinct clusters.
        - `Autonomous Teams`: If different engineering groups require completely different cluster-wide configurations, custom custom-resource definitions (CRDs), or different Kubernetes versions, giving them dedicated clusters prevents configuration conflicts.

3. Explain the Kubernetes control plane architecture and what happens when the API server becomes unavailable.
    - The Kubernetes control plane architecture acts as the "brain" of the cluster, maintaining the global desired state of your workloads. It functions as a distributed, decoupled system where components do not talk to each other directly; instead, every single control plane and worker node component communicates solely through the API server
    - Control Plane Components:
        - Kube API Server - The Gatekeeper & Hub. Exposes the REST API, handles authentication/authorization, and validates payload data. Only component that talks directly to etcd.
        - ETCD - The Source of Truth. Highly available, distributed key-value store holding the complete state configuration. Communicates exclusively with the API server.
        - Kube Controller Manager - The Reconciliation Loop. Runs controllers (Deployment, Node, StatefulSet) that track actual state vs. desired state. Watches API server; posts CRUD lifecycle requests back to it.
        - Kube Scheduler - The Placement Engine. Watches for newly created pods without an assigned node and selects the best worker node. Watches API server; updates pod specs with node assignments.
    - What Happens When the API Server Goes Down?
        - The Data Plane (Existing Applications Run Fine)
            - `Zero Container Outages`: Your running containers and pods continue executing exactly as they are. The container runtime (containderd) and the local kubelet agent on each worker node manage existing pods locally.
            - `Network Traffic Still Flows`: Traffic coming through external load balancers or internal service meshes still reaches your applications. The data plane routing rules are already programmed locally on the nodes by kube-proxy or your CNI plugin.
        - The Control Plane (Management Is Paralyzed)
            - Administrative Blindness: All kubectl commands, GitOps deployments, or API calls fail instantly. You cannot query cluster status, deploy code, or change configurations.
            - Total Loss of Self-Healing: If an existing worker node crashes or an application pod errors out, the cluster cannot recover. The kube-controller-manager cannot detect the dead node, and the kube-scheduler cannot schedule replacement pods.
            - Autoscaling Freezes: Horizontal Pod Autoscalers (HPA) and Cluster Autoscalers can neither read metrics nor trigger scale-up or scale-down events.
            - Stale Network Routing: If an active pod dies or is manually restarted at the container engine level, its replacement or missing status will not be propagated. Internal DNS and routing tables will become stale, leading to potential traffic drop-offs for that service. In a healthy cluster, pods are ephemeral—they frequently die, scale up, or move to different nodes, changing their internal IP addresses. The API server tracks these shifts and tells your network components to update their routing tables immediately. When the API server goes down, that communication chain is broken - `The Disconnect Between Pods and Services`
        - The Loophole: Creating "Static Pods" - You can still run containers on a specific worker node by taking advantage of Static Pods. These are pods managed directly by a specific node's kubelet daemon, completely independent of the master control plane. If a node is completely cut off from the network and you can't run kubectl, an SRE can SSH into the machine, drop a debugging tool's YAML into the manifests folder, and immediately have a local pod running to inspect host storage, local container logs, or run network diagnostics.
        - How to Access the Pod Anyway - Direct SSH + crictl.
4. How would you troubleshoot a pod stuck in Pending when the cluster has sufficient total CPU and memory?
    - The immediate first step is always to run kubectl describe pod <pod-name> and look at the "Events" section at the bottom. This will typically tell you exactly why the scheduler skipped certain nodes.
    - Causes:
        - Resource Fragmentation (No Single Node Fits)
        - Taints and Tolerations
        - Node Selectors and Affinity Rules(nodeAffinity, podAffinity, or podAntiAffinity).
        - Unfulfilled Persistent Volume Claims (PVC) - If your pod relies on persistent storage, the scheduler must find a node that can bind to or access that specific storage volume. The PVC might be bound to a specific availability zone (AZ), restricting the pod to nodes in that exact zone.
        - Node Port Conflicts (hostPort)
        - Namespace Quotas - A ResourceQuota applied to the namespace might be exhausted, preventing new pods from claiming resources in that specific namespace.

5. A pod is repeatedly restarting with OOMKilled, but its memory usage appears normal in Grafana. How would you investigate?
    - This discrepancy usually happens because Grafana dashboards typically show average memory usage over time or node-level metrics, which can smooth out and hide sudden, sharp memory spikes. Kubepods are terminated instantly the moment they hit their hard memory limit.
    - Verify the OOM Source
        - Check the Pod Termination Reason - If it says OOMKilled: true and shows an Exit Code: 137, the container hit its specific resource limit.
        - Check Node Events - Run kubectl get events -n <namespace> --sort-by='.metadata.creationTimestamp'. Look for SystemOOM events. If the node ran out of memory, the OS kernel will kill processes randomly to save the host, even if your container didn't hit its specific limit.
6. Explain the differences between requests, limits, QoS classes, and eviction thresholds. How do they affect scheduling and node stability?
    - `Requests`: The minimum amount of CPU or memory a container is `guaranteed` to get. Used by the scheduler to find a fitting node. Node capacity is tracked by the sum of all pod requests, not actual usage
    - `Limits`: The maximum amount of CPU or memory a container is allowed to consume. Prevents a single pod from starving other processes. Enforced at runtime by the container runtime via cgroups.
    - `Qos`: A classification (Guaranteed, Burstable, BestEffort) automatically assigned based on requests/limits. Determines pod priority during resource shortages. Dictates which pods are killed first when a node runs out of memory.
    - `Eviction Thresholds`: Settings defined on the kubelet (e.g., memory.available<100Mi) that trigger pod eviction. Proactively protects the node from crashing due to resource exhaustion. Reclaims node stability by gracefully or forcefully terminating pods.
    - `How They Affect Scheduling` - The Kubernetes scheduler (kube-scheduler) only looks at Requests, never Limits or actual usage. Overcommit: Because scheduling is based on requests, you can overcommit a node's limits. If five pods each have a request of 1GB but a limit of 4GB, they can all schedule onto a 5GB node. However, if they all try to use 4GB at the same time, the node will experience resource pressure.
    - `How They Affect Node Stability & Runtime` - When a node runs out of resources, Limits, QoS Classes, and Eviction Thresholds dictate how the node protects itself. Kubernetes automatically assigns one of three QoS classes to every pod.This class determines the pod's oom_score_adj (Out of Memory score adjustment), defining who gets killed first during a crisis.
        - Guaranteed - Highly stable. They are the last to be killed.
        - Burstable - Moderately stable. They can "burst" past their requests if capacity exists, but risk eviction if the node runs low.
        - BestEffort - Unstable. They consume whatever is free on the node, but they carry the highest oom_score_adj and are the first to be killed during resource pressure.
    - Eviction Thresholds (The Last Line of Defense)
        - The kubelet monitors node health. If a resource drops below an Eviction Threshold, the node enters a pressure state (e.g., MemoryPressure) and stops accepting new pods. It then begins evicting existing pods to free up resources, ranking them by their QoS class.
        - Hard Eviction Thresholds / Soft Eviction Thresholds 
        - Eviction Signals & Corresponding Node Conditions
            - memory.available -> MemoryPressure
            - nodefs.available / nodefs.inodesFree -> DiskPressure
            - imagefs.available / imagefs.inodesFree -> DiskPressure
            - pid.available -> PIDPressure
7. How would you design a safe Kubernetes version upgrade process across 12 or more production clusters?
    - Pre-Upgrade Validation & Automation
        - API Deprecation Checks: Pluto
        - Kube Addons Compatibility Checks: Ingress Controllers, Kube Dashboard, CSI Drivers, External DNS, Monitoring Stack, Etc.
    - The Phased Rollout Schedule
        - Divide your 12+ clusters into distinct deployment rings. Never upgrade more than one production environment at the same time within the same region.
    - Step-by-Step Cluster Upgrade Execution
        - State Backup - Ensure that the k8s backups are available. We use Velero.
        - Control Plane Upgrade (N → N+1)- If using managed services (EKS, GKE, AKS), trigger the control plane upgrade via Terraform/API. The cloud provider handles the rolling update of the API servers automatically.
        - Managed Add-ons: You upgrade kube-proxy, CoreDNS, and VPC CNI
        - Worker Node Pool Rollout - In our case, We update the AMI and this inturn triggers Karpenter to swap old nodes with new nodes.
8. How would you upgrade EKS worker nodes without disrupting workloads that have strict availability requirements?
    - Enforce Workload Resilience
        - Pod Disruption Budgets - Define minAvailable or maxUnavailable in a PDB. The EKS/Kubernetes eviction API will block node draining if evicting a pod violates these constraints.
        - Topology Spread Constraints / Anti-Affinity: Ensure replicas of the same application are distributed across different nodes and Availability Zones (AZs) so that losing one node doesn't take down the entire service.
        - Readiness Probes: Configure accurate readiness probes so the routing layer (Service/Ingress) only sends traffic to new pods once they are completely initialized.
        - Graceful Termination: Pods receive a SIGTERM signal. Your application must catch this signal, cease accepting new requests, finish in-flight requests, and exit cleanly within the terminationGracePeriodSeconds window.
    - Choose the Upgrade Method
        - Managed Node Groups (MNG) - Native Rolling Upgrade, EKS automatically creates a new EC2 instance, waits for it to become ready, drains an old node, and terminates it. You can configure maxUnavailable percentage to control the speed.
        - Karpenter - Drift Management - Update the AMI ID or version in your EC2NodeClass. Karpenter detects the "drift," automatically provisions upgraded nodes, gracefully drains the old ones, and terminates them.
        - Self-Managed Node Groups - Blue/Green Node Groups - Deploy a completely new Auto Scaling Group (ASG) with the new AMI. Cordon and drain the old ASG systematically, then terminate it once empty.
9. What strategies would you use to prevent a bad deployment from affecting every cluster simultaneously?
    - Progressive Delivery & Ring Deployments: Instead of deploying to all clusters at once, group your clusters into logical deployment rings (or phases) based on risk tolerance.
        - In a multi region setup, Treat each geographic region as an independent failure domain. Never deploy to paired regions (e.g., US-East and US-West) at the same time. Global Load Balancing - Use a global DNS or load balancer (like Cloudflare or AWS Route 53) that can automatically detect an unhealthy cluster and instantly route user traffic away from it to healthy clusters.
    - Automated Health Checks & Baking Periods: Do not move from one cluster ring to the next without validation.
10. How do PodDisruptionBudgets, readiness probes, liveness probes, and startup probes work together during maintenance?
    - PodDisruptionBudgets (PDBs) Limit Concurrent Voluntary Evictions - The PDB acts as a gatekeeper. It checks how many replicas must remain available (via minAvailable or maxUnavailable). Evictions are paused if taking a pod down would violate the budget. The system evicts pods in waves, waiting for new pods on other nodes to become fully healthy before evicting the next batch.
    - Startup Probes Handle Initial Initialization - When a new pod is scheduled on a different node to replace an evicted one, the startup probe runs first. `It disables all other probes` (liveness and readiness) while the application performs heavy lifting like loading caches or running database migrations. If the application takes 2 minutes to boot, the startup probe keeps the pod from being killed prematurely by an aggressive liveness probe.
    - Readiness Probes Signal When It Is Safe to Proceed - Once the startup probe passes, the readiness probe begins checking if the application can accept traffic. Traffic control: The pod is only added to the Service endpoints (and receives traffic) after the readiness probe succeeds. The PDB counts a pod as "Available" only when its readiness probe passes. Once the new pod is ready, the PDB budget is restored, giving the node-draining process the green light to evict the next pod.
    - Liveness Probes Protect Against Deadlocks - The liveness probe runs continuously after the startup probe passes to ensure the container hasn't deadlocked or crashed internally.Maintenance impact: If a newly scheduled pod experiences a deadlock during its high-load startup phase on the new node, the liveness probe restarts it. Because a restarting pod fails its readiness probe, the PDB will block further node evictions until the application stabilizes.
11. How would you design multi-tenancy in Kubernetes so teams can deploy independently without compromising platform security?
    - To design secure, independent multi-tenancy in Kubernetes, you should implement a  Hard Namespace Isolation strategy. For most enterprises, Namespace-as-a-Service (NaaS) augmented with network policies, RBAC, and admission controllers strikes the best balance between team autonomy and platform security.
    - Control Plane & Access Isolation
        - Strict RBAC (Role-Based Access Control) - Bind teams to a specific namespace using a RoleBinding instead of a ClusterRoleBinding. Team members should only have access to their designated namespaces.
        - Namespace Provisioning Automation - Implement GitOps (e.g., ArgoCD or Flux) or a portal (e.g., Capsule, Loft) to provision namespaces. Teams must never have permission to create their own namespaces directly.
    - Network & Traffic Isolation
        - Default-Deny Network Policies
        - Explicit Whitelisting
        - Service Mesh (Optional): For advanced requirements, use a service mesh (like Istio or Linkerd) to enforce mutual TLS (mTLS) and identity-based authorization between tenant workloads
    -  Compute & Resource Guardrails
        - Prevent a single team from accidentally or maliciously exhausting cluster resources (a "noisy neighbor" scenario)
        - ResourceQuotas / LimitRanges / Node Affinities & Taints
    - Runtime & Container Security
        - Admission Controllers: Deploy policy engines like Kyverno or OPA/Gatekeeper. Use them to block dangerous configurations before they reach the cluster.
    - Shared Infrastructure & Ingress
        - Ingress Sharding: Instead of a single shared Ingress Controller, deploy dedicated Ingress Controllers per tenant or group of tenants. This prevents one team’s misconfigured routing rule from breaking external access for other teams.
        - External DNS: Use tools like ExternalDNS coupled with a centralized, automated certificate provider (like cert-manager) so teams can manage their own domains safely without touching core DNS infrastructure.