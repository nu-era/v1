console.log("Service Worker Loaded...");

self.addEventListener("push", e => {
  const why = e.data.json();
  setTimeout(function(why) {
    console.log("Push Recieved...");
    self.registration.showNotification(why.title, {
      body: "Notified by Traversy Media!",
      icon: "http://image.ibb.co/frYOFd/tmlogo.png"
    });
  }, 3000);
});
