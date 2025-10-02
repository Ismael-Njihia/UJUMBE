<script>
  import { Router, Route } from 'svelte-routing';
  import Login from './routes/Login.svelte';
  import Register from './routes/Register.svelte';
  import Dashboard from './routes/Dashboard.svelte';
  import SendEmail from './routes/SendEmail.svelte';
  import Templates from './routes/Templates.svelte';
  import Domains from './routes/Domains.svelte';
  import Analytics from './routes/Analytics.svelte';
  import { setApiKey } from './lib/api';
  import { onMount } from 'svelte';

  let isAuthenticated = false;

  onMount(() => {
    const apiKey = localStorage.getItem('apiKey');
    if (apiKey) {
      setApiKey(apiKey);
      isAuthenticated = true;
    }
  });

  function handleLogin(apiKey) {
    localStorage.setItem('apiKey', apiKey);
    setApiKey(apiKey);
    isAuthenticated = true;
  }

  function handleLogout() {
    localStorage.removeItem('apiKey');
    setApiKey(null);
    isAuthenticated = false;
  }
</script>

<Router>
  <div class="app">
    {#if isAuthenticated}
      <nav class="navbar">
        <div class="nav-brand">
          <h1>UJUMBE</h1>
        </div>
        <ul class="nav-links">
          <li><a href="/">Dashboard</a></li>
          <li><a href="/send">Send Email</a></li>
          <li><a href="/templates">Templates</a></li>
          <li><a href="/domains">Domains</a></li>
          <li><a href="/analytics">Analytics</a></li>
          <li><button on:click={handleLogout} class="btn-logout">Logout</button></li>
        </ul>
      </nav>
      <main class="container">
        <Route path="/" component={Dashboard} />
        <Route path="/send" component={SendEmail} />
        <Route path="/templates" component={Templates} />
        <Route path="/domains" component={Domains} />
        <Route path="/analytics" component={Analytics} />
      </main>
    {:else}
      <Route path="/register" component={Register} props={{ onRegister: handleLogin }} />
      <Route path="/*" component={Login} props={{ onLogin: handleLogin }} />
    {/if}
  </div>
</Router>

<style>
  :global(body) {
    margin: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
    background: #f5f5f5;
  }

  .app {
    min-height: 100vh;
  }

  .navbar {
    background: #2c3e50;
    color: white;
    padding: 1rem 2rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  }

  .nav-brand h1 {
    margin: 0;
    font-size: 1.5rem;
  }

  .nav-links {
    display: flex;
    list-style: none;
    margin: 0;
    padding: 0;
    gap: 1.5rem;
    align-items: center;
  }

  .nav-links a {
    color: white;
    text-decoration: none;
    transition: color 0.3s;
  }

  .nav-links a:hover {
    color: #3498db;
  }

  .btn-logout {
    background: #e74c3c;
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    transition: background 0.3s;
  }

  .btn-logout:hover {
    background: #c0392b;
  }

  .container {
    max-width: 1200px;
    margin: 2rem auto;
    padding: 0 1rem;
  }
</style>
