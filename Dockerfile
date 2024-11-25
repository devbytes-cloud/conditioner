FROM bitnami/kubectl:1.31.3
COPY kubectl-conditioner /opt/bitnami/kubectl/bin/kubectl-conditioner
ENTRYPOINT ["kubectl"]
