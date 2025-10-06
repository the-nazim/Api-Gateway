import requests
import time

# URL to send the GET requests to
url = "http://127.0.0.1:8080/user/1"

# Number of requests to send
num_requests = 100

for i in range(num_requests):
    try:
        response = requests.get(url)
        print(f"[{i+1}] Status Code: {response.status_code}")
    except requests.RequestException as e:
        print(f"[{i+1}] Request failed: {e}")