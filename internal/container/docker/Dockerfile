FROM ubuntu:latest
RUN apt update && apt install -y g++ gcc bash perl
RUN DEBIAN_FRONTEND=noninteractive apt install -y openjdk-8-jdk
WORKDIR /tests
COPY . /tests
RUN apt install -y --install-recommends ./benchexec_*.deb
RUN chmod -R 700 /tests
CMD ["/bin/bash"]