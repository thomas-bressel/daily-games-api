#!/bin/bash

echo "🐳 Démarrage Docker Compose..."
sudo docker compose -f docker-compose.yml up -d --build --remove-orphans
sudo docker compose -f docker-compose.yml ps
echo "✅ Services démarrés !"
