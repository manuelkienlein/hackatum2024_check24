
import sys
import json

from collections import Counter
from utils import count_agr_actions
input_file = ""


def main():
    input_file_path = sys.argv[1]
    tests = []
    with open(input_file_path, 'r', encoding="utf-8") as file:
        for line_number, line in enumerate(file, start=1):
            # Parse the JSON object
            data = json.loads(line)
            tests.append(data)

    request_types = [test["requestType"] for test in tests]
    print(count_agr_actions(request_types))

    actual_logs = [test["log"] for test in tests]
    last_error = ""
    last_error_count = 0
    for index, log in enumerate(actual_logs):
        if "search_error" in log:
            if log["search_error"] == last_error:
                last_error_count += 1
            else:
                if last_error_count > 0: 
                    print("+", last_error_count)
                    last_error_count = 0
                print(f"Search ERROR at request {index}", log["search_error"])
                last_error = log["search_error"]
        if "write_error" in log:
            if log["wirte_error"] == last_error:
                last_error_count += 1
            else:
                if last_error_count > 0:
                    print("+", last_error_count)
                    last_error_count = 0
                print(f"Write ERROR at request {index}", log["write_error"])
                last_error = log["write_error"]



    if last_error_count > 0:
        print("+", last_error_count)

if __name__=="__main__":
    main()