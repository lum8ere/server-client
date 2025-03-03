import winreg
from clientV2.core.entities.installedApplication import InstalledApplication

def get_installed_applications_windows():
    """
    Получает список установленных приложений на Windows с информацией о версии и типе.
    Возвращает список объектов InstalledApplication.
    """
    apps = []
    uninstall_paths = [
        (winreg.HKEY_LOCAL_MACHINE, r"Software\Microsoft\Windows\CurrentVersion\Uninstall"),
        (winreg.HKEY_LOCAL_MACHINE, r"Software\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall"),
        (winreg.HKEY_CURRENT_USER, r"Software\Microsoft\Windows\CurrentVersion\Uninstall"),
    ]
    
    for root, path in uninstall_paths:
        try:
            key = winreg.OpenKey(root, path)
        except Exception:
            continue

        for i in range(winreg.QueryInfoKey(key)[0]):
            try:
                subkey_name = winreg.EnumKey(key, i)
                subkey = winreg.OpenKey(key, subkey_name)
                try:
                    display_name, _ = winreg.QueryValueEx(subkey, "DisplayName")
                    if not display_name:
                        continue
                    name = display_name.strip().replace("\n", " ")

                    # Пытаемся получить версию
                    version = ""
                    try:
                        ver_val, _ = winreg.QueryValueEx(subkey, "DisplayVersion")
                        version = ver_val.strip()
                    except Exception:
                        pass

                    # Определяем тип на основе значения DisplayIcon (если есть)
                    app_type = "Unknown"
                    try:
                        display_icon, _ = winreg.QueryValueEx(subkey, "DisplayIcon")
                        if display_icon:
                            if display_icon.lower().endswith(".exe"):
                                app_type = "Executable"
                            else:
                                app_type = "Other"
                    except Exception:
                        pass

                    apps.append(InstalledApplication(name=name, version=version, app_type=app_type))
                except Exception:
                    pass
                winreg.CloseKey(subkey)
            except Exception:
                continue

        winreg.CloseKey(key)
    return apps
