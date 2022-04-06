#!/bin/bash
set -x

bridge="br0"
cidr="192.168.1.0/24"
gatewayIP="192.168.1.1"
gateway="${gatewayIP}/24"
hostNetWork="eth0"

function addBridge() {
    # add bridge
    sudo ip link add name ${bridge} type bridge

    # add addr
    sudo ip addr add ${gateway} dev ${bridge}

    # setup
    sudo ip link set ${bridge} up

    echo "set bridge $bridge done"
}

function setupVeth() {
    ns=$1
    veth=$2
    vethbr=${veth}-br
    ip=$3

    # add netns
    sudo ip netns add $ns

    # add veth pair, link to bridge and netns
    sudo ip link add $veth type veth peer name $vethbr
    sudo ip link set $veth netns $ns
    sudo ip link set dev $vethbr master ${bridge}

    # add addr for veth
    sudo ip -n $ns addr add local $ip dev $veth

    # setup veth
    sudo ip -n $ns link set $veth up
    sudo ip link set $vethbr up
    sudo ip netns exec $ns ip link set lo up

    # add route
    sudo ip netns exec $ns ip route add default via $gatewayIP

    echo "set netns $ns done"
}

function setIptables() {
    # enable bridge forward
    sudo iptables -A FORWARD -i ${bridge} -j ACCEPT
}

function setOuterIptables() {
    # snat
    sudo iptables -t nat -A POSTROUTING -s $cidr -j MASQUERADE
    # enable forward for hostnetwork
    sudo iptables -A FORWARD -i $hostNetWork -j ACCEPT
}

function clearIptables() {
    sudo iptables -D FORWARD -i ${bridge} -j ACCEPT
    sudo iptables -t nat -D POSTROUTING -s $cidr -j MASQUERADE
    sudo iptables -D FORWARD -i $hostNetWork -j ACCEPT
}

function start() {
    echo "start run script"
    addBridge

    setupVeth "ns1" "vns1" "192.168.1.101/24"
    setupVeth "ns2" "vns2" "192.168.1.102/24"

    setIptables
}

function clear() {
    sudo ip netns del ns1
    sudo ip netns del ns2

    sudo ip link del ${bridge}
    clearIptables
}

while getopts 'rci' opt; do
    case $opt in
        r) start;;
        c) clear;;
        i) setOuterIptables;;
    esac
done
