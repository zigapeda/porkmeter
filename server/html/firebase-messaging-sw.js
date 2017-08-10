self.addEventListener('push', function(event) {
    console.log('[Service Worker] Push Received.');
    console.log(`[Service Worker] Push had this data:`);
    console.log(event);

    const title = 'Porkmeter';
    const options = {
        body: 'Temperaturen prüfen!',
        icon: 'appico.png',
        badge: 'badge.png'
    };

    event.waitUntil(self.registration.showNotification(title, options));
});

self.addEventListener('notificationclick', function(event) {
    console.log('[Service Worker] Notification click Received.');

    event.notification.close();

    event.waitUntil(
        clients.openWindow("https://porkmeter.maplpapl.de")
    );
});