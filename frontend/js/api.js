/* =============================================
   Mapa das Sementes do Brasil — API Client
   Camada de comunicação com o backend Go/Gin
   ============================================= */

// Detecta automaticamente o ambiente:
// - Em produção (Railway): frontend e API no mesmo domínio → string vazia
// - Em desenvolvimento local: aponta para Go na porta 8080
const API_BASE = (
  window.location.hostname === 'localhost' ||
  window.location.hostname === '127.0.0.1'
) ? 'http://localhost:8080' : '';

// Chaves usadas para guardar sessão no localStorage
const TOKEN_KEY   = 'msb_token';
const USUARIO_KEY = 'msb_usuario';

/* -----------------------------------------------
   Gestão de sessão / JWT
----------------------------------------------- */

function salvarSessao(token, usuario) {
  localStorage.setItem(TOKEN_KEY, token);
  localStorage.setItem(USUARIO_KEY, JSON.stringify(usuario));
}

function obterToken() {
  return localStorage.getItem(TOKEN_KEY);
}

function obterUsuario() {
  const dados = localStorage.getItem(USUARIO_KEY);
  return dados ? JSON.parse(dados) : null;
}

function encerrarSessao() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(USUARIO_KEY);
  window.location.href = '/frontend/login.html';
}

function estaLogado() {
  return !!obterToken();
}

/** Redireciona para login se não houver sessão. Use no topo de páginas protegidas. */
function exigirLogin() {
  if (!estaLogado()) {
    window.location.href = '/frontend/login.html';
  }
}

/* -----------------------------------------------
   Função base de requisição HTTP
----------------------------------------------- */

/**
 * @param {string}      rota    - ex: '/api/sementes'
 * @param {string}      metodo  - GET | POST | PUT | DELETE
 * @param {object|null} corpo   - body JSON (opcional)
 * @param {boolean}     publico - se true, não envia Authorization
 */
async function requisicao(rota, metodo = 'GET', corpo = null, publico = false) {
  const headers = { 'Content-Type': 'application/json' };

  if (!publico) {
    const token = obterToken();
    if (!token) {
      encerrarSessao();
      throw new Error('Sessão expirada. Faça login novamente.');
    }
    headers['Authorization'] = `Bearer ${token}`;
  }

  const opcoes = { method: metodo, headers };
  if (corpo && metodo !== 'GET') {
    opcoes.body = JSON.stringify(corpo);
  }

  const resposta = await fetch(`${API_BASE}${rota}`, opcoes);

  if (resposta.status === 204) return null;

  const dados = await resposta.json();

  if (!resposta.ok) {
    // A API Go retorna { "error": "mensagem" }
    throw new Error(dados.error || `Erro ${resposta.status}`);
  }

  return dados;
}

/* -----------------------------------------------
   Autenticação
----------------------------------------------- */

async function apiLogin(email, senha) {
  const resposta = await requisicao('/api/auth/login', 'POST', { email, senha }, true);
  // O backend retorna { message, data: { token, usuario } }
  const dados = resposta.data || resposta;
  salvarSessao(dados.token, dados.usuario);
  return dados;
}

// Rota real no main.go: POST /api/auth/cadastro
// Campos obrigatórios no backend: nome_completo, email, senha, telefone
async function apiCadastrar(usuario) {
  const resposta = await requisicao('/api/auth/cadastro', 'POST', usuario, true);
  return resposta.data || resposta;
}

/* -----------------------------------------------
   Espécies (públicas)
----------------------------------------------- */

async function apiListarEspecies() {
  return await requisicao('/api/especies', 'GET', null, true);
}

async function apiDetalheEspecie(id) {
  return await requisicao(`/api/especies/${id}`, 'GET', null, true);
}

async function apiCriarEspecie(dados) {
  return await requisicao('/api/especies', 'POST', dados);
}

async function apiEditarEspecie(id, dados) {
  return await requisicao(`/api/especies/${id}`, 'PUT', dados);
}

async function apiDeletarEspecie(id) {
  return await requisicao(`/api/especies/${id}`, 'DELETE');
}

/* -----------------------------------------------
   Sementes (públicas para leitura)
----------------------------------------------- */

async function apiListarSementes() {
  return await requisicao('/api/sementes', 'GET', null, true);
}

