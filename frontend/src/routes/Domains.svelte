<script>
  import { onMount } from 'svelte';
  import { domainAPI } from '../lib/api';

  let domains = [];
  let loading = true;
  let showAddForm = false;
  let newDomain = '';
  let error = '';
  let success = '';

  onMount(loadDomains);

  async function loadDomains() {
    loading = true;
    try {
      const response = await domainAPI.list();
      domains = response.data || [];
    } catch (err) {
      error = 'Failed to load domains';
    } finally {
      loading = false;
    }
  }

  async function handleAdd() {
    error = '';
    success = '';

    try {
      await domainAPI.add(newDomain);
      success = 'Domain added! Please verify it in your DNS settings.';
      showAddForm = false;
      newDomain = '';
      await loadDomains();
    } catch (err) {
      error = err.response?.data?.error || 'Failed to add domain';
    }
  }

  async function handleVerify(id) {
    try {
      await domainAPI.verify(id);
      success = 'Domain verified successfully!';
      await loadDomains();
    } catch (err) {
      error = 'Failed to verify domain. Please ensure DNS records are set up correctly.';
    }
  }

  async function handleDelete(id) {
    if (!confirm('Are you sure you want to delete this domain?')) return;

    try {
      await domainAPI.delete(id);
      success = 'Domain deleted successfully!';
      await loadDomains();
    } catch (err) {
      error = 'Failed to delete domain';
    }
  }
</script>

<div class="domains">
  <div class="header">
    <h1>Verified Domains</h1>
    <button class="btn btn-primary" on:click={() => showAddForm = !showAddForm}>
      {showAddForm ? 'Cancel' : 'Add Domain'}
    </button>
  </div>

  {#if success}
    <div class="alert alert-success">{success}</div>
  {/if}

  {#if error}
    <div class="alert alert-error">{error}</div>
  {/if}

  {#if showAddForm}
    <div class="add-form">
      <h2>Add New Domain</h2>
      <form on:submit|preventDefault={handleAdd}>
        <div class="form-group">
          <label for="domain">Domain Name</label>
          <input
            id="domain"
            type="text"
            bind:value={newDomain}
            required
            placeholder="yourdomain.com"
          />
          <small>Add your domain to verify ownership and send emails from it</small>
        </div>
        <button type="submit" class="btn btn-primary">Add Domain</button>
      </form>
    </div>
  {/if}

  {#if loading}
    <p>Loading domains...</p>
  {:else if domains.length === 0}
    <div class="empty-state">
      <p>No domains yet. Add a domain to start sending emails from your own address!</p>
    </div>
  {:else}
    <div class="domains-list">
      {#each domains as domain}
        <div class="domain-card">
          <div class="domain-info">
            <h3>{domain.domain}</h3>
            <span class="status" class:verified={domain.verified}>
              {domain.verified ? '✓ Verified' : '⚠ Pending Verification'}
            </span>
          </div>
          <div class="actions">
            {#if !domain.verified}
              <button class="btn btn-small btn-success" on:click={() => handleVerify(domain.id)}>
                Verify
              </button>
            {/if}
            <button class="btn btn-small btn-danger" on:click={() => handleDelete(domain.id)}>
              Delete
            </button>
          </div>
        </div>
      {/each}
    </div>
  {/if}

  <div class="info-box">
    <h3>How to verify your domain</h3>
    <ol>
      <li>Add your domain using the "Add Domain" button above</li>
      <li>Log in to your domain registrar (e.g., GoDaddy, Namecheap)</li>
      <li>Add the required DNS TXT records for AWS SES verification</li>
      <li>Wait for DNS propagation (usually 5-10 minutes)</li>
      <li>Click the "Verify" button to complete verification</li>
    </ol>
  </div>
</div>

<style>
  .domains {
    padding: 2rem 0;
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
  }

  h1 {
    color: #2c3e50;
    margin: 0;
  }

  .add-form {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
    margin-bottom: 2rem;
  }

  .add-form h2 {
    margin-top: 0;
    color: #2c3e50;
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

  small {
    color: #7f8c8d;
    font-size: 0.85rem;
  }

  .domains-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-bottom: 2rem;
  }

  .domain-card {
    background: white;
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .domain-info h3 {
    margin: 0 0 0.5rem 0;
    color: #2c3e50;
  }

  .status {
    padding: 0.25rem 0.75rem;
    border-radius: 12px;
    font-size: 0.85rem;
    font-weight: 500;
  }

  .status.verified {
    background: #d4edda;
    color: #155724;
  }

  .status:not(.verified) {
    background: #fff3cd;
    color: #856404;
  }

  .actions {
    display: flex;
    gap: 0.5rem;
  }

  .btn {
    padding: 0.75rem 1.5rem;
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

  .btn-primary:hover {
    background: #5568d3;
  }

  .btn-small {
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
  }

  .btn-success {
    background: #27ae60;
    color: white;
  }

  .btn-success:hover {
    background: #229954;
  }

  .btn-danger {
    background: #e74c3c;
    color: white;
  }

  .btn-danger:hover {
    background: #c0392b;
  }

  .alert {
    padding: 1rem;
    border-radius: 4px;
    margin-bottom: 1.5rem;
  }

  .alert-success {
    background: #d4edda;
    color: #155724;
    border: 1px solid #c3e6cb;
  }

  .alert-error {
    background: #fee;
    color: #c00;
    border: 1px solid #fcc;
  }

  .empty-state {
    background: white;
    padding: 3rem;
    border-radius: 8px;
    text-align: center;
    color: #7f8c8d;
    margin-bottom: 2rem;
  }

  .info-box {
    background: #e8f4fd;
    padding: 1.5rem;
    border-radius: 8px;
    border-left: 4px solid #3498db;
  }

  .info-box h3 {
    margin-top: 0;
    color: #2c3e50;
  }

  .info-box ol {
    margin: 0;
    padding-left: 1.5rem;
    color: #7f8c8d;
  }

  .info-box li {
    margin-bottom: 0.5rem;
  }
</style>
