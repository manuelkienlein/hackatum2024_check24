import json
import sys
from pathlib import Path
import os
from collections import defaultdict


class LogParser:
    def __init__(self, log_file):
        self.log_file = log_file
        self.push_entries = []
        self.read_entries = []
        self.error_entries = []
        self.current_state = defaultdict(list)
        self.compacted_logs = []

    def parse_logs(self):
        with open(self.log_file, 'r') as file:
            for line in file:
                try:
                    log_entry = json.loads(line.strip())
                    self._process_entry(log_entry)
                except json.JSONDecodeError:
                    print(f"Invalid JSON: {line.strip()}")

    def _process_entry(self, log_entry):
        request_type = log_entry.get("requestType")
        if request_type == "PUSH":
            self.push_entries.append(log_entry)
            self._process_push(log_entry)
        elif request_type == "READ":
            self.read_entries.append(log_entry)
        elif request_type == "COMPETITOR_ERROR":
            self.error_entries.append(log_entry)
        else:
            print(f"Unknown requestType: {request_type}")

    def _process_push(self, log_entry):
        offers = log_entry["log"]["write_config"]["Offers"]
        for offer in offers:
            self.current_state[offer["RegionID"]].append(offer)

    def compare_read_results(self):
        mismatches = []
        for entry in self.read_entries:
            expected = entry["log"]["expected_result"]
            actual = entry["log"]["actual_result"]
            if expected != actual:
                mismatches.append({
                    "id": entry["log"]["id"],
                    "expected": expected,
                    "actual": actual
                })
        return mismatches

    def print_summary(self):
        print("***** -----===== Log Analysis Summary: =====----- *****")
        print(f"Total PUSH entries: {len(self.push_entries)}")
        print(f"Total READ entries: {len(self.read_entries)}")
        print(f"Total ERROR entries: {len(self.error_entries)}\n")

        print("Summary of last 1000 entries:")
        print(", ".join(self.compacted_logs))

        mismatches = self.compare_read_results()
        print(f"\nMismatches found in READ entries: {len(mismatches)}")
        return # DEBUGGING: ONLY DO AVOID TERMINAL BUFFER OVERFLOW
        #for mismatch in mismatches:
        #    print(f"\nID: {mismatch['id']}")
        #    print(f"Expected: {json.dumps(mismatch['expected'], indent=2)}")
        #    print(f"Actual: {json.dumps(mismatch['actual'], indent=2)}\n")
        mismatch = mismatches[0]
        print(f"\nID: {mismatch['id']}")
        print(f"\033[92m Expected: {json.dumps(mismatch['expected'])}")
        print(f"\033[93m Actual: {json.dumps(mismatch['actual'])}\n")

        print("Current Database State:")
        for region, offers in self.current_state.items():
            print(f"Region {region}: {len(offers)} offers")
            for offer in offers:
                print(json.dumps(offer))

    def analyze_last_1000_entries(self):
        """Analyze the last 1000 entries in the log file."""
        with open(self.log_file, 'r') as file:
            lines = file.readlines()[-1000:]  # Get the last 1000 lines
            previous_type = None
            count = 0

            for line in lines:
                try:
                    log_entry = json.loads(line.strip())
                    request_type = log_entry.get("requestType", "UNKNOWN")

                    if request_type == previous_type:
                        count += 1
                    else:
                        if previous_type is not None:
                            # Append the compacted log entry for the previous type
                            self._add_compacted_entry(previous_type, count)
                        previous_type = request_type
                        count = 1

                except json.JSONDecodeError:
                    print(f"Invalid JSON: {line.strip()}")

            # Handle the last group
            if previous_type is not None:
                self._add_compacted_entry(previous_type, count)

    def _add_compacted_entry(self, request_type, count):
        """Add a compacted log entry to the list."""
        if count == 1:
            self.compacted_logs.append(request_type)
        else:
            self.compacted_logs.append(f"{request_type} ({count}x)")

# Usage

# ---Find latest log file in downloads folder
def get_latest_log_file(download_folder):
    # Suche nach der neuesten Datei im Download-Ordner mit der Endung '.log'
    log_files = [f for f in Path(download_folder).glob('*.log')]
    if not log_files:
        print("Keine Log-Dateien im Download-Ordner gefunden.")
        return None
    # Neueste Datei anhand des Änderungsdatums ermitteln
    latest_file = max(log_files, key=lambda f: f.stat().st_mtime)
    return latest_file


# Bestimme den Pfad zum Download-Ordner des Benutzers
download_folder = str(Path.home() / "Downloads")

# Finde die neueste Log-Datei im Download-Ordner
latest_log_file = get_latest_log_file(download_folder)

if latest_log_file:
    print(f"Neueste Log-Datei: {latest_log_file}")
    log_parser = LogParser(latest_log_file)
else:
    print("Keine Log-Datei gefunden.")

# ---Legacy via log file path argument
#log_parser = LogParser(sys.argv[1])
log_parser.analyze_last_1000_entries()
log_parser.parse_logs()
log_parser.print_summary()