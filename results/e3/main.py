import matplotlib.pyplot as plt
import numpy as np
import pandas as pd

delay_average = []
delay_err = []

p90_average = []
p90_err = []

dropped_average = []
dropped_err = []

x_axis = ['1s', '100ms', '10ms']
files = ['./s1.csv', './ms100.csv', './ms10.csv']

for result_file_name in files:
    print(result_file_name)

    r = pd.read_csv(result_file_name, header=None)

    print(r.head())

    r = np.array(r.values)

    # average delay
    avg = np.array([])
    for e in r:
        avg = np.append(avg, np.average(e[e > 0]))
        print(np.histogram(e[e > 0]))

    delay_average.append(np.average(avg))
    delay_err.append(np.std(avg))

    # maximum delay
    percentile = np.array([])
    for e in r:
        percentile = np.append(percentile, np.percentile(e[e > 0], 90))
    p90_average.append(np.average(percentile))
    p90_err.append(np.std(percentile))

    # dropped packets
    dropped = np.sum(r < 0, axis=1) / np.sum(r != 0, axis=1) * 100
    dropped_average.append(np.average(dropped))
    dropped_err.append(np.std(dropped))

print(f'Average Delay: {delay_average}')
print(f'Average Delay Standard Deviation: {delay_err}')

print(f'90% Percentile: {p90_average}')
print(f'90% Percentile Standard Deviation: {p90_err}')

print(f'Dropped Packets: {dropped_average}')
print(f'Dropped Packets Standard Deviation: {dropped_err}')

fig, ax = plt.subplots(figsize=(10, 10))
ax.errorbar(x=x_axis, y=delay_average, fmt='r--', yerr=delay_err)
ax.set_title('Latency')
ax.set(ylabel='Average Delay (s)', xlabel='Packet Rate (pps)')
fig.savefig('latency.png')

fig, ax = plt.subplots(figsize=(10, 10))
ax.errorbar(x=x_axis, y=p90_average, fmt='b--', yerr=p90_err)
ax.set_title('P90 Latency')
ax.set(ylabel='P90 Delay (s)', xlabel='Packet Rate (pps)')
fig.savefig('p90_latency.png')

fig, ax = plt.subplots(figsize=(10, 10))
ax.errorbar(x=x_axis, y=dropped_average, fmt='g--', yerr=dropped_err)
ax.set_title('Packet Drop Ratio')
ax.set(ylabel='Drop Ratio (%)', xlabel='Packet Rate (pps)')
fig.savefig('drop.png')
