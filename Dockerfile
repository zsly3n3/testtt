FROM scratch
COPY test_k8s_deploy test_k8s_deploy
EXPOSE 8180
CMD ["./test_k8s_deploy"]
