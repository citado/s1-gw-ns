"""
this modules load and process the forth experiment of
the chirpstack load testing scenario.
"""
from typing import Dict, List

import matplotlib.pyplot as plt
import numpy as np
import pandas as pd

TOTAL_MESSAGES = 1000

# name indicates column names
names: List[str] = ["device_id"] + [str(i) for i in range(0, 1000)]

# delivery_ratio for each examined rate.
# for each rate we have [min, mean, max]
delivery_ratio: Dict[str, np.ndarray] = {}
mean_latency: Dict[str, np.ndarray] = {}

for rate in ["s1", "ms500", "ms100"]:
    delivery_ratio_per_try: np.ndarray = np.zeros(10)
    mean_latency_per_try: np.ndarray = np.zeros(10)

    for t in range(0, 10):
        print(f"reading {rate} {t+1}")
        d = pd.DataFrame(
            pd.read_csv(
                f"{rate}_{t+1}.csv",
                header=None,
                names=names,
            )
        )

        # drop device_id if it exists
        if (x := d.drop("device_id", axis=1)) is not None:
            d = x

        delivery_ratio_per_device: np.ndarray = (
            d.notna().sum(axis=1) / TOTAL_MESSAGES * 100
        ).to_numpy()
        print("delivery ratio:")
        print(f"\t mean: {delivery_ratio_per_device.mean()}%")
        print(f"\t min: {delivery_ratio_per_device.min()}%")
        print(f"\t max: {delivery_ratio_per_device.max()}%")

        mean_latency_per_device: np.ndarray = np.nanmean(d.to_numpy(), axis=1)
        max_latency_per_device: np.ndarray = np.nanmax(d.to_numpy(), axis=1)
        min_latency_per_device: np.ndarray = np.nanmin(d.to_numpy(), axis=1)

        print("latency:")
        print(f"\t mean: {mean_latency_per_device.mean()}")
        print(f"\t min: {min_latency_per_device.min()}")
        print(f"\t max: {max_latency_per_device.max()}")

        print(f"\n{d.describe()}")

        delivery_ratio_per_try[t] = delivery_ratio_per_device.mean()
        mean_latency_per_try[t] = mean_latency_per_device.mean()

    delivery_ratio[rate] = np.array(
        [
            delivery_ratio_per_try.min(),
            delivery_ratio_per_try.mean(),
            delivery_ratio_per_try.max(),
        ]
    )

    mean_latency[rate] = np.array(
        [
            mean_latency_per_try.min(),
            mean_latency_per_try.mean(),
            mean_latency_per_try.max(),
        ]
    )

print(delivery_ratio)

fig, ax = plt.subplots(figsize=(10, 10))
ax.errorbar(
    x=[k for k in delivery_ratio.keys()],
    y=[v[1] for v in delivery_ratio.values()],
    fmt="g--",
    yerr=[
        [v[1] - v[0] for v in delivery_ratio.values()],
        [v[2] - v[1] for v in delivery_ratio.values()],
    ],
)
ax.set_title("Packet Delivery Ratio")
ax.set(ylabel="Delivery Ratio (%)", xlabel="Packet Rate (pps)")
fig.savefig("drop.png")

fig, ax = plt.subplots(figsize=(10, 10))
ax.errorbar(
    x=[k for k in mean_latency.keys()],
    y=[v[1] for v in mean_latency.values()],
    fmt="r--",
    yerr=[
        [v[1] - v[0] for v in mean_latency.values()],
        [v[2] - v[1] for v in mean_latency.values()],
    ],
)
ax.set_title("Latency")
ax.set(ylabel="Average Delay (s)", xlabel="Packet Rate (pps)")
fig.savefig("latency.png")
