# Istio, Envoy and API gateway migration

1. Why did you replace the standalone NGINX gateway with an Istio-native gateway? What problems were you solving?

Previously, external ingress traffic routed through our AWS Application Load Balancers (ALBs) had to cross two separate gateway tiers: the Envoy-based Istio Ingress Gateway and a standalone NGINX API Gateway. The NGINX layer was heavily utilized for host/path routing, regex-based URL rewrites with capture groups, header manipulations, and authorization token validation.
By eliminating the NGINX tier and absorbing its functionality directly into the mesh-native Istio Gateway, we resolved several `critical systemic problems across latency, observability, and regional resilience`.

1. High Latency and Architectural Redundancy (Multi-Hop Ingress)
    - The dual-gateway layout introduced unnecessary processing overhead. A standard request was forced to hop through the AWS ALB, the Envoy ingress gateway, and then the NGINX API Gateway before finally routing to the target application sidecar.
    - Benchmarks demonstrate that Envoy's modern architecture delivers significantly higher throughput and performance compared to NGINX, making the standalone NGINX layer an inefficient bottleneck.
2. Fragmented Observability and Elevated MTTR (Mean Time to Recovery)
    - Debugging transient errors (such as 5xx/502 responses) was highly complex. Incident commanders could not easily isolate whether a fault originated at the perimeter Envoy gateway, the NGINX gateway, the NGINX sidecar, or the upstream application sidecar.
    - Although telemetry data showed that the majority of historic 502 errors occurred within the NGINX container itself, the multi-tiered dependency forced engineering teams to inspect up to four distinct network components during active incidents, severely hurting our MTTR.
3. Inefficient Cross-Region Traffic Routing
    - The legacy NGINX setup lacked multi-region awareness. Consequently, cross-region traffic entering the secondary region ALB was suboptimally hairpinned back to the primary region's NGINX gateway before distribution, adding up to three cross-region hops for secondary-region workloads.
4. Lack of Automated Regional Failover
    - The routing dependency on the primary region's NGINX infrastructure introduced a single point of failure (SPOF) for global disaster recovery.
    - A primary region outage required manual, high-risk interventions to manipulate Route 53 DNS records and modify Istio routing weights, preventing clean, automated failovers.
5. Technical Debt and Maintenance Overhead
    - Managing two distinct gateway technologies doubled configuration complexity across our platform's Helm charts. Additionally, the open-source community support and feature velocity for traditional NGINX variants have stagnated compared to the cloud-native Envoy ecosystem.

We eliminated the standalone NGINX layer entirely and consolidated all core capabilities—including external custom authorization policies, dynamic header manipulations, and advanced routing workflows—directly into the mesh-native Istio Ingress Gateway. This structural reduction leaves us with a single, highly performant, unified control plane that optimizes traffic paths, unifies telemetry, and unlocks native cross-region mesh routing capabilities.



