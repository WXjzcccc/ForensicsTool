import base64
import random
import plyvel

random.seed(114514)
i = random.random()
print(i)

db = plyvel.DB(r'C:\Users\admin\AppData\Roaming\tabby\Local Storage\leveldb', create_if_missing=False)

for key,value in db.iterator():
    print(key,value)
db.close()