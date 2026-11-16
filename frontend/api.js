// All backend communication lives here.
//
// Request: GET {apiUrl}?regex=<pattern> (regex omitted when empty)
// Response: { "logs": [{ "id", "timestamp", "level", "source", "message" }], ... }

async function fetchLogs(regex) {
  const url = new URL(window.LOGSEEKER_CONFIG.apiUrl);
  if (regex) {
    url.searchParams.set("regex", regex);
  }

  const response = await fetch(url);
  if (!response.ok) {
    let detail = response.statusText;
    try {
      const error = await response.json();
      if (error && error.error) detail = error.error;
    } catch (_) {
      // Keep the HTTP status text when the server didn't return JSON.
    }
    throw new Error("API returned " + response.status + ": " + detail);
  }

  const data = await response.json();
  const logs = data && data.logs;
  if (!Array.isArray(logs)) {
    throw new Error("Unexpected API response format");
  }
  return logs;
}
