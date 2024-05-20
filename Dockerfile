FROM bitnami/kubectl:1.30.1
COPY kubectl-conditioner /opt/bitnami/kubectl/bin/kubectl-conditioner
ENTRYPOINT ["kubectl"]
