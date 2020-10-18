import json
import re

class MoveSet:
    def __init__(self, level, hp, attacks):
        self.Level = level
        self.HP = hp
        self.Attacks = attacks

class Attack:
    def __init__(self, name, cost, damage):
        self.Name = name
        self.Cost = cost
        self.Damage = damage

f = open("cards/base.json", "r", encoding="utf-8")
json_string = f.read()
original_cards = json.loads(json_string)

move_sets = []
for original_card in original_cards:
    if "mon" in original_card["supertype"]:
        cancelled = False
        attacks = []
        for original_attack in original_card["attacks"]:
            name = original_attack["name"]
            cost = original_attack["convertedEnergyCost"]
            damage = original_attack["damage"]
            if bool(re.search("^\\d+$", damage)):
                attack = Attack(name, cost, int(damage))
                attacks.append(attack)

        if len(attacks) > 0:
            level = int(original_card["level"])
            hp = int(original_card["hp"])
            move_set = MoveSet(level, hp, attacks)
            move_sets.append(move_set)

move_sets = sorted(move_sets, key=lambda p: p.Level, reverse=True)


# cards = []

# for original_card in original_cards:
#     if "mon" in original_card["supertype"]:
#         cards.append(Card(original_card["name"], int(original_card["hp"])))

o = open("cards/manipulatedPokemon.json", "w", encoding="utf-8")
o.write(json.dumps(move_sets, default=lambda o: o.__dict__, indent=4))