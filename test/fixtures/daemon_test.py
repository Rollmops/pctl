#!/usr/bin/env python3

import os
import sys
import time

pid = os.fork()
print(pid)
if pid != 0:
    sys.exit()

time.sleep(10)
