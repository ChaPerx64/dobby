import json
import http.client
import os

# Configuration
API_HOST = "api.dobby.chaianpar.dev"
API_TOKEN = "paste-actual-token"
ENVELOPE_ID = "2aeccd10-20f8-484f-80ab-fe920666e1da"
INPUT_FILE = "scripts/missing-transactions/missing_transactions.json"

def create_transaction(data):
    conn = http.client.HTTPSConnection(API_HOST)
    
    headers = {
        'Authorization': f'Bearer {API_TOKEN}',
        'Content-Type': 'application/json',
    }
    
    payload = {
        "envelopeId": ENVELOPE_ID,
        "amount": data["amount"],
        "description": data["description"],
        "date": data["date"],
        "category": "Food" # Defaulting to Food as per common usage in existing-transactions.json
    }
    
    # Simple heuristic for category
    desc = data["description"].upper()
    if "APOTEKA" in desc:
        payload["category"] = "Necessary"
    elif "BEOGRAD 411" in desc or "STAMPA" in desc:
        payload["category"] = "Necessary"
    elif "IMMANUEL" in desc:
        payload["category"] = "Education" # Pure guess
        
    json_payload = json.dumps(payload)
    
    conn.request("POST", "/transactions", json_payload, headers)
    response = conn.getresponse()
    
    print(f"Creating: {data['description']} ({data['amount']}) on {data['date']}")
    print(f"Response: {response.status} {response.reason}")
    
    if response.status >= 400:
        print(f"Error: {response.read().decode()}")
    else:
        print(f"Success: {response.read().decode()}")
    
    conn.close()

def main():
    if not os.path.exists(INPUT_FILE):
        print(f"File {INPUT_FILE} not found.")
        return

    with open(INPUT_FILE, 'r') as f:
        transactions = json.load(f)

    for tx in transactions:
        create_transaction(tx)

if __name__ == "__main__":
    main()
