# vmes

External Secrets not necessarily in Kubernetes.

This project uses [ESO](https://github.com/external-secrets/external-secrets) as a dependency and basicaly leverages its provider implementations to grab secrets and write them to an env file.

The initial need here is to run this project as a proccess on a VM that has some legacy applications, then we want to use secrets that are in a Secrets Manager. But we don't want to keep calling Secrets Managers from our applications, we want to just read values from env vars or env files.

![image](https://user-images.githubusercontent.com/2432275/148208852-03bf4422-e392-4fb0-86c4-a6166cb5bc67.png)


### Disclaimer ⚠️

This project is not stable or ready to be used. 

# How it works

Since we already implemented and tested multiple provider clients in the ESO project, we are just importing them here. Since we had some dependencies on Kubernetes resources for some of the configurations, we have built a client that just looks for yaml files locally instead of calling Kubernetes at all.

# Getting Started

Copy the example files provided to edit and configure them:

```
cp pkg/configdata/es.yml.example pkg/configdata/es.yml
cp pkg/configdata/ss.yml.example pkg/configdata/ss.yml
```

Open `pkg/configdata/es.yml` and edit the following fields:

- **spec.refreshInterval:** Choose a Time Duration interval used by the operator to fetch new secrets (1m = 1 minute, 1h = 1 hour, etc).
- **spec.target.name:** #TODO: chose the env file where this will end up? (for now writing to /etc/environment - you need to run as root to write to this file)
- **spec.data.secretKey:** The name of the Env Var injected in the machine.
- **spec.data.remoteRef.Key:** The name of the secret in the external provider.

Open `pkg/configdata/ss.yml` and edit the following fields:

- **spec.provider.aws.region:** Choose region where you created a secret.
- **spec.provider.aws.service:** Let us use SecretsManager for this example.
- **spec.provider.aws.auth:** Keep everything here the same for this example.


Export some credentials to be able to pull secrets:

```
export AWS_ACCESS_KEY_ID="******"
export AWS_SECRET_ACCESS_KEY="******"
```

Export the version that you want to use:

```
export VMES_VERSION=0.0.1
```

Get that release and put the executable in a directory in your PATH:

```
wget https://github.com/external-secrets/vmes/releases/download/${VMES_VERSION}/vmes_${VMES_VERSION}_linux_amd64.tar.gz
tar -xvf vmes_${VMES_VERSION}_linux_amd64.tar.gz
sudo cp vmes /usr/local/bin/
```

If you want you can build the executable locally:

```
go build
```

And run it:

```
./vmes
```

## Roadmap

- [] Only AWS provider working for now, need to reimplement schema here or have another way to grab the right provider
- [] Configure where secrets will sink in, being a arbitraty file, exported directly as env vars, or something else.
- [] Support ec2 assume role and other auth methods
- [] Test setup
