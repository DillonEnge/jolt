if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('./static/service-worker.js')
    .then(registration => {
      // Registration was successful
      console.log('ServiceWorker registration successful');
    }).catch(err => {
      // registration failed :(
      console.log(`ServiceWorker registration failed: ${err}`);
    });
  });
}
