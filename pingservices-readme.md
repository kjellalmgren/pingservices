# Project Docker container on raspberry PI 3 cluster

We need to update this documentation 2017-07-29 Kjell Almgren

##docker run images
	$ # not nessecary if you use docker-compose.yml
	$ docker run --publish 8443:8443 --name pingservices -t pingservices

##stop container
	$ docker stop <CONTAINER ID>
	$ docker ps	# list running container with CONTAINER ID
	$ docker ps -all # list all container (for rm)
	$ docker stats <CONTAINER ID>
	
##remove container
	$ docker rm <CONTAINER ID>
	
##Images
    #list all images
	$ docker images
	# remove images
	$ docker rmi <IMAGE ID> 

##remove all stopped containers
	$ docker rm $(docker ps -q -f status=exited)

##Entirely wipe out all containers
	$ docker rm $(docker ps -a -q)
	
#Dockerfile.builder

    # start from hypriot/rpi-alpine-scratch (nginx:alpine)
    #
    # -------------------------------------------------
    # FROM scratch
    # MAINTAINER kjell.almgren@tetracon.se
    # ADD pingservices /pingservices
    # ENTRYPOINT ["/pingservices"] 
    #
    # -------------------------------------------------
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
    # COPY pingservices pingservices
    ADD pingservices /pingservices

    # copy our self-signed certificate for now
    ##COPY tetracon-server.crt /go/src/server
    ##COPY tetracon-server.key /go/src/server

    # tell we are exposing our service on port 9000
    EXPOSE 9000

    # run it!
    CMD ["./pingservices"]
    #ENTRYPOINT ["/pingservices"]

To build it manually run this command to build it. 

	$ docker build -f Dockerfile.production -t pingservices:2.14 .
	
	
##Docker Comparing Containers and Virtual Machines

	http://www.docker.com/what-docker
	
##docker tutorial repro

	1. https://docs.docker.com/get-started/#container-diagram
	2. https://docs.docker.com/get-started/part2/#conclusion-of-part-one
	3. https://docs.docker.com/get-started/part3/#understanding-services
	4. 

#Docker on Raspberry PI3

	# We need to compile for 64 bits armv8
	$ GOOS=linux GOARCH=arm64 go build -v
	$ file pingservices
	
	# pingservices: ELF 64-bit LSB executable, ARM aarch64, version 1 (SYSV), statically linked, stripped
	
	$ docker build --file Dockerfile.production -t pingservices .
	
	
	#
	$ docker tag server tetracon/pingservices:2.14
	$ docker push tetracon/pingservices:2.14
	#if pulled from docker hub
	$ docker pull tetracon/pingservices:2.14
	
##Docker swarm visualizer
	
https://github.com/dockersamples/docker-swarm-visualizer

Should be run in the swarm manager, remember to login to hub.docker.com	

	<!-- -->
	
	docker service create \
	  --name=viz \
 	 --publish=4000:8080/tcp \
  	--constraint=node.role==manager \
  		--mount=type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock \
  	alexellis2/visualizer-arm:latest
  	
  	<!-- -->

    #command to be used for viz
	$ docker ps viz
	$ docker service ls
	$ docker node inspect self #find out ip-adress for manager
	$ docker service rm viz #remove current viz service
	
	<!-- -->

##Build docker image (pingservices)
	
	# build with tag pingservices (-t=tag)
	$ docker build --file Dockerfile.builder -t tetracon/pingservices:2.14 .
	# new from here
	$ docker login (add credentials)
	# push image to repository
	$ docker push tetracon/pingservices:2.14
	$ docker logoff
	
##Docker pull images from repository
	#
	$ docker login
	$ docker pull tetracon/pingservices:2.14
	$ docker images
	# docker service create or user docker-compose.yaml (see below)
	$ docker service create --name=pingservices --publish=80:9000 --with-registry-auth tetracon/pingservices:2.14
	# remove service
	$ docker service rm pingservices
	# find out wich node service is executing in
	$ docker service ps pingservices
	
		
##docker versions on HypriotOS 1.5.0

	<!-- -->
	
	
##docker-compose.yaml

	version: "3"

	services:
  	web:
    image: tetracon/pingservices:2.14
    deploy:
      replicas: 2
      resources:
        limits:
          cpus: "0.1"
          memory: 50M
      restart_policy:
        condition: on-failure
    ports:
      - "80:9000"
    networks:
     - webnet:
   networks:
     webnet:

	# for production we can use a docker-compose file, then we use docker stack deploy....      
	$ docker stack deploy -c docker-compose.yaml pingservices --with-registry-auth
	$ docker service ps pingservices
	$ docker service tasks pingservices
	$ docker stack ls
	
	# with docker run instead to be able to shell into the container
	$ docker run -it tetracon/pingservices:2.14 sh
	# stop the container when you are finnished
	$ docker stop <container-id>
	$ docker rm <container-id>
	$ sudo shutdown -h
	

##Docker exec

	$ docker ps # to find out container id
	# shell into the container
	$ docker exec -it <container_id> sh #remember to comment CMD["./PINGSERVICES"] in file Dockerfile.builder
	
##docker service create error

	$ docker service create --name=pingservices --publish=80:9000 tetracon/pingservices:2.14
	
	
==ERROR==

**image tetracon/pingservices:2.14 could not be accessed on a registry to record its digest. Each node will access tetracon/pingservices:2.14 independently, possibly leading to different nodes running different
versions of the image.**

This is solved by using ==**--with-registry-auth**== as a argument to docker service create

	$ docker service create --name=pingservices --publish=80:9000 --with-registry-auth tetracon/pingservices:2.14
	
**Explained by thaJeztah at github**

When updating services that need credentials to pull the image, you need to pass --with-registry-auth. Images pulled for a service take a different path than a regular docker pull, because the actual pull is performed on each node in the swarm where an instance is deployed. To pull an image, the swarm cluster needs to have the credentials stored (so that the credentials can be passed on to the node it performs the pull on).

Even though the "node" in this case is your local node, swarm takes the same approach (otherwise it would only be able to pull the image on the local node, but not on any of the other nodes).

Setting the --with-registry-auth option passes your locally stored credentials to the daemon, and stores them in the raft store. After that, the image digest is resolved (using those credentials), and the image is pulled on the node that the task is scheduled on (again, using the credentials that were stored in the service).
	
#Mac OSX known_hosts problem

	$ <user-id>/.ssh		#catalog on Mac

If you're removing hosts from the file, then you can just run the following on the command line

==**ssh-keygen -R HOSTNAME**==

You can search for a hostname with

==**ssh-keygen -F HOSTNAME**==
