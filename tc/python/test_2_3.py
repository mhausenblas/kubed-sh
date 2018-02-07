#!/usr/bin/env python

import socket

if __name__ == "__main__":
    hostname = socket.getfqdn()
    print("Hello from Python! I'm running in a container in the Kubernetes cluster, on host " + hostname)
