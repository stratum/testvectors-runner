# 
# Copyright 2019-present Open Networking Foundation
# 
# SPDX-License-Identifier: Apache-2.0
# 
# Docker image to run stratum_bmv2 switch with interfaces configured for data plane
# Use "make switch" to start the switch

FROM opennetworking/mn-stratum
RUN apt-get update && apt-get install -yq make vim
RUN mkdir -p tools/bmv2
COPY ./tools/bmv2 ./tools/bmv2/

WORKDIR /root/tools/bmv2

ENTRYPOINT make switch
