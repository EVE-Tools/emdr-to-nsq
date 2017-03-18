FROM alpine:latest

MAINTAINER zweizeichen@element-43.com

#
# Copy release to container and set command
#

# Add faster mirror and upgrade packages in base image, Erlang needs ncurses, NIF libstdc++
RUN printf "http://mirror.leaseweb.com/alpine/edge/main\nhttp://mirror.leaseweb.com/alpine/edge/community" > etc/apk/repositories && \
    apk update && \
    apk upgrade && \
    apk add zeromq && \
    rm -rf /var/cache/apk/*

# Copy build
COPY emdr-to-nsq emdr-to-nsq

CMD ["/emdr-to-nsq"]