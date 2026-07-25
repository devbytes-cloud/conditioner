FROM alpine/kubectl:1.36.3
COPY kubectl-conditioner /usr/local/bin/kubectl-conditioner
ENTRYPOINT ["kubectl"]
