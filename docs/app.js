/* tg-proxy-list — loads proxies.json and renders a searchable, sortable proxy list. */
(() => {
  "use strict";

  const $ = (sel) => document.querySelector(sel);
  const rowsEl = $("#rows");
  const loadingEl = $("#loading");
  const emptyEl = $("#empty");
  const shownEl = $("#shown");
  const searchEl = $("#search");
  const countryEl = $("#country");
  const sortEl = $("#sort");
  const hsOnlyEl = $("#hsonly");
  const toastEl = $("#toast");

  let ALL = [];

  const COUNTRY_NAMES = new Intl.DisplayNames(["en"], { type: "region" });

  function flag(cc) {
    if (!cc || cc.length !== 2 || cc === "??" || cc === "XX") return "🌐";
    const A = 0x1f1e6;
    return String.fromCodePoint(A + (cc.charCodeAt(0) - 65), A + (cc.charCodeAt(1) - 65));
  }
  function countryName(cc) {
    if (!cc || cc === "??" || cc === "XX") return "Unknown";
    try { return COUNTRY_NAMES.of(cc) || cc; } catch { return cc; }
  }
  function latClass(ms) {
    if (ms < 120) return "lat--fast";
    if (ms < 350) return "lat--mid";
    return "lat--slow";
  }
  function esc(s) {
    return String(s).replace(/[&<>"']/g, (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" }[c]));
  }
  function relTime(iso) {
    const t = Date.parse(iso);
    if (isNaN(t)) return "—";
    const mins = Math.round((Date.now() - t) / 60000);
    if (mins < 1) return "just now";
    if (mins < 60) return mins + "m ago";
    const h = Math.round(mins / 60);
    if (h < 24) return h + "h ago";
    return Math.round(h / 24) + "d ago";
  }

  function fillStats(data) {
    const hs = data.proxies.filter((p) => p.status === "handshake_ok").length;
    $('[data-stat="count"]').textContent = data.count.toLocaleString();
    $('[data-stat="handshake"]').textContent = hs.toLocaleString();
    $('[data-stat="countries"]').textContent = Object.keys(data.countries || {}).filter((c) => c !== "??").length;
    const upd = $('[data-stat="updated"]');
    upd.textContent = relTime(data.generated_at_utc);
    upd.title = data.generated_at_utc;
  }

  function fillCountryFilter(data) {
    const entries = Object.entries(data.countries || {})
      .filter(([cc]) => cc !== "??")
      .sort((a, b) => b[1] - a[1]);
    for (const [cc, n] of entries) {
      const o = document.createElement("option");
      o.value = cc;
      o.textContent = `${flag(cc)} ${countryName(cc)} (${n})`;
      countryEl.appendChild(o);
    }
  }

  function rowHTML(p) {
    const name = countryName(p.country);
    const typeLabel = { ee: "FakeTLS", dd: "Secure", plain: "Basic" }[p.type] || p.type;
    const statusLabel = p.status === "handshake_ok"
      ? '<span class="status-tag">● handshake</span>'
      : '<span class="status-tag status-tag--reach">● reachable</span>';
    return `<tr>
      <td class="col-country" data-label="Country"><span class="td-country"><span class="flag" aria-hidden="true">${flag(p.country)}</span><span>${esc(name)}</span><span class="cc">${esc(p.country)}</span></span></td>
      <td class="col-server" data-label="Server"><span class="server">${esc(p.server)}<span class="port">:${p.port}</span></span><br>${statusLabel}</td>
      <td class="col-type" data-label="Type"><span class="badge badge--${esc(p.type)}">${esc(typeLabel)}</span></td>
      <td class="col-num" data-label="Latency"><span class="lat ${latClass(p.latency_ms)}">${p.latency_ms} ms</span></td>
      <td class="col-num" data-label="Uptime"><span class="uptime">${p.uptime_pct ?? 0}%</span></td>
      <td class="col-actions" data-label="Connect"><div class="actions">
        <button class="btn btn--icon" type="button" data-copy="${esc(p.link)}" aria-label="Copy proxy link for ${esc(p.server)}" title="Copy link">
          <svg width="15" height="15" viewBox="0 0 24 24" aria-hidden="true"><path fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M9 9h10v10H9zM5 15H4a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1h10a1 1 0 0 1 1 1v1"/></svg>
        </button>
        <a class="btn btn--go" href="${esc(p.link)}" rel="noopener" aria-label="Connect to ${esc(p.server)} in Telegram">Connect</a>
      </div></td>
    </tr>`;
  }

  function applyFilters() {
    const q = searchEl.value.trim().toLowerCase();
    const cc = countryEl.value;
    const hsOnly = hsOnlyEl.checked;
    const sort = sortEl.value;

    let list = ALL.filter((p) => {
      if (cc && p.country !== cc) return false;
      if (hsOnly && p.status !== "handshake_ok") return false;
      if (q && !(p.server.toLowerCase().includes(q) || countryName(p.country).toLowerCase().includes(q) || (p.country || "").toLowerCase().includes(q))) return false;
      return true;
    });

    if (sort === "latency") list.sort((a, b) => a.latency_ms - b.latency_ms);
    else if (sort === "uptime") list.sort((a, b) => (b.uptime_pct ?? 0) - (a.uptime_pct ?? 0) || a.latency_ms - b.latency_ms);
    else if (sort === "country") list.sort((a, b) => countryName(a.country).localeCompare(countryName(b.country)) || a.latency_ms - b.latency_ms);

    rowsEl.innerHTML = list.map(rowHTML).join("");
    shownEl.textContent = list.length.toLocaleString();
    emptyEl.hidden = list.length !== 0;
  }

  function debounce(fn, ms) {
    let t;
    return (...a) => { clearTimeout(t); t = setTimeout(() => fn(...a), ms); };
  }

  function showToast(msg) {
    toastEl.textContent = msg;
    toastEl.hidden = false;
    requestAnimationFrame(() => toastEl.classList.add("show"));
    clearTimeout(showToast._t);
    showToast._t = setTimeout(() => {
      toastEl.classList.remove("show");
      setTimeout(() => (toastEl.hidden = true), 250);
    }, 1800);
  }

  rowsEl.addEventListener("click", async (e) => {
    const btn = e.target.closest("[data-copy]");
    if (!btn) return;
    const link = btn.getAttribute("data-copy");
    try {
      await navigator.clipboard.writeText(link);
      showToast("Proxy link copied");
    } catch {
      // Fallback for non-secure contexts
      const ta = document.createElement("textarea");
      ta.value = link; document.body.appendChild(ta); ta.select();
      try { document.execCommand("copy"); showToast("Proxy link copied"); } catch { showToast("Copy failed — long-press to copy"); }
      ta.remove();
    }
  });

  searchEl.addEventListener("input", debounce(applyFilters, 120));
  countryEl.addEventListener("change", applyFilters);
  sortEl.addEventListener("change", applyFilters);
  hsOnlyEl.addEventListener("change", applyFilters);

  fetch("proxies.json", { cache: "no-store" })
    .then((r) => { if (!r.ok) throw new Error("HTTP " + r.status); return r.json(); })
    .then((data) => {
      ALL = Array.isArray(data.proxies) ? data.proxies : [];
      loadingEl.hidden = true;
      fillStats(data);
      fillCountryFilter(data);
      applyFilters();
    })
    .catch((err) => {
      loadingEl.textContent = "Could not load the proxy list. Try the raw .txt or JSON links above.";
      console.error("proxies.json load failed:", err);
    });
})();
