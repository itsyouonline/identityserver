__author__ = 'Dylan'
from dronedeliveryConsumer.droneDeliveryClient.client import Client

client = Client()
client.url = "http://127.0.0.1:5000"

print(client.deliveries_get().json())
