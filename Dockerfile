FROM bitnami/kubectl:1.29.2
COPY kubectl-condition /opt/bitnami/kubectl/bin/kubectl-condition
ENTRYPOINT ["kubectl"]
