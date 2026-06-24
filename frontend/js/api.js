/* =============================================
   Mapa das Sementes do Brasil — API Client
   Camada de comunicação com o backend Go/Gin
   ============================================= */

// URL base da API — troque pelo domínio do Railway em produção
const API_BASE = 'http://localhost:8080';

// Chave usada para guardar o token no localStorage
const TOKEN_KEY = 'msb_token';
const USUARIO_KEY = 'msb_usuario';

/* -----------------------------------------------
   Autenticação — token JWT
----------------------------------------------- */

/**
 * Salva o token e os dados do usuário após login bem-sucedido.
 * @param {string} token - JWT retornado pela API
 * @param {object} usuario - dados do usuário logado
 */
function salvarSessao(token, usuario) {
  localStorage.setItem(TOKEN_KEY, token);
  localStorage.setItem(USUARIO_KEY, JSON.stringify(usuario));
}

/**
 * Retorna o token JWT salvo, ou null se não houver sessão.
 */
function obterToken() {
  return localStorage.getItem(TOKEN_KEY);
}

/**
 * Retorna os dados do usuário logado, ou null.
 */
function obterUsuario() {
  const dados = localStorage.getItem(USUARIO_KEY);
  return dados ? JSON.parse(dados) : null;
}

/**
 * Remove token e dados do usuário (logout).
 */
function encerrarSessao() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(USUARIO_KEY);
  window.location.href = 'login.html';
}

/**
 * Verifica se há uma sessão ativa.
 */
function estaLogado() {
  return !!obterToken();
}

/**
 * Redireciona para login se não houver sessão.
 * Use no topo de páginas protegidas.
 */
function exigirLogin() {
  if (!estaLogado()) {
    window.location.href = 'login.html';
  }
}

/* -----------------------------------------------
   Função base de requisição HTTP
----------------------------------------------- */

/**
 * Faz uma requisição autenticada para a API.
 *
 * @param {string} rota      - ex: '/api/sementes'
 * @param {string} metodo    - GET, POST, PUT, DELETE
 * @param {object|null} corpo - body JSON (opcional)
 * @param {boolean} publico  - se true, não envia o header Authorization
 * @returns {Promise<object>} - resposta JSON da API
 */
async function requisicao(rota, metodo = 'GET', corpo = null, publico = false) {
  const headers = {
    'Content-Type': 'application/json',
  };

  if (!publico) {
    const token = obterToken();
    if (!token) {
      encerrarSessao();
      throw new Error('Sessão expirada. Faça login novamente.');
    }
    headers['Authorization'] = `Bearer ${token}`;
  }

  const opcoes = {
    method: metodo,
    headers,
  };

  if (corpo && metodo !== 'GET') {
    opcoes.body = JSON.stringify(corpo);
  }

  const resposta = await fetch(`${API_BASE}${rota}`, opcoes);

  // Trata resposta sem corpo (ex: 204 No Content)
  if (resposta.status === 204) return null;

  const dados = await resposta.json();

  if (!resposta.ok) {
    // A API Go retorna { "error": "mensagem" }
    throw new Error(dados.error || `Erro ${resposta.status}`);
  }

  return dados;
}

/* -----------------------------------------------
   Endpoints de autenticação
----------------------------------------------- */

/**
 * Realiza login e salva a sessão.
 * @param {string} email
 * @param {string} senha
 */
async function apiLogin(email, senha) {
  const dados = await requisicao('/api/auth/login', 'POST', { email, senha }, true);
  salvarSessao(dados.token, dados.usuario);
  return dados;
}

/**
 * Cadastra um novo usuário.
 * @param {object} usuario - { nome, email, senha, papel }
 */
async function apiCadastrar(usuario) {
  return await requisicao('/api/auth/register', 'POST', usuario, true);
}

/* -----------------------------------------------
   Endpoints de sementes / espécies
----------------------------------------------- */

async function apiListarSementes() {
  return await requisicao('/api/sementes');
}

async function apiObterSemente(id) {
  return await requisicao(`/api/sementes/${id}`);
}

async function apiCriarSemente(dados) {
  return await requisicao('/api/sementes', 'POST', dados);
}

async function apiAtualizarSemente(id, dados) {
  return await requisicao(`/api/sementes/${id}`, 'PUT', dados);
}

async function apiDeletarSemente(id) {
  return await requisicao(`/api/sementes/${id}`, 'DELETE');
}

/* -----------------------------------------------
   Endpoints de registros (ocorrências / mapa)
----------------------------------------------- */

/**
 * Retorna todos os registros no formato GeoJSON-like
 * para uso no mapa Leaflet.
 */
async function apiRegistrosMapa() {
  return await requisicao('/api/registros/mapa', 'GET', null, true);
}

async function apiRegistrosPorEstado(uf) {
  return await requisicao(`/api/registros/estado/${uf}`);
}

async function apiRegistrosPorProximidade(lat, lng, raioKm) {
  return await requisicao(`/api/registros/proximidade?lat=${lat}&lng=${lng}&raio=${raioKm}`);
}

/* -----------------------------------------------
   Utilitários de UI
----------------------------------------------- */

/**
 * Exibe uma mensagem de erro ou sucesso em um elemento .msg.
 * @param {string} seletorId - id do elemento (sem #)
 * @param {string} texto
 * @param {'erro'|'sucesso'} tipo
 */
function exibirMensagem(seletorId, texto, tipo = 'erro') {
  const el = document.getElementById(seletorId);
  if (!el) return;
  el.textContent = texto;
  el.className = `msg msg--${tipo} visivel`;
}

function ocultarMensagem(seletorId) {
  const el = document.getElementById(seletorId);
  if (el) el.classList.remove('visivel');
}

/**
 * Ativa/desativa estado de carregamento em um botão.
 * @param {HTMLButtonElement} btn
 * @param {boolean} carregando
 * @param {string} textoOriginal
 */
function estadoBotao(btn, carregando, textoOriginal = '') {
  btn.disabled = carregando;
  btn.innerHTML = carregando
    ? '<span class="spinner"></span>'
    : textoOriginal;
}

/**
 * Formata papel do usuário para exibição.
 */
function formatarPapel(papel) {
  const papeis = {
    artesao:      'Artesão(ã)',
    pesquisador:  'Pesquisador(a)',
    especialista: 'Especialista',
    admin:        'Administrador(a)',
  };
  return papeis[papel] || papel;
}