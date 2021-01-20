FROM golang:1.15.6 AS go-builder
# RUN go get -d -v golang.org/x/net/html  
RUN GO111MODULE=on go get -u go.mozilla.org/sops/v3/cmd/sops@v3.6.1 && \
    cd $GOPATH/pkg/mod/go.mozilla.org/sops/v3@v3.6.1 && \
    make install


FROM python:3.8.5-alpine
COPY --from=go-builder /go/bin/sops /usr/bin/

ENV AZURE_CLI_VERSION=2.5.1
ENV TERRAFORM_VERSION=0.13.4
ENV TERRAFORM_EXTERNAL_VERSION=1.2.0
ENV TERRAFORM_AZURE_VERSION=2.29.0
ENV TERRAFORM_AWS_VERSION=3.8.0
ENV TERRAFORM_NULL_VERSION=2.1.2
ENV TERRAFORM_COMMERCETOOLS_VERSION=0.23.0
ENV TERRAFORM_PLUGINS_PATH=/root/.terraform.d/plugins/linux_amd64
RUN mkdir -p ${TERRAFORM_PLUGINS_PATH}

RUN apk update && \
    apk add --no-cache --virtual .build-deps g++ python3-dev libffi-dev openssl-dev && \
    apk add --no-cache --update python3 && \
    apk add bash curl tar ca-certificates git libc6-compat openssl jq unzip wget openssh-client make

# Install terraform
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/bin

# Install null provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-null/${TERRAFORM_NULL_VERSION}/terraform-provider-null_${TERRAFORM_NULL_VERSION}_linux_amd64.zip && \
    unzip terraform-provider-null_${TERRAFORM_NULL_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH}

# Install external provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-external/${TERRAFORM_EXTERNAL_VERSION}/terraform-provider-external_${TERRAFORM_EXTERNAL_VERSION}_linux_amd64.zip && \
    unzip terraform-provider-external_${TERRAFORM_EXTERNAL_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH}

# Install aws provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-aws/${TERRAFORM_AWS_VERSION}/terraform-provider-aws_${TERRAFORM_AWS_VERSION}_linux_amd64.zip && \
    unzip terraform-provider-aws_${TERRAFORM_AWS_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH}

# Install azure provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-azurerm/${TERRAFORM_AZURE_VERSION}/terraform-provider-azurerm_${TERRAFORM_AZURE_VERSION}_linux_amd64.zip && \
    unzip terraform-provider-azurerm_${TERRAFORM_AZURE_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH}


# Install commercetools provider
RUN cd /tmp && \
    wget https://github.com/labd/terraform-provider-commercetools/releases/download/v${TERRAFORM_COMMERCETOOLS_VERSION}/terraform-provider-commercetools_${TERRAFORM_COMMERCETOOLS_VERSION}_linux_amd64.zip && \
    unzip terraform-provider-commercetools_${TERRAFORM_COMMERCETOOLS_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH}

RUN rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

RUN pip --no-cache-dir install azure-cli==${AZURE_CLI_VERSION}

# TODO: use build containers to optimize this, for now this works ;^)
RUN mkdir /code
RUN mkdir /deployments
WORKDIR /code

ADD requirements.txt .
RUN pip install -r requirements.txt
COPY src /code/src/
ADD MANIFEST.in .
ADD setup.cfg . 
ADD setup.py . 
RUN python setup.py bdist_wheel && pip install dist/mach-0.0.0-py3-none-any.whl


ENTRYPOINT ["mach"]
