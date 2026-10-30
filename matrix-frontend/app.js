(function () {
  const apiUrl = (window.APP_CONFIG && window.APP_CONFIG.apiUrl) || '';
  const storageKey = 'qr_token';

  const loginSection = document.getElementById('login-section');
  const matrixSection = document.getElementById('matrix-section');
  const results = document.getElementById('results');
  const message = document.getElementById('message');
  const sessionStatus = document.getElementById('session-status');

  let token = sessionStorage.getItem(storageKey);

  function setMessage(text, isError) {
    message.textContent = text || '';
    message.className = isError ? 'error' : '';
  }

  function render() {
    loginSection.hidden = Boolean(token);
    matrixSection.hidden = !token;
    sessionStatus.textContent = token ? 'Sesión iniciada' : 'Sin sesión';
    if (!token) results.hidden = true;
  }

  function parseMatrix(text) {
    const trimmed = text.trim();
    if (trimmed.startsWith('[')) {
      return JSON.parse(trimmed);
    }
    return trimmed
      .split(/\r?\n/)
      .map((line) => line.trim())
      .filter((line) => line.length > 0)
      .map((line) => line.split(/[\s,;]+/).map(Number));
  }

  function format(value) {
    if (Number.isInteger(value)) return String(value);
    return value.toFixed(4);
  }

  function matrixTable(data) {
    const table = document.createElement('table');
    table.className = 'matrix';
    data.forEach((row) => {
      const tr = table.insertRow();
      row.forEach((value) => {
        tr.insertCell().textContent = format(value);
      });
    });
    return table;
  }

  function renderResults(data) {
    document.getElementById('q-table').replaceChildren(matrixTable(data.q));
    document.getElementById('r-table').replaceChildren(matrixTable(data.r));

    const s = data.statistics;
    const diagonal = s.matrices.filter((m) => m.isDiagonal).map((m) => m.name);
    const rows = [
      ['Valor máximo', format(s.max)],
      ['Valor mínimo', format(s.min)],
      ['Promedio', format(s.average)],
      ['Suma total', format(s.sum)],
      ['Cantidad de valores', String(s.count)],
      ['Alguna matriz diagonal', s.anyDiagonal ? 'Sí (' + diagonal.join(', ') + ')' : 'No'],
    ];

    const body = document.getElementById('stats-body');
    body.replaceChildren();
    rows.forEach(([label, value]) => {
      const tr = body.insertRow();
      const th = document.createElement('th');
      th.textContent = label;
      tr.appendChild(th);
      tr.insertCell().textContent = value;
    });

    results.hidden = false;
  }

  async function request(path, body) {
    const headers = { 'Content-Type': 'application/json' };
    if (token) headers.Authorization = 'Bearer ' + token;

    const res = await fetch(apiUrl + path, {
      method: 'POST',
      headers,
      body: JSON.stringify(body),
    });
    const data = await res.json().catch(() => ({}));
    if (!res.ok) {
      const error = new Error(data.error || 'Error ' + res.status);
      error.status = res.status;
      throw error;
    }
    return data;
  }

  document.getElementById('login-form').addEventListener('submit', async (event) => {
    event.preventDefault();
    const form = event.target;
    const button = form.querySelector('button');
    button.disabled = true;
    setMessage('');

    try {
      const data = await request('/api/auth/login', {
        username: form.username.value,
        password: form.password.value,
      });
      token = data.token;
      sessionStorage.setItem(storageKey, token);
      form.reset();
      render();
    } catch (err) {
      setMessage(err.status === 401 ? 'Usuario o contraseña incorrectos.' : err.message, true);
    } finally {
      button.disabled = false;
    }
  });

  document.getElementById('matrix-form').addEventListener('submit', async (event) => {
    event.preventDefault();
    const form = event.target;
    const button = form.querySelector('button[type="submit"]');
    setMessage('');

    let matrix;
    try {
      matrix = parseMatrix(form.matrix.value);
    } catch {
      setMessage('No se pudo leer la matriz.', true);
      return;
    }

    button.disabled = true;
    try {
      renderResults(await request('/api/qr', { matrix }));
    } catch (err) {
      if (err.status === 401) {
        token = null;
        sessionStorage.removeItem(storageKey);
        render();
        setMessage('La sesión expiró. Vuelve a ingresar.', true);
      } else {
        setMessage(err.message, true);
      }
    } finally {
      button.disabled = false;
    }
  });

  document.getElementById('logout').addEventListener('click', () => {
    token = null;
    sessionStorage.removeItem(storageKey);
    setMessage('');
    render();
  });

  render();
})();
