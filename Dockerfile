FROM alpine:latest as build

ARG DOCKER_METADATA_OUTPUT_JSON

RUN wget --no-check-certificate https://github.com/omerzamir/airflow-vars/releases/download/${DOCKER_METADATA_OUTPUT_JSON}/airflow-vars_Linux_arm64.tar.gz -O airflow-vars.tar.gz  \
    && tar -C /usr/local/bin -xzf airflow-vars.tar.gz \
    && rm airflow-vars.tar.gz

FROM gcr.io/distroless/static-debian11

COPY --from=build /usr/local/bin/airflow-vars /usr/local/bin/airflow-vars
ENTRYPOINT ["airflow-vars"]