#!/bin/bash

echo "🐳 Arrêt Docker Compose..."
sudo docker compose -f docker-compose.yml down

echo "✅ Services arrêtés !"
