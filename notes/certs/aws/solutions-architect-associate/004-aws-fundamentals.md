# AWS Fundamentals

`AWS Global Infrastructure`

`AWS Region` - Geographic Seperation (Isolated Fault Domain), GeoPolicital Seperation (Diff Governance), Location Control (Performance). Ex: us-east-1 (N.Virginia). Each AWS regions consists of a minimum of 3 Isolated and Physically seperate AZs.
`Availability Zone` - Group of logical data centers with independent power, cooling, physical security and is connected via redundant, ultra low latency networks. AWS Customer focused on HA can design their applications to run in multiple AZs to achieve greated Fault Tolerance. AZ's are physically seperated by a meaningful distance.

`VPC Basics`
1. VPC - A Virtual Network inside AWS.
2. A VPC is within 1 aws account and 1 aws region.
3. Private and Isolated unless you decide otherwise.
4. Two Types of VPC: Default VPC & Custom VPC.

Default VPC:
1. CIDR Range: 172.31.0.0/16
2. They always have the same CIDE and 1 subnet per AZ.
3. A default vps is created once per regions when an aws account is created.
4. CIDR range of the subnets cannot overlap with each other.
5. Consists of Internet gateway, NACL, SG.
6. Can be deleted and re-created.
7. Subnets by default assign Public IPv4 Addresses.

`EC2 Basics`
1. IAAS - Provides virtual machines called as Instances.
2. Private service by default - uses vpc networking.
3. AZ Resilient - Instance fails if AZ Fails.
4. Different Instance Sizes and Capabilities.
5. On-Demand Billing per Second.
6. Local on-host storage or EBS.
7. Instance State: Running , Stopped, Terminated.
8. An instance is charged for - CPU, Memory, Network, Storage.
9. When an instance is in Stoppped state, you are only charged for EBS Storage.

Amazon Machine Image - AMI
1. Permissions (Who can use the AMI) - Public: Everyone is allowed, Owner: Implicit allow, Explicit: Specific aws accounts are allowed.
2. Root Volume - The drive that boots the OS.
3. Block Device Mapping.

EC2 Status Checks
1. System Reachability Checks
2. Instance Reachability Checks

`S3 Basics`
1. Global Storage Platform - Regional Based / Resilient
2. Public Service, Unlimited Data and Multi User.
3. Used to store - media files like movies, audios, photos, other files, texts, large data sets.
4. Object based Buckets.

S3 Object
1. Key / Value
2. 0 Bytes -> 5 TB
3. Key is the object name and the value is the data.
4. Other metadata - Version ID, Tags, Access Control.

S3 Bucket
1. Bucket name must be globally unique across all s3 buckets in all aws accounts.
2. Blast Radius = Region
3. Unlimited objects in a Bucket.
4. Infinitely scalable storage
5. Flat Structure -  All objects are stored at the root level.
6. Object never leaves a region unless you configure replication.
7. Bucket Name must he 3 - 63 characters, all lower case, no underscores and it can start with a lowercase letter or number. It can't be IP formatted.
8. Per Account Limits - 100 Soft Limit and 1000 Hard Limit.
9. S3 is an object store and not a files or block based storage.
10. You cannot mount an s3 bucket to a server.
11. Block all Public Access - By default no body has access to the bucket apart from the user who creates it and the root user and IAM admin.

`Cloudformation Basics`
1. Cloudformation Template to create, update, delete resources.
2. Template Components: AWS Template Format Version, Descriptions, Metadata, Paramaters, Mappings, Conditions, Transform, Resources, Outputs.
3. Description - Is not a mandatory fiels, but of AWS Template Format Version field is present, then description fiels must alway be places followed by it. 

`CloudWatch Basics`
1. Collects and Manages operational data.
2. Metrics - AWS Products, On-Prem Servers, Apps (Custom Metrics).
3. Cloudwatch Logs
4. Cloudwatch Events - AWS Services and Schedules.
5. Concepts: Namespace - Logical grouping of metrics, DataPoint - Consists of a timestamp and metric data, Dimensions - seperate datapoint for the different resources or attributes within the same metric, Alarms - OK State or Alarm State. 

`High Availability`
1. AWS aimns to ensure an agreed level of operation performance usually uptime for a higher than norman period with minimal outages.
2. 99.9% (Three 9s) - 8.77 Hours per Year Downtime
3. 99.999 (Five 9s) - 5.26 Minutes per Year Downtime.

`Fault Tolerant - FT` - Is the property that enables a system to continue operatin properly in the event of the failure of some of its components.

`Disaster Recovery - DR` - A set of policies, tools and procedures to enable recovery or continuation of vital technology infra and systems following a natural or human induced disaster.

`R53 Fundamentals`
1. Register Domains
2. Hosted Zones & Managed NameServers
3. Global Service and Globally Resilient.
4. Hosted Zones: Public and Private(Linked to a VPC)
5. RecordSets are stored in a Hosted Zone.

DNS Record Types:
1. NS - Used for DNS Delegation
2. A or AAAA - Maps a hostname to an IP
3. CNAME - Maps a hostname to another hostname.
4. MX - Find an mail server for a domain.
5. TXT - Can add random facts to a recordset. Widely used to prove domain ownership.

Fact: A total of 13 DNS Root Server exists in the world managed by 12 large organications and the DNS Root Zone is manages by IANA.
