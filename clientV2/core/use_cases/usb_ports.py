import ctypes
import winreg
from clientV2.core.services.logger_service import LoggerService

logger = LoggerService()

def is_admin():
    try:
        return ctypes.windll.shell32.IsUserAnAdmin()
    except Exception as e:
        logger.error(f"Error checking admin rights: {e}")
        return False

def enable_usb_ports():
    try:
        reg_key = winreg.OpenKey(
            winreg.HKEY_LOCAL_MACHINE,
            r"SYSTEM\CurrentControlSet\Services\USBSTOR",
            0,
            winreg.KEY_SET_VALUE
        )
        winreg.SetValueEx(reg_key, "Start", 0, winreg.REG_DWORD, 3)
        winreg.CloseKey(reg_key)
        logger.info("USB ports have been enabled.")
        return "USB ports enabled"
    except Exception as e:
        logger.error(f"Failed to enable USB ports: {e}")
        return f"Error enabling USB ports: {e}"

def disable_usb_ports():
    try:
        reg_key = winreg.OpenKey(
            winreg.HKEY_LOCAL_MACHINE,
            r"SYSTEM\CurrentControlSet\Services\USBSTOR",
            0,
            winreg.KEY_SET_VALUE
        )
        winreg.SetValueEx(reg_key, "Start", 0, winreg.REG_DWORD, 4)
        winreg.CloseKey(reg_key)
        logger.info("USB ports have been disabled.")
        return "USB ports disabled"
    except Exception as e:
        logger.error(f"Failed to disable USB ports: {e}")
        return f"Error disabling USB ports: {e}"
