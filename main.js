const { app, BrowserWindow } = require('electron');
const path = require('path');
const { execFile } = require('child_process');
const os = require('os');

function createWindow() {
    const mainWindow = new BrowserWindow({
        width: 800,
        height: 600,
        webPreferences: {
            nodeIntegration: true,
            contextIsolation: false,
            preload: path.join(__dirname, 'preload.js')
        }
    });

    mainWindow.loadFile('index.html');
}

function startGoServer() {
    const goExecutable = path.join(__dirname, 'main');
    const goCmd = os.platform() === 'win32' ? `${goExecutable}.exe` : goExecutable;

    console.log(`Attempting to run: ${goCmd}`);

    execFile(goCmd, (error, stdout, stderr) => {
        if (error) {
            console.error(`execFile error: ${error}`);
            return;
        }
        console.log(`stdout: ${stdout}`);
        console.error(`stderr: ${stderr}`);
    });
}

app.whenReady().then(() => {
    startGoServer();
    createWindow();

    app.on('activate', () => {
        if (BrowserWindow.getAllWindows().length === 0) {
            createWindow();
        }
    });
});

app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') {
        app.quit();
    }
});
