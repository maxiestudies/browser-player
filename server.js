
const express = require('express');
const http = require('http');
const WebSocket = require('ws');

const app = express();
const server = http.createServer(app);
const wss = new WebSocket.Server({ server });

app.use(express.static('public')); // Serve static files from the "public" directory

wss.on('connection', (ws) => {
  console.log('Client connected');

  ws.on('message', (message) => {
    console.log(`Received: ${message}`);
  });

  // Example: send commands to the client at specific intervals
  setTimeout(() => ws.send(JSON.stringify({ action: 'startVideo' })), 5000); // 5 seconds
  setTimeout(() => ws.send(JSON.stringify({ action: 'stopVideo' })), 15000); // 10 seconds
  setTimeout(() => ws.send(JSON.stringify({ action: 'setImageSource', src: 'toxic-positivity-1.png' })), 15000); //
  setTimeout(() => ws.send(JSON.stringify({ action: 'showImage' })), 15000); // 15 seconds
  setTimeout(() => ws.send(JSON.stringify({ action: 'setVideoSource', src: 'file_example_MP4_640_3MG.mp4'})), 25000); // 25 seconds
  setTimeout(() => ws.send(JSON.stringify({ action: 'hideImage' })), 25010); //
  setTimeout(() => ws.send(JSON.stringify({ action: 'startVideo' })), 25010); //
  setTimeout(() => ws.send(JSON.stringify({ action: 'setImageSource', src: 'diselpunk-foosball-1.png' })), 27010); //
  setTimeout(() => ws.send(JSON.stringify({ action: 'stopVideo' })), 27010); //
  setTimeout(() => ws.send(JSON.stringify({ action: 'showImage'})), 27010); //

  ws.on('close', () => {
    console.log('Client disconnected');
  });
});

server.listen(8080, () => {
  console.log('Server is listening on port 8080');
});
