var conn = new WebSocket("ws://localhost:9292/connect");
conn.onmessage = function(m) { console.log("Received:", m.data); }
conn.onopen = function(m) {
  console.log("Connected.");
  conn.send(JSON.stringify({"action": "subscribe", "topic":"test"}))
}
conn.onerror = function (error) {
  console.log('Error ' + error);
};
conn.onclose = function (error) {
  console.log('Disconnected');
};

var conn = new WebSocket("ws://localhost:9292/connect");
conn.onmessage = function(e) { console.log("Received:", e.data); }
conn.onopen = function(m) {
  console.log("Connected.");
  conn.send(JSON.stringify({"action": "subscribe", "topic":"test"}))
  conn.send(JSON.stringify({"action": "publish", "topic":"test", "message": "hello"}))
}
conn.onerror = function (error) {
  console.log('Error ' + error);
};
conn.onclose = function (error) {
  console.log('Disconnected');
};

