

# load .env
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [ -f "$SCRIPT_DIR/.env" ]; then
    echo "Loading variables from $SCRIPT_DIR/.env"
    source "$SCRIPT_DIR/.env"
else
    echo "Error: .env file not found in $SCRIPT_DIR"
    exit 1
fi

# failsafe
set -e

echo "---------------------------------------------------"
echo "Running setup-airbyte.sh :)"
echo "UI WILL BE ON PORT $AIRBYTE_HOST:$AIRBYTE_WEBAPP_PORT"

if ! command -v abctl &> /dev/null; then
    echo "abctl not found. Installing abctl now..."
    curl -LsfS "$AIRBYTE_INSTALL_URL" | bash -
    echo "abctl installed successfully."
else
    echo "abctl is already installed. Skipping installation."
fi

echo "Cleaning up old installation..."
abctl local uninstall --verbose || echo "No installation to uninstall..."

echo "Installing abctl..."
abctl local install \
  --verbose \
  --port="$AIRBYTE_WEBAPP_PORT" \
  --no-browser \
  --insecure-cookies

echo "Setting up credentials..."
abctl local credentials \
  --verbose \
  --email "$AIRBYTE_EMAIL" \
  --password "$AIRBYTE_PASSWORD"

echo "setup-airbyte.sh complete :)"
echo "---------------------------------------------------"
