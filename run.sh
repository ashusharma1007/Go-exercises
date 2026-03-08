#!/bin/bash

set -e

# Stop and remove existing container if it exists
if podman ps -a | grep -q ot-editor; then
    echo "🛑 Stopping existing container..."
    podman stop ot-editor 2>/dev/null || true
    podman rm ot-editor 2>/dev/null || true
fi

# Run the container
echo "🚀 Starting OT Collaborative Editor..."
podman run -d -p 8080:8080 --name ot-editor ot-collaborative-editor:latest

echo "✅ Container started!"
echo ""
echo "📝 Access the editor at: http://localhost:8080"
echo ""
echo "Useful commands:"
echo "  podman logs ot-editor          # View logs"
echo "  podman logs -f ot-editor       # Follow logs"
echo "  podman stop ot-editor          # Stop container"
echo "  podman restart ot-editor       # Restart container"
