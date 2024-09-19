shop_items = {
    "book": 12,
    "pen": 2
}

item = input("> ")

if item not in shop_items:
    print("Not Found")
    exit(0)
    
q = int(input("Q> "))

print(shop_items.get(item) * q)