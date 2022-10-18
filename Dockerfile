###########################################################################
# Stage 1 Start
###########################################################################
FROM golang AS build-golang


RUN export GO111MODULE=on
RUN export GOPROXY=direct
RUN export GOSUMDB=off

################################
# Build Service:
################################
WORKDIR /usr/share/service/hypefast

COPY  . .

RUN make deploy

###########################################################################
# Stage 2 Start
###########################################################################
FROM ubuntu:18.04

# Change Repository ke kambing.ui:
RUN sed -i 's*archive.ubuntu.com*kambing.ui.ac.id*g' /etc/apt/sources.list

RUN apt-get update

RUN apt-get install -y ca-certificates

#Copy Desc File
COPY --from=build-golang /usr/share/service/hypefast/src/desc /usr/share/service/hypefast/src/

# Copy Binary
COPY --from=build-golang /usr/share/service/hypefast/bin /usr/share/service/hypefast/bin/

WORKDIR /usr/share/service/hypefast

# Create group and user to the group
RUN groupadd -r hypefast && useradd -r -s /bin/false -g hypefast hypefast

# Set ownership golang directory
RUN chown -R hypefast:hypefast /usr/share/service/hypefast

# Make docker container rootless
USER hypefast
