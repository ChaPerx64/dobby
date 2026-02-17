import json
import urllib.request
import urllib.error
import argparse
import os
import sys


def upload_transactions(file_path, base_url, token, negate=False):
    if not os.path.exists(file_path):
        print(f"Error: File {file_path} not found.")
        sys.exit(1)

    try:
        with open(file_path, "r") as f:
            transactions = json.load(f)
    except Exception as e:
        print(f"Error reading {file_path}: {e}")
        sys.exit(1)

    url = f"{base_url.rstrip('/')}/transactions"

    success_count = 0
    error_count = 0

    for tx in transactions:
        try:
            # Ensure amount is an integer and negate if requested
            if "amount" in tx:
                tx["amount"] = int(tx["amount"])
                if negate:
                    tx["amount"] = -tx["amount"]

            data = json.dumps(tx).encode("utf-8")
            req = urllib.request.Request(url, data=data, method="POST")
            req.add_header("Authorization", f"Bearer {token}")
            req.add_header("Content-Type", "application/json")

            with urllib.request.urlopen(req) as response:
                if response.status in (200, 201):
                    success_count += 1
                    print(
                        f"Uploaded: {tx.get('description', 'No description')} - {tx.get('amount')} cents"
                    )
                else:
                    error_count += 1
                    print(
                        f"Failed to upload: {tx.get('description', 'No description')}. Status: {response.status}"
                    )
        except urllib.error.HTTPError as e:
            error_count += 1
            error_body = e.read().decode("utf-8")
            print(
                f"Failed to upload: {tx.get('description', 'No description')}. Status: {e.code}, Error: {error_body}"
            )
        except Exception as e:
            error_count += 1
            print(
                f"Exception occurred while uploading {tx.get('description', 'No description')}: {str(e)}"
            )

    print("\nUpload complete!")
    print(f"Successfully uploaded: {success_count}")
    print(f"Failed: {error_count}")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Upload transactions to Dobby API")
    parser.add_argument(
        "--file", default="transactions.json", help="Path to transactions JSON file"
    )
    parser.add_argument(
        "--url", default="http://localhost:8080", help="Base URL of the API"
    )
    parser.add_argument("--token", help="Bearer token for authentication")
    parser.add_argument(
        "--negate",
        action="store_true",
        help="Negate all amounts (e.g. if expenses are positive in the file)",
    )

    args = parser.parse_args()

    token = args.token or os.environ.get("DOBBY_TOKEN")

    if not token:
        print(
            "Error: Bearer token must be provided via --token or DOBBY_TOKEN environment variable."
        )
        sys.exit(1)

    upload_transactions(args.file, args.url, token, args.negate)
