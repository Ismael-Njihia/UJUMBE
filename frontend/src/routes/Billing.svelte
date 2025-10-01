<script>
  import { onMount } from 'svelte';
  import api from '../lib/api';

  let transactions = [];
  let loading = true;
  let showModal = false;
  let topupForm = {
    amount: '',
    phone_number: ''
  };

  onMount(loadTransactions);

  async function loadTransactions() {
    try {
      const response = await api.getTransactions({ limit: 50, offset: 0 });
      transactions = response.data || [];
    } catch (error) {
      console.error('Failed to load transactions:', error);
    } finally {
      loading = false;
    }
  }

  async function handleTopup() {
    try {
      const response = await api.initiateTopup({
        amount: parseFloat(topupForm.amount),
        phone_number: topupForm.phone_number
      });
      alert(response.data.message);
      showModal = false;
      topupForm = { amount: '', phone_number: '' };
      await loadTransactions();
    } catch (error) {
      alert('Failed to initiate top-up');
    }
  }
</script>

<div class="billing">
  <div class="header">
    <h1>Billing & Transactions</h1>
    <button class="btn-primary" on:click={() => showModal = true}>💳 Top Up via M-Pesa</button>
  </div>

  <div class="pricing-info">
    <h2>Pricing</h2>
    <p>💌 100 free emails per month</p>
    <p>💰 KES 1.00 per email after free tier</p>
    <p>📱 Easy top-up with M-Pesa</p>
  </div>

  <h2>Transaction History</h2>
  {#if loading}
    <p>Loading...</p>
  {:else if transactions.length === 0}
    <div class="empty-state">
      <p>No transactions yet.</p>
    </div>
  {:else}
    <div class="transactions-table">
      <table>
        <thead>
          <tr>
            <th>Date</th>
            <th>Type</th>
            <th>Amount</th>
            <th>Emails</th>
            <th>Status</th>
            <th>Receipt</th>
          </tr>
        </thead>
        <tbody>
          {#each transactions as transaction}
            <tr>
              <td>{new Date(transaction.created_at).toLocaleString()}</td>
              <td>{transaction.type}</td>
              <td>{transaction.currency} {transaction.amount.toFixed(2)}</td>
              <td>{transaction.emails_purchased}</td>
              <td>
                <span class="status" class:completed={transaction.status === 'completed'}>
                  {transaction.status}
                </span>
              </td>
              <td>{transaction.mpesa_receipt_no || '-'}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

{#if showModal}
  <div class="modal-overlay" on:click={() => showModal = false}>
    <div class="modal" on:click|stopPropagation>
      <h2>Top Up with M-Pesa</h2>
      <form on:submit|preventDefault={handleTopup}>
        <div class="form-group">
          <label>Amount (KES)</label>
          <input type="number" bind:value={topupForm.amount} min="10" step="1" required />
          <small>Minimum: KES 10</small>
        </div>
        <div class="form-group">
          <label>M-Pesa Phone Number</label>
          <input type="tel" bind:value={topupForm.phone_number} placeholder="254712345678" required />
          <small>Format: 254XXXXXXXXX</small>
        </div>
        <div class="modal-actions">
          <button type="button" on:click={() => showModal = false}>Cancel</button>
          <button type="submit" class="btn-primary">Initiate Payment</button>
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
    background: #27ae60;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
  }

  .pricing-info {
    background: white;
    padding: 25px;
    border-radius: 10px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
    margin-bottom: 30px;
  }

  .pricing-info h2 {
    margin: 0 0 15px 0;
    color: #2c3e50;
  }

  .pricing-info p {
    margin: 8px 0;
    font-size: 16px;
    color: #495057;
  }

  h2 {
    color: #2c3e50;
    margin-bottom: 20px;
  }

  .transactions-table {
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

  .status {
    padding: 5px 12px;
    border-radius: 15px;
    font-size: 12px;
    background: #f39c12;
    color: white;
  }

  .status.completed {
    background: #27ae60;
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

  .form-group small {
    display: block;
    margin-top: 5px;
    color: #7f8c8d;
    font-size: 12px;
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

  .empty-state {
    text-align: center;
    padding: 60px 20px;
    color: #7f8c8d;
  }
</style>
