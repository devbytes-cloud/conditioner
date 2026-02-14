FROM alpine/kubectl:1.35.1
COPY kubectl-conditioner /usr/local/bin/kubectl-conditioner
ENTRYPOINT ["kubectl"]
