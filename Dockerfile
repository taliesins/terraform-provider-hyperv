FROM hashicorp/terraform:full

WORKDIR \
$GOPATH/src/github.com/taliesins/terraform-provider-hyperv

COPY . .

RUN apk add --update make bash && \
chmod -R +x ./scripts && \
make fmt && \
make build

WORKDIR $GOPATH
ENTRYPOINT ["terraform"]
