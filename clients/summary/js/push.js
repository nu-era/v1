const publicVapidKey =
  localStorage.getItem('pubKey');

if (Notification.permission === "default" || Notification.permission === "denied") {
    Notification.requestPermission();
}

// Check for service worker
if (Notification.permission === "granted" && "serviceWorker" in navigator) {
  send().catch((err) => console.error(err));
}

// Register SW, Register Push, Send Push
async function send() {
    // Register Service Worker
    const register = await navigator.serviceWorker.register("../js/service-worker.js");

    // Register Push
    const subscription = await register.pushManager.subscribe({
    userVisibleOnly: true,
    applicationServerKey: urlBase64ToUint8Array(publicVapidKey)
    });

    //register.active.postMessage(localStorage.getItem('device'))
    register.active.postMessage(JSON.stringify({
        'lat': localStorage.getItem('lat'),
        'long': localStorage.getItem('long')
    }))

    // Add user subscription to Back-End
    await fetch("https://api.bfranzen.me/subscribe", {
        method: "POST",
        body: JSON.stringify(subscription),
        headers: {
            "content-type": "application/json",
            "Authorization": localStorage.getItem('auth')
        }
    });
    console.log("User Registered for Push Notifications...");
}

function urlBase64ToUint8Array(base64String) {
    const padding = "=".repeat((4 - base64String.length % 4) % 4);
    const base64 = (base64String + padding)
    .replace(/\-/g, "+")
    .replace(/_/g, "/");

    const rawData = window.atob(base64);
    const outputArray = new Uint8Array(rawData.length);

    for (let i = 0; i < rawData.length; ++i) {
      outputArray[i] = rawData.charCodeAt(i);
    }
    return outputArray;
}