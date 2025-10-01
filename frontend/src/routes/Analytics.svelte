<script>
  import { onMount } from 'svelte';
  import api from '../lib/api';

  let analytics = [];
  let loading = true;

  onMount(loadAnalytics);

  async function loadAnalytics() {
    try {
      const endDate = new Date().toISOString().split('T')[0];
      const startDate = new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
      const response = await api.getAnalytics({ start_date: startDate, end_date: endDate });
      analytics = response.data || [];
    } catch (error) {
      console.error('Failed to load analytics:', error);
    } finally {
      loading = false;
    }
  }

  function getTotalSent() {
    return analytics.reduce((sum, a) => sum + a.emails_sent, 0);
  }

  function getTotalFailed() {
    return analytics.reduce((sum, a) => sum + a.failed, 0);
  }

  function getTotalBounced() {
    return analytics.reduce((sum, a) => sum + a.bounced, 0);
  }
</script>

<div class="analytics">
  <h1>Analytics</h1>

  {#if loading}
    <p>Loading...</p>
  {:else}
    <div class="summary-cards">
      <div class="summary-card">
        <h3>{getTotalSent()}</h3>
        <p>Total Sent (30 days)</p>
      </div>
      <div class="summary-card">
        <h3>{getTotalFailed()}</h3>
        <p>Total Failed</p>
      </div>
      <div class="summary-card">
        <h3>{getTotalBounced()}</h3>
        <p>Total Bounced</p>
      </div>
      <div class="summary-card">
        <h3>{analytics.length > 0 ? Math.round((getTotalSent() / (getTotalSent() + getTotalFailed())) * 100) : 0}%</h3>
        <p>Success Rate</p>
      </div>
    </div>

    <div class="chart-container">
      <h2>Daily Email Activity</h2>
      <div class="chart">
        {#each analytics.slice().reverse() as record}
          <div class="bar-group">
            <div class="bar sent" style="height: {(record.emails_sent / Math.max(...analytics.map(a => a.emails_sent), 1)) * 200}px"></div>
            <span class="label">{new Date(record.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}</span>
          </div>
        {/each}
      </div>
    </div>
  {/if}
</div>

<style>
  h1 {
    color: #2c3e50;
    margin-bottom: 30px;
  }

  .summary-cards {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 20px;
    margin-bottom: 40px;
  }

  .summary-card {
    background: white;
    padding: 25px;
    border-radius: 10px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
    text-align: center;
  }

  .summary-card h3 {
    margin: 0;
    font-size: 36px;
    color: #667eea;
  }

  .summary-card p {
    margin: 10px 0 0 0;
    color: #7f8c8d;
    font-size: 14px;
  }

  .chart-container {
    background: white;
    padding: 30px;
    border-radius: 10px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  }

  .chart-container h2 {
    margin: 0 0 30px 0;
    color: #2c3e50;
  }

  .chart {
    display: flex;
    align-items: flex-end;
    gap: 10px;
    height: 250px;
    padding-top: 20px;
    overflow-x: auto;
  }

  .bar-group {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 5px;
  }

  .bar {
    width: 30px;
    background: linear-gradient(to top, #667eea, #764ba2);
    border-radius: 5px 5px 0 0;
    min-height: 5px;
    transition: all 0.3s;
  }

  .bar:hover {
    opacity: 0.8;
  }

  .label {
    font-size: 10px;
    color: #7f8c8d;
    writing-mode: vertical-rl;
    text-orientation: mixed;
  }
</style>
