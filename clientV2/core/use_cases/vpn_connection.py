import subprocess
from clientV2.core.services.logger_service import LoggerService

logger = LoggerService()

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
        logger.info(f"VPN connection created successfully: {output}")
        return output
    except Exception as e:
        logger.error(f"Error creating VPN connection: {e}")
        return None
