CREATE TABLE IF NOT EXISTS currency (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT,
    codein TEXT,
    name TEXT,
    high TEXT,
    low TEXT,
    varBid TEXT,
    pctChange TEXT,
    bid TEXT,
    ask TEXT,
    timestamp TEXT,
    create_date TEXT
)