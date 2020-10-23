import json
import re

class Card:
    def __init__(self, identity, rank, attacks):
        self.Identity = identity
        self.Rank = rank
        self.Attacks = attacks

class IdentityInfo:
    def __init__(self, title):
        self.Title = title

class RankInfo:
    def __init__(self, level, hp):
        self.Level = level
        self.HP = hp

class Attack:
    def __init__(self, name, cost, damage):
        self.Name = name
        self.Cost = cost
        self.Damage = damage
    
f = open("cards/base.json", "r", encoding="utf-8")
json_string = f.read()
original_cards = json.loads(json_string)

cards = []
for original_card in original_cards:
    if "mon" in original_card["supertype"]:
        title = original_card["name"]
        identity = IdentityInfo(title)

        level = int(original_card["level"])
        hp = int(original_card["hp"])
        rank = RankInfo(level, hp)

        attacks = []
        for original_attack in original_card["attacks"]:
            name = original_attack["name"]
            cost = original_attack["convertedEnergyCost"]
            damage = original_attack["damage"]
            if bool(re.search("^\\d+$", damage)):
                attack = Attack(name, cost, int(damage))
                attacks.append(attack)

        card = Card(identity, rank, attacks)
        cards.append(card)

cards = sorted(cards, key=lambda p: p.Rank.Level, reverse=True)

print(f"Writing {len(cards)} cards...")
o = open("cards/manipulatedPokemon.json", "w", encoding="utf-8")
o.write(json.dumps(cards, default=lambda o: o.__dict__, indent=4))