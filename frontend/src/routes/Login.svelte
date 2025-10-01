<script>
  import { authAPI } from '../lib/api';
  import { navigate } from 'svelte-routing';

  export let onLogin;

  let email = '';
  let password = '';
  let error = '';
  let loading = false;

  async function handleSubmit() {
    error = '';
    loading = true;

    try {
      const response = await authAPI.login(email, password);
      if (response.data.success) {
        onLogin(response.data.api_key);
        navigate('/');
      } else {
        error = response.data.message;
      }
    } catch (err) {
      error = err.response?.data?.error || 'Login failed';
    } finally {
      loading = false;
    }
  }
</script>

<div class="auth-container">
  <div class="auth-card">
    <h2>Login to UJUMBE</h2>
    <p class="subtitle">Email delivery platform for Kenya</p>

    {#if error}
      <div class="alert alert-error">{error}</div>
    {/if}

    <form on:submit|preventDefault={handleSubmit}>
      <div class="form-group">
        <label for="email">Email</label>
        <input
          id="email"
          type="email"
          bind:value={email}
          required
          placeholder="your@email.com"
        />
      </div>

      <div class="form-group">
        <label for="password">Password</label>
        <input
          id="password"
          type="password"
          bind:value={password}
          required
          placeholder="••••••••"
        />
      </div>

      <button type="submit" class="btn btn-primary" disabled={loading}>
        {loading ? 'Logging in...' : 'Login'}
      </button>
    </form>

    <p class="auth-switch">
      Don't have an account? <a href="/register">Register here</a>
    </p>
  </div>
</div>

<style>
  .auth-container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  .auth-card {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 10px 25px rgba(0,0,0,0.2);
    width: 100%;
    max-width: 400px;
  }

  h2 {
    margin: 0 0 0.5rem 0;
    color: #2c3e50;
  }

  .subtitle {
    color: #7f8c8d;
    margin: 0 0 2rem 0;
  }

  .form-group {
    margin-bottom: 1.5rem;
  }

  label {
    display: block;
    margin-bottom: 0.5rem;
    color: #2c3e50;
    font-weight: 500;
  }

  input {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
    box-sizing: border-box;
  }

  input:focus {
    outline: none;
    border-color: #667eea;
  }

  .btn {
    width: 100%;
    padding: 0.75rem;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    transition: background 0.3s;
  }

  .btn-primary {
    background: #667eea;
    color: white;
  }

  .btn-primary:hover:not(:disabled) {
    background: #5568d3;
  }

  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .alert {
    padding: 0.75rem;
    border-radius: 4px;
    margin-bottom: 1rem;
  }

  .alert-error {
    background: #fee;
    color: #c00;
    border: 1px solid #fcc;
  }

  .auth-switch {
    text-align: center;
    margin-top: 1.5rem;
    color: #7f8c8d;
  }

  .auth-switch a {
    color: #667eea;
    text-decoration: none;
  }

  .auth-switch a:hover {
    text-decoration: underline;
  }
</style>
