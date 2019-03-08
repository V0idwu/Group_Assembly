#!/bin/bash

peer chaincode install -n group_activity -p github.com/group_activity/go -l golang -v 1.1
# sleep 2
sleep 2
peer chaincode upgrade -o orderer.example.com:7050 -C mychannel -n group_activity -l golang  -c '{"Args":[]}' -P "OR ('Org1MSP.member','Org2MSP.member')" -v 1.1