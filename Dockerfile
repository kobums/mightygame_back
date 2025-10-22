FROM       alpine
MAINTAINER missbsd@gmail.com

COPY ./bin/mighty.linux /usr/local/main/main
COPY ./config/config.json /usr/local/main/config/config.json
CMD  mkdir -p /usr/local/main/webdata
#ADD ./assets /usr/local/main/assets
#ADD ./views /usr/local/main/views

WORKDIR /usr/local/main
CMD    ./main
