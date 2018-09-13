#!/usr/bin/python

from websocket import create_connection
import json

ws = create_connection("ws://slither.local:4050/wscall")

# args = {'c': "ls", 'a': ["..", "-la" ]}

args = {'c': "sysType", 'a': ["", ]}
data = json.dumps(args).encode()
#print (data)

print ("Sending: {}".format( args ))
ws.send(data)

st = json.loads( ws.recv() )['R']
print("Received Response: {} ".format ( st ))

ws.close()
