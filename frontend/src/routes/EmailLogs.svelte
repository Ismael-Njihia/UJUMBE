<script>
  import { onMount } from 'svelte';
  import api from '../lib/api';

  let logs = [];
  let loading = true;

  onMount(loadLogs);

  async function loadLogs() {
    try {
      const response = await api.getEmailLogs({ limit: 100, offset: 0 });
      logs = response.data || [];
    } catch (error) {
      console.error('Failed to load logs:', error);
    } finally {
      loading = false;
    }
  }

  function getStatusColor(status) {
    switch(status) {
      case 'sent': return '#27ae60';
      case 'queued': return '#3498db';
      case 'failed': return '#e74c3c';
      case 'bounced': return '#f39c12';
      default: return '#95a5a6';
    }
  }
</script>

<div class="email-logs">
  <div class="header">
    <h1>Email Logs</h1>
    <button on:click={loadLogs}>🔄 Refresh</button>
  </div>

  {#if loading}
    <p>Loading...</p>
  {:else if logs.length === 0}
    <div class="empty-state">
      <p>No email logs yet.</p>
    </div>
  {:else}
    <div class="logs-table">
      <table>
        <thead>
          <tr>
            <th>Date</th>
            <th>From</th>
            <th>To</th>
            <th>Subject</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {#each logs as log}
            <tr>
              <td>{new Date(log.created_at).toLocaleString()}</td>
              <td>{log.from_email}</td>
              <td>{log.to_email}</td>
              <td>{log.subject}</td>
              <td>
                <span class="status-badge" style="background: {getStatusColor(log.status)}">
                  {log.status}
                </span>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

<style>
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 30px;
  }

  .header button {
    padding: 10px 20px;
    background: #3498db;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
  }

  .logs-table {
    background: white;
    border-radius: 10px;
    overflow: hidden;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  thead {
    background: #f8f9fa;
  }

  th {
    padding: 15px;
    text-align: left;
    font-weight: 600;
    color: #2c3e50;
    font-size: 14px;
  }

  td {
    padding: 15px;
    border-top: 1px solid #dee2e6;
    font-size: 14px;
    color: #495057;
  }

  tr:hover {
    background: #f8f9fa;
  }

  .status-badge {
    padding: 5px 12px;
    border-radius: 15px;
    color: white;
    font-size: 12px;
    font-weight: 500;
  }

  .empty-state {
    text-align: center;
    padding: 60px 20px;
    color: #7f8c8d;
  }
</style>
