#!/usr/bin/env python2
#encoding: utf−8

from pwn import *
import unittest
import os
import time

# Fonction utilitaires

port1 = b'0'
port2 = b'1'
port3 = b'2'
port4 = b'3'
getEluCmd = b'g'
t = 1.5 + 0.3
timeout = 2 + 3*t + 1
debugDelay = 4
config_file = "config.json"
config_2 = "config2.test.json"
config_2_dbg = "config2.test.dbg.json"
config_4 = "config4departage.test.json"
config_4_dbg = "config4.test.dbg.json"
script_dir = os.path.dirname(__file__)
dir = script_dir.replace('/e2e', '') + '/configs'
send_msg_hint = b'>>'

def createProcess(port):
    p = process('/bin/sh')
    send(p, b'cd .. && go run process.go ' + port)
    return p

def closeProcess(p):
    p.kill()

def closeProcesses(ps):
    for p in ps:
        closeProcess(p)
    time.sleep(2)

def getOneLineDebug(p):
    output = p.recvline(timeout=timeout)
    print(output)

def iWasWaitingFor(p, this, self):
    print(p)
    print(b'Waiting for: ' + this)
    output = p.recvuntil(this, timeout=timeout)
    self.assertFalse(b'' == output)

def waitInit(p):
    p.recvuntil(b'Election results are known after 3s', timeout=timeout)

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

def sleepInit():
    time.sleep(2)

class Test2Nodes(unittest.TestCase):

    def test_both_process_determine_the_same_initial_winner(self):
        configureWith(config_2)

        p1 = createProcess(port1)
        p2 = createProcess(port2)
        waitInit(p1)
        waitInit(p2)

        initialElectionMessage = b'L\'elu de l\'election initiale est le processus: 1'

        iWasWaitingFor(p1, initialElectionMessage, self)
        iWasWaitingFor(p2, initialElectionMessage, self)

        closeProcesses([p1, p2])

    def test_4_process_determine_the_same_initial_winner(self):
        configureWith(config_4)

        p1 = createProcess(port1)
        p2 = createProcess(port2)
        p3 = createProcess(port3)
        p4 = createProcess(port4)

        initialElectionMessage = b'Le processus 2 est l\'elu!'

        for p in [p1, p2, p3, p4]:
            waitInit(p)
            send(p, getEluCmd)
            iWasWaitingFor(p, initialElectionMessage, self)

        closeProcesses([p1, p2, p3, p4])


if __name__ == '__main__':
    unittest.main()
