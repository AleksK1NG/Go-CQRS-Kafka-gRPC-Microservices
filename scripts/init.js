use products

db.products.stats()

db.products.createIndex({ name: 1, description: 1 });
db.products.createIndex({ '$**': 'text' });

db.products.getIndexes();