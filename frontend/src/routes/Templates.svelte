<script>
  import { onMount } from 'svelte';
  import api from '../lib/api';

  let templates = [];
  let loading = true;
  let showModal = false;
  let editingTemplate = null;

  let form = {
    name: '',
    subject: '',
    html_body: '',
    text_body: '',
    variables: {}
  };

  onMount(loadTemplates);

  async function loadTemplates() {
    try {
      const response = await api.getTemplates();
      templates = response.data || [];
    } catch (error) {
      console.error('Failed to load templates:', error);
    } finally {
      loading = false;
    }
  }

  function openModal(template = null) {
    if (template) {
      editingTemplate = template;
      form = { ...template };
    } else {
      editingTemplate = null;
      form = { name: '', subject: '', html_body: '', text_body: '', variables: {} };
    }
    showModal = true;
  }

  async function handleSubmit() {
    try {
      if (editingTemplate) {
        await api.updateTemplate(editingTemplate.id, form);
      } else {
        await api.createTemplate(form);
      }
      showModal = false;
      await loadTemplates();
    } catch (error) {
      alert('Failed to save template');
    }
  }

  async function deleteTemplate(id) {
    if (confirm('Are you sure you want to delete this template?')) {
      try {
        await api.deleteTemplate(id);
        await loadTemplates();
      } catch (error) {
        alert('Failed to delete template');
      }
    }
  }
</script>

<div class="templates">
  <div class="header">
    <h1>Email Templates</h1>
    <button class="btn-primary" on:click={() => openModal()}>Create Template</button>
  </div>

  {#if loading}
    <p>Loading...</p>
  {:else if templates.length === 0}
    <div class="empty-state">
      <p>No templates yet. Create your first template!</p>
    </div>
  {:else}
    <div class="templates-grid">
      {#each templates as template}
        <div class="template-card">
          <h3>{template.name}</h3>
          <p class="subject">{template.subject}</p>
          <p class="date">Created: {new Date(template.created_at).toLocaleDateString()}</p>
          <div class="actions">
            <button on:click={() => openModal(template)}>Edit</button>
            <button class="delete" on:click={() => deleteTemplate(template.id)}>Delete</button>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

{#if showModal}
  <div class="modal-overlay" on:click={() => showModal = false}>
    <div class="modal" on:click|stopPropagation>
      <h2>{editingTemplate ? 'Edit' : 'Create'} Template</h2>
      <form on:submit|preventDefault={handleSubmit}>
        <div class="form-group">
          <label>Name</label>
          <input type="text" bind:value={form.name} required />
        </div>
        <div class="form-group">
          <label>Subject</label>
          <input type="text" bind:value={form.subject} required />
        </div>
        <div class="form-group">
          <label>HTML Body</label>
          <textarea bind:value={form.html_body} rows="8" required></textarea>
        </div>
        <div class="form-group">
          <label>Text Body (optional)</label>
          <textarea bind:value={form.text_body} rows="4"></textarea>
        </div>
        <div class="modal-actions">
          <button type="button" on:click={() => showModal = false}>Cancel</button>
          <button type="submit" class="btn-primary">Save</button>
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

  h1 {
    color: #2c3e50;
    margin: 0;
  }

  .btn-primary {
    padding: 10px 20px;
    background: #667eea;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    font-size: 14px;
  }

  .btn-primary:hover {
    background: #5568d3;
  }

  .templates-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 20px;
  }

  .template-card {
    background: white;
    padding: 20px;
    border-radius: 10px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  }

  .template-card h3 {
    margin: 0 0 10px 0;
    color: #2c3e50;
  }

  .subject {
    color: #7f8c8d;
    font-size: 14px;
    margin: 5px 0;
  }

  .date {
    color: #95a5a6;
    font-size: 12px;
    margin: 10px 0;
  }

  .actions {
    display: flex;
    gap: 10px;
    margin-top: 15px;
  }

  .actions button {
    flex: 1;
    padding: 8px;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    font-size: 14px;
  }

  .actions button:first-child {
    background: #3498db;
    color: white;
  }

  .actions button.delete {
    background: #e74c3c;
    color: white;
  }

  .empty-state {
    text-align: center;
    padding: 60px 20px;
    color: #7f8c8d;
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
    max-width: 600px;
    max-height: 90vh;
    overflow-y: auto;
  }

  .modal h2 {
    margin: 0 0 20px 0;
    color: #2c3e50;
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

  .form-group input,
  .form-group textarea {
    width: 100%;
    padding: 10px;
    border: 1px solid #dee2e6;
    border-radius: 5px;
    font-size: 14px;
    box-sizing: border-box;
    font-family: inherit;
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
    font-size: 14px;
  }

  .modal-actions button:first-child {
    background: #95a5a6;
    color: white;
  }
</style>
