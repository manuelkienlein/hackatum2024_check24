import json

DEEPNESS = 4

with open("regions.json", "r", encoding="utf-8") as file:
    read_json:dict = json.load(file)


next_dicts = [(read_json, 0)] 
id_name_dict = {}


while len(next_dicts) > 0:
    # print(next_dicts)
    dic, deepness = next_dicts[0]
    if deepness >= DEEPNESS:
        next_dicts.pop(0)
        continue
    id_name_dict[dic["name"]] = dic["id"]
    if "subregions" in dic:
        subregions = dic["subregions"]
        for sub in subregions:
            next_dicts.append((sub, deepness+1))
    next_dicts.pop(0)


print(json.dumps(id_name_dict))