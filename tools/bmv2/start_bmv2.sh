#!/bin/bash
# 
# Copyright 2019-present Open Networking Foundation
# 
# SPDX-License-Identifier: Apache-2.0
# 

CHASSIS_CONFIG_FILE=/root/tools/bmv2/chassis-config.txt
CPU_PORT=253
GRPC_PORT=50001
LOG_LEVEL=debug

# Setup veth interfaces
for idx in 0 1; do
    intf0="veth$(($idx*2))"
    intf1="veth$(($idx*2+1))"
    if ! ip link show $intf0 &> /dev/null; then
        ip link add name $intf0 type veth peer name $intf1
        ip link set dev $intf0 up
        ip link set dev $intf1 up

        # Set the MTU of these interfaces to be larger than default of
        # 1500 bytes, so that P4 behavioral-model testing can be done
        # on jumbo frames.
        # Note: ifconfig is deprecated, and no longer installed by
        # default in Ubuntu Linux minimal installs starting with
        # Ubuntu 18.04.  The ip command is installed in Ubuntu
        # versions since at least 16.04, and probably older versions,
        # too.
        ip link set $intf0 mtu 9500
        ip link set $intf1 mtu 9500

        # Disable IPv6 on the interfaces, so that the Linux kernel
        # will not automatically send IPv6 MDNS, Router Solicitation,
        # and Multicast Listener Report packets on the interface,
        # which can make P4 program debugging more confusing.
        #
        # Testing indicates that we can still send IPv6 packets across
        # such interfaces, both from scapy to simple_switch, and from
        # simple_switch out to scapy sniffing.
        #
        # https://superuser.com/questions/356286/how-can-i-switch-off-ipv6-nd-ra-transmissions-in-linux
        sysctl net.ipv6.conf.${intf0}.disable_ipv6=1
        sysctl net.ipv6.conf.${intf1}.disable_ipv6=1
    fi
done

# Start stratum_bmv2
stratum_bmv2 -device_id=1 -chassis_config_file=${CHASSIS_CONFIG_FILE} -forwarding_pipeline_configs_file=/tmp/s1/pipe.txt -persistent_config_dir=/tmp/s1 -initial_pipeline=/root/dummy.json -cpu_port=${CPU_PORT} -external_stratum_urls=0.0.0.0:${GRPC_PORT} -local_stratum_url=localhost:49343 -max_num_controllers_per_node=10 -write_req_log_file=/tmp/s1/write-reqs.txt -logtosyslog=false -logtostderr=true -bmv2_log_level=${LOG_LEVEL}