async function apiDetalheSemente(id) {
  return await requisicao(`/api/sementes/${id}`, 'GET', null, true);
}

async function apiCriarSemente(dados) {
  return await requisicao('/api/sementes', 'POST', dados);
}

async function apiEditarSemente(id, dados) {
  return await requisicao(`/api/sementes/${id}`, 'PUT', dados);
}

async function apiDeletarSemente(id) {
  return await requisicao(`/api/sementes/${id}`, 'DELETE');
}

/* -----------------------------------------------
   Registros / Mapa
----------------------------------------------- */

/** Pontos para o Leaflet — rota pública: GET /api/registros/mapa */
async function apiRegistrosMapa() {
  return await requisicao('/api/registros/mapa', 'GET', null, true);
}

async function apiListarRegistros() {
  return await requisicao('/api/registros', 'GET', null, true);
}

async function apiDetalheRegistro(id) {
  return await requisicao(`/api/registros/${id}`, 'GET', null, true);
}

async function apiCriarRegistro(dados) {
  return await requisicao('/api/registros', 'POST', dados);
}

async function apiEditarRegistro(id, dados) {
  return await requisicao(`/api/registros/${id}`, 'PUT', dados);
}

async function apiDeletarRegistro(id) {
  return await requisicao(`/api/registros/${id}`, 'DELETE');
}

/* -----------------------------------------------
   Busca avançada (pública) — /api/busca/...
----------------------------------------------- */

/** Estatísticas gerais — usadas na home */
async function apiEstatisticas() {
  return await requisicao('/api/busca/estatisticas', 'GET', null, true);
}

/** Busca por proximidade geográfica */
async function apiBuscaProximidade(lat, lng, raioKm) {
  return await requisicao(
    `/api/busca/mapa?lat=${lat}&lng=${lng}&raio=${raioKm}`,
    'GET', null, true
  );
}

/** Busca por estado (UF) */
async function apiBuscaPorEstado(uf) {
  return await requisicao(`/api/busca/estado/${uf}`, 'GET', null, true);
}

/** Busca geral por termo */
async function apiBuscaGeral(termo) {
  return await requisicao(`/api/busca?q=${encodeURIComponent(termo)}`, 'GET', null, true);
}

/** Busca de espécies por nome */
async function apiBuscaEspecies(termo) {
  return await requisicao(`/api/busca/especies?q=${encodeURIComponent(termo)}`, 'GET', null, true);
}

/* -----------------------------------------------
   Perfil (autenticado)
----------------------------------------------- */

async function apiMeuPerfil() {
  return await requisicao('/api/perfil');
}

async function apiEditarPerfil(dados) {
  return await requisicao('/api/perfil', 'PUT', dados);
}

async function apiMinhasContribuicoes() {
  return await requisicao('/api/perfil/contribuicoes');
}

/* -----------------------------------------------
   Conhecimento Tradicional
----------------------------------------------- */

async function apiListarConhecimentos() {
  return await requisicao('/api/conhecimentos', 'GET', null, true);
}

async function apiDetalheConhecimento(id) {
  return await requisicao(`/api/conhecimentos/${id}`, 'GET', null, true);
}

async function apiCriarConhecimento(dados) {
  return await requisicao('/api/conhecimentos', 'POST', dados);
}

async function apiCurtirConhecimento(id) {
  return await requisicao(`/api/conhecimentos/${id}/curtir`, 'POST');
}

async function apiValidarConhecimento(id) {
  return await requisicao(`/api/conhecimentos/${id}/validar`, 'POST');
}

/* -----------------------------------------------
   Utilitários de UI
----------------------------------------------- */

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

function estadoBotao(btn, carregando, textoOriginal = '') {
  btn.disabled = carregando;
  btn.innerHTML = carregando
    ? '<span class="spinner"></span>'
    : textoOriginal;
}

function formatarPapel(papel) {
  const papeis = {
    artesao:            'Artesão(ã)',
    pesquisador:        'Pesquisador(a)',
    especialista:       'Especialista',
    estudante:          'Estudante',
    agente_territorial: 'Agente Territorial Cultural',
    admin:              'Administrador(a)',
  };
  return papeis[papel] || papel;
}
