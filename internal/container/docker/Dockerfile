FROM ubuntu:latest
RUN apt update && apt install -y g++ gcc bash perl
WORKDIR /tests
COPY . /tests
RUN apt install -y --install-recommends ./benchexec_*.deb
CMD ["/bin/bash"]