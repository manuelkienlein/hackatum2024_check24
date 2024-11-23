def count_agr_actions(actions:list):
    result = []
    current_action = actions[0]
    count = 0

    for action in actions:
        if action == current_action:
            count += 1
        else:
            result.append(f"{count}x {current_action}")
            current_action = action
            count = 1
    result.append(f"{count}x {current_action}")  # Add the last group

    # Format the result
    formatted_result = ", ".join(result)
    return formatted_result

def compare_nested_structures(structure1, structure2, path=""):
    """
    Compares two nested structures (dictionaries or lists) and prints the differences.

    :param structure1: The first structure to compare.
    :param structure2: The second structure to compare.
    :param path: The path to the current key/element being compared.
    """
    if isinstance(structure1, dict) and isinstance(structure2, dict):
        all_keys = set(structure1.keys()).union(set(structure2.keys()))
        for key in all_keys:
            new_path = f"{path}.{key}" if path else str(key)
            if key not in structure1:
                print(f"Key '{new_path}' is missing in the expected result.")
            elif key not in structure2:
                print(f"Key '{new_path}' is missing in the actual result.")
            else:
                compare_nested_structures(structure1[key], structure2[key], new_path)
    elif isinstance(structure1, list) and isinstance(structure2, list):
        max_length = max(len(structure1), len(structure2))
        for i in range(max_length):
            new_path = f"{path}[{i}]"
            if i >= len(structure1):
                print(f"Index '{new_path}' is missing in the expected result.")
            elif i >= len(structure2):
                print(f"Index '{new_path}' is missing in the acutal result.")
            else:
                compare_nested_structures(structure1[i], structure2[i], new_path)
    else:
        if structure1 != structure2:
            print(f"Difference at '{path}': {structure1} != {structure2}")

