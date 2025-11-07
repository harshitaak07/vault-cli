const state = {
    authenticated: false,
    requiresPassword: false,
};

const els = {
    status: document.getElementById('status'),
    loginPanel: document.getElementById('login-panel'),
    loginForm: document.getElementById('login-form'),
    loginFeedback: document.getElementById('login-feedback'),
    uploadForm: document.getElementById('upload-form'),
    uploadFile: document.getElementById('upload-file'),
    uploadFeedback: document.getElementById('upload-feedback'),
    filesTableBody: document.querySelector('#files-table tbody'),
    secretsList: document.getElementById('secrets-list'),
    secretForm: document.getElementById('secret-form'),
    secretFeedback: document.getElementById('secret-feedback'),
    secretCat: document.getElementById('secret-category'),
    secretName: document.getElementById('secret-name'),
    secretValue: document.getElementById('secret-value'),
    refreshSecrets: document.getElementById('refresh-secrets'),
    auditList: document.getElementById('audit-list'),
    refreshAudit: document.getElementById('refresh-audit'),
    loginPassword: document.getElementById('password'),
    secretTemplate: document.getElementById('secret-item-template'),
};

async function api(path, options = {}) {
    try {
        const res = await fetch(path, Object.assign({
            headers: {
                'Accept': 'application/json',
            },
        }, options));

        if (res.status === 401) {
            handleAuthRequired();
            return null;
        }

        if (res.status === 204) {
            return null;
        }

        const text = await res.text();
        if (!text) {
            return null;
        }
        return JSON.parse(text);
    } catch (err) {
        console.error('api error', err);
        toast(els.status, 'API error — see console', true);
        return null;
    }
}

function handleAuthRequired() {
    state.authenticated = false;
    state.requiresPassword = true;
    els.loginPanel.hidden = false;
    els.loginFeedback.textContent = 'Session required. Enter master password.';
    toast(els.status, '🔒 Locked — login to continue', true);
}

function toast(el, message, warn = false) {
    if (!el) return;
    el.textContent = message;
    el.classList.toggle('warn', warn);
    if (message) {
        clearTimeout(el._timer);
        el._timer = setTimeout(() => {
            el.textContent = '';
            el.classList.remove('warn');
        }, 4000);
    }
}

async function loadAll() {
    const files = await api('/api/files');
    if (!files) {
        return;
    }
    state.authenticated = true;
    els.loginPanel.hidden = true;
    renderFiles(files);

    const [secrets, audits] = await Promise.all([
        api('/api/secrets'),
        api('/api/audit?limit=50'),
    ]);

    if (secrets) renderSecrets(secrets);
    if (audits) renderAudit(audits);
    toast(els.status, '✅ Synced just now');
}

