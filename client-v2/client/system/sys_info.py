import subprocess
import getpass
import platform
import psutil
import winreg

def current_user_has_password() -> bool:
    username = getpass.getuser()
    ps_command = f'Get-LocalUser -Name "{username}" | Select-Object -ExpandProperty PasswordRequired'
    command = ['powershell', '-Command', ps_command]
    try:
        output = subprocess.check_output(command, stderr=subprocess.STDOUT, text=True, timeout=10).strip().lower()
        return output == 'true'
    except Exception:
        return False

def get_min_password_length() -> int:
    try:
        output = subprocess.check_output(['net', 'accounts'], stderr=subprocess.STDOUT, text=True, timeout=10)
        for line in output.splitlines():
            if "Minimum password length" in line:
                return int(line.split()[-1])
        return -1
    except Exception:
        return -1

def get_running_processes():
    processes = []
    for proc in psutil.process_iter(attrs=['pid', 'name']):
        try:
            processes.append(proc.info)
        except (psutil.NoSuchProcess, psutil.AccessDenied):
            continue
    return processes

def get_running_services():
    services = []
    try:
        for svc in psutil.win_service_iter():
            try:
                info = svc.as_dict()
                services.append({
                    "name": info.get("name"),
                    "display_name": info.get("display_name"),
                    "status": info.get("status")
                })
            except psutil.NoSuchProcess:
                continue
    except Exception as e:
        services.append({"error": str(e)})
    return services

def create_vpn_connection():
    vpn_name = "MyVPN"
    server_address = "vpn.example.com"
    command = [
        "powershell",
        "-Command",
        f"Add-VpnConnection -Name \"{vpn_name}\" -ServerAddress \"{server_address}\" -TunnelType L2tp -Force -PassThru"
    ]
    try:
        output = subprocess.check_output(command, stderr=subprocess.STDOUT, text=True)
        print("VPN connection created successfully:", output)
    except Exception as e:
        print("Error creating VPN connection:", e)
