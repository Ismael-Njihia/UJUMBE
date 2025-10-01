<script>
  import { onMount } from 'svelte';
  import api from '../lib/api';

  let stats = {};
  let loading = true;

  onMount(async () => {
    try {
      const response = await api.getDashboardStats();
      stats = response.data;
    } catch (error) {
      console.error('Failed to load stats:', error);
    } finally {
      loading = false;
    }
  });
</script>

<div class="dashboard">
  <h1>Dashboard</h1>

  {#if loading}
    <p>Loading...</p>
  {:else}
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">📧</div>
        <div class="stat-content">
          <h3>{stats.emails_sent || 0}</h3>
          <p>Emails Sent This Month</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">✅</div>
        <div class="stat-content">
          <h3>{stats.quota_remaining || 0}</h3>
          <p>Quota Remaining</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">🎯</div>
        <div class="stat-content">
          <h3>{stats.total_sent || 0}</h3>
          <p>Total Delivered</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">💰</div>
        <div class="stat-content">
          <h3>KES {stats.balance || 0}</h3>
          <p>Account Balance</p>
        </div>
      </div>

      <div class="stat-card error">
        <div class="stat-icon">❌</div>
        <div class="stat-content">
          <h3>{stats.total_failed || 0}</h3>
          <p>Failed</p>
        </div>
      </div>

      <div class="stat-card warning">
        <div class="stat-icon">↩️</div>
        <div class="stat-content">
          <h3>{stats.total_bounced || 0}</h3>
          <p>Bounced</p>
        </div>
      </div>
    </div>

    <div class="quota-info">
      <h2>Monthly Quota</h2>
      <div class="progress-bar">
        <div 
          class="progress-fill" 
          style="width: {(stats.emails_sent / stats.email_quota * 100)}%"
        ></div>
      </div>
      <p>{stats.emails_sent} / {stats.email_quota} emails used ({Math.round(stats.emails_sent / stats.email_quota * 100)}%)</p>
      <p class="info-text">Quota resets on the 1st of each month. You get 100 free emails monthly!</p>
    </div>
  {/if}
</div>

<style>
  .dashboard {
    max-width: 1200px;
  }

  h1 {
    color: #2c3e50;
    margin-bottom: 30px;
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 20px;
    margin-bottom: 40px;
  }

  .stat-card {
    background: white;
    padding: 25px;
    border-radius: 10px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
    display: flex;
    align-items: center;
    gap: 20px;
  }

  .stat-card.error {
    border-left: 4px solid #e74c3c;
  }

  .stat-card.warning {
    border-left: 4px solid #f39c12;
  }

  .stat-icon {
    font-size: 48px;
  }

  .stat-content h3 {
    margin: 0;
    font-size: 32px;
    color: #2c3e50;
  }

  .stat-content p {
    margin: 5px 0 0 0;
    color: #7f8c8d;
    font-size: 14px;
  }

  .quota-info {
    background: white;
    padding: 30px;
    border-radius: 10px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  }

  .quota-info h2 {
    margin: 0 0 20px 0;
    color: #2c3e50;
  }

  .progress-bar {
    width: 100%;
    height: 20px;
    background: #ecf0f1;
    border-radius: 10px;
    overflow: hidden;
    margin-bottom: 10px;
  }

  .progress-fill {
    height: 100%;
    background: linear-gradient(90deg, #667eea, #764ba2);
    transition: width 0.3s ease;
  }

  .quota-info p {
    margin: 5px 0;
    color: #7f8c8d;
  }

  .info-text {
    margin-top: 15px;
    padding: 10px;
    background: #e8f4f8;
    border-left: 4px solid #3498db;
    border-radius: 5px;
  }
</style>
