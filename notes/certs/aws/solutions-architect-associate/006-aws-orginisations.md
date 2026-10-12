# AWS Organisations

Accounts types in an AWS Organisation:
1. Managemnt or Master Account
2. Single Payer Account
3. Member Accounts

An AWS Org can only has 1 management or master account and 0 or more member accounts. Member accounts can either be created in an aaws org or existing aws accounts can be invited to become an member account in an aws org.
AWS Org offers consolidated billing, reservations and volume discounts.

`Identity Account` - AWS Org best practice is to have a single aws account for all of its Identities. This can be IAM Users or Identity Federated Accounts. This allows identities to `Role Switch` from this account to other member account - (Centralized IAM).

`Service Control Policies - SCPs`
1. SCPs are Account Permission Boundaries. They limit what the aws account is allowed to do (including the root user). 
2. Fact: You cannot affect access to the root user directly, but by limiting access to the account by SCP we can indirectly limit root user access.
3. SCP's don't grant any permissions. The best practice is to follow the DENY list policy. By default there is an Implicit DENY of all resources and aws adds Full Access AWS Policy which negates the implicit deny. So as an Ord Admin, we can add an SCP with Explicit DENY's to limit access in aws accounts.
4. SCP's can be attached at the org root or org units or individual aws accounts.
5. A SCP does not limit or control the management account.

`Cloudwatch Logs`
1. Store, Monitor & Acess Log Data
2. AWS Integrations - EC2, VPC Flow Logs, Lambda, CloudTrail, R53, etc.
3. Can generate custom metrics based on the logs.
4. `Log Groups` consists of multiple log streams, metric filters, log retention policies.

`Cloud Trail`
1. Logs API Calls / Account Activities as a Cloud Trail Event.
2. 90 Days retention by default in the event history.
3. Enabled by default - No cost for 90 day history.
4. Cloud Trail events can be stored in Cloudwatch Logs and S3 can be used for long term storage beyond the defauls 90 days history.
5. IAM , STS, CloudFront are Global Services, So it generates Global Service Events. These events are logged to us-east-1 and a trail needs to be enabled to capture these events.
6. Cloud Trail os not real-time and there is some delay.

Cloud Trail Event Types
1. Management Events - These events are management operations like creation of an EC2 Instance or deletion of an EC2 Instance. By default cloudtrail stores Management events only, Data events can be enabled at an additional cost.
2. Data Events - These events are data operations like access an S3 object or lamnda invocations. By default, Data events are not enabled since it could be too much of operations to handle and store.
3. Insights Events.


`AWS Control Tower`
1. Quick and Easy setupo of Multi Account Environments.
2. Orchestrates other AWS Service to provide this functionality.
3. AWS Control Tower provides - Organizations, IAM Identiry Center (SSO), Cloud Formation, AWS Config and more.
4. `Landing Zone` - In a Multi Account Env, landing zone provides SSO / ID Federation, Centralised LOgging and Auditing.
5. `Guard Rails` - Detect / Mandate rules and standards accross all accounts in a landing zone.
6. `Account Factory` - Automated and Standardises new account creation.
7. `Dashboard` - Single page overview of the entire env.

