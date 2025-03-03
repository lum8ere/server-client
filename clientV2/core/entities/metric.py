from dataclasses import dataclass, field
from typing import List

@dataclass
class Metric:
    public_ip: str
    hostname: str
    os_info: str
    disk_total: int
    disk_used: int
    disk_free: int
    memory_total: int
    memory_used: int
    memory_available: int
    process_count: int

    def to_dict(self):
        return {
            "public_ip": self.public_ip,
            "hostname": self.hostname,
            "os": self.os_info,
            "disk_total": self.disk_total,
            "disk_used": self.disk_used,
            "disk_free": self.disk_free,
            "memory_total": self.memory_total,
            "memory_used": self.memory_used,
            "memory_available": self.memory_available,
            "process_count": self.process_count,
        }
