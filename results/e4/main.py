"""
this modules load and process the forth experiment of
the chirpstack load testing scenario.
"""
from typing import List

import numpy as np
import pandas as pd

# name indicates column names
names: List[str] = ["device_id"] + [str(i) for i in range(0, 1000)]

for rate in ["s1", "ms500", "ms100"]:
    for t in range(1, 10):
        print(f"reading {rate} {t}")
        d = pd.DataFrame(
            pd.read_csv(f"{rate}_{t}.csv", header=None, names=names)
        )

        # drop device_id if it exists
        if (x := d.drop("device_id", axis=1)) is not None:
            d = x

        delivery_ratio_per_device: np.ndarray = (
            d.notna().sum(axis=1) / 1000 * 100
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

        print(f"\n{d.head()}")
