#Docker container

## Docker test versions

	$ curl -fsSL https:test.docker.com/| sh #pipe character alt-7 on OSX

##docker run images
	$ # not nessecary if you use docker-compose.yml
	$ docker run --publish 8443:8443 --name pingservices -t pingservices

#run as service via docker-compose

##Build docker image
	# build with tag server (-t=tag)
	$ docker build --file Dockerfile.builder -t server .
	# new from here
	$ docker login
	$ docker tag server tetracon/golang:srv
	$ docker push tetracon/golang:srv	
	$ docker run -p4000:8080 tetracon/golang:srv
	

	
	$ docker-compose scale server=2
		Stopping and removing server_server_3 ... done
		Stopping and removing server_server_4 ... done
	
	$ docker-compose stop
	# stop all services
	$ docker-compose rm		#will remove all running services
	

	$ docker images
		REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
		server_web          latest              95f750d73235        40 minutes ago      717 MB
		server              latest              bbfb03d3c251        2 hours ago         708 MB
		golang              onbuild             ad323f40f596        2 weeks ago         703 MB
		golang              1.8                 c0ccf5f2c036        2 weeks ago         703 MB
		# server_web can be removed
	

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
	# remove image (rmi)
	$ docker rmi <IMAGE ID> 

##remove all stopped containers
	$ docker rm $(docker ps -q -f status=exited)
##Entirely wipe out all containers
	$ docker rm $(docker ps -a -q)
	
###Dockerfile.production

	# start with a scratch (no layers)
	FROM scratch

	# copy our static linked library
	COPY server server

	# tell we are exposing our service on port 8080
	EXPOSE 8080

	# run it!
	CMD ["./server"]
	
To build it manually run this command to build it. 

	$ docker build -f Dockerfile.production -t server:latest .
	
	
##Docker Comparing Containers and Virtual Machines

	http://www.docker.com/what-docker
	
##dock tutorial repro

	1. https://docs.docker.com/get-started/#container-diagram
	2. https://docs.docker.com/get-started/part2/#conclusion-of-part-one
	3. https://docs.docker.com/get-started/part3/#understanding-services
	4. 

#Docker on Raspberry PI2

	# We need to compile for 32 bits armv7
	$ GOOS=linux GOARCH=arm GOARM=7 go build -v
	$ file server
	
	# server: ELF 32-bit LSB executable, ARM, EABI5 version 1 (SYSV), statically linked, not stripped
	
	$ docker build --file Dockerfile.production -t server .
	
##Dockerfile.production
	
	# start with a scratch (no layers)
	#FROM scratch
	FROM scratch
	# copy our static linked library
	COPY server server
	
	# compile static linked binary
	$ GOOS=linux GOARCH=arm go build -a --ldflags '-extldflags "static"' -tags pingservices -installsuffix pingservices .
	# tell we are exposing our service on port 8080
	EXPOSE 8080

	# run it!
	CMD ["./server"]
	#
	$ docker tag server tetracon/server:srv
	$ docker push tetracon/server
	#if pulled from docker hub
	$ docker pull tetracon/server:srv
	
##On Raspberry PI2

	HypriotOS/armv7: pirate@black-pearl in ~
	$ docker run --publish 8081:8080 srv -t tetracon/server:srv
	

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
  	
-
	<!-- -->
	
	$ docker ps viz
	$ docker service ls
	$ docker node inspect self #find out ip-adress for manager
	$ docker service rm viz #remove current viz service
	
	<!-- -->

##Build docker image (pingservices)
	
	# build with tag pingservices (-t=tag)
	$ docker build --file Dockerfile.builder -t tetracon/pingservices:2.14 .
	# new from here
	$ docker login
	$ # docker tag pingservices tetracon/pingservices:2.14
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
	
	$ docker-compose -v
	# docker-compose version 1.14.0, build c7bdf9e
	$ docker-machine -v
	# docker-machine version 0.12.0, build 45c69ad
	$ docker -v
	# Docker version 17.05.0-ce, build 89658be
	#
	$ docker-compose version
	# docker-compose version 1.14.0, build c7bdf9e
	# docker-py version: 2.4.0
	# CPython version: 2.7.9
	# OpenSSL version: OpenSSL 1.0.1t  3 May 2016
	#
	$ docker version
	# Client:
	# Version:      17.05.0-ce
 	# API version:  1.29
 	# Go version:   go1.7.5
 	# Git commit:   89658be
 	# Built:        Thu May  4 22:30:54 2017
 	# OS/Arch:      linux/arm

	# Server:
 	# Version:      17.05.0-ce
 	# API version:  1.29 (minimum version 1.12)
 	# Go version:   go1.7.5
 	# Git commit:   89658be
 	# Built:        Thu May  4 22:30:54 2017
 	# OS/Arch:      linux/arm
 	# Experimental: false
 	#
 	# update docker hypriotOS
 	$ # DONT WORK -- apt-get update apt-get install docker-hypriot docker-compose
	$ apt-get update && apt-get install docker-compose
	
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

	# for production we can use docker-compose file, then we use docker stack deploy....      
	$ docker stack deploy -c docker-compose.yaml pingservices --with-registry-auth
	$ docker service ps pingservices
	$ docker service tasks pingservices
	$ docker stack ls
	
	# with docker run instead to be able to shell into the container
	$ docker run -it tetracon/pingservices:2.14 sh
	# stop the container when you are finnsihed
	$ docker stop <container-id>
	$ docker rm <container-id>
	$ sudo shutdown -h
	

##Docker exec

	$ docker ps # to find out container id
	# shell into the container
	$ docker exec -it <container_id> sh
	
##docker service create error

	$ docker service create --name=pingservices --publish=80:9000 tetracon/pingservices:2.14
	
	
==ERROR==

**image tetracon/pingservices:1.7 could not be accessed on a registry to record its digest. Each node will access tetracon/pingservices:1.7 independently, possibly leading to different nodes running different
versions of the image.**

This is solved by using ==**--with-registry-auth**== as a argument to docker service create

	$ docker service create --name=pingservices --publish=80:9000 --with-registry-auth tetracon/pingservices:1.9
	
**Explained by thaJeztah at github**

When updating services that need credentials to pull the image, you need to pass --with-registry-auth. Images pulled for a service take a different path than a regular docker pull, because the actual pull is performed on each node in the swarm where an instance is deployed. To pull an image, the swarm cluster needs to have the credentials stored (so that the credentials can be passed on to the node it performs the pull on).

Even though the "node" in this case is your local node, swarm takes the same approach (otherwise it would only be able to pull the image on the local node, but not on any of the other nodes).

Setting the --with-registry-auth option passes your locally stored credentials to the daemon, and stores them in the raft store. After that, the image digest is resolved (using those credentials), and the image is pulled on the node that the task is scheduled on (again, using the credentials that were stored in the service).

Can you confirm if passing --with-registry-auth makes the problem go away?
	
#Mac OSX known_hosts

	$ <user-id>/.ssh		#catalog on Mac

If you're removing hosts from the file, then you can just run the following on the command line

==**ssh-keygen -R HOSTNAME**==

You can search for a hostname with

==**ssh-keygen -F HOSTNAME**==
