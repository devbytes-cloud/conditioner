FROM alpine/kubectl:1.35.2
COPY kubectl-conditioner /usr/local/bin/kubectl-conditioner
ENTRYPOINT ["kubectl"]
