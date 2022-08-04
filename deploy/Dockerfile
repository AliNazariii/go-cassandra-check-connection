FROM frolvlad/alpine-glibc:glibc-2.33

ARG docker_version=0.0.1

COPY ./build/release/cassandra-connection-check_$docker_version /bin/cassandra-connection-check

CMD ["/bin/cassandra-connection-check"]