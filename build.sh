#!/bin/bash

set -e

echo "🚀 Building OT Collaborative Editor..."

# Build with Podman
echo "📦 Building container image..."
podman build -t ot-collaborative-editor:latest .

echo "✅ Build complete!"
echo ""
echo "To run the container:"
echo "  podman run -d -p 8080:8080 --name ot-editor ot-collaborative-editor:latest"
echo ""
echo "Then open: http://localhost:8080"
