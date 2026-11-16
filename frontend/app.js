const form = document.getElementById("search-form");
const input = document.getElementById("regex-input");
const regexError = document.getElementById("regex-error");
const status = document.getElementById("status");
const table = document.getElementById("logs-table");
const body = document.getElementById("logs-body");

function setStatus(text, isError) {
  status.textContent = text;
  status.className = isError ? "error" : "";
}

function renderLogs(logs) {
  body.replaceChildren();

  for (const log of logs) {
    const row = document.createElement("tr");

    const time = document.createElement("td");
    time.className = "timestamp";
    time.textContent = log.timestamp;

    const levelCell = document.createElement("td");
    const level = document.createElement("span");
    const levelName = String(log.level || "").toUpperCase();
    level.className = "level level-" + levelName.toLowerCase();
    level.textContent = levelName;
    levelCell.appendChild(level);

    const source = document.createElement("td");
    source.textContent = log.source || "";

    const message = document.createElement("td");
    message.className = "message";
    message.textContent = log.message;

    row.append(time, levelCell, source, message);
    body.appendChild(row);
  }
}

async function search() {
  const regex = input.value.trim();

  // Validate the pattern syntax before calling the API.
  // Note: JS regex syntax may differ slightly from the Go backend's (RE2).
  regexError.hidden = true;
  if (regex) {
    try {
      new RegExp(regex);
    } catch (err) {
      regexError.textContent = "Invalid regex: " + err.message;
      regexError.hidden = false;
      return;
    }
  }

  table.hidden = true;
  setStatus("Loading...", false);

  try {
    const logs = await fetchLogs(regex);
    if (logs.length === 0) {
      setStatus("No logs found.", false);
      return;
    }
    renderLogs(logs);
    table.hidden = false;
    setStatus(logs.length + " log(s)", false);
  } catch (err) {
    setStatus("Failed to load logs: " + err.message, true);
  }
}

form.addEventListener("submit", (event) => {
  event.preventDefault();
  search();
});

search();
