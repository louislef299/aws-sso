+++
date = '2025-09-05T15:55:03-05:00'
draft = false
title = 'Plugin Architecture'
+++

The `aws-sso` tool currently is *attempting* to use a plugin architecture for
each step of the authentication process. This gives developers(mostly just me)
the flexibility to try out newer features in isolation without breaking the
current tool. As of v1.6, there are three main plugins:

1. OIDC Plugin (Always Enabled - Core Auth)
    - Handles AWS SSO authentication
    - Manages token retrieval

2. EKS Plugin (Optional - Can be disabled)
    - Configures kubectl credentials
    - Creates kubeconfig entries
    - Can be disabled with `eks.disableEKSLogin` config value or
      `--disableEKSLogin`

3. ECR Plugin (Optional - Can be disabled)  
    - Configures Docker credentials
    - Allows container image pulls
    - Can be disabled with `ecr.disableECRLogin` config value or
      `--disableECRLogin`

```bash
                        +-------------------+
                        |                   |
                        |    aws-sso CLI    |
                        |                   |
                        +-------------------+
                                |
                                v
        +-----------------------------------------------------+
        |                                                     |
        |                  Login Command                      |
        |                                                     |
        +-----------------------------------------------------+
                                |
                                v
    +-------------------------------------------------------------+
    |                       Plugin Registry                       |
    |      (dlogin package - manages plugin registration)         |
    +-------------------------------------------------------------+
                |               |                 |
                v               v                 v
+------------------+    +------------------+    +------------------+
|                  |    |                  |    |                  |
|   OIDC Plugin    |    |    EKS Plugin    |    |    ECR Plugin    |
|   (Required)     |    |    (Optional)    |    |    (Optional)    |
|                  |    |                  |    |                  |
+------------------+    +------------------+    +------------------+
        |                       |                       |
        v                       v                       v
+------------------+    +------------------+    +------------------+
|                  |    |                  |    |                  |
| AWS SSO/OIDC     |    | Kubernetes       |    | Docker Registry  |
| Authentication   |    | Authentication   |    | Authentication   |
|                  |    |                  |    |                  |
+------------------+    +------------------+    +------------------+
        |
        v
+------------------+
|                  |
|   AWS Console    |
|   Access         |
|                  |
+------------------+
```
