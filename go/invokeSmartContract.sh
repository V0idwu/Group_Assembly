#!/bin/bash

peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n group_activity -c '{"function":$1,"Args":[$2]}'