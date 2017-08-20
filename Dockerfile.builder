# start from hypriot/rpi-alpine-scratch (nginx:alpine)
#
# -------------------------------------------------
# FROM resin/raspbian
# MAINTAINER kjell.almgren@tetracon.se
# -------------------------------------------------
#
FROM resin/rpi-raspbian

MAINTAINER kjell.almgren@tetracon.se

# make some update to the OS in the container
#RUN apk update && \
#apk upgrade && \
#apk add bash && \
#rm -rf /var/cache/apk/*

#make some changes to the container images (docker dns-bugs)
#COPY docker-compose.yml docker-compose.yaml
#switch to our app directory (/pingservices)
RUN mkdir -p /pingservices
WORKDIR /pingservices

#create our sub directories
RUN mkdir -p dist/css/
RUN mkdir -p dist/fonts
RUN mkdir -p dist/js
RUN mkdir -p images
RUN mkdir -p templates
RUN mkdir -p assets/js/vendor

#copy distributions files
COPY dist/css dist/css
COPY dist/fonts dist/fonts
COPY dist/js dist/js

#COPY images files
COPY images images

#vendor files
COPY assets assets
COPY assets/js assets/js
COPY assets/js/vendor assets/js/vendor

#copy main template files
COPY templates templates

#copy the main services
COPY main.css main.css
COPY services-prod.json services-prod.json
COPY services-qa.json services-qa.json
# COPY executable pingservices /pingservices
COPY pingservices /pingservices

# copy our self-signed certificate
##COPY tetracon-server.crt /go/src/pingservices
##COPY tetracon-server.key /go/src/pingservices

# tell we are exposing our service on port 9000
EXPOSE 9000

# run it!

ENTRYPOINT ["./pingservices"]
#CMD ["./pingservices"]