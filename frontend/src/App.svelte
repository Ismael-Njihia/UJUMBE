<script>
  import { onMount } from 'svelte';
  import Login from './routes/Login.svelte';
  import Dashboard from './routes/Dashboard.svelte';
  import Templates from './routes/Templates.svelte';
  import Domains from './routes/Domains.svelte';
  import APIKeys from './routes/APIKeys.svelte';
  import EmailLogs from './routes/EmailLogs.svelte';
  import Analytics from './routes/Analytics.svelte';
  import Billing from './routes/Billing.svelte';

  let currentRoute = 'login';
  let isAuthenticated = false;
  let user = null;

  onMount(() => {
    const token = localStorage.getItem('token');
    if (token) {
      isAuthenticated = true;
      currentRoute = 'dashboard';
    }
  });

  function handleLogin(event) {
    isAuthenticated = true;
    user = event.detail.user;
    currentRoute = 'dashboard';
  }

  function handleLogout() {
    localStorage.removeItem('token');
    isAuthenticated = false;
    user = null;
    currentRoute = 'login';
  }

  function navigate(route) {
    currentRoute = route;
  }
</script>

<main>
  {#if !isAuthenticated}
    <Login on:login={handleLogin} />
  {:else}
    <div class="app-container">
      <nav class="sidebar">
        <div class="logo">
          <h1>UJUMBE</h1>
          <p>Email Delivery Platform</p>
        </div>
        
        <ul class="nav-menu">
          <li class:active={currentRoute === 'dashboard'}>
            <button on:click={() => navigate('dashboard')}>Dashboard</button>
          </li>
          <li class:active={currentRoute === 'templates'}>
            <button on:click={() => navigate('templates')}>Templates</button>
          </li>
          <li class:active={currentRoute === 'domains'}>
            <button on:click={() => navigate('domains')}>Domains</button>
          </li>
          <li class:active={currentRoute === 'api-keys'}>
            <button on:click={() => navigate('api-keys')}>API Keys</button>
          </li>
          <li class:active={currentRoute === 'logs'}>
            <button on:click={() => navigate('logs')}>Email Logs</button>
          </li>
          <li class:active={currentRoute === 'analytics'}>
            <button on:click={() => navigate('analytics')}>Analytics</button>
          </li>
          <li class:active={currentRoute === 'billing'}>
            <button on:click={() => navigate('billing')}>Billing</button>
          </li>
        </ul>

        <div class="user-section">
          <p>{user?.email}</p>
          <button on:click={handleLogout} class="logout-btn">Logout</button>
        </div>
      </nav>

      <div class="content">
        {#if currentRoute === 'dashboard'}
          <Dashboard />
        {:else if currentRoute === 'templates'}
          <Templates />
        {:else if currentRoute === 'domains'}
          <Domains />
        {:else if currentRoute === 'api-keys'}
          <APIKeys />
        {:else if currentRoute === 'logs'}
          <EmailLogs />
        {:else if currentRoute === 'analytics'}
          <Analytics />
        {:else if currentRoute === 'billing'}
          <Billing />
        {/if}
      </div>
    </div>
  {/if}
</main>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
    background: #f5f7fa;
  }

  main {
    width: 100%;
    min-height: 100vh;
  }

  .app-container {
    display: flex;
    min-height: 100vh;
  }

  .sidebar {
    width: 260px;
    background: #2c3e50;
    color: white;
    padding: 20px;
    display: flex;
    flex-direction: column;
  }

  .logo {
    margin-bottom: 40px;
  }

  .logo h1 {
    margin: 0;
    font-size: 28px;
    color: #3498db;
  }

  .logo p {
    margin: 5px 0 0 0;
    font-size: 12px;
    color: #95a5a6;
  }

  .nav-menu {
    list-style: none;
    padding: 0;
    margin: 0;
    flex-grow: 1;
  }

  .nav-menu li {
    margin-bottom: 10px;
  }

  .nav-menu button {
    width: 100%;
    padding: 12px 15px;
    background: transparent;
    border: none;
    color: #ecf0f1;
    text-align: left;
    cursor: pointer;
    border-radius: 5px;
    font-size: 15px;
    transition: background 0.2s;
  }

  .nav-menu li.active button,
  .nav-menu button:hover {
    background: #34495e;
  }

  .user-section {
    margin-top: auto;
    padding-top: 20px;
    border-top: 1px solid #34495e;
  }

  .user-section p {
    margin: 0 0 10px 0;
    font-size: 14px;
    color: #bdc3c7;
  }

  .logout-btn {
    width: 100%;
    padding: 10px;
    background: #e74c3c;
    border: none;
    color: white;
    border-radius: 5px;
    cursor: pointer;
    font-size: 14px;
  }

  .logout-btn:hover {
    background: #c0392b;
  }

  .content {
    flex: 1;
    padding: 30px;
    overflow-y: auto;
  }
</style>
