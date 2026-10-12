# AWS Accounts

```
An AWS Account is a container for Identities and Resources.
```

`Root User` - The email address provided at the time of account creation becomes the root user and it has full control over all of this account and any resources created within it and cannot be restricted.

`Pay As You Go Model` - AWS Accounts can contain the impact of admin errors or exploits by bad actors/ Use seperate accounts for seperate things like dev, test, prod, teams, products and clients.

`Best Practices`
1. Create `AWS Budgets` to avoid any un-intended costs (cost budget).
2. Enable `MFA` for the Root Account.
3. Activate `IAM User and Role access to Billing Information` - This setting alone doesnt grant IAM User and Roles the necessary permissions for AWS Billing console pages. In adition to activating the IAM Access, you must also attach the required IAM Policies to these users and roles.
4. Enable `Billing Perferences` to receive billing alerts and receive PDF invoice by email and receive free tier usage alerts.
5. Create `Account Aliases` to make the account sign-in URL user friendly.

`AWS Support Plans`
1. Basic Support     - Free
2. Developer Support - 29$ Per Month
3. Business Support  - 100$ Per Month

`IAM Basics` - LEAST PREVILAGE ACCESS
1. Control Access - ALLOW or DENY access to Identities in an AWS Account.
2. Global Service and not region specific.
3. 3 Types of Identity Objects - User / Groups / Role.
4. IAM is free (No cost).
5. No direct contrl on external accounts or users.
6. Identity Federationa and MFA.

`IAM Functionalities`
1. Manage Identities - An ID Provider (IDP)
2. Authenticate - Prove who you claim to be.
3. Authorize - Allow or deny access.

`Task that require ROOT User Credentials`
1. Change your account settings.
2. Restore IAM user permissions.
3. Activate IAM access to the billing cost management console.
4. Close your AWS Account.
5. Sign-up for Gov-Cloud.
6. Register as a seller in the RI Marketplace.
7. Configure S3 Bucket to enable MFA.
8. Edit or delete an S3 bucket policy that incluides an invalid VPC id or endpoint.

`Long Term Credentials`
1. UserName and PWD
2. IAM Access Keys - These key don't change automatically and an IAM User can only have 2 sets of access key whether they are active or incative.



