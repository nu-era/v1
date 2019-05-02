const {app, BrowserWindow, protocol} = require('electron')
const path = require('path')
const PROTOCOL = 'newEra://'
const PROTOCOL_PREFIX = 'newEra'
const BASE = "https://api.bfranzen.me/"
// this file should handle all system events i.e. handle window being
// closed, re-created

let win

function formatURL(inPath) {
    return `${base}${inPath}`
} 

function getPath(inUrl) {
    if (!inUrl) {
        return
    } else if (!inUrl.startsWith(PROTOCOL)) {
        return
    } else { // return everything after newera://
        return inUrl.substr(9)
    }
}

function createWindow() {
     let win = new BrowserWindow({width: 800, height: 600})

     protocol.registerHttpProtocol(PROTOCOL_PREFIX, (req, cb) => {
        let url = req.url
        let tmp = getPath(url)
        const fullUrl = formatURL(tmp)
        win.loadURL(fullUrl)
     }, (err) => {
         if (!err) {
             console.log("protocol registered...")
         } else {
             console.log("error generating protocol, " + err)
         }
     })

     win.loadFile('index.html')

     //win.webContents.openDevTools()

     win.on('closed', () => {
          win = null
     })
}

app.on('ready', () => {
    createWindow()
});

app.on('window-all-closed', () => {
     if (process.platform !== 'darwin') {
           app.quit()
     }
});

app.on('activate', () => {
     if (win === null) {
          createWindow()
     }
});
