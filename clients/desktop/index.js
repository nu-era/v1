const {app, BrowserWindow} = require('electron')

// this file should handle all system events i.e. handle window being
// closed, re-created

let win

function createWindow() {
     let win = new BrowserWindow({width: 800, height: 600})

     win.loadFile('index.html')

     win.webContents.openDevTools()

     win.on('closed', () => {
          win = null
     })
}

app.on('ready', createWindow)

app.on('window-all-closed', () => {
     if (process.platform !== 'darwin') {
           app.quit()
     }
})

app.on('activate', () => {
     if (win === null) {
          createWindow()
     }
})
