# Deployment Guide

## Prerequisites
- Podman (or Docker) installed
- Go 1.22+ (for local development)

## Building with Podman

### 1. Build the container image
```bash
podman build -t ot-collaborative-editor:latest .
```

### 2. Run the container
```bash
podman run -d -p 8080:8080 --name ot-editor ot-collaborative-editor:latest
```

### 3. Access the application
Open your browser to: http://localhost:8080

## Container Management

### View running containers
```bash
podman ps
```

### View logs
```bash
podman logs ot-editor
```

### Stop the container
```bash
podman stop ot-editor
```

### Remove the container
```bash
podman rm ot-editor
```

### Remove the image
```bash
podman rmi ot-collaborative-editor:latest
```

## Custom Port Configuration

Run on a different port (e.g., 3000):
```bash
podman run -d -p 3000:8080 --name ot-editor ot-collaborative-editor:latest
```

Or change the internal port via environment variable:
```bash
podman run -d -p 9000:9000 -e PORT=9000 --name ot-editor ot-collaborative-editor:latest
```

## Production Deployment Tips

1. **Use a reverse proxy** (nginx/traefik) for HTTPS
2. **Set resource limits**:
   ```bash
   podman run -d -p 8080:8080 \
     --memory="256m" \
     --cpus="0.5" \
     --name ot-editor \
     ot-collaborative-editor:latest
   ```

3. **Enable restart policy**:
   ```bash
   podman run -d -p 8080:8080 \
     --restart=unless-stopped \
     --name ot-editor \
     ot-collaborative-editor:latest
   ```

## Kubernetes/OpenShift Deployment

Create a deployment manifest:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ot-editor
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ot-editor
  template:
    metadata:
      labels:
        app: ot-editor
    spec:
      containers:
      - name: ot-editor
        image: ot-collaborative-editor:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
---
apiVersion: v1
kind: Service
metadata:
  name: ot-editor
spec:
  selector:
    app: ot-editor
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

## Local Development (without container)

```bash
go run .
```

Server will start on http://localhost:8080

## Building for Production (static binary)

```bash
CGO_ENABLED=0 GOOS=linux go build -o ot-collaborative-editor .
```

## Troubleshooting

### Port already in use
```bash
# Check what's using port 8080
sudo lsof -i :8080

# Or use a different port
podman run -d -p 8081:8080 --name ot-editor ot-collaborative-editor:latest
```

### Container won't start
```bash
# Check logs
podman logs ot-editor

# Run interactively for debugging
podman run -it --rm -p 8080:8080 ot-collaborative-editor:latest
```

### Build fails
```bash
# Clean build cache
podman build --no-cache -t ot-collaborative-editor:latest .
```
