#!/usr/bin/env python3
"""Seed performance_metrics with sample data for UI verification."""

import math
import random
import subprocess
import sys
import time

MYSQL = ["mysql", "-h127.0.0.1", "-P3306", "-uroot", "-p123456", "im_db"]
NODES = ["node-a", "node-b"]
MINUTES = 120  # last 2 hours, 1 point per minute
INTERVAL_MS = 60 * 1000

METRIC_BASES = {
    "cpu.usage_percent": ("percent", 35, 25),
    "load.load1": ("load", 1.2, 0.8),
    "load.load5": ("load", 1.0, 0.5),
    "load.load15": ("load", 0.8, 0.3),
    "memory.total_bytes": ("fixed", 16 * 1024**3, 0),
    "memory.used_bytes": ("mem_used", 8 * 1024**3, 2 * 1024**3),
    "memory.free_bytes": ("mem_free", 4 * 1024**3, 1 * 1024**3),
    "memory.available_bytes": ("mem_avail", 6 * 1024**3, 1.5 * 1024**3),
    "memory.usage_percent": ("percent", 55, 15),
    "memory.swap_total_bytes": ("fixed", 2 * 1024**3, 0),
    "memory.swap_used_bytes": ("swap_used", 256 * 1024**2, 128 * 1024**2),
    "memory.swap_free_bytes": ("swap_free", 1792 * 1024**2, 128 * 1024**2),
    "memory.swap_usage_percent": ("percent", 12, 8),
    "disk.total_bytes": ("fixed", 500 * 1024**3, 0),
    "disk.used_bytes": ("disk_used", 280 * 1024**3, 40 * 1024**3),
    "disk.free_bytes": ("disk_free", 220 * 1024**3, 40 * 1024**3),
    "disk.usage_percent": ("percent", 56, 8),
    "go_runtime.goroutine_count": ("count", 320, 120),
    "go_runtime.gomaxprocs": ("fixed", 8, 0),
    "go_runtime.cgo_call_count": ("count", 42, 10),
    "go_runtime.alloc_bytes": ("bytes", 48 * 1024**2, 12 * 1024**2),
    "go_runtime.total_alloc_bytes": ("grow", 8 * 1024**3, 0),
    "go_runtime.sys_bytes": ("bytes", 96 * 1024**2, 8 * 1024**2),
    "go_runtime.heap_alloc_bytes": ("bytes", 40 * 1024**2, 10 * 1024**2),
    "go_runtime.heap_sys_bytes": ("bytes", 72 * 1024**2, 6 * 1024**2),
    "go_runtime.heap_inuse_bytes": ("bytes", 52 * 1024**2, 8 * 1024**2),
    "go_runtime.stack_inuse_bytes": ("bytes", 4 * 1024**2, 1 * 1024**2),
    "go_runtime.next_gc_bytes": ("bytes", 64 * 1024**2, 8 * 1024**2),
    "go_runtime.last_gc_time_unix_nano": ("grow_ts", 0, 0),
    "go_runtime.num_gc": ("grow_slow", 120, 0),
    "go_runtime.pause_total_ns": ("grow_slow", 50 * 1e6, 0),
}


def aligned_minute_ms(offset_minutes: int) -> int:
    now_ms = int(time.time() * 1000)
    base = now_ms // INTERVAL_MS * INTERVAL_MS
    return base - offset_minutes * INTERVAL_MS


def metric_value(metric_type: str, kind: str, base: float, amp: float, minute_idx: int, node_idx: int) -> float:
    phase = minute_idx / 18.0 + node_idx * 1.7
    wave = math.sin(phase) * amp
    noise = random.uniform(-amp * 0.08, amp * 0.08)

    if kind == "fixed":
        return base
    if kind == "percent":
        return max(0, min(100, base + wave + noise))
    if kind == "load":
        return max(0.01, base + wave * 0.5 + noise * 0.3)
    if kind == "count":
        return max(0, base + wave + noise)
    if kind == "bytes":
        return max(0, base + wave + noise)
    if kind == "mem_used":
        total = 16 * 1024**3
        pct = max(20, min(90, 55 + math.sin(phase) * 15))
        return total * pct / 100
    if kind == "mem_free":
        total = 16 * 1024**3
        used = total * max(20, min(90, 55 + math.sin(phase) * 15)) / 100
        return max(0, total - used - 2 * 1024**3)
    if kind == "mem_avail":
        total = 16 * 1024**3
        used = total * max(20, min(90, 55 + math.sin(phase) * 15)) / 100
        return max(0, total - used)
    if kind == "swap_used":
        total = 2 * 1024**3
        pct = max(0, min(80, 12 + math.sin(phase) * 8))
        return total * pct / 100
    if kind == "swap_free":
        total = 2 * 1024**3
        used = total * max(0, min(80, 12 + math.sin(phase) * 8)) / 100
        return total - used
    if kind == "disk_used":
        total = 500 * 1024**3
        pct = max(30, min(85, 56 + math.sin(phase) * 8))
        return total * pct / 100
    if kind == "disk_free":
        total = 500 * 1024**3
        used = total * max(30, min(85, 56 + math.sin(phase) * 8)) / 100
        return total - used
    if kind == "grow":
        return base + minute_idx * 5 * 1024**2 + node_idx * 1024**3
    if kind == "grow_slow":
        return base + minute_idx * (2 + node_idx) + node_idx * 10
    if kind == "grow_ts":
        return (aligned_minute_ms(MINUTES - 1 - minute_idx) - 300_000) * 1_000_000
    return base


def escape_sql(s: str) -> str:
    return s.replace("\\", "\\\\").replace("'", "''")


def main() -> int:
    random.seed(42)
    now_aligned = aligned_minute_ms(0)

    delete_sql = (
        "DELETE FROM performance_metrics WHERE node_name IN ('node-a', 'node-b');"
    )
    subprocess.run(MYSQL + ["-e", delete_sql], check=True)

    batch_size = 500
    values = []
    total = 0

    for node_idx, node in enumerate(NODES):
        for minute_idx in range(MINUTES):
            collect_time = aligned_minute_ms(MINUTES - 1 - minute_idx)
            for metric_type, (kind, base, amp) in METRIC_BASES.items():
                val = metric_value(metric_type, kind, base, amp, minute_idx, node_idx)
                values.append(
                    f"('{escape_sql(node)}',{collect_time},'{escape_sql(metric_type)}',{val:.6f})"
                )
                if len(values) >= batch_size:
                    sql = (
                        "INSERT INTO performance_metrics "
                        "(node_name, collect_time, metric_type, metric_value) VALUES "
                        + ",".join(values)
                        + ";"
                    )
                    subprocess.run(MYSQL + ["-e", sql], check=True)
                    total += len(values)
                    values.clear()

    if values:
        sql = (
            "INSERT INTO performance_metrics "
            "(node_name, collect_time, metric_type, metric_value) VALUES "
            + ",".join(values)
            + ";"
        )
        subprocess.run(MYSQL + ["-e", sql], check=True)
        total += len(values)

    verify = subprocess.run(
        MYSQL
        + [
            "-N",
            "-e",
            "SELECT COUNT(*), MIN(collect_time), MAX(collect_time) "
            "FROM performance_metrics WHERE node_name IN ('node-a','node-b');",
        ],
        check=True,
        capture_output=True,
        text=True,
    )
    print(f"Inserted {total} rows for nodes {NODES}.")
    print(f"Verify: count/min_ts/max_ts = {verify.stdout.strip()}")
    print(f"Latest aligned minute (ms): {now_aligned}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
