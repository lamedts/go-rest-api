FROM ubuntu:16.04

# create user
RUN groupadd zakkaya
RUN useradd -ms /bin/bash yay

# create config directory
RUN mkdir -p /var/yay/config
RUN mkdir -p /var/log/yay
RUN chown -R yay /var/log/yay
RUN chown -R yay /var/yay

RUN apt-get update
# RUN apt-get install -y sudo
RUN apt-get install -y ca-certificates

# copy binary bundle
RUN mkdir - p /opt/yay
WORKDIR /opt/yay
COPY ./build/linux_x64/go-rest-api /opt/yay/


# copy zakkaya configuration
COPY ./build/linux_x64/config.yaml /opt/yay/
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.3.0/wait /wait
RUN chmod +x /wait

# set environment variables
# ENV APP_ENV=DEBUG
EXPOSE 8080
# ENTRYPOINT /opt/yay/go-rest-api 
# CMD ["tail", "-f", "/dev/null"]
CMD /wait && /opt/yay/go-rest-api 

