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