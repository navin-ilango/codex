# Go, API design and engineering automation

Walk me through the architecture of your Kubernetes upgrade tracker, from data collection through notification delivery.



Why did you choose Go for this service? What alternatives did you consider?



How would you design an API that reports the upgrade status of thousands of clusters and their workloads?



How would you model cluster state, version compatibility, upgrade history, and status transitions?



Why did you choose S3 as the backing store? What trade-offs would change if the service needed frequent updates and complex queries?



How would you handle Jira or Slack API rate limits, intermittent failures, and service outages?



How would you make notification processing idempotent so retries do not generate duplicate alerts?



How would you implement authentication, authorization, audit logging, and secrets management for this API?



How would you test a Go service that integrates with AWS, Jira, Slack, and Kubernetes?



How would you instrument the service so you can monitor its own availability and diagnose failures?



What would you change if the tracker became a business-critical service used by every engineering team?



How would you evolve the service from a single deployment to a highly available architecture without overengineering it?

Coding and design exercises to practice
- Write a Go function that retries an HTTP request using exponential backoff, a maximum retry count, and context cancellation.
- Design a worker pool that processes upgrade-status checks concurrently while respecting API rate limits.
- Design a reconciliation loop that compares desired and actual Kubernetes upgrade states.
- Review a sample Go service and identify concurrency bugs, resource leaks, error-handling issues, and missing tests.