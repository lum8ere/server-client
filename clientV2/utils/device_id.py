import os
import uuid

def get_device_id():
    """
    Возвращает сохранённый идентификатор устройства, или генерирует новый, если его нет.
    Сохраняет идентификатор в файле в домашней директории пользователя.
    """
    device_id_file = os.path.join(os.path.expanduser("~"), ".clientv2_device_id")
    if os.path.exists(device_id_file):
        with open(device_id_file, "r") as f:
            device_id = f.read().strip()
            if device_id:
                return device_id

    device_id = str(uuid.uuid4())
    with open(device_id_file, "w") as f:
        f.write(device_id)
    return device_id

if __name__ == "__main__":
    print(f"Device ID: {get_device_id()}")
