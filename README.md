# Azure Security Hub

Azure Security Hub is a CLI tool designed to inspect and monitor the security posture of Azure resources.

<br/>

### CIS Microsoft Azure Foundations Benchmark v3.0.0 (09-05-2024)

**2 Identity**

- **2.1 Security Defaults (Per-User MFA)**
  - (Manual) 2.1.1 Ensure Security Defaults is enabled on Microsoft Entra ID
  - (Manual) 2.1.2 Ensure that 'Multi-Factor Auth Status' is 'Enabled' for all Privileged Users
  - (Manual) 2.1.3 Ensure that 'Multi-Factor Auth Status' is 'Enabled' for all Non-Privileged Users
  - (Manual) 2.1.4 Ensure that 'Allow users to remember multi-factor authentication on devices they trust' is Disabled
- **2.2 Conditional Access**
  - (Manual) 2.2.1 Ensure Trusted Locations Are Defined
  - (Manual) 2.2.2 Ensure that an exclusionary Geographic Access Policy is considered
  - (Manual) 2.2.3 Ensure that an exclusionary Device code flow policy is considered
  - (Manual) 2.2.4 Ensure that A Multi-factor Authentication Policy Exists for Administrative Groups
  - (Manual) 2.2.5 Ensure that A Multi-factor Authentication Policy Exists for All Users
  - (Manual) 2.2.6 Ensure Multi-factor Authentication is Required for Risky Sign-ins
  - (Manual) 2.2.7 Ensure Multi-factor Authentication is Required for Windows Azure Service Management API
  - (Manual) 2.2.8 Ensure Multi-factor Authentication is Required to access Microsoft Admin Portals
- [x] 2.3 Ensure that 'Restrict non-admin users from creating tenants' is set to 'Yes'
- (Manual) 2.4 Ensure Guest Users Are Reviewed on a Regular Basis
- (Manual) 2.5 Ensure That 'Number of methods required to reset' is set to '2'
- (Manual) 2.6 Ensure that account 'Lockout Threshold' is less than or equal to '10'
- (Manual) 2.7 Ensure that account 'Lockout duration in seconds' is greater than or equal to '60'
- (Manual) 2.8 Ensure that a Custom Bad Password List is set to 'Enforce' for your Organization
- (Manual) 2.9 Ensure that 'Number of days before users are asked to re-confirm their authentication information' is not set to '0'
- (Manual) 2.10 Ensure that 'Notify users on password resets?' is set to 'Yes'
- (Manual) 2.11 Ensure That 'Notify all admins when other admins reset their password?' is set to 'Yes'
- (Manual) 2.12 Ensure `User consent for applications` is set to `Do not allow user consent`
- (Manual) 2.13 Ensure 'User consent for applications' Is Set To 'Allow for Verified Applications'
- [x] 2.14 Ensure That 'Users Can Register Application' Is Set to 'No'

**3 Security**

- **3.1 Microsoft Defender for cloud**
  - **3.1.1 Microsoft Cloud Security Posture Management (CSPM)**
    - (Deprecated) 3.1.1.1 Ensure that Auto provisioning of 'Log Analytics agent for Azure VM's is Set to 'On'
    - [ ] 3.1.1.2 Ensure that Microsoft Defender for Cloud Apps integration with Microsoft Defender for Cloud is Selected
  - **3.1.2 Defender Plan: APIs**
  - **3.1.3 Defender Plan: Servers**
    - [ ] 3.1.3.1 Ensure That Microsoft Defender for Servers Is Set to 'On'

**4 Storage Accounts**

- [x] 4.1 Ensure that 'Secure transfer required' is set to 'Enabled'
- [x] 4.2 Ensure that 'Enable Infrastructure Encryption' for Each Storage Account in Azure Storage is Set to 'enabled'

**5 Database Services**

**6 Logging and Monitoring**

**7 Networking**

**8 Virtual Machines**

**9 AppService**

**10 Miscellaneous**
