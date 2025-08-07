import pandas as pd
from clickhouse_connect import get_client

# Run
'''
anirudhrv@MacBook-Pro zhuana-trading % kubectl port-forward svc/algonexus-clickhouse 8123:8123

Forwarding from 127.0.0.1:8123 -> 8123
Forwarding from [::1]:8123 -> 8123
Handling connection for 8123

'''

# Connect to ClickHouse
client = get_client(
    host='localhost',
    port=8123,
    username='default',
    password='password',
    database='algonexus',
)


# Drop the old table (optional during dev)
client.command("DROP TABLE IF EXISTS OHLC")

# Create OHLC table with symbol and market
client.command("""
CREATE TABLE OHLC (
    Date_Time DateTime('UTC'),
    Open Float64,
    High Float64,
    Low Float64,
    Close Float64,
    Volume UInt64,
    Day UInt8,
    Weekday UInt8,
    Week UInt8,
    Month UInt8,
    Year UInt16,
    Symbol String,
    Market String
) ENGINE = MergeTree()
ORDER BY (Symbol, Market, Date_Time)
""")

# Load CSV
df = pd.read_csv("spy.csv")  # Replace with actual path

# Convert 'Date' column to full DateTime
df['Date_Time'] = pd.to_datetime(df['Date'])

# Add constant values for Symbol and Market
df['Symbol'] = 'SPY'
df['Market'] = 'NYSEARCA'

# Drop original Date column
df = df.drop(columns=['Date'])

# Reorder columns to match table schema
df = df[[
    'Date_Time', 'Open', 'High', 'Low', 'Close',
    'Volume', 'Day', 'Weekday', 'Week', 'Month', 'Year',
    'Symbol', 'Market'
]]

# Insert data
client.insert_df('OHLC', df)

print("Data inserted successfully.")
