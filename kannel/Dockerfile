FROM ubuntu:20.04

# Install Kannel and cleanup
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
        kannel \
        curl \
        nano \
        inetutils-ping && \
    apt-get clean

# Expose Kannel ports
EXPOSE 13000 13001 13013

# Default command to start bearerbox and smsbox
CMD ["sh", "-c", "bearerbox /etc/kannel/kannel.conf & sleep 2 && smsbox /etc/kannel/kannel.conf -v 1"]