1. Explain the difference between an API gateway, an ingress controller, a service mesh, and an east-west gateway.
    - `Ingress Controller` - An Ingress Controller is a specialized load balancer designed specifically to get traffic into a Kubernetes cluster. It sits at the edge of the cluster, watches the Kubernetes API for Ingress or Gateway resources, and translates them into routing rules.
        - Traffic Direction: North-South(External client -> Internal cluster service).
        - Primary Scope: Handles the transition from public/private corporate networks to private container networks. It takes care of TLS termination, host-based routing (e.g., https://company.com), and path-based routing (e.g., /api/v1).
    - `API Gateway` - An API Gateway is an application-aware proxy layer. While an Ingress Controller cares about networks and pods, an API Gateway cares about business logic, APIs, and developers. It sits in front of business applications to shield them from client-facing complexities.
        - Traffic Direction: North-South (Client -> API endpoints).
        - Primary Scope: High-level application functionalities. This includes rate limiting (per user/API key), user authentication/authorization (OAuth2, JWT validation), request/response body mutation, etc.
    - `Service Mesh` - A Service Mesh manages communication inside the cluster between your microservices once traffic has already cleared the front door. In Istio’s case, this is achieved by injecting a sidecar proxy (Envoy) next to every single application container, controlled by a centralized control plane.
        - Traffic Direction: East-West (Service-to-Service within the mesh).
        - Primary Scope: Microservice reliability, security, and observability. It enforces mutual TLS (mTLS) for all internal communication, handles circuit breaking, retries, and outlier detection, and generates distributed tracing maps without requiring changes to application code.
    - `East-West Gateway` - An East-West Gateway is a specific architectural component used to route traffic between different service meshes or different clusters/regions. Unlike a standard Ingress gateway that exposes apps to the open internet, an East-West gateway is dedicated strictly to secure, mesh-to-mesh infrastructure routing.
        - Traffic Direction: East-West Across Boundaries (Cluster A Service -> Cluster B Service).
        - Primary Scope: Cross-cluster replication, disaster recovery, and multi-region failover. It acts as an entry/exit point that preserves internal mesh identities (like SPIFFE/SPIRE) and keeps traffic strictly encrypted via mTLS as it transits over the WAN between environments.


1. Walk me through the request path before and after your gateway migration.

1. `Before Migration`: 
    - `Public Ingress`: Internet -> AWS Global Accelerator -> ALB -> Envoy Gateway -> **Nginx API Gateway** -> Application Pod
    - `Private Ingress`: Internet -> ALB -> Envoy Gateway -> **Nginx API Gateway** -> Application Pod
    - `Secondary Region Ingress`: Internet -> AWS Global Accelerator -> ALB -> Envoy Gateway -> **Nginx API Gateway (Secondary Region) -> Nginx API Gateway (Primary Region)** -> Application Pod (Primary Region)
2. `After Migration`: 
    - `Public Ingress`: Internet -> AWS Global Accelerator -> ALB -> Envoy Gateway -> Application Pod
    - `Private Ingress`: Internet -> ALB -> Envoy Gateway -> Application Pod
    - `Secondary Region Ingress`: -> Internet -> AWS Global Accelerator -> ALB -> Envoy Gateway -> Application Pod (Secondary Region)

3. How did you achieve feature parity between NGINX and Istio? Which features were the hardest to reproduce?

We achieved feature parity with a help if Istio Virtual Services Per Gateway and Per Application Routes, Destination Rules Per Applicatin, Authorization Policy Per Application Routes.

Feature which were challenging to implement:
- `Maximum request body size`
    - NGINX handles limits instantly using the native client_max_body_size directive or (nginx.ingress.kubernetes.io/proxy-body-size: " annotation), easily overridden on individual location paths. It caps RAM usage to a tiny footprint (via client_body_buffer_size) and transparently streams larger uploads onto local disk storage. It reads the Content-Length header immediately, closing connections and returning a 413 before wasting system memory or disk I/O.
    - By default, Istio (via its underlying Envoy proxies) streams HTTP requests directly to upstream applications without checking or limiting the total body size. Because native Istio resources like VirtualService or Gateway do not provide a built-in payload size parameter, you must use an EnvoyFilter to explicitly inject Envoy's envoy.filters.http.buffer filter
        - Ingress Gateway Level (Edge Protection): Applying the filter at the istio-system namespace blocks oversized requests before they ever enter your service mesh, protecting downstream compute and memory. You can deploy an EnvoyFilter targeted at the istio-system namespace with context: GATEWAY and a max_request_bytes configuration (e.g., 10485760 for 10MB).
        -  Route Overrides Directly at the Ingress Gateway Edge - Service/Sidecar Level (Fine-Grained Protection): If you want to apply specific payload size limits to particular backend microservices, apply the EnvoyFilter in the application namespace to grant exceptions (e.g., 20MB) to specific backend apps, you must configure Per-Route Overrides directly on the Ingress Gateway.
        - Completely Disabling the Limit on a Route: If you have a dedicated large-upload endpoint, you can bypass the buffer restriction entirely for that specific route path.
            - Handle Limits Safely at the Application Level: Once Envoy streams the raw TCP packets directly into your application, you must protect your pod from being overwhelmed. Rather than reading the entire payload into the application container's RAM, read it as a stream and enforce the size limit on the fly.
    - Challenges:
        - Memory Footprint: The buffer filter forces Envoy to keep the entire request body in memory before forwarding it. If you set max_request_bytes to 50MB and handle 100 concurrent requests, your sidecar container can consume up to an extra 5GB of RAM (max_request_bytes * concurrent_requests). Be sure to adjust your Kubernetes resources.limits.memory on targeted pods accordingly!
- Authourization Header Validation 
    - A shift from the current current In-Pod Sidecar / Local Gateway Authorization Pattern to Centralized Service-Mesh Sidecar Authorization Pattern.
    - Enforement Point: We moved away from a Local NGINX wrapper container inside the exact same pod boundary to Istio Sidecar Proxy (Envoy) attached to your application container.
    - This pattern allowed our Identity team to autoscale and update policies completely independently.

3. How did you migrate production traffic without dropping requests?

To migrate production traffic across public and private edges in two regions without dropping requests,` a phased, Blue/Green DNS shift was executed during a low-traffic maintenance window`.

`Blue Green Deployment of the Gateways`:

- `Parallel Deployment`: Co-located new gateway clusters directly alongside the active legacy gateways across all public and private edges in both operating regions.
- `Upstream Alignment`: Coordinated with application teams to redeploy and configure backends so they could seamlessly accept and process traffic identically from both gateway generations.
- `TTL Optimization`: Lowered all DNS Time-To-Live (TTL) timeouts to 60 seconds prior to the cutover window to guarantee rapid client-side cache expiration.
- `Off-Peak Timing`: Scheduled the actual switch during a designated low-traffic maintenance window to minimize the blast radius of any transient issues.
- `Phased DNS Execution`: Initiated a structured, multi-stage DNS routing shift to transition traffic gradually rather than all at once:
	- Private Host-Based Endpoints moved first.
	- Private Path-Based Routing moved second.
    - Public Edge Traffic followed the same routing sequence once private layers stabilized.
- `Regional Isolation`: Executed the entire phased DNS sequence within the secondary region first to validate the configuration in production before repeating the exact same steps in the primary region.

`Challenges`:

- Stale TCP Connections: Active clients using HTTP Keep-Alive will continue sending data over existing, open TCP pipes to the old gateway, completely ignoring the 60-second DNS TTL update.
    - Fix: Extended Co-Existence Window: Keep the old gateway infrastructure running for a defined soak period (e.g., 24–48 hours) to monitor connection decay, and use access logs to identify and manually remediate stubborn clients.
- Non-Compliant DNS Caching: Certain client application stacks (such as legacy Java JVM configurations or mobile devices) cache DNS records indefinitely at the application level, completely ignoring your 60s TTL.
    - Fix: Same as above.
- During the rolling DNS transition:  A client's initial request (like a login) might hit the old gateway, while a subsequent stateful request hits the new gateway.
    - the user experience will break mid-session, resulting in scattered errors across your system during the DNS propagation window.
        - Session Invalidation (401 Unauthorized) 
            - • The Failure: A client authenticates via the old gateway. The DNS record expires 60 seconds later, and their next API call lands on the new gateway. Because the new gateway has no record of that active session or cryptographic handshake context, it rejects the request outright, typically returning a 413 (if fine-grained protection state mismatches) or a 401 Unauthorized / 403 Forbidden.
            - Standard HTTP clients do not treat a 401 or 403 as a network glitch. They treat it as a hard authentication failure. The client app will stop retrying and will likely force the user back to a login screen.
        - The TCP / HTTP/2 Connection Breakage
            - A WebSocket connection pinned to an application pod through the old gateway cannot be smoothly handed over. If the old gateway is forcefully terminated or drops its connection without a shared state, the socket snaps.
            - A dropped TCP connection or an HTTP 5xx gateway error will trigger standard network retry logic in almost all modern HTTP clients. Because the connection broke, the client will re-resolve DNS, see the new gateway IP, and establish a fresh session on the new infrastructure.

- Cross-Region Hairpinning: Triggering the secondary region's DNS switch can cause temporary routing imbalances or hair-pinning across regions if global traffic profiles alter latency
    - Fix: `Locality-Aware Traffic Policie`s: Enforce strict regional boundaries at your Global Traffic Manager (GTM) level, and ensure application target groups in the secondary region are fully warmed and scaled.


4. What rollback strategy did you design, and what conditions would trigger a rollback?

The Rollback Strategy:
- Kept the old gateway clusters fully powered on, scaled, and configured during the entire migration window and for a 48-hour "soak period" afterward. If a rollback is triggered, execute an automated script or pipeline to immediately point your DNS record or global accelerator routing configurations back to the old gateway load balancers.

Critical Rollback Triggers:
- L7 Telemetry: Edge Error Spikes - Global HTTP 5xx (Gateway Error) or HTTP 4xx (Auth/State Mismatch) error rates 
- Network Telemetry: Connection Drops & Timeouts & Latency -  A sudden drop in successful TCP handshakes or an exponential spike in upstream connection timeouts (HTTP 504) on the new gateway.
    - Envoy Connection DownStream Metrics.
- Blast Radius Guardrails: Resource Exhaustion.
    - Track the memory / cpu usage of the new gateways.
    - Ensure HPA scales up and down accordingly.
- ALB Access Logs
    - Tracking TLS Negotiation Spikes
    - Tracking Dropouts After Connection

In our case, We did not rollback, instead we waited for about for about 2 hours after migration and ensure that the error levels we same as before before calling it a successfull migration.

We left the old gateways for around the weekend and then eventually took them down after no new connections or active connections through them.

Minor data loss occurred during the transition, and we subsequently replayed the missing documents to achieve full recovery.


5. How do Istio's Gateway, VirtualService, and DestinationRule resources work together?

In Istio, the Gateway, VirtualService, and DestinationRule Custom Resource Definitions (CRDs) work together as a pipeline to manage how traffic enters the cluster, where it gets directed, and how it behaves when it arrives.

The Request Flow Pipeline: `External Request -> Gateway -> Virtual Service -> Destination Rule -> Application Pods`

- Gateway: The Entry Point - Handles edge ingress or egress configuration. It acts strictly at the layer 4-6 level. It specifies which ports to expose, which protocols to accept (HTTP, HTTPS, gRPC, TCP), and which hostnames are allowed to pass through
- VirtualService: The Traffic Director - Defines the actual routing rules and matches incoming traffic using Layer 7 logic. It links directly to a Gateway (for external traffic) or binds to the internal mesh (for service-to-service communication). Virtual capabilities include - Layer 7 Request Matching (URI Paths, HTTP Headers, HTTP Methods, Query Parameters), URI and Authority Rewriting, HTTP Redirects & Direct Responses, Header Manipulation, Timeouts & Retries, CORS Policies.
- DestinationRule: The Policy Enforcer - Defines policies that apply after routing decisions have been made by the VirtualService. Traffic Policies - Workload Subsets, Load Balancing Policies (RoundRobin, Least Connection), Consistent Hashing & Sticky Sessions, Locality & Cross-Region Routing (Weighted and Failover), Connection Management & Circuit Breaking (TCP Connection Pools, HTTP/1.1 and HTTP/2 Tuning), Outlier Detection (Ejects faulty pods), Transport Layer Security (TLS) Settings (TLS Modes: Mutual TLS, SNI & SAN Verification), 

6. Explain how Envoy handles routing, retries, circuit breaking, connection pooling, and timeouts. 

TODO

7. How can misconfigured retries amplify an outage or cause duplicate transactions?

- Retry Storm: If an application normally receives 1,000 requests per second and has a global retry limit of 3 without backoff, a brief 503 error spike from a database can suddenly balloon inbound network traffic to 4,000 requests per second. The backend, which was already resource-starved, is completely crushed under the artificial load spike.

- How we mitigated: 
    - Implement Strict Concurrency Limits (Retry Budgets)
    - In most cases, we let the application code handle retries on their own instead of the mesh.


8. How would you troubleshoot intermittent 503 errors that appear only between certain services in the mesh?

- To pinpoint exactly which services are throwing intermittent 503 errors, you can run a targeted PromQL query using the istio_requests_total metric. Once the PromQL query reveals the specific source_workload and destination_workload causing the errors, jump straight into the application and sidecar logs for those specific pods (application container logs or istio proxy logs).


9. Explain Istio mTLS. What is the difference between STRICT, PERMISSIVE, and DISABLE modes?

Istio Mutual TLS (mTLS) is a core security feature that automatically encrypts and authenticates communication between microservices within a service mesh.

Unlike standard TLS (where only the client verifies the server's identity), mTLS requires both the client and the server to verify each other's certificates before establishing a connection. Istio handles this transparently using sidecar proxies (Envoy) without requiring changes to your application code.

The peer authentication policy in Istio controls how strictly these mTLS requirements are enforced on the receiving (server) side through three distinct modes: STRICT (Only Istio mTLS), PERMISSIVE (Plain text and Istio mTLS), and DISABLE (Only plain text).

By default, modern versions of Istio implicitly enforce PERMISSIVE mode across the entire service mesh - This means that out of the box—even if you haven't explicitly created or applied any PeerAuthentication custom resources—Istio will automatically upgrade traffic to mTLS between sidecars wherever possible, while still accepting plain text traffic from outside the mesh.


10. How would you rotate Istio root and intermediate certificate authorities across multiple clusters without an outage?

- The 4-Phase Zero-Downtime Rotation Plan: `Phase 1: Dual Trust Bundle (Trust Root A + Root B) -> Phase 2: Switch Intermediates (Sign with Intermediate B) -> Phase 3: Wait for Workload Expiry -> Phase 4: Clean Up (Trust Root B Only)`

- Phase 1: Establish a Combined Root Trust Bundle : Concatenate the new and old roots (cat new-root-cert.pem old-root-cert.pem > combined-root-cert.pem) and update the cacerts secret across all clusters with this combined root while retaining the old intermediate CA. Restart istiod (kubectl rollout restart deployment/istiod -n istio-system) to distribute the expanded bundle via the istio-ca-root-cert ConfigMap
- Phase 2: Pivot Intermediate CAs to the New Root: Sequentially update each cluster with a new intermediate CA signed by Root B and the combined root certificate, followed by restarting istiod. This allows existing and new workloads to communicate seamlessly during the transition.
- Phase 3: Wait for Workload Expiration: Allow enough time (typically matching your 24-hour default certificate TTL plus a buffer) for all workload proxies to automatically request new certificates via the SDS API without manual restarts.
- Phase 4: Evict the Old Root Certificate: Once verification confirms no proxies use old leaf certificates, update the cacerts secrets across all clusters to use only the new root and execute a final istiod rollout
    - Metric to watch: IstioWorkloadCertRotationFailing - `(rate(citadel_server_csr_count[5m]) - rate(citadel_server_success_cert_issuance_count[5m])) > 0`


11. What would happen if one intermediate certificate expired while workloads were still using it?

- The Immediate Impact: Total Data Plane Outage: The moment the intermediate CA certificate crosses its expiration timestamp, all new cross-workload mTLS connections will instantly fail - When Client Pod A attempts to establish an mTLS session with Server Pod B, they exchange their certificate chains. Envoy proxies validate the incoming leaf certificate by parsing the entire chain up to the root CA. When Envoy detects that the signing intermediate CA in that chain is expired, cryptographic validation fails
- The Control Plane Impact: Brick-Walled Leaf Renewals: Because Istio leaf certificates only last 24 hours and attempt renewal every 12 hours, the background istio-agent on your pods will continue to shoot CSR requests to istiod. istiod cannot sign certificates using an expired intermediate private key/cert pair. Leaf renewals will lock up mesh-wide, causing citadel_server_csr_sign_err_count to spike massively across your control plane metrics

Metrics to track:

- IstioRootCAExpiringSoon: Alert if the Root CA is expiring within 30 days: `(citadel_server_root_cert_expiry_timestamp - time()) < 30 * 24 * 3600`
- IstioIntermediateCAExpiringCritical: Critical Alert if the Intermediate CA is expiring within 7 days`(pilot_cert_expiry_seconds - time()) < 7 * 24 * 3600`


How would you safely upgrade Istio across hundreds of services with different owners and release schedules?

This is achieved by deploying Canary Control Planes in tandem with abstraction mechanisms called Revision Tags.

The Multi-Phased Upgrade Execution Plan:
- Phase 1: Deploy the Parallel Control Plane (No Traffic Impact) 
- Phase 2: Establish the Canary Tag Hook: istioctl x revision tag set istio-1-25-0 --revision 1-25-0
- Phase 3: Gradual Decentralized Workload Migration: Test sample workloads with the new revision tags `kubectl label namespace sample-workloads istio.io/rev=istio-1-25-0 --overwrite`
- Phase 4: Overwrite the Default Production Tag: `istioctl x revision tag set prod-stable --revision 1-25-0` - Any application namespace utilizing prod-stable will gracefully pick up the new proxy control track on its very next scheduled deployment rollout or as and when new pods are rolled.
- Phase 5: Deprecate and Purge - Once logging metrics show zero active connections hitting the old control plane, remove it cleanly. `count(istio_build{component="proxy", version="1.25.0"})`

12. When would you choose ambient mesh over sidecar-based Istio, and what trade-offs would influence that decision?

The fundamental difference lies in how proxying responsibilities are distributed across your cluster:

- Sidecar Mode: Embeds a single, heavy Envoy proxy instance right inside every application pod. This proxy handles both Layer 4 (mTLS/TCP) and Layer 7 (HTTP routing, headers, and policies) processing locally inside the same network namespace.
- Ambient Mode: Splits the data plane into two decoupled infrastructure layers:
	1. ztunnel (Layer 4): A lightweight, secure tunnel agent running as a DaemonSet on each node. It transparently handles mTLS and TCP routing for all local pods.
	2. Waypoint Proxies (Layer 7): On-demand Envoy-based proxies deployed independently at the namespace or service level. They are strictly spun up only if your applications require intensive L7 traffic routing or authorization

When to choose Ambient / Sidecar Mode ? TODO


