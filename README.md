# AutoTyper for IMS (Go + JS)

**Automatic typing tool for Typewriter.ch**, optimized for the **Swiss German** keyboard layout.

---

## Installation

### 1. Run `AutoTyper.exe`

OR Build it yourself:

### 2. Open Typewriter.ch

and Open any typing exercise.

### 3. Inject the Browser Script

1. Open DevTools Console
Press: F12

2. Copy this script:

```
(function () {
    const panel = document.createElement("div");
    panel.style = `
        position:fixed;
        top:20px;
        right:20px;
        z-index:999999;
        background:rgba(0,0,0,0.85);
        color:#fff;
        padding:12px;
        border-radius:10px;
        font-family:Arial,sans-serif;
        font-size:13px;
    `;
    panel.innerHTML = `
        <b>AutoTyper (Go)</b><br>
        <br>
        Speed (ms): <input id="at_speed" type="number" value="70" min="20" max="500" style="width:60px"><br><br>
        <button id="at_start">Start</button>
        <button id="at_stop">Stop</button>
        <button id="at_close">Close</button>
    `;
    document.body.appendChild(panel);

    let intervalId = null;

    function getVisibleChar() {
        const box = document.querySelector("#text_todo_1");
        if (!box) return null;
        const span = box.querySelector("span");
        if (!span) return null;
        return span.textContent || null;
    }

    function sendChar() {
        const ch = getVisibleChar();
        if (!ch) return;
        fetch("http://localhost:9090/type", {
            method: "POST",
            body: ch
        }).catch(err => console.error("sendChar error", err));
    }

    document.getElementById("at_start").onclick = () => {
        const speed = parseInt(document.getElementById("at_speed").value) || 70;
        if (intervalId) clearInterval(intervalId);
        intervalId = setInterval(sendChar, speed);
        document.getElementById("at_status").innerText = "Running";
    };

    document.getElementById("at_stop").onclick = () => {
        if (intervalId) clearInterval(intervalId);
        intervalId = null;
        document.getElementById("at_status").innerText = "Stopped";
    };

    document.getElementById("at_close").onclick = () => {
        if (intervalId) clearInterval(intervalId);
        panel.remove();
    };

    console.log("AutoTyper (Go) client initialized");
})();

```

4. Paste into the browser console and press Enter

### As Result: You will see a small AutoTyper panel appear in the top-right corner.
