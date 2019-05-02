# Copyright 2019 Robert Andrei STEFAN
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from os import listdir
from os.path import isfile, join
from collections import namedtuple
import numpy as np

PATH = "metrics/logs"
SYNCLOG = "synced buffer with message"
NOROUNDS = 50
NOLOSS = 10
NOBETA = 10
NOPEERS = 50

PeerStruct = namedtuple("PeerStruct", ["min", "max"])

# initialize two-dimensioanl arrays
# We can use: % peers = matrix[beta][loss][round]
matrix = np.empty((NOBETA, NOLOSS, NOROUNDS,), dtype = PeerStruct)


def init():
	for i in range (NOBETA):
		for j in range (NOLOSS):
			for k in range (NOROUNDS):
				matrix[i][j][k] = PeerStruct(min = -1.0, max = -1.0)

def parseLogFiles():
	# get files from path
	files = [f for f in listdir(PATH) if isfile(join(PATH, f))]

	lenFiles = len(files)
	for f in range(lenFiles):
		try:
			name = files[f].replace(".", "_").split("_")
			loss = int(round(float(name[4]) / 10))
			beta = int(round(float(name[6]) / 10))
		except:
			continue

		# get lines (logs) from file
		fileName = PATH + "/" + files[f]
		with open(fileName) as file:
			lines = file.readlines()

		countPeers = 0
		rounds = [ 0 for i in range(NOROUNDS) ]

		# parse synchronization logs and update rounds array with number of nodes
		lenLines = len(lines)
		for l in range(lenLines):
			if SYNCLOG in lines[l]:
				try:
					rrr = int(lines[l].split(" ")[9])
					rounds[rrr] = rounds[rrr] + 1
					countPeers = countPeers + 1
				except:
					print("Invalid log message: ", lines[l])
					continue

		# rounds[i] -> no of nodes with synced message after round i
		for r in range(1, NOROUNDS):
			rounds[r] = rounds[r - 1] + rounds[r]

		# Do not take into account if enough nodes have not been synchronized.
		# This may mean that nodes have not started
		if countPeers >= int(NOPEERS) - 10 and countPeers <= NOPEERS:
			for rr in range(NOROUNDS):
				# procent of nodes
				m = 100.0 * rounds[rr] / countPeers

				if matrix[beta][loss][rr] == PeerStruct(min = -1.0, max = -1.0):
					matrix[beta][loss][rr] = PeerStruct(min = m, max = m)
				else:
					if m < matrix[beta][loss][rr].min:
						matrix[beta][loss][rr] = PeerStruct(min = m, max = matrix[beta][loss][rr].max)
					if m > matrix[beta][loss][rr].max:
						matrix[beta][loss][rr] = PeerStruct(min = matrix[beta][loss][rr].min, max = m)
		else:
			print("Invalid log file: ", fileName)

		file.close()

	# for b in range(len(matrix)):
	# 	for l in range(len(matrix[b])):
	# 		for r in range(len(matrix[b][l])):
	# 			print("for beta = ", b, ", loss = ", l, " and round ", r)
	# 			print("min: ", matrix[b][l][r].min, " max: ", matrix[b][l][r].max)

	# print(matrix)

def main():
	init()
	parseLogFiles()

if __name__ == "__main__":
	main()
