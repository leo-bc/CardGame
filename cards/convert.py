import json

class Card:
    def __init__(self, title, hp):
        self.Title = title
        self.HP = hp
    
    def toString(self):
        return f'Pokemon: {self.Title}! [{self.HP}]'

f = open("cards/base.json", "r")
json_string = f.read()
original_cards = json.loads(json_string)
cards = []

for original_card in original_cards:
    if "mon" in original_card["supertype"]:
        cards.append(Card(original_card["name"], int(original_card["hp"])))

o = open("cards/converted.json", "w")
o.write(json.dumps(cards, default=lambda o: o.__dict__, indent=4))