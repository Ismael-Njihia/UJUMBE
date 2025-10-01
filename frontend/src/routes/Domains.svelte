<script>
  import { onMount } from 'svelte';
  import api from '../lib/api';

  let domains = [];
  let loading = true;
  let showModal = false;
  let newDomain = '';

  onMount(loadDomains);

  async function loadDomains() {
    try {
      const response = await api.getDomains();
      domains = response.data || [];
    } catch (error) {
      console.error('Failed to load domains:', error);
    } finally {
      loading = false;
    }
  }

  async function addDomain() {
    try {
      await api.addDomain({ domain: newDomain });
      showModal = false;
      newDomain = '';
      await loadDomains();
    } catch (error) {
      alert('Failed to add domain');
    }
  }

  async function verifyDomain(id) {
    try {
      await api.verifyDomain(id);
      await loadDomains();
      alert('Domain verified successfully!');
    } catch (error) {
      alert('Failed to verify domain');
    }
  }

  async function deleteDomain(id) {
    if (confirm('Are you sure?')) {
      try {
        await api.deleteDomain(id);
        await loadDomains();
      } catch (error) {
        alert('Failed to delete domain');
      }
    }
  }
</script>

<div class="domains">
  <div class="header">
    <h1>Verified Sender Domains</h1>
    <button class="btn-primary" on:click={() => showModal = true}>Add Domain</button>
  </div>

  {#if loading}
    <p>Loading...</p>
  {:else if domains.length === 0}
    <div class="empty-state">
      <p>No domains yet. Add your first domain to start sending emails!</p>
    </div>
  {:else}
    <div class="domains-list">
      {#each domains as domain}
        <div class="domain-card">
          <div class="domain-info">
            <h3>{domain.domain}</h3>
            <span class="status" class:verified={domain.is_verified}>
              {domain.is_verified ? '✓ Verified' : '⏳ Pending'}
            </span>
          </div>
          {#if !domain.is_verified}
            <div class="verification">
              <p>Add this TXT record to verify:</p>
              <code>{domain.verification_code}</code>
              <button on:click={() => verifyDomain(domain.id)}>Check Verification</button>
            </div>
          {/if}
          <button class="delete" on:click={() => deleteDomain(domain.id)}>Delete</button>
        </div>
      {/each}
    </div>
  {/if}
</div>

{#if showModal}
  <div class="modal-overlay" on:click={() => showModal = false}>
    <div class="modal" on:click|stopPropagation>
      <h2>Add Domain</h2>
      <form on:submit|preventDefault={addDomain}>
        <div class="form-group">
          <label>Domain Name</label>
          <input type="text" bind:value={newDomain} placeholder="example.com" required />
        </div>
        <div class="modal-actions">
          <button type="button" on:click={() => showModal = false}>Cancel</button>
          <button type="submit" class="btn-primary">Add</button>
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
    margin-bottom: 30px;
  }

  .btn-primary {
    padding: 10px 20px;
    background: #667eea;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
  }

  .domains-list {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .domain-card {
    background: white;
    padding: 20px;
    border-radius: 10px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  }

  .domain-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15px;
  }

  .domain-info h3 {
    margin: 0;
    color: #2c3e50;
  }

  .status {
    padding: 5px 15px;
    border-radius: 20px;
    font-size: 12px;
    background: #f39c12;
    color: white;
  }

  .status.verified {
    background: #27ae60;
  }

  .verification {
    background: #f8f9fa;
    padding: 15px;
    border-radius: 5px;
    margin-bottom: 15px;
  }

  .verification p {
    margin: 0 0 10px 0;
    font-size: 14px;
    color: #7f8c8d;
  }

  .verification code {
    display: block;
    padding: 10px;
    background: white;
    border: 1px solid #dee2e6;
    border-radius: 5px;
    font-size: 12px;
    word-break: break-all;
    margin-bottom: 10px;
  }

  .verification button {
    padding: 8px 15px;
    background: #3498db;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
  }

  .delete {
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
    font-size: 14px;
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
