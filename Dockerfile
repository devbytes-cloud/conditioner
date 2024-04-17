FROM bitnami/kubectl:1.29.4
COPY kubectl-conditioner /opt/bitnami/kubectl/bin/kubectl-conditioner
ENTRYPOINT ["kubectl"]
