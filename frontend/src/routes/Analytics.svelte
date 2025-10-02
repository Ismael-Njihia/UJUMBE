<script>
  import { onMount } from 'svelte';
  import { emailAPI } from '../lib/api';

  let analytics = null;
  let loading = true;

  onMount(async () => {
    try {
      const response = await emailAPI.getAnalytics();
      analytics = response.data;
    } catch (err) {
      console.error('Failed to load analytics:', err);
    } finally {
      loading = false;
    }
  });
</script>

<div class="analytics">
  <h1>Email Analytics</h1>

  {#if loading}
    <p>Loading analytics...</p>
  {:else if analytics}
    <div class="stats-grid">
      <div class="stat-card">
        <h3>Emails Sent</h3>
        <p class="stat-number success">{analytics.total_emails_sent}</p>
      </div>

      <div class="stat-card">
        <h3>Emails Failed</h3>
        <p class="stat-number danger">{analytics.total_emails_failed}</p>
      </div>

      <div class="stat-card">
        <h3>Emails Pending</h3>
        <p class="stat-number warning">{analytics.total_emails_pending}</p>
      </div>

      <div class="stat-card">
        <h3>Success Rate</h3>
        <p class="stat-number">{analytics.success_rate.toFixed(2)}%</p>
      </div>
    </div>

    <div class="quota-section">
      <h2>Current Quota</h2>
      <div class="quota-grid">
        <div class="quota-card">
          <h3>Free Emails Remaining</h3>
          <p class="quota-number">{analytics.free_emails_remaining}</p>
          <div class="progress-bar">
            <div
              class="progress-fill"
              style="width: {(analytics.free_emails_remaining / 100) * 100}%"
            ></div>
          </div>
          <p class="quota-label">out of 100 monthly</p>
        </div>

        <div class="quota-card">
          <h3>Paid Email Balance</h3>
          <p class="quota-number">{analytics.paid_emails_balance}</p>
          <p class="quota-label">Top up with M-Pesa for more</p>
        </div>
      </div>
    </div>
  {:else}
    <p class="error">Failed to load analytics</p>
  {/if}
</div>

<style>
  .analytics {
    padding: 2rem 0;
  }

  h1 {
    color: #2c3e50;
    margin-bottom: 2rem;
  }

  h2 {
    color: #2c3e50;
    margin-bottom: 1.5rem;
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1.5rem;
    margin-bottom: 3rem;
  }

  .stat-card {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
    text-align: center;
  }

  .stat-card h3 {
    margin: 0 0 1rem 0;
    color: #7f8c8d;
    font-size: 0.9rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .stat-number {
    font-size: 2.5rem;
    font-weight: bold;
    margin: 0;
    color: #2c3e50;
  }

  .stat-number.success {
    color: #27ae60;
  }

  .stat-number.danger {
    color: #e74c3c;
  }

  .stat-number.warning {
    color: #f39c12;
  }

  .quota-section {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  }

  .quota-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 2rem;
  }

  .quota-card {
    text-align: center;
  }

  .quota-card h3 {
    margin: 0 0 1rem 0;
    color: #7f8c8d;
    font-size: 0.9rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .quota-number {
    font-size: 3rem;
    font-weight: bold;
    color: #667eea;
    margin: 0 0 1rem 0;
  }

  .progress-bar {
    width: 100%;
    height: 12px;
    background: #ecf0f1;
    border-radius: 6px;
    overflow: hidden;
    margin: 1rem 0;
  }

  .progress-fill {
    height: 100%;
    background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
    transition: width 0.3s ease;
  }

  .quota-label {
    color: #95a5a6;
    font-size: 0.85rem;
    margin: 0.5rem 0 0 0;
  }

  .error {
    color: #e74c3c;
  }
</style>
