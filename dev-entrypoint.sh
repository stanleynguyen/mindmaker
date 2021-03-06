ngrok http -log=stdout 8080  > ngrok.log &
until [ $(curl localhost:4040/api/tunnels | jq ".tunnels" | jq length) -gt 0 ]
do
  echo "Waiting for ngrok to start..."
  sleep 3
done
sed -i '/WEBHOOK_URL/ c\' .env
echo \ >> .env
echo "WEBHOOK_URL=$(curl http://localhost:4040/api/tunnels | jq ".tunnels[] | select (.proto == \"https\")" | jq ".public_url")" >> .env
fresh
