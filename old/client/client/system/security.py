import ctypes
import winreg

def is_admin():
    try:
        return ctypes.windll.shell32.IsUserAnAdmin()
    except Exception:
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
        print("USB ports have been enabled.")
    except Exception as e:
        print("Failed to enable USB ports:", e)

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
        print("USB ports have been disabled.")
    except Exception as e:
        print("Failed to disable USB ports:", e)
