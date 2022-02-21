# vmes

External Secrets not necessarily in Kubernetes.

This project uses [ESO](https://github.com/external-secrets/external-secrets) as a dependency and basicaly leverages its provider implementations to grab secrets and write them to an env file.

The initial need here is to run this project as a proccess on a VM that has some legacy applications, then we want to use secrets that are in a Secrets Manager. But we don't want to keep calling Secrets Managers from our applications, we want to just read values from env vars or env files.

There is a way to encrypt and decrypt data to hide secrets a bit more. We are still experimenting with it. The idea is that you should generate private and public keys and use vmes together with [saferun](https://github.com/ContainerSolutions/saferun) if you would like asymetric encryption in your machine. No changes are needed to apps reading env vars/env files.

### Disclaimer ⚠️

This project is in early access stage. Use at your own discretion.

# How it works

Since we already implemented and tested multiple provider clients in the ESO project, we are just importing them here. Since we had some dependencies on Kubernetes resources for some of the configurations, we have built a client that just looks for yaml files locally instead of calling Kubernetes at all.

# Getting Started

Copy the example files provided to edit and configure them:

```
cp pkg/configdata/es.yml.example ~/.vmes/es.yml
cp pkg/configdata/ss.yml.example ~/.vmes/ss.yml
```

Open `~/.vmes/es.yml` and edit the following fields:

- **spec.refreshInterval:** Choose a Time Duration interval used by the operator to fetch new secrets (1m = 1 minute, 1h = 1 hour, etc).
- **spec.target.name:** Chose the env file where this will end up (default is /etc/environment - You need to run as root to write there).
- **spec.data.secretKey:** The name of the Env Var injected in the machine.
- **spec.data.remoteRef.Key:** The name of the secret in the external provider.

Open `~/.vmes/ss.yml` and edit the following fields:

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

If you are using vmes encryption and saferun, generate a key pair:

```
openssl genrsa -out myuser.key 2048
openssl rsa -in myuser.key -out myuser.pub -pubout -outform PEM
```

To run the installed release just call vmes anywhere (if you don't want asymmetric encryption, just omit `--public-key-path /home/youruser/.vmes/test.pub`):

```
vmes --config-path /home/youruser/.vmes --public-key-path /home/youruser/.vmes/test.pub
```

If you want you can build the executable locally:

```
go build
```

And run it (if you don't want asymmetric encryption, just omit `--public-key-path /home/youruser/.vmes/test.pub`):

```
./vmes --config-path /home/youruser/.vmes --public-key-path /home/youruser/.vmes/test.pub
```

To get values in your app with saferun, you can use:

```
./bin/saferun run --private-key=test.key  --only-encrypted /bin/env

or with your app

./bin/saferun run --private-key=test.key  --only-encrypted /path/to/app
```

If you are getting vars from /etc/environment and don't want to re-login to check you can run:

```
for line in $( cat /etc/environment ) ; do export $line ; done
```

## Systemd config

You probably want to to run this tool as a service in a machine. For that you can create a Systemd unit file and let Systemd manage it. Run these commands:

```
sudo cat > /etc/systemd/system/vmes.service <<EOF
[Unit]
Description=vmes
After=network.target

[Service]
Type=idle
User=root
Group=keycloak
ExecStart=vmes --config-path /home/youruser/.vmes --public-key-path /home/youruser/.vmes/test.pub
TimeoutStartSec=600
TimeoutStopSec=600

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl start vmes
sudo systemctl enable vmes
```

## Docs

You can find examples and other docs at [docs](docs)

## Roadmap

- [] Only AWS provider working for now, need to reimplement schema here or have another way to grab the right provider
- [] Add option to use multiple public keys and multiple files to sink in (read more ESs)
- [] Configure where secrets will sink in
    - ✅ being a arbitraty file, 
    - [] exported directly as env vars, or something else.
- ✅ Support assume role and other auth methods
- [] Test setup
- ✅ Provide a way to configure different paths for where yaml files could be
- ✅ Integrate with saferun
- ✅ Add option to enable/disable encryption (also adds/removes SAFE_RUN_ prefix to envs in source files)
- [] Upgrade to new ESO CRDS