function renderFiles(files) {
    els.filesTableBody.innerHTML = '';
    if (!files.length) {
        const row = document.createElement('tr');
        const cell = document.createElement('td');
        cell.colSpan = 6;
        cell.textContent = 'No files recorded yet.';
        row.appendChild(cell);
        els.filesTableBody.appendChild(row);
        return;
    }
    files.forEach((file) => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${escapeHtml(file.filename)}</td>
            <td>${formatBytes(file.size)}</td>
            <td>${escapeHtml(file.location || '')}</td>
            <td>${escapeHtml(file.mode || '')}</td>
            <td>${formatDate(file.uploaded_at)}</td>
            <td><button data-download="${encodeURIComponent(file.filename)}">Download</button></td>
        `;
        els.filesTableBody.appendChild(row);
    });
}

function renderSecrets(items) {
    const groups = {};
    items.forEach((item) => {
        if (!groups[item.Category || item.category]) {
            groups[item.Category || item.category] = [];
        }
        groups[item.Category || item.category].push(item);
    });
    els.secretsList.innerHTML = '';
    if (!items.length) {
        const empty = document.createElement('p');
        empty.textContent = 'No secrets stored yet.';
        empty.className = 'muted';
        els.secretsList.appendChild(empty);
        return;
    }
    Object.keys(groups).forEach((category) => {
        const wrapper = document.createElement('div');
        wrapper.className = 'secret-group';
        const title = document.createElement('h3');
        title.textContent = category || 'uncategorized';
        wrapper.appendChild(title);

        groups[category].forEach((item) => {
            const tpl = els.secretTemplate.content.cloneNode(true);
            tpl.querySelector('.name').textContent = item.Name || item.name;
            tpl.querySelector('.timestamp').textContent = item.UpdatedAt ? `Updated ${formatDate(item.UpdatedAt)}` : '';
            const viewBtn = tpl.querySelector('.view');
            const delBtn = tpl.querySelector('.delete');
            viewBtn.dataset.category = item.Category || item.category;
            viewBtn.dataset.name = item.Name || item.name;
            delBtn.dataset.category = item.Category || item.category;
            delBtn.dataset.name = item.Name || item.name;
            wrapper.appendChild(tpl);
        });

        els.secretsList.appendChild(wrapper);
    });
}

function renderAudit(items) {
    els.auditList.innerHTML = '';
    if (!items.length) {
        const li = document.createElement('li');
        li.textContent = 'No audit events yet.';
        els.auditList.appendChild(li);
        return;
    }
    items.forEach((a) => {
        const li = document.createElement('li');
        li.className = a.success ? 'success' : 'error';
        li.innerHTML = `
            <strong>${escapeHtml(a.action)}</strong> • ${escapeHtml(a.filename || '')} → ${escapeHtml(a.target || '')}<br>
            <small>${formatDate(a.timestamp || a.ts)}${a.error ? ' — ' + escapeHtml(a.error) : ''}</small>
        `;
        els.auditList.appendChild(li);
    });
}

function escapeHtml(str = '') {
    return str.replace(/[&<>"]+/g, (c) => ({
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
    }[c] || c));
}

function formatBytes(bytes) {
    if (!Number.isFinite(bytes)) return '—';
    if (bytes === 0) return '0 B';
    const units = ['B', 'KB', 'MB', 'GB'];
    const exponent = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1);
    const value = bytes / Math.pow(1024, exponent);
    return `${value.toFixed(value < 10 && exponent > 0 ? 1 : 0)} ${units[exponent]}`;
}

function formatDate(input) {
    if (!input) return '—';
    const date = new Date(input);
    if (Number.isNaN(date.getTime())) return input;
    return date.toLocaleString();
}

async function onLogin(event) {
    event.preventDefault();
    const password = els.loginPassword.value;
    if (!password) return;
    els.loginFeedback.textContent = 'Verifying...';
    const res = await api('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ password }),
    });
    if (!res || res.error) {
        els.loginFeedback.textContent = res && res.error ? res.error : 'Login failed';
        toast(els.status, 'Login failed', true);
        return;
    }
    els.loginFeedback.textContent = '';
    els.loginPassword.value = '';
    toast(els.status, 'Unlocked');
    await loadAll();
}

async function onUpload(event) {
    event.preventDefault();
    const file = els.uploadFile.files[0];
    if (!file) return;
    const form = new FormData();
    form.append('file', file);
    els.uploadFeedback.textContent = 'Uploading...';
    const res = await api('/api/upload', {
        method: 'POST',
        body: form,
    });
    if (res && !res.error) {
        els.uploadFeedback.textContent = '✅ Upload complete';
        els.uploadFile.value = '';
        await loadAll();
    } else {
        els.uploadFeedback.textContent = res && res.error ? res.error : 'Upload failed';
    }
}

async function onSecretSubmit(event) {
    event.preventDefault();
    const payload = {
        category: els.secretCat.value,
        name: els.secretName.value,
        value: els.secretValue.value,
    };
    els.secretFeedback.textContent = 'Saving...';
    const res = await api('/api/secrets', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
    });
    if (res && !res.error) {
        els.secretFeedback.textContent = '✅ Stored';
        els.secretForm.reset();
        await loadAll();
    } else {
        els.secretFeedback.textContent = res && res.error ? res.error : 'Failed to store secret';
    }
}

async function onSecretsClick(event) {
    const target = event.target;
    if (target.matches('button[data-download]')) {
        const name = target.dataset.download;
        window.location = `/api/download?name=${name}`;
        return;
    }
    if (target.classList.contains('view')) {
        const { category, name } = target.dataset;
        const res = await api(`/api/secrets/value?category=${encodeURIComponent(category)}&name=${encodeURIComponent(name)}`);
        if (res && res.value !== undefined) {
            navigator.clipboard?.writeText(res.value).catch(() => {});
            alert(`${category}/${name}:\n${res.value}`);
        }
    }
    if (target.classList.contains('delete')) {
        const { category, name } = target.dataset;
        if (!confirm(`Delete secret ${category}/${name}?`)) return;
        const res = await api(`/api/secrets?category=${encodeURIComponent(category)}&name=${encodeURIComponent(name)}`, {
            method: 'DELETE',
        });
        if (res && !res.error) {
            await loadAll();
        }
    }
}

document.addEventListener('DOMContentLoaded', () => {
    els.loginForm?.addEventListener('submit', onLogin);
    els.uploadForm?.addEventListener('submit', onUpload);
    els.secretForm?.addEventListener('submit', onSecretSubmit);
    document.body.addEventListener('click', onSecretsClick);
    els.refreshSecrets?.addEventListener('click', loadAll);
    els.refreshAudit?.addEventListener('click', loadAll);

    loadAll();
});


