"""
this modules load and process the sixth experiment of
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
# for each number of gateways we have [min, mean, max]
delivery_ratio: Dict[str, np.ndarray] = {}
mean_latency: Dict[str, np.ndarray] = {}
p90_latency: Dict[str, np.ndarray] = {}

for number_of_gateways in ["g1", "g2", "g5", "g10"]:
    delivery_ratio_per_try: np.ndarray = np.zeros(10)
    mean_latency_per_try: np.ndarray = np.zeros(10)
    p90_latency_per_try: np.ndarray = np.zeros(10)

    for t in range(0, 10):
        print(f"reading {number_of_gateways} {t+1}")
        d = pd.DataFrame(
            pd.read_csv(
                f"{number_of_gateways}_{t+1}.csv",
                header=None,
                names=names,
                engine="python",
                on_bad_lines=lambda bad_line: bad_line[:TOTAL_MESSAGES]
                if len(bad_line) > TOTAL_MESSAGES
                else None,
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
        p90_latency_per_device: np.ndarray = np.nanquantile(
            d.to_numpy(), axis=1, q=0.90
        )

        print("latency:")
        print(f"\t mean: {mean_latency_per_device.mean()}")
        print(f"\t min: {min_latency_per_device.min()}")
        print(f"\t max: {max_latency_per_device.max()}")

        print(f"\n{d.describe()}")

        delivery_ratio_per_try[t] = delivery_ratio_per_device.mean()
        mean_latency_per_try[t] = mean_latency_per_device.mean()
        p90_latency_per_try[t] = p90_latency_per_device.mean()

    delivery_ratio[number_of_gateways] = np.array(
        [
            delivery_ratio_per_try.min(),
            delivery_ratio_per_try.mean(),
            delivery_ratio_per_try.max(),
        ]
    )

    mean_latency[number_of_gateways] = np.array(
        [
            mean_latency_per_try.min(),
            mean_latency_per_try.mean(),
            mean_latency_per_try.max(),
        ]
    )

    p90_latency[number_of_gateways] = np.array(
        [
            p90_latency_per_try.min(),
            p90_latency_per_try.mean(),
            p90_latency_per_try.max(),
        ]
    )

fig, axs = plt.subplots(figsize=(10, 10), ncols=1, nrows=3)
axs[0].errorbar(
    x=list(delivery_ratio.keys()),
    y=[v[1] for v in delivery_ratio.values()],
    fmt="g--",
    yerr=[
        [v[1] - v[0] for v in delivery_ratio.values()],
        [v[2] - v[1] for v in delivery_ratio.values()],
    ],
)
axs[0].set_title("Packet Delivery Ratio")
axs[0].set(ylabel="Delivery Ratio (%)", xlabel="Number of Gateways")

axs[1].errorbar(
    x=list(mean_latency.keys()),
    y=[v[1] for v in mean_latency.values()],
    fmt="r--",
    yerr=[
        [v[1] - v[0] for v in mean_latency.values()],
        [v[2] - v[1] for v in mean_latency.values()],
    ],
)
axs[1].set_title("Latency")
axs[1].set(ylabel="Average Delay (s)", xlabel="Number of Gateways")

axs[2].errorbar(
    x=list(p90_latency.keys()),
    y=[v[1] for v in p90_latency.values()],
    fmt="b--",
    yerr=[
        [v[1] - v[0] for v in p90_latency.values()],
        [v[2] - v[1] for v in p90_latency.values()],
    ],
)
axs[2].set_title("P90 Latency")
axs[2].set(ylabel="P90 Delay (s)", xlabel="Number of Gateways")

fig.tight_layout()
fig.savefig('e6.png')
