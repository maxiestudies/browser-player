const video = document.getElementById('video');
const image = document.getElementById('image');
const buttonContainer = document.getElementById('buttons');

const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
  console.log('Connected to server');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  handleServerMessage(message);
};

function handleServerMessage(message) {
  switch (message.action) {
    case 'setVideoSource':
      setVideoSource(message.src);
      break;
    case 'startVideo':
      requestFullscreen(video);
      video.style.display = 'block';
      video.play();
      break;
    case 'stopVideo':
      video.pause();
      video.style.display = 'none';
      break;
    case 'showImage':
      image.style.display = 'block';
      break;
    case 'hideImage':
      image.style.display = 'none';
    case 'setImageSource':
      setImageSource(message.src);
      break;
    default:
      console.log('Unknown action:', message.action);
  }
}

function setVideoSource(src) {
  video.src = src;
  video.load(); // Reload the video with the new source
}

function setImageSource(src) {
  image.src = src;
}

function requestFullscreen(element) {
  if (element.mozRequestFullScreen) {
    element.mozRequestFullScreen();
  }
}

document.addEventListener('fullscreenchange', () => {
  if (!document.fullscreenElement) {
    console.log('Exited fullscreen mode');
  }
});

ws.onclose = () => {
  console.log('Disconnected from server');
};

// Get all buttons on the page
const buttons = document.querySelectorAll('button');

// Add a click event listener to each button
buttons.forEach(button => {
    button.addEventListener('click', (event) => {
        // Get the id of the clicked button
        const buttonId = event.target.id;

        // Print the id to the console
        console.log(`Last clicked button ID: ${buttonId}`);
        ws.send(buttonId);
        buttonContainer.style.display = 'none';
    });
});
