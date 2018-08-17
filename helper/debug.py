import socket
import sys
import time
import json
import decimal

import xmlrpclib
import os
import time

def fetchMeasurements():

    "Detects smart power sockets by sending a UDP broadcast."

    # Create a UDP socket
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.bind(('', 0))
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)
    sock.settimeout(15)

    server_address = ('192.168.178.255', 8888)
    message = '00dv=all,2016-01-31,19:27:45,13;'

    no_of_expected_devices = 2
    data_dict = {}

    try:

        # Send data
        # print >>sys.stderr, 'sending "%s"' % message
        sent = sock.sendto(message, server_address)

        timeout = time.time() + 15   # 5 minutes from now

        msg_received = 0
        while True:
            if msg_received == no_of_expected_devices or time.time() > timeout:
                break

            # Receive response
            # print >> sys.stderr, 'waiting to receive'
            data, server = sock.recvfrom(4096)
            print >> sys.stderr, 'received from %s: %s' % (server, data)

            data_dict[server[0]] = data
            msg_received = msg_received + 1

    except socket.timeout:
        # do nothing as
        print >> sys.stderr, "Timeout occured."
        return

    finally:
        sock.close()

fetchMeasurements()
