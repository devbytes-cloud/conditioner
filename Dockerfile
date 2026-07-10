FROM alpine/kubectl:1.36.2
COPY kubectl-conditioner /usr/local/bin/kubectl-conditioner
ENTRYPOINT ["kubectl"]
