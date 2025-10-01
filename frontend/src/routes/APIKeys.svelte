<script>
  import { onMount } from 'svelte';
  import api from '../lib/api';

  let apiKeys = [];
  let loading = true;
  let showModal = false;
  let keyName = '';

  onMount(loadAPIKeys);

  async function loadAPIKeys() {
    try {
      const response = await api.getAPIKeys();
      apiKeys = response.data || [];
    } catch (error) {
      console.error('Failed to load API keys:', error);
    } finally {
      loading = false;
    }
  }

  async function createKey() {
    try {
      await api.createAPIKey({ name: keyName });
      showModal = false;
      keyName = '';
      await loadAPIKeys();
    } catch (error) {
      alert('Failed to create API key');
    }
  }

  async function revokeKey(id) {
    if (confirm('Are you sure you want to revoke this API key?')) {
      try {
        await api.revokeAPIKey(id);
        await loadAPIKeys();
      } catch (error) {
        alert('Failed to revoke API key');
      }
    }
  }

  function copyToClipboard(text) {
    navigator.clipboard.writeText(text);
    alert('Copied to clipboard!');
  }
</script>

<div class="api-keys">
  <div class="header">
    <h1>API Keys</h1>
    <button class="btn-primary" on:click={() => showModal = true}>Create API Key</button>
  </div>

  <div class="info-box">
    <p>Use API keys to authenticate your API requests. Keep them secure!</p>
    <code>X-API-Key: your-api-key-here</code>
  </div>

  {#if loading}
    <p>Loading...</p>
  {:else if apiKeys.length === 0}
    <div class="empty-state">
      <p>No API keys yet. Create your first API key!</p>
    </div>
  {:else}
    <div class="keys-list">
      {#each apiKeys as key}
        <div class="key-card">
          <div class="key-info">
            <h3>{key.name}</h3>
            <p class="key-value">
              <code>{key.key}</code>
              <button class="copy-btn" on:click={() => copyToClipboard(key.key)}>📋 Copy</button>
            </p>
            <p class="date">Created: {new Date(key.created_at).toLocaleDateString()}</p>
            <span class="status" class:active={key.is_active}>
              {key.is_active ? 'Active' : 'Revoked'}
            </span>
          </div>
          {#if key.is_active}
            <button class="revoke" on:click={() => revokeKey(key.id)}>Revoke</button>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</div>

{#if showModal}
  <div class="modal-overlay" on:click={() => showModal = false}>
    <div class="modal" on:click|stopPropagation>
      <h2>Create API Key</h2>
      <form on:submit|preventDefault={createKey}>
        <div class="form-group">
          <label>Key Name</label>
          <input type="text" bind:value={keyName} placeholder="e.g., Production Server" required />
        </div>
        <div class="modal-actions">
          <button type="button" on:click={() => showModal = false}>Cancel</button>
          <button type="submit" class="btn-primary">Create</button>
        </div>
      </form>
    </div>
  </div>
{/if}

<style>
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  .btn-primary {
    padding: 10px 20px;
    background: #667eea;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
  }

  .info-box {
    background: #e8f4f8;
    padding: 15px;
    border-radius: 5px;
    margin-bottom: 30px;
    border-left: 4px solid #3498db;
  }

  .info-box p {
    margin: 0 0 10px 0;
    color: #2c3e50;
  }

  .info-box code {
    display: block;
    padding: 10px;
    background: white;
    border-radius: 3px;
    font-size: 13px;
  }

  .keys-list {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .key-card {
    background: white;
    padding: 20px;
    border-radius: 10px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  }

  .key-info h3 {
    margin: 0 0 10px 0;
    color: #2c3e50;
  }

  .key-value {
    display: flex;
    align-items: center;
    gap: 10px;
    margin: 10px 0;
  }

  .key-value code {
    flex: 1;
    padding: 10px;
    background: #f8f9fa;
    border: 1px solid #dee2e6;
    border-radius: 5px;
    font-size: 12px;
    overflow-x: auto;
  }

  .copy-btn {
    padding: 8px 12px;
    background: #3498db;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    font-size: 12px;
  }

  .date {
    margin: 10px 0;
    color: #7f8c8d;
    font-size: 12px;
  }

  .status {
    padding: 5px 15px;
    border-radius: 20px;
    font-size: 12px;
    background: #95a5a6;
    color: white;
  }

  .status.active {
    background: #27ae60;
  }

  .revoke {
    margin-top: 15px;
    padding: 8px 15px;
    background: #e74c3c;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
  }

  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .modal {
    background: white;
    border-radius: 10px;
    padding: 30px;
    width: 90%;
    max-width: 500px;
  }

  .form-group {
    margin-bottom: 20px;
  }

  .form-group label {
    display: block;
    margin-bottom: 5px;
    font-weight: 500;
  }

  .form-group input {
    width: 100%;
    padding: 10px;
    border: 1px solid #dee2e6;
    border-radius: 5px;
    box-sizing: border-box;
  }

  .modal-actions {
    display: flex;
    gap: 10px;
    justify-content: flex-end;
  }

  .modal-actions button {
    padding: 10px 20px;
    border: none;
    border-radius: 5px;
    cursor: pointer;
  }

  .modal-actions button:first-child {
    background: #95a5a6;
    color: white;
  }
</style>
