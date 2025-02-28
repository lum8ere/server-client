import unittest
from client.services.metrics_service import get_init_info

class TestMetricsService(unittest.TestCase):
    def test_get_init_info(self):
        data = get_init_info()
        # Проверяем, что возвращается словарь и есть ключ "metrics"
        self.assertIsInstance(data, dict)
        self.assertIn("metrics", data)
        
if __name__ == "__main__":
    unittest.main()
