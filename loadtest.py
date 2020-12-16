import time 
from locust import HttpUser, task 

class QuickstartUser(HttpUser):
    @task
#   def access_model(self):
#self.client.post("http://dominio/url")
#    self.client.post("/prueba"\
#    ,json={"data":1})
#self.client.get("http://dominio/url")
    def on_start(self):
        self.client.get("/")