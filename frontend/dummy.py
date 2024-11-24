import json
import time

with open("dummy_response2.json", "r", encoding="utf-8") as file:
    response = json.load(file)
    

def get_dummy_response():
    time.sleep(2)
    return response