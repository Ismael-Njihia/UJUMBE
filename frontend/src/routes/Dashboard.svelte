<script>
  import { onMount } from 'svelte';
  import { emailAPI } from '../lib/api';

  let quota = null;
  let loading = true;

  onMount(async () => {
    try {
      const response = await emailAPI.getQuota();
      quota = response.data;
    } catch (err) {
      console.error('Failed to load quota:', err);
    } finally {
      loading = false;
    }
  });
</script>

<div class="dashboard">
  <h1>Dashboard</h1>
  
  {#if loading}
    <p>Loading...</p>
  {:else if quota}
    <div class="cards">
      <div class="card">
        <h3>Free Emails Remaining</h3>
        <p class="stat">{quota.free_emails_remaining}</p>
        <p class="subtitle">Resets monthly</p>
      </div>

      <div class="card">
        <h3>Paid Email Balance</h3>
        <p class="stat">{quota.paid_emails_balance}</p>
        <p class="subtitle">Top up with M-Pesa</p>
      </div>

      <div class="card">
        <h3>Total Available</h3>
        <p class="stat">{quota.free_emails_remaining + quota.paid_emails_balance}</p>
        <p class="subtitle">Ready to send</p>
      </div>
    </div>

    <div class="quick-actions">
      <h2>Quick Actions</h2>
      <div class="action-buttons">
        <a href="/send" class="btn btn-primary">Send Email</a>
        <a href="/templates" class="btn btn-secondary">Manage Templates</a>
        <a href="/domains" class="btn btn-secondary">Verify Domain</a>
        <a href="/analytics" class="btn btn-secondary">View Analytics</a>
      </div>
    </div>
  {:else}
    <p class="error">Failed to load dashboard data</p>
  {/if}
</div>

<style>
  .dashboard {
    padding: 2rem 0;
  }

  h1 {
    color: #2c3e50;
    margin-bottom: 2rem;
  }

  .cards {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1.5rem;
    margin-bottom: 3rem;
  }

  .card {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
    text-align: center;
  }

  .card h3 {
    margin: 0 0 1rem 0;
    color: #7f8c8d;
    font-size: 0.9rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .stat {
    font-size: 3rem;
    font-weight: bold;
    color: #667eea;
    margin: 0;
  }

  .subtitle {
    color: #95a5a6;
    font-size: 0.85rem;
    margin-top: 0.5rem;
  }

  .quick-actions {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  }

  .quick-actions h2 {
    margin: 0 0 1.5rem 0;
    color: #2c3e50;
  }

  .action-buttons {
    display: flex;
    gap: 1rem;
    flex-wrap: wrap;
  }

  .btn {
    padding: 0.75rem 1.5rem;
    border-radius: 4px;
    text-decoration: none;
    transition: all 0.3s;
    display: inline-block;
  }

  .btn-primary {
    background: #667eea;
    color: white;
  }

  .btn-primary:hover {
    background: #5568d3;
  }

  .btn-secondary {
    background: #ecf0f1;
    color: #2c3e50;
  }

  .btn-secondary:hover {
    background: #dfe6e9;
  }

  .error {
    color: #e74c3c;
  }
</style>
