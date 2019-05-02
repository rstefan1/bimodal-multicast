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
import matplotlib.pyplot as plt
from matplotlib.widgets import Slider

PATH = "metrics/logs"
SYNCLOG = "synced buffer with message"
NOROUNDS = 50
NOLOSS = 10
NOBETA = 10
NOPEERS = 50

betaPlotMin = 0    # the minimial value of the paramater beta
betaPlotMax = 0.9   # the maximal value of the paramater beta
betaPlotInit = 0.5   # the value of the parameter beta to be used initially, when the graph is created

lossPlotMin = 0    # the minimial value of the paramater loss
lossPlotMax = 0.9   # the maximal value of the paramater loss
lossPlotInit = 0   # the value of the parameter loss to be used initially, when the graph is created

PeerStruct = namedtuple("PeerStruct", ["min", "max"])

# % procent of peers = matrix[beta][loss][round]
matrix = np.empty((NOBETA, NOLOSS, NOROUNDS,), dtype = PeerStruct)

# global values for beta and loss used for plotting
globalBeta = 0
globalLoss = 0
globalRoundNumber = 25

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
					r = int(lines[l].split(" ")[9])
					rounds[r] = rounds[r] + 1
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
			for r in range(NOROUNDS):
				# procent of nodes
				m = 100.0 * rounds[r] / countPeers

				if matrix[beta][loss][r] == PeerStruct(min = -1.0, max = -1.0):
					matrix[beta][loss][r] = PeerStruct(min = m, max = m)
				else:
					if m < matrix[beta][loss][r].min:
						matrix[beta][loss][r] = PeerStruct(min = m, max = matrix[beta][loss][r].max)
					if m > matrix[beta][loss][r].max:
						matrix[beta][loss][r] = PeerStruct(min = matrix[beta][loss][r].min, max = m)
		else:
			print("Invalid log file: ", fileName)

		file.close()



def plot():
	fig = plt.figure(figsize=(8,3))

	global globalBeta
	globalBeta = int(round(betaPlotInit * 10))

	global globalLoss
	globalLoss = int(round(lossPlotInit * 10))

	plotRounds = range(globalRoundNumber)
	plotPeersMin = []
	plotPeersMax = []
	for i in range(globalRoundNumber):
		plotPeersMin.append(matrix[globalBeta][globalLoss][i].min)
		plotPeersMax.append(matrix[globalBeta][globalLoss][i].max)

	# first we create the general layout of the figure
	# with two axes objects: one for the plot of the function
	# and the other for the slider
	figAx = plt.axes([0.1, 0.27, 0.8, 0.65])
	betaSliderAx = plt.axes([0.1, 0.12, 0.8, 0.05])
	lossSliderAx = plt.axes([0.1, 0.05, 0.8, 0.05])

	# in plot_ax we plot the function with the initial value of the parameter beta
	plt.sca(figAx) # select sin_ax
	figPlotMin, = plt.plot(plotRounds, plotPeersMin, 'r', label='Min bound')
	figPlotMax, = plt.plot(plotRounds, plotPeersMax, 'g', label='Max bound')
	plt.xlabel('Rounds')
	plt.ylabel('% peers')
	plt.legend()

	# here we create the beta slider
	betaSlider = Slider(betaSliderAx,	# the axes object containing the slider
		'Beta',							# the name of the slider parameter
		betaPlotMin,					# minimal value of the parameter
		betaPlotMax,					# maximal value of the parameter
		valinit = betaPlotInit			# initial value of the parameter
		)

	# here we create the loss slider
	lossSlider = Slider(lossSliderAx,	# the axes object containing the slider
		'Loss',							# the name of the slider parameter
		lossPlotMin,					# minimal value of the parameter
		lossPlotMax,					# maximal value of the parameter
		valinit = lossPlotInit			# initial value of the parameter
		)

	# We define a function that will be executed each time the value
	# indicated by the slider changes. The variable of this function will
	# be assigned the value of the beta slider.
	def updateBeta(beta):
		beta = int(round(beta * 10))
		global globalBeta
		globalBeta = beta

		plotPeersMin = []
		plotPeersMax = []

		for i in range(globalRoundNumber):
			plotPeersMin.append(matrix[beta][globalLoss][i].min)
			plotPeersMax.append(matrix[beta][globalLoss][i].max)

		figPlotMax.set_ydata(plotPeersMax)	# set new y-coordinates of the plotted points
		figPlotMin.set_ydata(plotPeersMin)	# set new y-coordinates of the plotted points
		fig.canvas.draw_idle()				# redraw the plot

	# We define a function that will be executed each time the value
	# indicated by the slider changes. The variable of this function will
	# be assigned the value of the loss slider.
	def updateLoss(loss):
		loss = int(round(loss * 10))
		global globalLoss
		globalLoss = loss

		plotPeersMin = []
		plotPeersMax = []

		for i in range(globalRoundNumber):
			plotPeersMin.append(matrix[globalBeta][loss][i].min)
			plotPeersMax.append(matrix[globalBeta][loss][i].max)

		figPlotMax.set_ydata(plotPeersMax)	# set new y-coordinates of the plotted points
		figPlotMin.set_ydata(plotPeersMin)	# set new y-coordinates of the plotted points
		fig.canvas.draw_idle()				# redraw the plot

	# specify that the sliders needs to
	# execute the above function when its value changes
	betaSlider.on_changed(updateBeta)
	lossSlider.on_changed(updateLoss)

	plt.show()

def main():
	init()
	parseLogFiles()
	plot()

if __name__ == "__main__":
	main()
