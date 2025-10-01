<script>
  import { createEventDispatcher } from 'svelte';
  import api from '../lib/api';

  const dispatch = createEventDispatcher();

  let isLogin = true;
  let email = '';
  let password = '';
  let name = '';
  let error = '';
  let loading = false;

  async function handleSubmit() {
    error = '';
    loading = true;

    try {
      if (isLogin) {
        const response = await api.login({ email, password });
        localStorage.setItem('token', response.data.token);
        dispatch('login', { user: response.data.user });
      } else {
        await api.register({ email, password, name });
        // Auto login after registration
        const response = await api.login({ email, password });
        localStorage.setItem('token', response.data.token);
        dispatch('login', { user: response.data.user });
      }
    } catch (err) {
      error = err.response?.data?.error || 'An error occurred';
    } finally {
      loading = false;
    }
  }
</script>

<div class="login-container">
  <div class="login-card">
    <div class="login-header">
      <h1>UJUMBE</h1>
      <p>Email Delivery Platform for Kenya</p>
    </div>

    <div class="tabs">
      <button class:active={isLogin} on:click={() => isLogin = true}>Login</button>
      <button class:active={!isLogin} on:click={() => isLogin = false}>Register</button>
    </div>

    <form on:submit|preventDefault={handleSubmit}>
      {#if !isLogin}
        <div class="form-group">
          <label for="name">Name</label>
          <input
            type="text"
            id="name"
            bind:value={name}
            required
            placeholder="Your name"
          />
        </div>
      {/if}

      <div class="form-group">
        <label for="email">Email</label>
        <input
          type="email"
          id="email"
          bind:value={email}
          required
          placeholder="your@email.com"
        />
      </div>

      <div class="form-group">
        <label for="password">Password</label>
        <input
          type="password"
          id="password"
          bind:value={password}
          required
          placeholder="••••••••"
        />
      </div>

      {#if error}
        <div class="error">{error}</div>
      {/if}

      <button type="submit" class="submit-btn" disabled={loading}>
        {loading ? 'Please wait...' : (isLogin ? 'Login' : 'Register')}
      </button>
    </form>

    <div class="info">
      <p>✨ 100 free emails monthly</p>
      <p>💳 Pay-as-you-go with M-Pesa</p>
      <p>🚀 Fast & reliable delivery</p>
    </div>
  </div>
</div>

<style>
  .login-container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  .login-card {
    background: white;
    border-radius: 10px;
    padding: 40px;
    width: 100%;
    max-width: 400px;
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
  }

  .login-header {
    text-align: center;
    margin-bottom: 30px;
  }

  .login-header h1 {
    margin: 0;
    font-size: 32px;
    color: #667eea;
  }

  .login-header p {
    margin: 5px 0 0 0;
    color: #6c757d;
    font-size: 14px;
  }

  .tabs {
    display: flex;
    gap: 10px;
    margin-bottom: 20px;
  }

  .tabs button {
    flex: 1;
    padding: 10px;
    border: none;
    background: #f8f9fa;
    color: #6c757d;
    cursor: pointer;
    border-radius: 5px;
    font-size: 14px;
    transition: all 0.3s;
  }

  .tabs button.active {
    background: #667eea;
    color: white;
  }

  .form-group {
    margin-bottom: 20px;
  }

  .form-group label {
    display: block;
    margin-bottom: 5px;
    font-size: 14px;
    color: #495057;
    font-weight: 500;
  }

  .form-group input {
    width: 100%;
    padding: 12px;
    border: 1px solid #dee2e6;
    border-radius: 5px;
    font-size: 14px;
    box-sizing: border-box;
  }

  .form-group input:focus {
    outline: none;
    border-color: #667eea;
  }

  .error {
    background: #fee;
    color: #c33;
    padding: 10px;
    border-radius: 5px;
    margin-bottom: 15px;
    font-size: 14px;
  }

  .submit-btn {
    width: 100%;
    padding: 12px;
    background: #667eea;
    color: white;
    border: none;
    border-radius: 5px;
    font-size: 16px;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.3s;
  }

  .submit-btn:hover:not(:disabled) {
    background: #5568d3;
  }

  .submit-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .info {
    margin-top: 30px;
    text-align: center;
  }

  .info p {
    margin: 8px 0;
    font-size: 13px;
    color: #6c757d;
  }
</style>
