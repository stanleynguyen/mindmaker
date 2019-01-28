ngrok http -log=stdout 8080  > ngrok.log &
echo \ >> .env
echo "WEBHOOK_URL=$(curl http://localhost:4040/api/tunnels | jq ".tunnels[] | select (.proto == \"https\")" | jq ".public_url")" >> .env
fresh
