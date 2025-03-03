from dataclasses import dataclass

@dataclass
class InstalledApplication:
    name: str
    version: str
    app_type: str

    def to_dict(self):
        return {
            "name": self.name,
            "version": self.version,
            "app_type": self.app_type,
        }