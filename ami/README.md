[Packer](https://www.packer.io/)を使ってAMIの作成をする。

# How to use

## Prerequisites
* [Session Manager plugin for the AWS CLI](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)

## Steps

```
$ cd ishocon1/ami
$ cp _secret_vars.hcl secret_vars.hcl
# Update your values at `secret_vars.hcl`

$ cd ../ && tar -zcvf ami/webapp.tar.gz ./webapp && cd ami
$ cd ../admin/ && tar -zcvf ../ami/benchmarker.tar.gz ./benchmarker && cd ../ami

$ packer init .
$ packer validate -var-file=shared_vars.hcl -var-file=secret_vars.hcl .

# This will take around 15 minutes to run
$ packer build -var-file=shared_vars.hcl -var-file=secret_vars.hcl ami.pkr.hcl

# You can see the AMI ID from the output
```
