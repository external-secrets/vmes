# Making vmes assume a role

You probably don't want to insert credentials into your machine that will have big permissions across your aws account. To be able to write a very basic credential to your machine, but still be able to fetch secrets, you will need to make vmes assume a role. Vmes supports this and this doc explains an example workflow.

## Create a service account with assume role permissions

Create a service account that can perfomr sts:AssumeRole action. You will also probably want to limit to which resources it can assume a role (for example  'arn:aws:iam::account-id:role/Test*' to only allow this service account to assume roles starting with 'Test'). 

## Create a new policy that only has secret read actions

Name it `SecretReader` and add following actions:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "secretsmanager:GetResourcePolicy",
                "secretsmanager:GetSecretValue",
                "secretsmanager:DescribeSecret",
                "secretsmanager:ListSecretVersionIds"
            ],
            "Resource": [
                "arn:aws:secretsmanager:us-west-2:111122223333:secret:Project-Env-*"
            ]
        }
    ]
}
```

This policy is also limiting secrets starting with `Project-Env-` prefix, adapt accordingly.

## Use this Service Account's credentials to authenticate

Simply export, as in previous instructions. 

```
export AWS_ACCESS_KEY_ID="******"
export AWS_SECRET_ACCESS_KEY="******"
```

## Edit the ss.yml file to tell vmes which role to assume

Vmes will use AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY to authenticate, and later assume the role that you set on ss.yml.

```
[...]
spec:
  provider:
    aws:
      role: arn:aws:iam::479664765532:role/vmes-role ## This is the role that vmes will assume
      service: SecretsManager
[...]
```

Just run vmes as before, passing --config-path and --public-key-path. It will authenticate and assume the role. Only available after version 0.0.3