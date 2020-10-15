FROM ubuntu:latest
RUN apt update && apt install -y g++ gcc bash perl sudo
WORKDIR /tests
COPY . /tests
CMD ["/bin/bash"]