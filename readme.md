# Awsom - AWS applications made easy

[![Version](https://img.shields.io/badge/Awsom-0.3.0-blue.svg)](https://github.com/hekonsek/awsom/releases)

Awsom is a toolkit providing opinionated setup of infrastructure, applications and CI/CD based on AWS. This project
is intended to cover majority of common AWS deployments by following convention-over-configuration approach. Think about
it like "Rails but for AWS".

## Usage

```
docker run --net=host -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID -e AWS_DEFAULT_REGION=$AWS_DEFAULT_REGION -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY -it hekonsek/awsom:0.3.0
```

## License

This project is distributed under Apache 2.0 license.