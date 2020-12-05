#!/usr/bin/env python2
#encoding: utf−8

from pwn import *
import unittest
import os

# Fonction utilitaires

port1 = b'0'
port2 = b'1'
port3 = b'2'
port4 = b'3'
debugDelay = 8
config_file = "nodes.json"
config_2 = "nodes2.test.json"
config_2_dbg = "nodes2.test.dbg.json"
config_4 = "nodes4.test.json"
config_4_dbg = "nodes4.test.dbg.json"
script_dir = os.path.dirname(__file__)
dir = script_dir.replace('/e2e', '') + '/nodes'
send_msg_hint = b'>>'

def createNodeAndConnect(port):
    node = process('sh')
    send(node, b'go run nodes/node.go ' + port + b' e2e')
    return node

def closeNode(node):
    node.close()
    time.sleep(0.05)

def closeNodes(nodes):
    for node in nodes:
        closeNode(node)

def iWasWaitingFor(node, this, self):
    print(node)
    print(b'Waiting for: ' + this)
    output = node.recvuntil(this, timeout=1 + debugDelay)
    print(output)
    self.assertFalse(b'' == output)

def send(io, line):
    print(b'You sent: ' + line)
    io.sendline(line)

def configureWith(fileToUse):
    f = open(dir + '/' + fileToUse, "r")
    text = f.read(-1)
    f2 = open(dir + '/' + config_file, "w")
    f2.write(text)
    f.close()
    f2.close()

def waitCommunication():
    time.sleep(0.5)

class Test2Nodes(unittest.TestCase):
    def test_2nodes(self):
        configureWith(config_2)

        test_name = b'test_2nodes'

        node1 = createNodeAndConnect(port1)
        node2 = createNodeAndConnect(port2)

        time.sleep(1)

        send(node1, test_name + b'1')
        time.sleep(1)
        send(node2, test_name + b'2')

        iWasWaitingFor(node2, b'Hi there', self)
        iWasWaitingFor(node1, b'Hi there', self)

        closeNodes([node1, node2])


if __name__ == '__main__':
    unittest.main()
