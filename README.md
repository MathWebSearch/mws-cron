# mws-cron

[![Build Status](https://travis-ci.org/MathWebSearch/mws-cron.svg?branch=master)](https://travis-ci.org/MathWebSearch/mws-cron)

A Docker Container responsible for automatically updating a MathWebSearch instance.  

## Building

The code can be built using standard `golang` tools. 
Furthermore, a Makefile is provided which can be used as:

```bash
make
```

A static `mws-cron` executable (both for the current architecture and cross-compiled for others) can then be found in the `out` directory. 

## Usage

The executable can be run as follows:

```
Usage of mws-cron:
  -label string
        Label for MathWebSearch daemon (default "org.mathweb.mwsd")
  -pidfile string
        Pidfile to use
  -schedule string
        Cronline representing time to run job on (default "@midnight")
  -trigger
        Trigger manually running a cron job in running instance
```

This executable is intended to be run from within a Docker Container. 
When run, it will regularly update MWS indexes, that is it will start the [MathWebSearch/mws-indexer](https://github.com/MathWebSearch/mws-indexer) image with appropriate arguments and afterwards restart the mwsd container (which is found by the appropriate label). 

It is also possible to run the re-indexing process manually. 
This can be achieved by executing

```
/mws-cron --trigger
```

within the container. 

## Docker Image

This image can be found as the automated build [mathwebsearch/mws-cron](https://hub.docker.com/r/mathwebsearch/mws-cron) on DockerHub. 

## LICENSE

Licensed under GPL 3.0. 