FROM alpine:latest

#
# Copy release to container and set command
#

# Add faster mirror and upgrade packages in base image, Erlang needs ncurses, NIF libstdc++
RUN printf "http://mirror.leaseweb.com/alpine/3.6/main\nhttp://mirror.leaseweb.com/alpine/3.6/community" > etc/apk/repositories && \
    apk update && \
    apk upgrade && \
    apk add zeromq && \
    rm -rf /var/cache/apk/*

# Copy build
COPY emdr-to-nsq emdr-to-nsq

CMD ["/emdr-to-nsq"]