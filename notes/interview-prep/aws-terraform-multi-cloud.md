# AWS, Terraform provisioning, multi-cloud platform engineering, cost optimization.

Walk me through how you identified and delivered the $365K in annual infrastructure savings.



How did you calculate the savings baseline, account for workload growth, and verify that the savings were sustainable?



Explain how Karpenter works. How does its scheduling and provisioning model differ from Cluster Autoscaler?



What happens when Karpenter cannot provision a suitable node for a pending pod?



How would you prevent Karpenter from creating too many expensive nodes during a traffic spike?



How would you choose between Spot and On-Demand instances for a production EKS fleet?



What are the risks of migrating workloads from x86 AMD instances to ARM-based Graviton instances?



How would you validate application compatibility, performance, and cost benefits during an AMD-to-ARM migration?



Explain Kubernetes bin-packing. How can aggressive bin-packing increase the risk of correlated failures?



How would you design Availability Zone-aware traffic routing and capacity placement?



How do you structure reusable Terraform modules for networking, IAM, EKS, and cluster add-ons?



How do you detect and remediate Terraform drift in production without accidentally destroying critical resources?



How would you manage Terraform state, state locking, module versioning, and concurrent changes across multiple teams?



How would you implement policy enforcement and security checks in a Terraform CI/CD pipeline?



Your company wants to reduce cloud costs by another 20%, but application teams refuse to accept lower availability. How would you proceed?

How would you design a Terraform module that provisions a production-ready EKS cluster with networking, IAM, and add-ons?



How do you avoid creating overly generic Terraform modules that are difficult to maintain?



How would you version and roll out breaking changes to shared infrastructure modules?



How would you structure Terraform repositories for multiple environments, accounts, and regions?



How do you prevent secrets from being exposed through Terraform state or CI/CD logs?



You have AWS and GCP platforms with similar application requirements. What should be standardized, and what should remain cloud-specific?



How would you migrate an application from EKS to GKE while minimizing downtime and avoiding vendor lock-in?



How would you design self-service infrastructure provisioning while preserving platform governance?



How would you detect and recover from a partially failed infrastructure deployment?



How do you decide whether to automate a manual process, build an internal platform service, or use an existing managed service?