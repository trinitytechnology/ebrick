import http from 'k6/http';
import { sleep } from 'k6';

export let options = {
  vus: 1, // Number of virtual users
  iterations: 15, // Number of iterations
};

export default function () {
  let server_list = ["localhost:8080", "localhost:8082"]; // Add more servers
  let endpoint_list = ["/api/tenants"];
  let names = ["linkify", "example", "test"]; // Add more names if needed

  server_list.forEach(function(server) {
    endpoint_list.forEach(function(endpoint) {
      let randomName = names[Math.floor(Math.random() * names.length)];
      http.post(`http://${server}${endpoint}`, JSON.stringify({ name: randomName }), { headers: { "Content-Type": "application/json" } });
    });
  });

  sleep(0.5);
}
