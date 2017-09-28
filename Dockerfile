FROM alpine:3.6

#
# Copy release to container and set command
#

# Add faster mirror and upgrade packages in base image, Erlang needs ncurses, NIF libstdc++
RUN printf "http://mirror.leaseweb.com/alpine/v3.6/main\nhttp://mirror.leaseweb.com/alpine/v3.6/community" > etc/apk/repositories && \
    apk update && \
    apk upgrade && \
    apk add zeromq && \
    rm -rf /var/cache/apk/*

# Do not run as root
USER element43:element43

# Copy build
COPY emdr-to-nsq emdr-to-nsq

CMD ["/emdr-to-nsq"]