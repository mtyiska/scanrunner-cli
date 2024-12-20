#!/bin/bash
set -e

echo "Starting installation of dependencies and tools..."

# Verify Go is installed
command -v go >/dev/null 2>&1 || { echo "Error: Go is not installed. Please install Go and try again."; exit 1; }

echo "Installing Go libraries..."
# Ensure Go module is initialized
if [ ! -f "go.mod" ]; then
    echo "Initializing Go module..."
    go mod init scanrunner || echo "Go module already initialized."
fi

# Add necessary libraries to go.mod
go get k8s.io/apimachinery@latest
go get k8s.io/client-go@latest
go get github.com/aquasecurity/trivy-operator@latest
go get github.com/sirupsen/logrus@latest
go get github.com/spf13/cobra@latest
go get k8s.io/cli-runtime@latest
go get github.com/moby/buildkit@v0.11.5

# Tidy up dependencies
go mod tidy

echo "Go libraries installed successfully."

echo "Installing third-party tools..."

# Trivy
echo "Installing Trivy..."
curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh
if [ -f "./bin/trivy" ]; then
    sudo mv ./bin/trivy /usr/local/bin/
fi
if ! command -v trivy >/dev/null 2>&1; then
    echo "Error: Trivy installation failed. Please check your environment."
    exit 1
fi
echo "Trivy installed successfully."

# Hadolint
echo "Installing Hadolint..."
curl -LO https://github.com/hadolint/hadolint/releases/latest/download/hadolint-$(uname -s)-$(uname -m)
chmod +x hadolint-$(uname -s)-$(uname -m)
sudo mv hadolint-$(uname -s)-$(uname -m) /usr/local/bin/hadolint
if ! command -v hadolint >/dev/null 2>&1; then
    echo "Error: Hadolint installation failed. Please check your environment."
    exit 1
fi
echo "Hadolint installed successfully."

# Actionlint
echo "Installing Actionlint..."
curl -LO https://github.com/rhysd/actionlint/releases/latest/download/actionlint-$(uname -s)-$(uname -m)
chmod +x actionlint-$(uname -s)-$(uname -m)
sudo mv actionlint-$(uname -s)-$(uname -m) /usr/local/bin/actionlint
if ! command -v actionlint >/dev/null 2>&1; then
    echo "Error: Actionlint installation failed. Please check your environment."
    exit 1
fi
echo "Actionlint installed successfully."

# Kubectl
echo "Installing Kubectl..."
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl
sudo mv kubectl /usr/local/bin/
if ! command -v kubectl >/dev/null 2>&1; then
    echo "Error: Kubectl installation failed. Please check your environment."
    exit 1
fi
echo "Kubectl installed successfully."

echo "All tools and libraries have been installed successfully!"
