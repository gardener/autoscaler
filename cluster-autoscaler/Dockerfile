
FROM gcr.io/distroless/static:nonroot

ARG TARGETARCH

COPY cluster-autoscaler-${TARGETARCH} /cluster-autoscaler
WORKDIR /
CMD ["/cluster-autoscaler"]
